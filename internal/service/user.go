package service

import (
	"github.com/WenkeZhou/flash-sale/global"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/WenkeZhou/flash-sale/pkg/gredis"
	"strconv"
)

type User struct {
	ID       uint32 `json:"id"`
	UserName string `json:"username"`
}

func (u *User) AddUserVisitCount() error {
	k := global.VerifySetting.UserVisitCountPrefix + strconv.Itoa(int(u.ID))
	err := gredis.IncrWithExpiry(global.RedisConn, k, 1)
	if err != nil {
		return err
	}
	//fmt.Printf("用户 UserId [%v], 访问次数 [%v]\n", u.ID, v)
	return nil
}

func (u *User) GetUserIsBanded() (bool, error) {
	k := global.VerifySetting.UserVisitCountPrefix + strconv.Itoa(int(u.ID))
	v, err := gredis.GetInt(global.RedisConn, k)
	if err != nil {
		if err != errcode.NotFound {
			return false, err
		}
		v = 0
	}
	if v > global.VerifySetting.MaxUserBuyCount {
		return false, nil
	}
	return true, nil
}
