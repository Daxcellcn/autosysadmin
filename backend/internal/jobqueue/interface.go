// backend/internal/jobqueue/interface.go
package jobqueue

import "context"

type JobQueue interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (*Job, error)
	CompleteJob(ctx context.Context, jobID string, result string) error
	FailJob(ctx context.Context, jobID string, errorMsg string) error
	GetJob(ctx context.Context, jobID string) (*Job, error)
	Close() error
}