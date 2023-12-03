package middleware

import (
	"GinProject/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var tokenBuckets = make(map[string]*utils.TokenBucket)

func LimiterMiddleWare(cap int, rate int, unit time.Duration, pattern int, timeOut time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		tb, exists := tokenBuckets[c.FullPath()]
		if !exists {
			tb = utils.NewTokenBucket(cap, rate, unit)
			tokenBuckets[c.FullPath()] = tb
		}
		switch pattern {
		case 1:
			tb.Acquire()
			break
		case 2:
			if ok := tb.TryAcquire(); !ok {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "too many requests,reject",
				})
				c.Abort()
				return
			}
			break
		case 3:
			if ok := tb.AcquireWithTimeOut(timeOut); !ok {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "too many requests,reject",
				})
				c.Abort()
				return
			}
			break
		default:
			tb.TryAcquire()
		}
		c.Next()
	}
}
