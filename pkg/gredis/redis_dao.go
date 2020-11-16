package gredis

import (
	"encoding/json"
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/errcode"
	"github.com/WenkeZhou/flash-sale/pkg/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)

func InitRedisConn(rds *setting.RedisSettingS) (*redis.Pool, error) {
	RedisConn := &redis.Pool{
		MaxIdle:     rds.MaxIdle,
		MaxActive:   rds.MaxActive,
		IdleTimeout: rds.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(rds.LinkType, rds.Address)
			if err != nil {
				return nil, err
			}

			if rds.Password != "" {
				if _, err := c.Do("AUTH", rds.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return RedisConn, nil
}

func SetCommon(RedisConn *redis.Pool, key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	fmt.Printf("Redis设置 key[%v], value[%v], expire[%v]\n", key, data, time)
	return nil
}

func Set(RedisConn *redis.Pool, key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func Get(RedisConn *redis.Pool, key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	//reply, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func Incr(RedisConn *redis.Pool, key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func IncrWithExpiry(RedisConn *redis.Pool, key string, t int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, t)
	r, err := conn.Do("EXEC")

	//reply, err := redis.Int(conn.Do("INCR", key))
	fmt.Printf("IncrWithExpiry:result:[%v] \n", r)
	if err != nil {
		return err
	}
	return nil
}

func Exists(RedisConn *redis.Pool, key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Delete(RedisConn *redis.Pool, key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(RedisConn *redis.Pool, key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, v := range keys {
		_, err := Delete(RedisConn, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCommon(RedisConn *redis.Pool, key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	//reply, err := conn.Do("GET", key)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func GetInt(RedisConn *redis.Pool, key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("GET", key))
	//reply, err := conn.Do("GET", key)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return reply, err
		} else {
			return reply, errcode.NotFound
		}
	}
	return reply, nil
}
