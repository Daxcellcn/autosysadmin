// backend/internal/billing/service.go
package billing

import (
	"context"
	"time"

	"github.com/autosysadmin/backend/internal/subscriptions"
	"github.com/autosysadmin/backend/internal/usage"
)

type BillingService interface {
	CreateSubscription(userID, planID string) (*Subscription, error)
	CancelSubscription(subscriptionID string) error
	GetSubscription(subscriptionID string) (*Subscription, error)
	ListPlans() ([]Plan, error)
	GetUsage(subscriptionID string, start, end time.Time) (*UsageSummary, error)
	ProcessPayment(subscriptionID string, amount float64) (*Payment, error)
}

type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PlanID    string    `json:"plan_id"`
	Status    string    `json:"status"` // active, canceled, expired
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Plan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Features    []string `json:"features"`
}

type Payment struct {
	ID            string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // pending, completed, failed
	CreatedAt     time.Time `json:"created_at"`
}

type UsageSummary struct {
	SubscriptionID string    `json:"subscription_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	ServerCount    int       `json:"server_count"`
	CPUHours       float64   `json:"cpu_hours"`
	MemoryGBHours  float64   `json:"memory_gb_hours"`
	NetworkGB      float64   `json:"network_gb"`
	StorageGB      float64   `json:"storage_gb"`
}

type billingService struct {
	subscriptionService subscriptions.Service
	usageTracker       usage.Tracker
	plans              []Plan
}

func NewBillingService(subscriptionService subscriptions.Service, usageTracker usage.Tracker) BillingService {
	// Initialize with default plans
	plans := []Plan{
		{
			ID:          "free",
			Name:        "Free Tier",
			Description: "Basic monitoring for small setups",
			Price:       0,
			Currency:    "USD",
			Features: []string{
				"5 servers max",
				"Basic monitoring (5-minute intervals)",
				"Email notifications",
				"Community support",
			},
		},
		{
			ID:          "pro",
			Name:        "Professional",
			Description: "For growing businesses",
			Price:       49,
			Currency:    "USD",
			Features: []string{
				"50 servers",
				"1-minute monitoring intervals",
				"Slack/Webhook notifications",
				"Basic automation workflows",
				"Email support",
			},
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise",
			Description: "For large scale deployments",
			Price:       299,
			Currency:    "USD",
			Features: []string{
				"Unlimited servers",
				"15-second monitoring",
				"Advanced RBAC",
				"API access",
				"Audit logging",
				"Priority support",
				"On-premise option",
			},
		},
	}

	return &billingService{
		subscriptionService: subscriptionService,
		usageTracker:       usageTracker,
		plans:              plans,
	}
}

func (s *billingService) CreateSubscription(userID, planID string) (*Subscription, error) {
	// Validate plan exists
	var selectedPlan *Plan
	for _, p := range s.plans {
		if p.ID == planID {
			selectedPlan = &p
			break
		}
	}
	if selectedPlan == nil {
		return nil, errors.New("invalid plan ID")
	}

	// Create subscription in database
	sub := &Subscription{
		ID:        generateUUID(),
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0), // 1 year
	}

	if err := s.subscriptionService.Create(context.Background(), sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *billingService) CancelSubscription(subscriptionID string) error {
	return s.subscriptionService.Cancel(context.Background(), subscriptionID)
}

func (s *billingService) GetSubscription(subscriptionID string) (*Subscription, error) {
	return s.subscriptionService.Get(context.Background(), subscriptionID)
}

func (s *billingService) ListPlans() ([]Plan, error) {
	return s.plans, nil
}

func (s *billingService) GetUsage(subscriptionID string, start, end time.Time) (*UsageSummary, error) {
	// Get usage data from tracker
	usageData, err := s.usageTracker.GetUsage(context.Background(), subscriptionID, start, end)
	if err != nil {
		return nil, err
	}

	// Convert to summary format
	summary := &UsageSummary{
		SubscriptionID: subscriptionID,
		StartDate:      start,
		EndDate:        end,
		ServerCount:    usageData.ServerCount,
		CPUHours:       usageData.CPUHours,
		MemoryGBHours:  usageData.MemoryGBHours,
		NetworkGB:      usageData.NetworkGB,
		StorageGB:      usageData.StorageGB,
	}

	return summary, nil
}

func (s *billingService) ProcessPayment(subscriptionID string, amount float64) (*Payment, error) {
	// In a real implementation, this would integrate with Stripe or similar
	payment := &Payment{
		ID:            generateUUID(),
		SubscriptionID: subscriptionID,
		Amount:        amount,
		Currency:      "USD",
		Status:        "completed",
		CreatedAt:     time.Now(),
	}

	return payment, nil
}

func generateUUID() string {
	// In production, use github.com/google/uuid
	return "generated-uuid"
}