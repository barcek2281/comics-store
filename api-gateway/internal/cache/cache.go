package cache

import "github.com/go-redis/redis/v8"

func NewRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        DB:   0,
    })
}