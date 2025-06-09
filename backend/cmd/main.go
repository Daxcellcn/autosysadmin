// backend/cmd/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/autosysadmin/backend/internal/agent"
	"github.com/autosysadmin/backend/internal/api"
	"github.com/autosysadmin/backend/internal/auth"
	"github.com/autosysadmin/backend/internal/billing"
	"github.com/autosysadmin/backend/internal/jobqueue"
	"github.com/autosysadmin/backend/internal/monitoring"
	"github.com/autosysadmin/backend/internal/patching"
	"github.com/autosysadmin/backend/internal/security"
	"github.com/autosysadmin/backend/internal/subscriptions"
	"github.com/autosysadmin/backend/internal/usage"
)

func main() {
	// Initialize all components
	authService := auth.NewAuthService()
	jobQueue := jobqueue.NewRedisJobQueue()
	agentManager := agent.NewManager(jobQueue)
	monitoringService := monitoring.NewMonitor()
	patchingService := patching.NewPatchManager()
	securityScanner := security.NewVulnerabilityScanner()
	billingService := billing.NewBillingService()
	subscriptionService := subscriptions.NewService()
	usageTracker := usage.NewTracker()

	// Start the API server
	apiServer := api.NewServer(
		authService,
		agentManager,
		monitoringService,
		patchingService,
		securityScanner,
		billingService,
		subscriptionService,
		usageTracker,
	)

	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	apiServer.Stop()
	jobQueue.Close()
	log.Println("Server exited properly")
}