// backend/internal/api/middleware/ratelimit.go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visits map[string][]time.Time
	mu     sync.Mutex
	limit  int
	window time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		visits: make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	visits := rl.visits[ip]

	// Remove visits outside the time window
	var recentVisits []time.Time
	for _, visit := range visits {
		if now.Sub(visit) <= rl.window {
			recentVisits = append(recentVisits, visit)
		}
	}

	if len(recentVisits) >= rl.limit {
		return false
	}

	recentVisits = append(recentVisits, now)
	rl.visits[ip] = recentVisits
	return true
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}