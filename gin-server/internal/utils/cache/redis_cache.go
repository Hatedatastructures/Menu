package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"shack/internal/global"
)

// RedisCache Redis缓存操作工具
type RedisCache struct {
	prefix string // key前缀
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(prefix string) *RedisCache {
	return &RedisCache{
		prefix: prefix,
	}
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if global.GVA_REDIS == nil {
		return "", fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Get(ctx, r.buildKey(key)).Result()
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Set(ctx, r.buildKey(key), value, expiration).Err()
}

// Del 删除缓存
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	if len(keys) == 0 {
		return nil
	}
	buildKeys := make([]string, len(keys))
	for i, key := range keys {
		buildKeys[i] = r.buildKey(key)
	}
	return global.GVA_REDIS.Del(ctx, buildKeys...).Err()
}

// Exists 检查key是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if global.GVA_REDIS == nil {
		return false, fmt.Errorf("redis not initialized")
	}
	count, err := global.GVA_REDIS.Exists(ctx, r.buildKey(key)).Result()
	return count > 0, err
}

// Expire 设置过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Expire(ctx, r.buildKey(key), expiration).Err()
}

// GetInt 获取整数类型缓存
func (r *RedisCache) GetInt(ctx context.Context, key string) (int, error) {
	if global.GVA_REDIS == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Get(ctx, r.buildKey(key)).Int()
}

// GetInt64 获取int64类型缓存
func (r *RedisCache) GetInt64(ctx context.Context, key string) (int64, error) {
	if global.GVA_REDIS == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Get(ctx, r.buildKey(key)).Int64()
}

// SetInt 设置整数类型缓存
func (r *RedisCache) SetInt(ctx context.Context, key string, value int, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration)
}

// SetInt64 设置int64类型缓存
func (r *RedisCache) SetInt64(ctx context.Context, key string, value int64, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration)
}

// GetJSON 获取JSON类型缓存并反序列化
func (r *RedisCache) GetJSON(ctx context.Context, key string, result interface{}) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	val, err := global.GVA_REDIS.Get(ctx, r.buildKey(key)).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

// SetJSON 序列化为JSON并设置缓存
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return global.GVA_REDIS.Set(ctx, r.buildKey(key), data, expiration).Err()
}

// Incr 递增
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	if global.GVA_REDIS == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Incr(ctx, r.buildKey(key)).Result()
}

// Decr 递减
func (r *RedisCache) Decr(ctx context.Context, key string) (int64, error) {
	if global.GVA_REDIS == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.Decr(ctx, r.buildKey(key)).Result()
}

// HGet 获取哈希字段
func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	if global.GVA_REDIS == nil {
		return "", fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.HGet(ctx, r.buildKey(key), field).Result()
}

// HSet 设置哈希字段
func (r *RedisCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.HSet(ctx, r.buildKey(key), field, value).Err()
}

// HDel 删除哈希字段
func (r *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	if len(fields) == 0 {
		return nil
	}
	return global.GVA_REDIS.HDel(ctx, r.buildKey(key), fields...).Err()
}

// HGetAll 获取所有哈希字段
func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if global.GVA_REDIS == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.HGetAll(ctx, r.buildKey(key)).Result()
}

// ZAdd 添加到有序集合
func (r *RedisCache) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.ZAdd(ctx, r.buildKey(key), members...).Err()
}

// ZRange 获取有序集合范围
func (r *RedisCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if global.GVA_REDIS == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.ZRange(ctx, r.buildKey(key), start, stop).Result()
}

// ZRem 从有序集合删除
func (r *RedisCache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}
	return global.GVA_REDIS.ZRem(ctx, r.buildKey(key), members...).Err()
}

// buildKey 构建完整的key
func (r *RedisCache) buildKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s%s", r.prefix, key)
}

// KeysPattern 根据模式删除key
func KeysPattern(ctx context.Context, pattern string) error {
	if global.GVA_REDIS == nil {
		return fmt.Errorf("redis not initialized")
	}

	iter := global.GVA_REDIS.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := global.GVA_REDIS.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// DeleteByPrefix 根据前缀删除key
func DeleteByPrefix(ctx context.Context, prefix string) error {
	return KeysPattern(ctx, fmt.Sprintf("%s*", prefix))
}
