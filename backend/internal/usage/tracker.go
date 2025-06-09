// backend/internal/usage/tracker.go
package usage

import (
	"context"
	"sync"
	"time"
)

type Tracker interface {
	RecordUsage(ctx context.Context, record *UsageRecord) error
	GetUsage(ctx context.Context, subscriptionID string, start, end time.Time) (*UsageRecord, error)
	AggregateUsage(ctx context.Context, subscriptionID string, period string) ([]UsageRecord, error)
}

type UsageRecord struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	Timestamp      time.Time `json:"timestamp"`
	ServerCount    int       `json:"server_count"`
	CPUHours       float64   `json:"cpu_hours"`
	MemoryGBHours  float64   `json:"memory_gb_hours"`
	NetworkGB      float64   `json:"network_gb"`
	StorageGB      float64   `json:"storage_gb"`
}

type usageTracker struct {
	repo Repository
	mu   sync.RWMutex
}

func NewTracker(repo Repository) Tracker {
	return &usageTracker{
		repo: repo,
	}
}

func (t *usageTracker) RecordUsage(ctx context.Context, record *UsageRecord) error {
	record.ID = generateUUID()
	record.Timestamp = time.Now()
	return t.repo.Save(ctx, record)
}

func (t *usageTracker) GetUsage(ctx context.Context, subscriptionID string, start, end time.Time) (*UsageRecord, error) {
	records, err := t.repo.GetForSubscription(ctx, subscriptionID, start, end)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &UsageRecord{
			SubscriptionID: subscriptionID,
		}, nil
	}

	// Aggregate the records into a single summary
	var summary UsageRecord
	summary.SubscriptionID = subscriptionID
	serverCounts := make(map[int]bool)

	for _, record := range records {
		summary.CPUHours += record.CPUHours
		summary.MemoryGBHours += record.MemoryGBHours
		summary.NetworkGB += record.NetworkGB
		summary.StorageGB += record.StorageGB
		serverCounts[record.ServerCount] = true
	}

	// Get max server count
	maxServers := 0
	for count := range serverCounts {
		if count > maxServers {
			maxServers = count
		}
	}
	summary.ServerCount = maxServers

	return &summary, nil
}

func (t *usageTracker) AggregateUsage(ctx context.Context, subscriptionID string, period string) ([]UsageRecord, error) {
	var start, end time.Time
	now := time.Now()

	switch period {
	case "day":
		start = now.AddDate(0, 0, -1)
		end = now
	case "week":
		start = now.AddDate(0, 0, -7)
		end = now
	case "month":
		start = now.AddDate(0, -1, 0)
		end = now
	default:
		return nil, errors.New("invalid period")
	}

	return t.repo.GetForSubscription(ctx, subscriptionID, start, end)
}

func generateUUID() string {
	// In production, use github.com/google/uuid
	return "generated-uuid"
}