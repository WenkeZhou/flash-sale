package dao

import "github.com/WenkeZhou/flash-sale/internal/model"

type User struct {
	ID       uint32 `json:"id"`
	UserName string `json:"user_name"`
}

func (d *Dao) GetUser(id uint32) (model.User, error) {
	user := model.User{ID: id}
	return user.Get(d.engine)
}
