// backend/internal/agent/stats.go
package agent

import (
	"time"
)

type SystemStats struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	NetworkIn   float64   `json:"network_in"`
	NetworkOut  float64   `json:"network_out"`
	Timestamp   time.Time `json:"timestamp"`
}

type ProcessStats struct {
	PID         int     `json:"pid"`
	Name        string  `json:"name"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type AgentStats struct {
	System   SystemStats    `json:"system"`
	Processes []ProcessStats `json:"processes"`
	AgentID  string         `json:"agent_id"`
}

func (a *Agent) CollectStats() (*AgentStats, error) {
	// In a real implementation, this would collect actual stats from the system
	// This is a simplified version for demonstration
	return &AgentStats{
		System: SystemStats{
			CPUUsage:    25.5,
			MemoryUsage: 45.2,
			DiskUsage:   60.1,
			NetworkIn:   1024,
			NetworkOut:  512,
			Timestamp:   time.Now(),
		},
		Processes: []ProcessStats{
			{PID: 1, Name: "init", CPUUsage: 0.1, MemoryUsage: 0.5},
			{PID: 123, Name: "nginx", CPUUsage: 5.2, MemoryUsage: 10.3},
		},
		AgentID: a.ID,
	}, nil
}