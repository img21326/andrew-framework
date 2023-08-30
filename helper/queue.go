package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue struct{}

var queue *Queue = &Queue{}

func GetQueueInstance() *Queue {
	return queue
}

type Job struct {
	RetryCount  int         `json:"retry_count"`
	MaxRetry    int         `json:"max_retry"`
	LastRunTime int64       `json:"last_run_time"`
	JobType     string      `json:"job_type"`
	JobData     interface{} `json:"job_data"`
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
	return GetRedisInstance().LPush(ctx, "job_queue", job).Err()
}

func (q *Queue) Work(ctx context.Context) error {
	var job Job
	err := GetRedisInstance().LPop(ctx, "job_queue").Scan(&job)
	if err != nil {
		if err == redis.Nil {
			time.Sleep(5 * time.Second)
			return nil
		}
		return err
	}
	err = jobWorkMap[job.JobType](job)
	if err != nil {
		job.RetryCount++
		if job.RetryCount < job.MaxRetry {
			return q.PushJob(ctx, job)
		}
	}
	return nil
}

func StartWorker(ctx context.Context) {
	logger := NewLogger("queue")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := GetQueueInstance().Work(ctx)
			if err != nil {
				logger.Error(ctx, "work error: %s", err.Error())
			}
		}
	}
}
