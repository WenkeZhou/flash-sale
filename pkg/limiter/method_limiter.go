package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"strings"
)

type MethodLimiter struct {
	*Limiter
}

func (m MethodLimiter) Key(c *gin.Context) string {
	url := c.Request.RequestURI
	index := strings.Index(url, "?")
	if index == -1 {
		return url
	}
	return url[:index]
}

func (m MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := m.LimiterBuckets[key]
	return bucket, ok
}

func (m MethodLimiter) AddBucket(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		if _, ok := m.LimiterBuckets[rule.Key]; !ok {
			m.LimiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}
	return m
}

func NewMethodLimiter() LimiterIface {
	return MethodLimiter{
		Limiter: &Limiter{LimiterBuckets: make(map[string]*ratelimit.Bucket)},
	}
}
