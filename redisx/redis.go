package redisx

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// 127.0.0.1:6379
func GetRedisPool(connStr string, maxActive int, maxIdle int, passport string) *redis.Pool {
	RedisPool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 3600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connStr)
			if err != nil {
				return nil, err
			}
			c.Do("auth", passport)
			return c, nil
		},
	}

	return RedisPool
}
