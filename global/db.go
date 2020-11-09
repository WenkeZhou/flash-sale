package global

import (
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

var (
	DBEngine  *gorm.DB
	RedisConn *redis.Pool
)
