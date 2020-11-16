package main

import (
	"github.com/WenkeZhou/flash-sale/global"
	"github.com/WenkeZhou/flash-sale/internal/model"
	"github.com/WenkeZhou/flash-sale/internal/routers"
	"github.com/WenkeZhou/flash-sale/pkg/gredis"
	"github.com/WenkeZhou/flash-sale/pkg/setting"
	"log"
	"net/http"
	"time"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init setupDBEngine err: %v", err)
	}

	err = setupRedisConn()
	if err != nil {
		log.Fatalf("init setupRedisConn err: %v", err)
	}

}

func setupSetting() error {
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("RedisDatabase", &global.RedisSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("VerifyData", &global.VerifySetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("BusinessSetting", &global.BusinessSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.RedisSetting.IdleTimeout *= time.Second

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupRedisConn() error {
	var err error
	global.RedisConn, err = gredis.InitRedisConn(global.RedisSetting)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//gin.SetMode("debug")
	router := routers.NewRouter()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
