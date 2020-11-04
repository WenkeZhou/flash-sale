package service

import (
	"errors"
	"github.com/WenkeZhou/flash-sale/internal/dao"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"time"
)

type BuyRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type Stock struct {
	ID      uint32 `json:"id"`
	Name    string `json:"string"`
	Count   uint32 `json:"count"`
	Sale    uint32 `json:"sale"`
	Version uint32 `json:"version"`
}

func (svc *Service) Buy(param *BuyRequest) (*StockOrder, error) {
	stock, err := svc.dao.GetStock(param.ID)
	if err != nil {
		return nil, err
	}

	// 检查库存是否充足
	if stock.Sale == stock.Count {
		return nil, errcode.NotFound
	}

	// 扣库存
	err = svc.saleStock(stock.ID, stock.Sale)

	// 创建订单
	stockOrder, err := svc.dao.CreateStockOrder(&dao.StockOrder{
		Sid:        stock.ID,
		Name:       stock.Name,
		CreateTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, errors.New("创建订单失败")
	}
	return &StockOrder{
		ID:         stockOrder.ID,
		Sid:        stockOrder.Sid,
		Name:       stockOrder.Name,
		CreateTime: stockOrder.CreateTime,
	}, nil
}

//func (svc *Service) checkStock(stock Stock) bool {
//	if stock.Sale == stock.Count {
//		return false
//	}
//	return true
//}

func (svc *Service) saleStock(id uint32, sale uint32) error {
	err := svc.dao.UpdateStock(&dao.Stock{
		ID:   id,
		Sale: sale + 1,
	})
	if err != nil {
		return err
	}
	return nil
}
