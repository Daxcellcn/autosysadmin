// backend/internal/agent/agent.go
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/autosysadmin/backend/internal/jobqueue"
)

type Agent struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Hostname      string    `json:"hostname"`
	IPAddress     string    `json:"ip_address"`
	OS            string    `json:"os"`
	Architecture  string    `json:"architecture"`
	Version       string    `json:"version"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Status        string    `json:"status"` // online, offline, degraded
	Tags          []string  `json:"tags"`
}

type AgentCommand struct {
	Command string        `json:"command"`
	Args    []string      `json:"args"`
	Timeout time.Duration `json:"timeout"`
}

func (a *Agent) ExecuteCommand(ctx context.Context, cmd AgentCommand, queue jobqueue.JobQueue) ([]byte, error) {
	job := jobqueue.Job{
		ID:        fmt.Sprintf("%s-%d", a.ID, time.Now().Unix()),
		AgentID:   a.ID,
		Command:   cmd.Command,
		Args:      cmd.Args,
		Timeout:   cmd.Timeout,
		CreatedAt: time.Now(),
	}

	if err := queue.Enqueue(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	// In a real implementation, we'd wait for the response
	// This is simplified for demonstration
	return json.Marshal(map[string]interface{}{
		"job_id":    job.ID,
		"status":    "queued",
		"timestamp": time.Now(),
	})
}