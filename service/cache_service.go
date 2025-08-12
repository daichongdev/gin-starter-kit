package service

import (
	"context"
	"encoding/json"
	"gin-demo/database"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	rdb *redis.Client
}

func NewCacheService() *CacheService {
	return &CacheService{
		rdb: database.GetRedis(),
	}
}

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, expiration).Err()
}

func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.rdb.Exists(ctx, key).Result()
	return count > 0, err
}