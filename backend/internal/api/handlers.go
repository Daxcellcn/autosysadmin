// backend/internal/api/handlers.go
package api

import (
	"net/http"

	"github.com/autosysadmin/backend/internal/agent"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleLogin(c *gin.Context) {
	// Implementation would use s.authService
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint"})
}

func (s *Server) handleRegister(c *gin.Context) {
	// Implementation would use s.authService
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint"})
}

func (s *Server) handleRefreshToken(c *gin.Context) {
	// Implementation would use s.authService
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token endpoint"})
}

func (s *Server) listAgents(c *gin.Context) {
	agents := s.agentManager.ListAgents()
	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

func (s *Server) registerAgent(c *gin.Context) {
	var newAgent agent.Agent
	if err := c.ShouldBindJSON(&newAgent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.agentManager.RegisterAgent(&newAgent)
	c.JSON(http.StatusCreated, gin.H{"agent": newAgent})
}

func (s *Server) getAgent(c *gin.Context) {
	agentID := c.Param("id")
	agent, exists := s.agentManager.GetAgent(agentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (s *Server) runCommand(c *gin.Context) {
	// Implementation would run commands on agents
	c.JSON(http.StatusOK, gin.H{"message": "Command execution endpoint"})
}

func (s *Server) getAgentStats(c *gin.Context) {
	agentID := c.Param("id")
	agent, exists := s.agentManager.GetAgent(agentID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	stats, err := agent.CollectStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// ... other handler implementations would follow the same pattern