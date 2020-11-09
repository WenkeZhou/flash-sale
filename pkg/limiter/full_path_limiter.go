package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type FullPathLimiter struct {
	*Limiter
}

func (fp FullPathLimiter) Key(c *gin.Context) string {
	url := c.FullPath()
	return url
}

func (fp FullPathLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := fp.LimiterBuckets[key]
	return bucket, ok
}

func (fp FullPathLimiter) AddBucket(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		if _, ok := fp.LimiterBuckets[rule.Key]; !ok {
			fp.LimiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}
	return fp
}

func NewFullPathLimiter() LimiterIface {
	return FullPathLimiter{
		Limiter: &Limiter{LimiterBuckets: make(map[string]*ratelimit.Bucket)},
	}
}
