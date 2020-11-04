package dao

import "github.com/WenkeZhou/flash-sale/internal/model"

type StockOrder struct {
	ID         uint32 `json:"id"`
	Sid        uint32 `json:"sid"`
	Name       string `json:"string"`
	CreateTime int64  `json:"create_time"`
}

func (d *Dao) CreateStockOrder(param *StockOrder) (*model.StockOrder, error) {
	stockOrder := model.StockOrder{
		Sid:        param.Sid,
		Name:       param.Name,
		CreateTime: param.CreateTime,
	}
	return stockOrder.Create(d.engine)
}
