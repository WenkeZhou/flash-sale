package model

import (
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"gorm.io/gorm"
)

type User struct {
	ID       uint32 `gorm:"primary_key" json:"id"`
	UserName string `json:"username"`
}

func (u User) TableName() string {
	return "user"
}

func (u User) Get(db *gorm.DB) (User, error) {
	var user User
	err := db.Where("id = ? ", u.ID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}

	if err == gorm.ErrRecordNotFound {
		return user, errcode.NotFoundStock
	}

	return user, nil
}
