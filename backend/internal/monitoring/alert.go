// backend/internal/monitoring/alert.go
package monitoring

import (
	"context"
	"fmt"
	"time"
)

type AlertNotifier interface {
	Notify(alert Alert) error
}

type AlertManager struct {
	notifiers []AlertNotifier
	alerts    map[string]Alert // alertID -> alert
	mu        sync.RWMutex
}

func NewAlertManager(notifiers ...AlertNotifier) *AlertManager {
	return &AlertManager{
		notifiers: notifiers,
		alerts:    make(map[string]Alert),
	}
}

func (m *AlertManager) AddAlert(alert Alert) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.alerts[alert.ID] = alert

	// Notify all notifiers
	for _, notifier := range m.notifiers {
		go func(n AlertNotifier) {
			if err := n.Notify(alert); err != nil {
				fmt.Printf("Failed to send alert notification: %v\n", err)
			}
		}(notifier)
	}
}

func (m *AlertManager) ResolveAlert(alertID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	alert, exists := m.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found")
	}

	alert.Status = "resolved"
	m.alerts[alertID] = alert
	return nil
}

func (m *AlertManager) GetActiveAlerts() []Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var activeAlerts []Alert
	for _, alert := range m.alerts {
		if alert.Status == "active" {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	return activeAlerts
}

// EmailNotifier implements AlertNotifier for email notifications
type EmailNotifier struct {
	SMTPConfig SMTPConfig
}

type SMTPConfig struct {
	Server   string
	Port     int
	Username string
	Password string
	From     string
}

func (n *EmailNotifier) Notify(alert Alert) error {
	// In a real implementation, this would send an email
	fmt.Printf("Sending email alert for %s: %s\n", alert.AgentID, alert.Message)
	return nil
}

// SlackNotifier implements AlertNotifier for Slack notifications
type SlackNotifier struct {
	WebhookURL string
}

func (n *SlackNotifier) Notify(alert Alert) error {
	// In a real implementation, this would send a Slack message
	fmt.Printf("Sending Slack alert for %s: %s\n", alert.AgentID, alert.Message)
	return nil
}

// WebhookNotifier implements AlertNotifier for generic webhook notifications
type WebhookNotifier struct {
	URL string
}

func (n *WebhookNotifier) Notify(alert Alert) error {
	// In a real implementation, this would POST to a webhook
	fmt.Printf("Sending webhook alert for %s: %s\n", alert.AgentID, alert.Message)
	return nil
}