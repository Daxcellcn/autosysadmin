// backend/internal/subscriptions/repository.go
package subscriptions

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Repository interface {
	Save(ctx context.Context, subscription *Subscription) error
	Get(ctx context.Context, id string) (*Subscription, error)
	ListByUser(ctx context.Context, userID string) ([]Subscription, error)
}

type inMemoryRepository struct {
	subscriptions map[string]Subscription
	mu            sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		subscriptions: make(map[string]Subscription),
	}
}

func (r *inMemoryRepository) Save(ctx context.Context, subscription *Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscriptions[subscription.ID] = *subscription
	return nil
}

func (r *inMemoryRepository) Get(ctx context.Context, id string) (*Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, exists := r.subscriptions[id]
	if !exists {
		return nil, errors.New("subscription not found")
	}

	return &sub, nil
}

func (r *inMemoryRepository) ListByUser(ctx context.Context, userID string) ([]Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userSubs []Subscription
	for _, sub := range r.subscriptions {
		if sub.UserID == userID {
			userSubs = append(userSubs, sub)
		}
	}

	return userSubs, nil
}