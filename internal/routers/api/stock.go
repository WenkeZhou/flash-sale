package api

import (
	"github.com/WenkeZhou/flash-sale/internal/service"
	"github.com/WenkeZhou/flash-sale/pkg/app"
	"github.com/WenkeZhou/flash-sale/pkg/convert"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/gin-gonic/gin"
	"log"
)

type Stock struct{}

func NewStock() Stock {
	return Stock{}
}

func (s Stock) Buy(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.Buy(&param)
	if err != nil {
		log.Printf("Buy||err:%v", err)
		response.ToErrorResponse(errcode.ErrorBuyStock)
		return
	}
	response.ToResponse(stockOrder)
	return
}
