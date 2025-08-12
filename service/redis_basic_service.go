package service

import (
	"context"
	"gin-demo/database"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBasicService struct {
	rdb *redis.Client
}

func NewRedisBasicService() *RedisBasicService {
	return &RedisBasicService{
		rdb: database.GetRedis(),
	}
}

func (r *RedisBasicService) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r *RedisBasicService) GetString(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisBasicService) Increment(ctx context.Context, key string) (int64, error) {
	return r.rdb.Incr(ctx, key).Result()
}

func (r *RedisBasicService) SetHash(ctx context.Context, key, field, value string) error {
	return r.rdb.HSet(ctx, key, field, value).Err()
}

func (r *RedisBasicService) GetHash(ctx context.Context, key, field string) (string, error) {
	return r.rdb.HGet(ctx, key, field).Result()
}

func (r *RedisBasicService) PushToList(ctx context.Context, key string, values ...interface{}) error {
	return r.rdb.LPush(ctx, key, values...).Err()
}

func (r *RedisBasicService) GetFromList(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.rdb.LRange(ctx, key, start, stop).Result()
}

func (r *RedisBasicService) Delete(ctx context.Context, keys ...string) error {
	return r.rdb.Del(ctx, keys...).Err()
}

func (r *RedisBasicService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.rdb.Exists(ctx, key).Result()
	return count > 0, err
}
