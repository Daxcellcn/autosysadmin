// backend/internal/agent/manager.go
package agent

import (
	"context"
	"sync"
	"time"

	"github.com/autosysadmin/backend/internal/jobqueue"
)

type Manager struct {
	agents  map[string]*Agent
	mu      sync.RWMutex
	queue   jobqueue.JobQueue
	timeout time.Duration
}

func NewManager(queue jobqueue.JobQueue) *Manager {
	return &Manager{
		agents:  make(map[string]*Agent),
		queue:   queue,
		timeout: 30 * time.Second,
	}
}

func (m *Manager) RegisterAgent(agent *Agent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.agents[agent.ID] = agent
}

func (m *Manager) GetAgent(id string) (*Agent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, exists := m.agents[id]
	return agent, exists
}

func (m *Manager) ListAgents() []*Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agents := make([]*Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents
}

func (m *Manager) RunCommandOnAgent(ctx context.Context, agentID string, cmd AgentCommand) ([]byte, error) {
	agent, exists := m.GetAgent(agentID)
	if !exists {
		return nil, fmt.Errorf("agent not found")
	}

	return agent.ExecuteCommand(ctx, cmd, m.queue)
}

func (m *Manager) RunCommandOnAll(ctx context.Context, cmd AgentCommand) map[string]interface{} {
	results := make(map[string]interface{})
	agents := m.ListAgents()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, agent := range agents {
		wg.Add(1)
		go func(a *Agent) {
			defer wg.Done()
			result, err := a.ExecuteCommand(ctx, cmd, m.queue)
			mu.Lock()
			if err != nil {
				results[a.ID] = map[string]interface{}{
					"error": err.Error(),
				}
			} else {
				results[a.ID] = map[string]interface{}{
					"result": string(result),
				}
			}
			mu.Unlock()
		}(agent)
	}

	wg.Wait()
	return results
}