// backend/internal/monitoring/metrics.go
package monitoring

import (
	"context"
	"sync"
	"time"
)

type MetricsCollector interface {
	Collect(agentID string) ([]Metric, error)
}

type MetricsStorage interface {
	Store(metrics []Metric) error
	Query(agentID string, start, end time.Time) ([]Metric, error)
}

type PrometheusMetricsCollector struct {
	// Configuration for Prometheus client
}

func (c *PrometheusMetricsCollector) Collect(agentID string) ([]Metric, error) {
	// In a real implementation, this would query Prometheus
	return []Metric{
		{Name: "cpu", Value: 25.5, Timestamp: time.Now()},
		{Name: "memory", Value: 45.2, Timestamp: time.Now()},
	}, nil
}

type InMemoryMetricsStorage struct {
	metrics map[string][]Metric // agentID -> metrics
	mu      sync.RWMutex
}

func NewInMemoryMetricsStorage() *InMemoryMetricsStorage {
	return &InMemoryMetricsStorage{
		metrics: make(map[string][]Metric),
	}
}

func (s *InMemoryMetricsStorage) Store(metrics []Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, metric := range metrics {
		s.metrics[metric.AgentID] = append(s.metrics[metric.AgentID], metric)
	}
	return nil
}

func (s *InMemoryMetricsStorage) Query(agentID string, start, end time.Time) ([]Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agentMetrics, exists := s.metrics[agentID]
	if !exists {
		return nil, fmt.Errorf("no metrics found for agent %s", agentID)
	}

	var result []Metric
	for _, metric := range agentMetrics {
		if (metric.Timestamp.After(start) || metric.Timestamp.Equal(start)) &&
			(metric.Timestamp.Before(end) || metric.Timestamp.Equal(end)) {
			result = append(result, metric)
		}
	}

	return result, nil
}

type MetricsService struct {
	collector MetricsCollector
	storage   MetricsStorage
}

func NewMetricsService(collector MetricsCollector, storage MetricsStorage) *MetricsService {
	return &MetricsService{
		collector: collector,
		storage:   storage,
	}
}

func (s *MetricsService) CollectAndStore(agentID string) error {
	metrics, err := s.collector.Collect(agentID)
	if err != nil {
		return err
	}
	return s.storage.Store(metrics)
}

func (s *MetricsService) QueryMetrics(agentID string, start, end time.Time) ([]Metric, error) {
	return s.storage.Query(agentID, start, end)
}