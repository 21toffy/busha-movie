package cache

import (
	"context"
	"encoding/json"

	"fmt"

	"time"

	"github.com/21toffy/busha-movie/internal/customerror"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

func InitRedisCache() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         viper.GetString("redis.addr"),
			Password:     viper.GetString("redis.password"),
			DB:           viper.GetInt("redis.db"),
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     10,
			PoolTimeout:  30 * time.Second,
			IdleTimeout:  5 * time.Minute,
		})

		_, err := redisClient.Ping(redisClient.Context()).Result()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
	}

	return redisClient
}

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: InitRedisCache(),
		ctx:    context.Background(),
	}
}

// function to get a particular data with an associated key from redis
func (c *RedisCache) Get(key string, result interface{}) error {
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return customerror.ErrCacheMiss
		}
		return customerror.OtherCacheError
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return customerror.UnmarshalingError
	}

	return nil
}

// function to load/store data into redis
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return customerror.UnmarshalingError
	}

	err = c.client.Set(c.ctx, key, data, expiration).Err()
	if err != nil {
		return customerror.CacheSetError
	}

	return nil
}
