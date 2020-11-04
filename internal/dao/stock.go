package dao

import "github.com/WenkeZhou/flash-sale/internal/model"

type Stock struct {
	ID      uint32 `json:"id"`
	Name    string `json:"string"`
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
