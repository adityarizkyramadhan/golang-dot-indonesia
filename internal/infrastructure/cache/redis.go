package cache

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	redis *redis.Client
}

func NewRedis() *Redis {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &Redis{client}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expired time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.redis.WithContext(ctx).Set(key, jsonBytes, expired).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	val, err := r.redis.WithContext(ctx).Get(key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), value)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	err := r.redis.WithContext(ctx).Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}
