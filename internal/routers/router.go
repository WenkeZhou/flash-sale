package routers

import (
	"github.com/WenkeZhou/flash-sale/internal/middleware"
	"github.com/WenkeZhou/flash-sale/internal/routers/api"
	"github.com/WenkeZhou/flash-sale/pkg/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/buywithoptlocklimiter",
	FillInterval: time.Millisecond * 100,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimiter(methodLimiters))

	stock := api.NewStock()

	r.POST("/buy/:id", stock.Buy)
	r.POST("/buywithpesslock/:id", stock.BuyWithPessimisticLock)
	r.POST("/buywithoptlock/:id", stock.BuyWithOptimisticLock)
	r.POST("/buywithoptlocklimiter/:id", stock.BuyWithOptimisticLock)

	return r
}
