// backend/internal/api/server.go
package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/autosysadmin/backend/internal/agent"
	"github.com/autosysadmin/backend/internal/auth"
	"github.com/autosysadmin/backend/internal/billing"
	"github.com/autosysadmin/backend/internal/monitoring"
	"github.com/autosysadmin/backend/internal/patching"
	"github.com/autosysadmin/backend/internal/security"
	"github.com/autosysadmin/backend/internal/subscriptions"
	"github.com/autosysadmin/backend/internal/usage"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	router            *gin.Engine
	httpServer        *http.Server
	authService       auth.AuthService
	agentManager      *agent.Manager
	monitoringService monitoring.Monitor
	patchingService   patching.PatchManager
	securityScanner   security.VulnerabilityScanner
	billingService    billing.BillingService
	subscriptionService subscriptions.Service
	usageTracker      usage.Tracker
}

func NewServer(
	authService auth.AuthService,
	agentManager *agent.Manager,
	monitoringService monitoring.Monitor,
	patchingService patching.PatchManager,
	securityScanner security.VulnerabilityScanner,
	billingService billing.BillingService,
	subscriptionService subscriptions.Service,
	usageTracker usage.Tracker,
) *Server {
	router := gin.Default()
	server := &Server{
		router:            router,
		authService:       authService,
		agentManager:      agentManager,
		monitoringService: monitoringService,
		patchingService:   patchingService,
		securityScanner:   securityScanner,
		billingService:    billingService,
		subscriptionService: subscriptionService,
		usageTracker:      usageTracker,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", s.handleLogin)
			authGroup.POST("/register", s.handleRegister)
			authGroup.POST("/refresh", s.handleRefreshToken)
		}

		// Agent routes
		agentGroup := api.Group("/agents")
		agentGroup.Use(s.authMiddleware())
		{
			agentGroup.GET("/", s.listAgents)
			agentGroup.POST("/", s.registerAgent)
			agentGroup.GET("/:id", s.getAgent)
			agentGroup.POST("/:id/command", s.runCommand)
			agentGroup.GET("/:id/stats", s.getAgentStats)
		}

		// Monitoring routes
		monitorGroup := api.Group("/monitoring")
		monitorGroup.Use(s.authMiddleware())
		{
			monitorGroup.GET("/alerts", s.listAlerts)
			monitorGroup.POST("/alerts", s.createAlert)
			monitorGroup.GET("/metrics", s.getMetrics)
		}

		// Patching routes
		patchGroup := api.Group("/patching")
		patchGroup.Use(s.authMiddleware())
		{
			patchGroup.GET("/updates", s.listAvailableUpdates)
			patchGroup.POST("/apply", s.applyUpdates)
		}

		// Security routes
		securityGroup := api.Group("/security")
		securityGroup.Use(s.authMiddleware())
		{
			securityGroup.GET("/scan", s.runSecurityScan)
			securityGroup.GET("/compliance", s.checkCompliance)
		}

		// Billing routes
		billingGroup := api.Group("/billing")
		billingGroup.Use(s.authMiddleware())
		{
			billingGroup.GET("/plans", s.listPlans)
			billingGroup.GET("/usage", s.getUsage)
			billingGroup.POST("/subscribe", s.createSubscription)
		}
	}

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	var g errgroup.Group
	g.Go(func() error {
		log.Println("Starting HTTP server on :8080")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	return g.Wait()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}