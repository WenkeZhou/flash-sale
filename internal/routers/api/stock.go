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
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) BuyWithPessimisticLock(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyWithPessimisticLock(&param)
	if err != nil {
		log.Printf("BuyWithPessimisticLock||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) BuyWithOptimisticLock(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyWithOptimisticLock(&param)
	if err != nil {
		log.Printf("BuyWithOptimisticLock||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) GetVerifyHash(c *gin.Context) {
	param := service.GetVerifyHashRequest{
		SID:    convert.StrTo(c.Param("sid")).MustInt32(),
		UserID: convert.StrTo(c.Param("userid")).MustInt32(),
	}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	verifyHash, err := svc.GetVerifyHash(&param)
	if err != nil {
		log.Printf("GetVerifyHash||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(map[string]interface{}{"verifyHash": verifyHash})
	return
}

func (s Stock) BuyMd5(c *gin.Context) {
	param := service.UserByRequest{
		ID:         convert.StrTo(c.Param("sid")).MustInt32(),
		UserID:     convert.StrTo(c.Param("userid")).MustInt32(),
		VerifyHash: c.Param("verifyhash"),
	}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyMd5(&param)
	if err != nil {
		log.Printf("BuyMd5||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) GetStockByDB(c *gin.Context) {
	param := service.GetStock{
		ID: convert.StrTo(c.Param("sid")).MustInt32(),
	}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stock, err := svc.GetStockByDB(&param)
	if err != nil {
		log.Printf("GetStockByDB||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stock)
	return
}

func (s Stock) GetStockByCache(c *gin.Context) {
	param := service.GetStock{
		ID: convert.StrTo(c.Param("sid")).MustInt32(),
	}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stock, err := svc.GetStockByCache(&param)
	if err != nil {
		log.Printf("GetStockByCache||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stock)
	return
}

func (s Stock) BuyWithCacheV1(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyWithCacheV1(&param)
	if err != nil {
		log.Printf("BuyWithCacheV1||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) BuyWithCacheV2(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyWithCacheV2(&param)
	if err != nil {
		log.Printf("BuyWithCacheV2||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}

func (s Stock) BuyWithCacheV3(c *gin.Context) {
	param := service.BuyRequest{ID: convert.StrTo(c.Param("id")).MustInt32()}
	response := app.NewResponse(c)

	svc := service.New(c.Request.Context())
	stockOrder, err := svc.BuyWithCacheV3(&param)
	if err != nil {
		log.Printf("BuyWithCacheV3||err:%v", err)
		switch err.(type) {
		case *errcode.Error:
			response.ToErrorResponse(err.(*errcode.Error))
		default:
			response.ToErrorResponse(errcode.ErrorBuyStock)
		}
		return
	}
	response.ToResponse(stockOrder)
	return
}
