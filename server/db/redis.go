package db

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

func InitPool(address string, maxIdle, maxActive, idleTimeout int) {
	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}