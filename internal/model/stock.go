package model

import (
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"gorm.io/gorm"
)

type Stock struct {
	ID      uint32 `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Count   uint32 `json:"count"`
	Sale    uint32 `json:"sale"`
	Version uint32 `json:"version"`
}

func (s Stock) TableName() string {
	return "stock"
}

func (s Stock) Get(db *gorm.DB) (Stock, error) {
	var stock Stock
	err := db.Where("id = ? ", s.ID).First(&stock).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return stock, err
	}

	if err == gorm.ErrRecordNotFound {
		return stock, errcode.NotFoundStock
	}

	return stock, nil
}

func (s Stock) Update(db *gorm.DB, values interface{}) error {
	return db.Model(&s).Updates(values).Where("id = ? ", s.ID).Error
}
