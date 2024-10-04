package cache

import (
	"os"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func NewRedis() (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	redisClient = client
	return redisClient, nil
}
