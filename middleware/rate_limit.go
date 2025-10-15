package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var visitTimestamps = struct {
	sync.Mutex
	data map[string]time.Time
}{data: make(map[string]time.Time)}

// RateLimitMiddleware limits requests per IP every 10 seconds
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		visitTimestamps.Lock()
		defer visitTimestamps.Unlock()

		if last, exists := visitTimestamps.data[ip]; exists && now.Sub(last) < 10*time.Second {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests — please wait a few seconds",
			})
			return
		}

		visitTimestamps.data[ip] = now
		c.Next()
	}
}