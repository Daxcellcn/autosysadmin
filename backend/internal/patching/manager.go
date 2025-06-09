// backend/internal/patching/manager.go
package patching

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/autosysadmin/backend/internal/agent"
	"github.com/autosysadmin/backend/internal/jobqueue"
)

type PatchManager interface {
	CheckForUpdates(agentID string) ([]Update, error)
	ApplyUpdates(agentID string, updates []string) (string, error)
	GetPatchHistory(agentID string) ([]PatchRecord, error)
	SchedulePatch(agentID string, updates []string, when time.Time) (string, error)
}

type Update struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // low, medium, high, critical
	Size        int64  `json:"size"`     // in bytes
}

type PatchRecord struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Updates   []string  `json:"updates"`
	Status    string    `json:"status"` // pending, in-progress, completed, failed
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Logs      string    `json:"logs"`
}

type patchManager struct {
	agentManager *agent.Manager
	jobQueue     jobqueue.JobQueue
	updates      map[string][]Update // agentID -> updates
	history      map[string][]PatchRecord // agentID -> history
	mu           sync.RWMutex
}

func NewPatchManager() PatchManager {
	return &patchManager{
		updates: make(map[string][]Update),
		history: make(map[string][]PatchRecord),
	}
}

func (m *patchManager) CheckForUpdates(agentID string) ([]Update, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	updates, exists := m.updates[agentID]
	if !exists {
		return nil, fmt.Errorf("no updates found for agent %s", agentID)
	}

	return updates, nil
}

func (m *patchManager) ApplyUpdates(agentID string, updates []string) (string, error) {
	agent, exists := m.agentManager.GetAgent(agentID)
	if !exists {
		return "", fmt.Errorf("agent not found")
	}

	// Create patch record
	record := PatchRecord{
		ID:        fmt.Sprintf("patch-%s-%d", agentID, time.Now().Unix()),
		AgentID:   agentID,
		Updates:   updates,
		Status:    "pending",
		StartedAt: time.Now(),
	}

	m.mu.Lock()
	m.history[agentID] = append(m.history[agentID], record)
	m.mu.Unlock()

	// Create job to apply updates
	cmd := agent.AgentCommand{
		Command: "apply-updates",
		Args:    updates,
		Timeout: 30 * time.Minute,
	}

	ctx := context.Background()
	result, err := agent.ExecuteCommand(ctx, cmd, m.jobQueue)
	if err != nil {
		m.updatePatchStatus(agentID, record.ID, "failed", err.Error())
		return "", err
	}

	go m.monitorPatchJob(agentID, record.ID, string(result))
	return record.ID, nil
}

func (m *patchManager) monitorPatchJob(agentID, patchID, jobID string) {
	// In a real implementation, this would monitor the job status
	time.Sleep(5 * time.Second) // Simulate job running
	m.updatePatchStatus(agentID, patchID, "completed", "Updates applied successfully")
}

func (m *patchManager) updatePatchStatus(agentID, patchID, status, logs string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, record := range m.history[agentID] {
		if record.ID == patchID {
			m.history[agentID][i].Status = status
			m.history[agentID][i].EndedAt = time.Now()
			m.history[agentID][i].Logs = logs
			break
		}
	}
}

func (m *patchManager) GetPatchHistory(agentID string) ([]PatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	history, exists := m.history[agentID]
	if !exists {
		return nil, fmt.Errorf("no patch history found for agent %s", agentID)
	}

	return history, nil
}

func (m *patchManager) SchedulePatch(agentID string, updates []string, when time.Time) (string, error) {
	// In a real implementation, this would schedule a job for future execution
	record := PatchRecord{
		ID:        fmt.Sprintf("patch-%s-%d", agentID, time.Now().Unix()),
		AgentID:   agentID,
		Updates:   updates,
		Status:    "scheduled",
		StartedAt: when,
	}

	m.mu.Lock()
	m.history[agentID] = append(m.history[agentID], record)
	m.mu.Unlock()

	return record.ID, nil
}