// backend/internal/billing/pricing.go
package billing

import "time"

type PricingCalculator interface {
	CalculateCost(usage *UsageRecord) float64
	EstimateCost(planID string, servers int, duration time.Duration) float64
}

type pricingCalculator struct {
	plans map[string]PlanPricing
}

type PlanPricing struct {
	BasePrice       float64
	PricePerServer  float64
	IncludedServers int
}

func NewPricingCalculator() PricingCalculator {
	// Define pricing models for each plan
	plans := map[string]PlanPricing{
		"free": {
			BasePrice:       0,
			PricePerServer:  0,
			IncludedServers: 5,
		},
		"pro": {
			BasePrice:       49,
			PricePerServer:  0.5,
			IncludedServers: 10,
		},
		"enterprise": {
			BasePrice:       299,
			PricePerServer:  0,
			IncludedServers: 0, // unlimited
		},
	}

	return &pricingCalculator{
		plans: plans,
	}
}

func (c *pricingCalculator) CalculateCost(usage *UsageRecord) float64 {
	pricing, exists := c.plans[usage.PlanID]
	if !exists {
		return 0
	}

	if usage.PlanID == "enterprise" {
		return pricing.BasePrice
	}

	// Calculate cost based on server count
	serverCount := usage.ServerCount
	if serverCount <= pricing.IncludedServers {
		return pricing.BasePrice
	}

	additionalServers := serverCount - pricing.IncludedServers
	return pricing.BasePrice + (float64(additionalServers) * pricing.PricePerServer)
}

func (c *pricingCalculator) EstimateCost(planID string, servers int, duration time.Duration) float64 {
	pricing, exists := c.plans[planID]
	if !exists {
		return 0
	}

	if planID == "enterprise" {
		return pricing.BasePrice
	}

	// Calculate monthly cost
	var monthlyCost float64
	if servers <= pricing.IncludedServers {
		monthlyCost = pricing.BasePrice
	} else {
		additionalServers := servers - pricing.IncludedServers
		monthlyCost = pricing.BasePrice + (float64(additionalServers) * pricing.PricePerServer)
	}

	// Adjust for duration (assuming duration is in months)
	months := duration.Hours() / (24 * 30)
	return monthlyCost * months
}