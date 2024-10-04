package cache

import (
	"os"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func NewRedis() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	redisClient = client
	return redisClient
}
