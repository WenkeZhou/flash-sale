package model

import (
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"gorm.io/gorm"
	"time"
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

func (s Stock) BuyWithPessimisticLock(db *gorm.DB) (*StockOrder, error) {
	var stock Stock
	var stockOrder StockOrder

	err := db.Transaction(func(tx *gorm.DB) error {

		// 查询库存
		err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? ", s.ID).First(&stock).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			return errcode.NotFoundStock
		}

		// 判断库存
		if stock.Sale == stock.Count {
			return errcode.SellOutStock
		}

		//扣除库存
		err = tx.Model(&stock).Updates(map[string]interface{}{"sale": stock.Sale + 1}).Where("id = ? ", stock.ID).Error
		if err != nil {
			return err
		}
		stockOrder = StockOrder{Sid: stock.ID, Name: stock.Name, CreateTime: time.Now().Unix()}
		//创建订单
		err = tx.Create(&stockOrder).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	//stockOrder = StockOrder{Sid: stock.ID, Name: stock.Name, CreateTime: time.Now().Unix()}
	////创建订单
	//err = db.Create(&stockOrder).Error
	//if err != nil {
	//	return nil, err
	//}

	return &stockOrder, nil
}

func (s Stock) BuyWithPessimisticLock2(db *gorm.DB) (*StockOrder, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var stock Stock
	var stockOrder StockOrder

	// 查询库存
	//err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? ", s.ID).First(&stock).Error
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&stock, s.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, errcode.NotFoundStock
	}

	// 判断库存
	if stock.Sale == stock.Count {
		tx.Rollback()
		return nil, errcode.SellOutStock
	}

	//扣除库存
	stock.Sale = stock.Sale + 1
	if err = tx.Save(&stock).Error; err != nil {
		tx.Rollback()
		return nil, err
	} else {
		stockOrder = StockOrder{Sid: stock.ID, Name: stock.Name, CreateTime: time.Now().Unix()}
		//创建订单
		err = tx.Create(&stockOrder).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	//result := tx.Model(Stock{}).Where("id = ? ", stock.ID).Update("sale", gorm.Expr("sale + ?", 1))
	//if result.RowsAffected == 0 {
	//	tx.Rollback()
	//	return nil, errcode.SellOutStock
	//}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &stockOrder, nil
}

func (s Stock) BuyWithPessimisticLock3(db *gorm.DB) (*StockOrder, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var stock Stock
	var stockOrder StockOrder

	// 查询库存
	result := tx.Model(Stock{}).Where("id = ? AND sale < count ", s.ID).Update("sale", gorm.Expr("sale + ?", 1))
	if result.RowsAffected > 0 {
		tx.First(&stock, s.ID)
		stockOrder = StockOrder{Sid: stock.ID, Name: stock.Name, CreateTime: time.Now().Unix()}
		//创建订单
		err := tx.Create(&stockOrder).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		tx.Rollback()
		return nil, errcode.ErrorBuyStock
	}

	tx.Commit()

	return &stockOrder, nil
}

func (s Stock) BuyWithOptimisticLock(db *gorm.DB) (*StockOrder, error) {
	tx := db.Begin()
	var stock Stock
	var stockOrder StockOrder

	// 查询库存
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? ", s.ID).First(&stock).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, errcode.NotFoundStock
	}

	// 判断库存
	if stock.Sale == stock.Count {
		tx.Rollback()
		return nil, errcode.SellOutStock
	}

	updateValues := map[string]interface{}{
		"sale":    stock.Sale + 1,    //扣除库存
		"version": stock.Version + 1, // 乐观锁递增
	}
	result := tx.Model(Stock{}).Where(map[string]interface{}{"id": stock.ID, "version": stock.Version}).Updates(updateValues)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errcode.SellOutStock
	}

	stockOrder = StockOrder{Sid: stock.ID, Name: stock.Name, CreateTime: time.Now().Unix()}
	//创建订单
	err = tx.Create(&stockOrder).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &stockOrder, nil
}
