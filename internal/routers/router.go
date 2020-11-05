package routers

import (
	"github.com/WenkeZhou/flash-sale/internal/routers/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	stock := api.NewStock()

	r.POST("/buy/:id", stock.Buy)
	r.POST("/buywithpesslock/:id", stock.BuyWithPessimisticLock)
	r.POST("/buywithoptlock/:id", stock.BuyWithOptimisticLock)

	return r
}
