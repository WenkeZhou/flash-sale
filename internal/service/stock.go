package service

import (
	"errors"
	"fmt"
	"github.com/WenkeZhou/flash-sale/global"
	"github.com/WenkeZhou/flash-sale/internal/dao"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/WenkeZhou/flash-sale/pkg/gredis"
	"github.com/WenkeZhou/flash-sale/pkg/util"
	"strconv"
	"time"
)

type GetStock struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type BuyRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type UserByRequest struct {
	ID         uint32 `form:"id" binding:"required,gte=1"`
	UserID     uint32 `form:"userid" binding:"required,gte=1"`
	VerifyHash string `form:"verifyhash" binding:"required"`
}

type GetVerifyHashRequest struct {
	SID    uint32 `form:"sid" binding:"required,gte=1"`
	UserID uint32 `form:"userid" binding:"required,gte=1"`
}

type Stock struct {
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Count   uint32 `json:"count"`
	Sale    uint32 `json:"sale"`
	Version uint32 `json:"version"`
}

type StockStorage struct {
	ID      uint32 `json:"id"`
	Storage uint32 `json:"storage"`
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

func (svc *Service) BuyWithCacheV1(param *BuyRequest) (*StockOrder, error) {
	// 删除库存缓存
	if err := svc.DeleteStockStorageCache(param.ID); err != nil {
		return nil, err
	}

	// 完成扣库存 下单
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

func (svc *Service) BuyWithCacheV2(param *BuyRequest) (*StockOrder, error) {
	// 删除库存缓存
	if err := svc.DeleteStockStorageCache(param.ID); err != nil {
		return nil, err
	}

	// 完成扣库存 下单
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

func (svc *Service) DeleteStockStorageCache(sid uint32) error {
	key := global.BusinessSetting.StockCachePrefix + strconv.Itoa(int(sid))
	_, err := gredis.Delete(global.RedisConn, key)
	if err != nil {
		return err
	} else {
		fmt.Println("DeleteStockStorageCache||sid=", sid)
	}
	return nil
}

func (svc *Service) BuyWithCacheV3(param *BuyRequest) (*StockOrder, error) {
	// 完成扣库存 下单
	stockOrder, err := svc.dao.BuyWithOptimisticLock(param.ID)

	if err != nil {
		return nil, err
	}

	// 删除库存缓存
	if err = svc.DeleteStockStorageCache(param.ID); err != nil {
		return nil, err
	}

	go svc.DelayDeleteStockStorageCache(param.ID, 200)

	return &StockOrder{
		ID:         stockOrder.ID,
		Sid:        stockOrder.Sid,
		Name:       stockOrder.Name,
		CreateTime: stockOrder.CreateTime,
	}, nil
}

func (svc *Service) DelayDeleteStockStorageCache(sid uint32, st int) {
	// sid 是商品 id
	// st 是延迟时间，单位毫秒
	ticker := time.Tick(time.Duration(st) * time.Millisecond)
	for {
		select {
		case <-ticker:
			// 删除库存缓存
			err := svc.DeleteStockStorageCache(sid)
			if err != nil {
				fmt.Println("DelayDeleteStockStorageCache||err:", err)
			}
			return
		}
	}
}

func (svc *Service) GetVerifyHash(param *GetVerifyHashRequest) (string, error) {
	sId := param.SID
	userId := param.UserID

	// 判断用户存不存在
	user, err := svc.dao.GetUser(userId)
	if err != nil {
		return "", err
	}

	// 判断商品存不存在
	stock, err := svc.dao.GetStock(sId)
	if err != nil {
		return "", err
	}

	// 生成 hash
	verify := global.VerifySetting.VerifySalt + string(user.ID) + string(stock.ID)
	verifyHash := util.EncodeMD5(verify)
	tmpKey := global.VerifySetting.UserHashKeyPrefix + "_" + strconv.Itoa(int(user.ID)) + "_" + strconv.Itoa(int(stock.ID))
	err = gredis.SetCommon(
		global.RedisConn,
		tmpKey,
		verifyHash,
		3600,
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Redis [写入Key: %v][写入Value: %v]\n", tmpKey, verifyHash)

	return verifyHash, nil
}

func (svc *Service) BuyMd5(param *UserByRequest) (*StockOrder, error) {
	stockId := param.ID
	userId := param.UserID
	verifyHash := param.VerifyHash
	fmt.Printf("UserByRequest_params:%v, %v, %v\n", stockId, userId, verifyHash)

	// 验证 has 值合法性
	hashKey := global.VerifySetting.UserHashKeyPrefix + "_" + strconv.Itoa(int(userId)) + "_" + strconv.Itoa(int(stockId))
	redisValueHash, err := gredis.Get(global.RedisConn, hashKey)
	if err != nil {
		return nil, errcode.RedisGetVerifyHashError
	}

	if verifyHash != redisValueHash {
		return nil, errcode.VerifyHashNotEqual
	}

	// 添加用户访问次数
	user := &User{ID: userId}
	err = user.AddUserVisitCount()
	if err != nil {
		return nil, err
	}

	// 判断用户是否超出频次
	v, err := user.GetUserIsBanded()
	if err != nil {
		return nil, err
	}
	if v == false {
		return nil, errcode.UserRequestFrequently
	}

	// 用户是否存在
	_, err = svc.dao.GetUser(userId)
	if err != nil {
		return nil, err
	}

	stockOrder, err := svc.dao.BuyMd5(stockId, userId)

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

func (svc *Service) GetStockByDB(param *GetStock) (*StockStorage, error) {
	stock, err := svc.dao.GetStock(param.ID)
	if err != nil {
		return nil, err
	}

	return &StockStorage{
		ID:      stock.ID,
		Storage: stock.Count - stock.Sale,
	}, nil
}

func (svc *Service) GetStockByCache(param *GetStock) (*StockStorage, error) {
	k := global.BusinessSetting.StockCachePrefix + strconv.Itoa(int(param.ID))
	storage, err := gredis.GetInt(global.RedisConn, k)
	var count int
	if err != nil {
		if err != errcode.NotFound {
			return nil, err
		}
		stock, err := svc.dao.GetStock(param.ID)
		if err != nil {
			return nil, err
		}
		if err = gredis.SetCommon(global.RedisConn, k, stock.Count-stock.Sale, 60*2); err != nil {
			return nil, err
		}
		count = int(stock.Count - stock.Sale)
	} else {
		count = storage
		fmt.Println("从 Redis cache 中获取缓存, storage:", storage)
	}

	return &StockStorage{
		ID:      param.ID,
		Storage: uint32(count),
	}, nil
}
