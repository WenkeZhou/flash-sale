package dao

import "github.com/WenkeZhou/flash-sale/internal/model"

type Stock struct {
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Count   uint32 `json:"count"`
	Sale    uint32 `json:"sale"`
	Version uint32 `json:"version"`
}

func (d *Dao) GetStock(id uint32) (model.Stock, error) {
	stock := model.Stock{ID: id}
	return stock.Get(d.engine)
}

func (d *Dao) UpdateStock(param *Stock) error {
	stock := model.Stock{ID: param.ID}
	values := map[string]interface{}{}
	values["sale"] = param.Sale
	return stock.Update(d.engine, values)
}

func (d *Dao) BuyWithPessimisticLock(id uint32) (*model.StockOrder, error) {
	stock := model.Stock{ID: id}
	return stock.BuyWithPessimisticLock3(d.engine)
}

func (d *Dao) BuyWithOptimisticLock(id uint32) (*model.StockOrder, error) {
	stock := model.Stock{ID: id}
	return stock.BuyWithOptimisticLock(d.engine)
}

func (d *Dao) BuyMd5(sid, userId uint32) (*model.UserStockOrder, error) {
	stock := model.Stock{ID: sid}
	return stock.BuyWithUser(d.engine, userId)
}
