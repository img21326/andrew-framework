package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var queueLogger *Logger

type Queue struct{}

var queue *Queue = &Queue{}

func GetQueueInstance() *Queue {
	if queueLogger == nil {
		queueLogger = NewLogger("queue")
	}
	return queue
}

type Job struct {
	RetryCount  int         `json:"retry_count"`
	MaxRetry    int         `json:"max_retry"`
	LastRunTime *time.Time  `json:"last_run_time"`
	JobType     string      `json:"job_type"`
	JobDataRaw  string      `json:"job_data_raw"`
	JobData     interface{} `json:"-"`
}

func NewJob(jobType string, jobData interface{}) Job {
	return Job{
		JobType: jobType,
		JobData: jobData,
	}
}

var jobWorkMap = map[string]func(job Job) error{}

func (q *Queue) RegisterJobWork(jobType string, work func(job Job) error) {
	jobWorkMap[jobType] = work
}

func (q *Queue) PushJob(ctx context.Context, job Job) error {
	var find bool = false
	for k := range jobWorkMap {
		if k == job.JobType {
			find = true
			break
		}
	}
	if !find {
		return fmt.Errorf("job type %s not found", job.JobType)
	}
	data, _ := json.Marshal(job.JobData)
	job.JobDataRaw = string(data)
	raw, _ := json.Marshal(job)
	return GetRedisInstance().LPush(ctx, "job_queue", raw).Err()
}

func (q *Queue) Work(ctx context.Context) error {
	var jobString string
	err := GetRedisInstance().LPop(ctx, "job_queue").Scan(&jobString)
	if err != nil {
		if err == redis.Nil {
			time.Sleep(5 * time.Second)
			return nil
		}
		return err
	}
	var job Job
	err = json.Unmarshal([]byte(jobString), &job)
	if err != nil {
		return err
	}
	err = jobWorkMap[job.JobType](job)
	if err != nil {
		now := time.Now()
		job.LastRunTime = &now
		if job.MaxRetry == 0 {
			job.MaxRetry = 5
		}
		job.RetryCount++
		if job.RetryCount < job.MaxRetry {
			return q.PushJob(ctx, job)
		} else {
			queueLogger.Error(ctx, "job: %+v, err: %s", job, err.Error())
		}
	}
	return nil
}

func StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := GetQueueInstance().Work(ctx)
			if err != nil {
				queueLogger.Error(ctx, "work error: %s", err.Error())
			}
		}
	}
}
