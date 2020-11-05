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

// 不做任何限制
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

// 悲观锁购买
func (svc *Service) BuyWithPessimisticLock(param *BuyRequest) (*StockOrder, error) {
	stockOrder, err := svc.dao.BuyWithPessimisticLock(param.ID)

	if err != nil {
		return nil, err
	}
	return &StockOrder{
		ID:         stockOrder.ID,
		Sid:        stockOrder.Sid,
		Name:       stockOrder.Name,
		CreateTime: stockOrder.CreateTime,
	}, nil
}

// 乐观锁购买

func (svc *Service) BuyWithOptimisticLock(param *BuyRequest) (*StockOrder, error) {
	stockOrder, err := svc.dao.BuyWithOptimisticLock(param.ID)

	if err != nil {
		return nil, err
	}
	return &StockOrder{
		ID:         stockOrder.ID,
		Sid:        stockOrder.Sid,
		Name:       stockOrder.Name,
		CreateTime: stockOrder.CreateTime,
	}, nil
}
