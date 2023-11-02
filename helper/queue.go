package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
	JobID       string      `json:"job_id"`
	RetryCount  int         `json:"retry_count"`
	MaxRetry    int         `json:"max_retry"`
	LastRunTime *time.Time  `json:"last_run_time"`
	JobType     string      `json:"job_type"`
	JobDataRaw  string      `json:"job_data_raw"`
	JobData     interface{} `json:"-"`
}

func NewJob(jobType string, jobData interface{}) Job {
	id := uuid.NewString()
	return Job{
		JobID:   id,
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
		panic(fmt.Sprintf("job type %s not found", job.JobType))
	}
	if job.JobDataRaw != "" {
		data, _ := json.Marshal(job.JobData)
		job.JobDataRaw = string(data)
	}
	raw, _ := json.Marshal(job)
	return GetRedisInstance().LPush(ctx, "job_queue", raw).Err()
}

func (q *Queue) Work(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			queueLogger.Error(ctx, "panic: %v", err)
			queueLogger.Error(ctx, "stack: %s", string(debug.Stack()))
		}
	}()
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
	queueLogger.Info(ctx, "strating job id: %s", job.JobID)
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
			if mailHelper := GetEmailHelper(); mailHelper != nil && gin.Mode() == gin.ReleaseMode {
				viper := viper.GetViper()
				adminEmail := viper.GetStringSlice("ADMIN_EMAIL")
				if len(adminEmail) == 0 {
					return nil
				}
				body := fmt.Sprintf("job: %+v, err: %s", job, err.Error())
				mailHelper.SendEmail(EmailSendOption{
					To:      adminEmail,
					Subject: "Job Error",
					Body:    body,
				})
			}
		}
	}
	queueLogger.Info(ctx, "job id: %s done", job.JobID)
	return nil
}

func StartWorker(ctx context.Context, done chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			done <- struct{}{}
		default:
			err := GetQueueInstance().Work(ctx)
			if err != nil {
				queueLogger.Error(ctx, "work error: %s", err.Error())
			}
		}
	}
}
