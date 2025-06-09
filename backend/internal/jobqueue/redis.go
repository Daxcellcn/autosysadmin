// backend/internal/jobqueue/redis.go
package jobqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisJobQueue struct {
	client *redis.Client
	prefix string
}

type Job struct {
	ID        string        `json:"id"`
	AgentID   string        `json:"agent_id"`
	Command   string        `json:"command"`
	Args      []string      `json:"args"`
	Timeout   time.Duration `json:"timeout"`
	CreatedAt time.Time     `json:"created_at"`
	Status    string        `json:"status"` // queued, running, completed, failed
	Result    string        `json:"result"`
}

func NewRedisJobQueue(redisAddr string, prefix string) *RedisJobQueue {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisJobQueue{
		client: client,
		prefix: prefix,
	}
}

func (q *RedisJobQueue) Enqueue(ctx context.Context, job Job) error {
	job.Status = "queued"
	jobKey := fmt.Sprintf("%s:jobs:%s", q.prefix, job.ID)
	queueKey := fmt.Sprintf("%s:queue", q.prefix)

	// Serialize job
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Store job in hash and add to queue
	pipe := q.client.TxPipeline()
	pipe.HSet(ctx, jobKey, "data", jobData)
	pipe.LPush(ctx, queueKey, job.ID)
	_, err = pipe.Exec(ctx)
	return err
}

func (q *RedisJobQueue) Dequeue(ctx context.Context) (*Job, error) {
	queueKey := fmt.Sprintf("%s:queue", q.prefix)
	jobID, err := q.client.RPop(ctx, queueKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // No jobs available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	jobKey := fmt.Sprintf("%s:jobs:%s", q.prefix, jobID)
	data, err := q.client.HGet(ctx, jobKey, "data").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Update job status to running
	job.Status = "running"
	jobData, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated job: %w", err)
	}

	if err := q.client.HSet(ctx, jobKey, "data", jobData).Err(); err != nil {
		return nil, fmt.Errorf("failed to update job status: %w", err)
	}

	return &job, nil
}

func (q *RedisJobQueue) CompleteJob(ctx context.Context, jobID string, result string) error {
	jobKey := fmt.Sprintf("%s:jobs:%s", q.prefix, jobID)
	data, err := q.client.HGet(ctx, jobKey, "data").Result()
	if err != nil {
		return fmt.Errorf("failed to get job data: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	job.Status = "completed"
	job.Result = result
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal completed job: %w", err)
	}

	return q.client.HSet(ctx, jobKey, "data", jobData).Err()
}

func (q *RedisJobQueue) FailJob(ctx context.Context, jobID string, errorMsg string) error {
	jobKey := fmt.Sprintf("%s:jobs:%s", q.prefix, jobID)
	data, err := q.client.HGet(ctx, jobKey, "data").Result()
	if err != nil {
		return fmt.Errorf("failed to get job data: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	job.Status = "failed"
	job.Result = errorMsg
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal failed job: %w", err)
	}

	return q.client.HSet(ctx, jobKey, "data", jobData).Err()
}

func (q *RedisJobQueue) GetJob(ctx context.Context, jobID string) (*Job, error) {
	jobKey := fmt.Sprintf("%s:jobs:%s", q.prefix, jobID)
	data, err := q.client.HGet(ctx, jobKey, "data").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

func (q *RedisJobQueue) Close() error {
	return q.client.Close()
}