// backend/internal/usage/repository.go
package usage

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Repository interface {
	Save(ctx context.Context, record *UsageRecord) error
	GetForSubscription(ctx context.Context, subscriptionID string, start, end time.Time) ([]UsageRecord, error)
}

type inMemoryUsageRepository struct {
	records []UsageRecord
	mu      sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryUsageRepository{
		records: make([]UsageRecord, 0),
	}
}

func (r *inMemoryUsageRepository) Save(ctx context.Context, record *UsageRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, *record)
	return nil
}

func (r *inMemoryUsageRepository) GetForSubscription(ctx context.Context, subscriptionID string, start, end time.Time) ([]UsageRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []UsageRecord
	for _, record := range r.records {
		if record.SubscriptionID == subscriptionID &&
			!record.Timestamp.Before(start) &&
			!record.Timestamp.After(end) {
			results = append(results, record)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no usage records found")
	}

	return results, nil
}