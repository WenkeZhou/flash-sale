package model

import "gorm.io/gorm"

type StockOrder struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	Sid        uint32 `json:"sid"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
}

func (so StockOrder) TableName() string {
	return "stock_order"
}

func (so StockOrder) Create(db *gorm.DB) (*StockOrder, error) {
	if err := db.Create(&so).Error; err != nil {
		return nil, err
	}
	return &so, nil
}
