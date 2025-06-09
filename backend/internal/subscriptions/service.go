// backend/internal/subscriptions/service.go
package subscriptions

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	Create(ctx context.Context, subscription *Subscription) error
	Get(ctx context.Context, id string) (*Subscription, error)
	Cancel(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]Subscription, error)
	UpdatePlan(ctx context.Context, id, planID string) (*Subscription, error)
}

type Subscription struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	PlanID       string    `json:"plan_id"`
	Status       string    `json:"status"` // active, canceled, expired
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	RenewalDate  time.Time `json:"renewal_date"`
	PaymentMethodID string `json:"payment_method_id"`
}

type subscriptionService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &subscriptionService{
		repo: repo,
	}
}

func (s *subscriptionService) Create(ctx context.Context, subscription *Subscription) error {
	if subscription.UserID == "" {
		return errors.New("user ID is required")
	}
	if subscription.PlanID == "" {
		return errors.New("plan ID is required")
	}

	subscription.ID = generateUUID()
	subscription.Status = "active"
	subscription.StartDate = time.Now()
	subscription.EndDate = time.Now().AddDate(1, 0, 0) // 1 year
	subscription.RenewalDate = time.Now().AddDate(1, 0, 0)

	return s.repo.Save(ctx, subscription)
}

func (s *subscriptionService) Get(ctx context.Context, id string) (*Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *subscriptionService) Cancel(ctx context.Context, id string) error {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	sub.Status = "canceled"
	return s.repo.Save(ctx, sub)
}

func (s *subscriptionService) ListByUser(ctx context.Context, userID string) ([]Subscription, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *subscriptionService) UpdatePlan(ctx context.Context, id, planID string) (*Subscription, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	sub.PlanID = planID
	if err := s.repo.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func generateUUID() string {
	// In production, use github.com/google/uuid
	return "generated-uuid"
}