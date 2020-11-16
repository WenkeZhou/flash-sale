package routers

import (
	"github.com/WenkeZhou/flash-sale/internal/middleware"
	"github.com/WenkeZhou/flash-sale/internal/routers/api"
	"github.com/WenkeZhou/flash-sale/pkg/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBucket(limiter.LimiterBucketRule{
	Key:          "/buywithoptlocklimiter/1",
	FillInterval: time.Millisecond * 200,
	Capacity:     10,
	Quantum:      10,
})

var fullPathLimiters = limiter.NewFullPathLimiter().AddBucket(limiter.LimiterBucketRule{
	Key:          "/buywithoptlocklimiter/:id",
	FillInterval: time.Second * 5,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.RateLimiter(fullPathLimiters))

	stock := api.NewStock()

	r.POST("/buy/:id", stock.Buy)
	r.POST("/buywithpesslock/:id", stock.BuyWithPessimisticLock)
	r.POST("/buywithoptlock/:id", stock.BuyWithOptimisticLock)
	r.POST("/buywithoptlocklimiter/:id", stock.BuyWithOptimisticLock)
	r.GET("/getverifyhash/stock/:sid/user/:userid", stock.GetVerifyHash)
	r.POST("/buymd5/stock/:sid/user/:userid/verifyhash/:verifyhash", stock.BuyMd5)
	r.GET("/getstockbydb/:sid", stock.GetStockByDB)
	r.GET("/getstockbycache/:sid", stock.GetStockByCache)

	r.POST("/buywithcachev1/:id", stock.BuyWithCacheV1)
	r.POST("/buywithcachev2/:id", stock.BuyWithCacheV2)
	r.POST("/buywithcachev3/:id", stock.BuyWithCacheV3)

	return r
}
