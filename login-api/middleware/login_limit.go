package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)
type attemptInfo struct {
	Count    int
	LastTry time.Time
}
var (
	attempts = make(map[string]*attemptInfo)
	mu       sync.Mutex
)
const (
	maxAttempts = 5
	blockTime   = 2 * time.Minute
)
func LoginLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()

		mu.Lock()
		info, exists := attempts[ip]
		if !exists {
			attempts[ip] = &attemptInfo{
				Count:    1,
				LastTry: time.Now(),
			}
			mu.Unlock()
			c.Next()
			return
		}
		if info.Count >= maxAttempts {
			if time.Since(info.LastTry) < blockTime {
				mu.Unlock()
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Too many login attempts. Try again later.",
				})
				c.Abort()
				return
			}
		
			info.Count = 1
			info.LastTry = time.Now()
			mu.Unlock()
			c.Next()
			return
		}
		info.Count++
		info.LastTry = time.Now()
		mu.Unlock()

		c.Next()
	}
}