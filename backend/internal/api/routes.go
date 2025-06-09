package api

import (
	"github.com/gin-gonic/gin"
	"github.com/autosysadmin/backend/internal/api/middleware"
)

func (s *Server) setupRoutes() {
	// Public routes
	public := s.router.Group("/api/v1")
	{
		authGroup := public.Group("/auth")
		{
			authGroup.POST("/login", s.handleLogin)
			authGroup.POST("/register", s.handleRegister)
			authGroup.POST("/refresh", s.handleRefreshToken)
		}
	}

	// Protected routes (require authentication)
	protected := s.router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(s.authService))
	{
		// Agent routes
		agentGroup := protected.Group("/agents")
		{
			agentGroup.GET("/", s.listAgents)
			agentGroup.POST("/", s.registerAgent)
			agentGroup.GET("/:id", s.getAgent)
			agentGroup.POST("/:id/command", s.runCommand)
			agentGroup.GET("/:id/stats", s.getAgentStats)
			agentGroup.GET("/:id/updates", s.listAvailableUpdates)
			agentGroup.POST("/:id/updates", s.applyUpdates)
		}

		// Monitoring routes
		monitorGroup := protected.Group("/monitoring")
		{
			monitorGroup.GET("/", s.getMonitoringDashboard)
			monitorGroup.GET("/alerts", s.listAlerts)
			monitorGroup.POST("/alerts", s.createAlert)
			monitorGroup.GET("/metrics", s.getMetrics)
			monitorGroup.GET("/metrics/:agent_id", s.getAgentMetrics)
		}

		// Security routes
		securityGroup := protected.Group("/security")
		{
			securityGroup.POST("/scan/:agent_id", s.runSecurityScan)
			securityGroup.GET("/scans/:agent_id", s.getScanResults)
			securityGroup.GET("/compliance/:standard", s.getComplianceReport)
			securityGroup.POST("/ssh-keys", s.addSSHKey)
			securityGroup.GET("/ssh-keys/:agent_id", s.listSSHKeys)
		}

		// Billing routes
		billingGroup := protected.Group("/billing")
		{
			billingGroup.GET("/plans", s.listPlans)
			billingGroup.GET("/subscription", s.getSubscription)
			billingGroup.POST("/subscription", s.createSubscription)
			billingGroup.PUT("/subscription", s.updateSubscription)
			billingGroup.GET("/usage", s.getUsage)
			billingGroup.POST("/payment", s.processPayment)
			billingGroup.GET("/payment-history", s.getPaymentHistory)
		}

		// User routes
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/profile", s.getUserProfile)
			userGroup.PUT("/profile", s.updateProfile)
			userGroup.PUT("/password", s.changePassword)
			userGroup.GET("/settings", s.getUserSettings)
			userGroup.PUT("/settings", s.updateUserSettings)
		}
	}

	// Health check route
	s.router.GET("/health", s.healthCheck)
}

func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	db, err := s.db.DB()
	if err != nil {
		c.JSON(503, gin.H{"status": "unhealthy", "error": "database connection failed"})
		return
	}

	if err := db.Ping(); err != nil {
		c.JSON(503, gin.H{"status": "unhealthy", "error": "database ping failed"})
		return
	}

	// Check Redis connection
	if _, err := s.redis.Ping(c).Result(); err != nil {
		c.JSON(503, gin.H{"status": "unhealthy", "error": "redis connection failed"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
	})
}