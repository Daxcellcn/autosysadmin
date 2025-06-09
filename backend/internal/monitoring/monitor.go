// backend/internal/monitoring/monitor.go
package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/autosysadmin/backend/internal/agent"
)

type Monitor interface {
	StartMonitoring(agentID string, interval time.Duration) error
	StopMonitoring(agentID string) error
	GetMetrics(agentID string) ([]Metric, error)
	GetAlerts(agentID string) ([]Alert, error)
	SetAlertThreshold(agentID, metric string, threshold float64) error
}

type Metric struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Alert struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"` // active, resolved
}

type monitor struct {
	agentManager *agent.Manager
	metrics      map[string][]Metric // agentID -> metrics
	alerts       map[string][]Alert // agentID -> alerts
	thresholds   map[string]map[string]float64 // agentID -> metric -> threshold
	mu           sync.RWMutex
	cancelFuncs  map[string]context.CancelFunc // agentID -> cancelFunc
}

func NewMonitor() Monitor {
	return &monitor{
		metrics:     make(map[string][]Metric),
		alerts:      make(map[string][]Alert),
		thresholds:  make(map[string]map[string]float64),
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

func (m *monitor) StartMonitoring(agentID string, interval time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.cancelFuncs[agentID]; exists {
		return fmt.Errorf("monitoring already started for agent %s", agentID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFuncs[agentID] = cancel

	go m.monitorAgent(ctx, agentID, interval)
	return nil
}

func (m *monitor) StopMonitoring(agentID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cancel, exists := m.cancelFuncs[agentID]
	if !exists {
		return fmt.Errorf("no monitoring active for agent %s", agentID)
	}

	cancel()
	delete(m.cancelFuncs, agentID)
	return nil
}

func (m *monitor) monitorAgent(ctx context.Context, agentID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			agent, exists := m.agentManager.GetAgent(agentID)
			if !exists {
				continue
			}

			stats, err := agent.CollectStats()
			if err != nil {
				continue
			}

			// Record metrics
			metrics := []Metric{
				{Name: "cpu", Value: stats.System.CPUUsage, Timestamp: time.Now()},
				{Name: "memory", Value: stats.System.MemoryUsage, Timestamp: time.Now()},
				{Name: "disk", Value: stats.System.DiskUsage, Timestamp: time.Now()},
				{Name: "network_in", Value: stats.System.NetworkIn, Timestamp: time.Now()},
				{Name: "network_out", Value: stats.System.NetworkOut, Timestamp: time.Now()},
			}

			m.mu.Lock()
			m.metrics[agentID] = append(m.metrics[agentID], metrics...)
			
			// Check thresholds and generate alerts
			for _, metric := range metrics {
				if thresholds, ok := m.thresholds[agentID]; ok {
					if threshold, ok := thresholds[metric.Name]; ok {
						if metric.Value > threshold {
							alert := Alert{
								ID:        fmt.Sprintf("%s-%s-%d", agentID, metric.Name, time.Now().Unix()),
								AgentID:   agentID,
								Metric:    metric.Name,
								Value:     metric.Value,
								Threshold: threshold,
								Message:   fmt.Sprintf("%s exceeded threshold (%.2f > %.2f)", metric.Name, metric.Value, threshold),
								Timestamp: time.Now(),
								Status:    "active",
							}
							m.alerts[agentID] = append(m.alerts[agentID], alert)
						}
					}
				}
			}
			m.mu.Unlock()
		}
	}
}

func (m *monitor) GetMetrics(agentID string) ([]Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.metrics[agentID]
	if !exists {
		return nil, fmt.Errorf("no metrics found for agent %s", agentID)
	}

	return metrics, nil
}

func (m *monitor) GetAlerts(agentID string) ([]Alert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	alerts, exists := m.alerts[agentID]
	if !exists {
		return nil, fmt.Errorf("no alerts found for agent %s", agentID)
	}

	return alerts, nil
}

func (m *monitor) SetAlertThreshold(agentID, metric string, threshold float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.thresholds[agentID]; !exists {
		m.thresholds[agentID] = make(map[string]float64)
	}

	m.thresholds[agentID][metric] = threshold
	return nil
}