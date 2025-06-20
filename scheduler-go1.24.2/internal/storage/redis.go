package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

// InitRedis 初始化Redis
func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{Addr: addr})
	if err := RDB.Ping(Ctx).Err(); err != nil {
		panic("failed to connect redis: " + err.Error())
	}
}
