package middleware

import (
	"github.com/WenkeZhou/flash-sale/pkg/app"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/WenkeZhou/flash-sale/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if bucket, ok := l.GetBucket(key); ok {
			count := bucket.TakeAvailable(1)
			if count == 0 {
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.ErrorTooManyRequest)
			}
		}
	}
}
