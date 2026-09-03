package online

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"shack/internal/global"
)

const (
	// 用户在线状态相关常量
	UserOnlineKeyPrefix    = "chat:online:"
	UserOnlineExpire      = 300 // 5分钟
)

// UserOnlineManager 用户在线状态管理器
type UserOnlineManager struct {
	redis redis.UniversalClient
}

// NewUserOnlineManager 创建用户在线状态管理器
func NewUserOnlineManager() *UserOnlineManager {
	return &UserOnlineManager{
		redis: global.GVA_REDIS,
	}
}

// SetOnline 设置用户在线状态
func (m *UserOnlineManager) SetOnline(ctx context.Context, userID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}
	key := m.buildKey(userID)
	return m.redis.Set(ctx, key, 1, UserOnlineExpire*time.Second).Err()
}

// SetOffline 设置用户离线状态
func (m *UserOnlineManager) SetOffline(ctx context.Context, userID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}
	key := m.buildKey(userID)
	return m.redis.Del(ctx, key).Err()
}

// IsOnline 检查用户是否在线
func (m *UserOnlineManager) IsOnline(ctx context.Context, userID int64) (bool, error) {
	if m.redis == nil {
		return false, fmt.Errorf("redis not initialized")
	}
	key := m.buildKey(userID)
	exists, err := m.redis.Exists(ctx, key).Result()
	return exists > 0, err
}

// RefreshOnline 刷新用户在线状态（重置过期时间）
func (m *UserOnlineManager) RefreshOnline(ctx context.Context, userID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}
	key := m.buildKey(userID)
	exists, err := m.redis.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return m.redis.Expire(ctx, key, UserOnlineExpire*time.Second).Err()
	}
	return nil
}

// GetOnlineUsers 获取所有在线用户ID列表
func (m *UserOnlineManager) GetOnlineUsers(ctx context.Context) ([]int64, error) {
	if m.redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	pattern := fmt.Sprintf("%s*", UserOnlineKeyPrefix)
	var userIDs []int64

	iter := m.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		// 从key中提取userID
		var userID int64
		_, err := fmt.Sscanf(key, fmt.Sprintf("%s%%d", UserOnlineKeyPrefix), &userID)
		if err == nil {
			userIDs = append(userIDs, userID)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

// GetOnlineCount 获取在线用户数量
func (m *UserOnlineManager) GetOnlineCount(ctx context.Context) (int64, error) {
	if m.redis == nil {
		return 0, fmt.Errorf("redis not initialized")
	}

	pattern := fmt.Sprintf("%s*", UserOnlineKeyPrefix)
	var count int64

	iter := m.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		count++
	}

	if err := iter.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// buildKey 构建Redis key
func (m *UserOnlineManager) buildKey(userID int64) string {
	return fmt.Sprintf("%s%d", UserOnlineKeyPrefix, userID)
}

// BatchSetOnline 批量设置用户在线状态
func (m *UserOnlineManager) BatchSetOnline(ctx context.Context, userIDs []int64) error {
	if m.redis == nil || len(userIDs) == 0 {
		return nil
	}

	pipe := m.redis.Pipeline()
	for _, userID := range userIDs {
		key := m.buildKey(userID)
		pipe.Set(ctx, key, 1, UserOnlineExpire*time.Second)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// BatchSetOffline 批量设置用户离线状态
func (m *UserOnlineManager) BatchSetOffline(ctx context.Context, userIDs []int64) error {
	if m.redis == nil || len(userIDs) == 0 {
		return nil
	}

	pipe := m.redis.Pipeline()
	for _, userID := range userIDs {
		key := m.buildKey(userID)
		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// BatchIsOnline 批量检查用户是否在线
func (m *UserOnlineManager) BatchIsOnline(ctx context.Context, userIDs []int64) (map[int64]bool, error) {
	if m.redis == nil || len(userIDs) == 0 {
		return make(map[int64]bool), nil
	}

	pipe := m.redis.Pipeline()
	cmds := make(map[int64]*redis.IntCmd)

	for _, userID := range userIDs {
		key := m.buildKey(userID)
		cmds[userID] = pipe.Exists(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[int64]bool)
	for userID, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			result[userID] = val > 0
		} else {
			result[userID] = false
		}
	}

	return result, nil
}
