package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count     int
	resetTime time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
	limit    = 100
	window   = 60 * time.Second
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var identifier string

		userIDVal, exists := c.Get("userID")
		if exists {
			identifier = "user:" + userIDVal.(string)
		} else {
			identifier = "ip:" + c.ClientIP()
		}

		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		v, ok := visitors[identifier]
		if !ok || now.After(v.resetTime) {
			visitors[identifier] = &visitor{
				count:     1,
				resetTime: now.Add(window),
			}
			c.Next()
			return
		}

		if v.count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}
		v.count++
		c.Next()
	}
}
