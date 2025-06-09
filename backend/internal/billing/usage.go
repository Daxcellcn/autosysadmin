// backend/internal/billing/usage.go
package billing

import (
	"context"
	"time"
)

type UsageRecord struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	PlanID         string    `json:"plan_id"`
	Timestamp      time.Time `json:"timestamp"`
	ServerCount    int       `json:"server_count"`
	CPUHours       float64   `json:"cpu_hours"`
	MemoryGBHours  float64   `json:"memory_gb_hours"`
	NetworkGB      float64   `json:"network_gb"`
	StorageGB      float64   `json:"storage_gb"`
}

type UsageRepository interface {
	Save(ctx context.Context, record *UsageRecord) error
	GetForSubscription(ctx context.Context, subscriptionID string, start, end time.Time) ([]UsageRecord, error)
	Aggregate(ctx context.Context, subscriptionID string, start, end time.Time) (*UsageSummary, error)
}

type usageRepository struct {
	// Would be backed by database in real implementation
	records []UsageRecord
}

func NewUsageRepository() UsageRepository {
	return &usageRepository{
		records: make([]UsageRecord, 0),
	}
}

func (r *usageRepository) Save(ctx context.Context, record *UsageRecord) error {
	r.records = append(r.records, *record)
	return nil
}

func (r *usageRepository) GetForSubscription(ctx context.Context, subscriptionID string, start, end time.Time) ([]UsageRecord, error) {
	var results []UsageRecord
	for _, record := range r.records {
		if record.SubscriptionID == subscriptionID &&
			!record.Timestamp.Before(start) &&
			!record.Timestamp.After(end) {
			results = append(results, record)
		}
	}
	return results, nil
}

func (r *usageRepository) Aggregate(ctx context.Context, subscriptionID string, start, end time.Time) (*UsageSummary, error) {
	records, err := r.GetForSubscription(ctx, subscriptionID, start, end)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &UsageSummary{
			SubscriptionID: subscriptionID,
			StartDate:      start,
			EndDate:        end,
		}, nil
	}

	var summary UsageSummary
	summary.SubscriptionID = subscriptionID
	summary.StartDate = start
	summary.EndDate = end

	// Find the most recent plan ID
	summary.PlanID = records[len(records)-1].PlanID

	// Calculate averages and totals
	var totalCPU, totalMemory, totalNetwork, totalStorage float64
	serverCounts := make(map[int]bool)

	for _, record := range records {
		totalCPU += record.CPUHours
		totalMemory += record.MemoryGBHours
		totalNetwork += record.NetworkGB
		totalStorage += record.StorageGB
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
	summary.CPUHours = totalCPU
	summary.MemoryGBHours = totalMemory
	summary.NetworkGB = totalNetwork
	summary.StorageGB = totalStorage

	return &summary, nil
}