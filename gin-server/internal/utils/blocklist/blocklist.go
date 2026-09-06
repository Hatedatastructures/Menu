package blocklist

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"shack/internal/global"
)

const (
	// 黑名单缓存相关常量
	BlacklistKeyPrefix = "chat:blacklist:"
	BlacklistExpire    = 1800 // 30分钟
)

// BlocklistManager 黑名单管理器
type BlocklistManager struct {
	redis redis.UniversalClient
}

// NewBlocklistManager 创建黑名单管理器
func NewBlocklistManager() *BlocklistManager {
	return &BlocklistManager{
		redis: global.GVA_REDIS,
	}
}

// IsBlocked 检查是否被拉黑
func (m *BlocklistManager) IsBlocked(ctx context.Context, userID, targetUserID int64) (bool, error) {
	if m.redis == nil {
		return false, fmt.Errorf("redis not initialized")
	}

	// 先从缓存获取
	blocked, err := m.checkFromCache(ctx, userID, targetUserID)
	if err == nil {
		return blocked, nil
	}

	// 缓存未命中，返回false（默认未拉黑）
	return false, nil
}

// AddToBlacklist 添加到黑名单
func (m *BlocklistManager) AddToBlacklist(ctx context.Context, userID, targetUserID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	// 添加到Redis集合
	key := m.buildKey(userID)
	return m.redis.SAdd(ctx, key, targetUserID).Err()
}

// RemoveFromBlacklist 从黑名单移除
func (m *BlocklistManager) RemoveFromBlacklist(ctx context.Context, userID, targetUserID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	// 从Redis集合移除
	key := m.buildKey(userID)
	return m.redis.SRem(ctx, key, targetUserID).Err()
}

// GetBlacklist 获取黑名单列表
func (m *BlocklistManager) GetBlacklist(ctx context.Context, userID int64) ([]int64, error) {
	if m.redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	key := m.buildKey(userID)
	members, err := m.redis.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 转换字符串为int64
	result := make([]int64, 0, len(members))
	for _, member := range members {
		var id int64
		if _, err := fmt.Sscanf(member, "%d", &id); err == nil {
			result = append(result, id)
		}
	}

	return result, nil
}

// ClearBlacklist 清空黑名单
func (m *BlocklistManager) ClearBlacklist(ctx context.Context, userID int64) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	key := m.buildKey(userID)
	return m.redis.Del(ctx, key).Err()
}

// InvalidateCache 使缓存失效
func (m *BlocklistManager) InvalidateCache(ctx context.Context, userID int64) error {
	if m.redis == nil {
		return nil
	}

	key := m.buildKey(userID)
	return m.redis.Del(ctx, key).Err()
}

// BatchIsBlocked 批量检查是否被拉黑
func (m *BlocklistManager) BatchIsBlocked(ctx context.Context, userID int64, targetUserIDs []int64) (map[int64]bool, error) {
	if m.redis == nil || len(targetUserIDs) == 0 {
		result := make(map[int64]bool)
		for _, id := range targetUserIDs {
			result[id] = false
		}
		return result, nil
	}

	key := m.buildKey(userID)
	members, err := m.redis.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 构建被拉黑用户的集合
	blockedSet := make(map[int64]bool)
	for _, member := range members {
		var id int64
		if _, err := fmt.Sscanf(member, "%d", &id); err == nil {
			blockedSet[id] = true
		}
	}

	// 检查每个目标用户是否被拉黑
	result := make(map[int64]bool)
	for _, targetID := range targetUserIDs {
		result[targetID] = blockedSet[targetID]
	}

	return result, nil
}

// GetBlacklistCount 获取黑名单数量
func (m *BlocklistManager) GetBlacklistCount(ctx context.Context, userID int64) (int64, error) {
	if m.redis == nil {
		return 0, fmt.Errorf("redis not initialized")
	}

	key := m.buildKey(userID)
	return m.redis.SCard(ctx, key).Result()
}

// AddToBlacklistWithReason 添加到黑名单（带原因）
func (m *BlocklistManager) AddToBlacklistWithReason(ctx context.Context, userID, targetUserID int64, reason string) error {
	if m.redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	// 添加到Redis集合
	key := m.buildKey(userID)
	if err := m.redis.SAdd(ctx, key, targetUserID).Err(); err != nil {
		return err
	}

	// 如果有原因，额外存储
	if reason != "" {
		reasonKey := m.buildReasonKey(userID, targetUserID)
		return m.redis.Set(ctx, reasonKey, reason, BlacklistExpire*time.Second).Err()
	}

	return nil
}

// GetBlockedReason 获取拉黑原因
func (m *BlocklistManager) GetBlockedReason(ctx context.Context, userID, targetUserID int64) (string, error) {
	if m.redis == nil {
		return "", fmt.Errorf("redis not initialized")
	}

	reasonKey := m.buildReasonKey(userID, targetUserID)
	reason, err := m.redis.Get(ctx, reasonKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	return reason, err
}

// SyncBlacklistFromDB 从数据库同步黑名单到Redis
func (m *BlocklistManager) SyncBlacklistFromDB(ctx context.Context, userID int64, blockedUserIDs []int64) error {
	if m.redis == nil || len(blockedUserIDs) == 0 {
		return nil
	}

	key := m.buildKey(userID)
	// 清空现有黑名单
	if err := m.redis.Del(ctx, key).Err(); err != nil {
		return err
	}

	// 批量添加新黑名单
	members := make([]interface{}, len(blockedUserIDs))
	for i, id := range blockedUserIDs {
		members[i] = fmt.Sprintf("%d", id)
	}

	return m.redis.SAdd(ctx, key, members...).Err()
}

// checkFromCache 从缓存检查是否被拉黑
func (m *BlocklistManager) checkFromCache(ctx context.Context, userID, targetUserID int64) (bool, error) {
	key := m.buildKey(userID)
	exists, err := m.redis.SIsMember(ctx, key, fmt.Sprintf("%d", targetUserID)).Result()
	return exists, err
}

// buildKey 构建Redis key
func (m *BlocklistManager) buildKey(userID int64) string {
	return fmt.Sprintf("%s%d", BlacklistKeyPrefix, userID)
}

// buildReasonKey 构建拉黑原因的Redis key
func (m *BlocklistManager) buildReasonKey(userID, targetUserID int64) string {
	return fmt.Sprintf("%s%d:reason:%d", BlacklistKeyPrefix, userID, targetUserID)
}

// ExportBlacklist 导出黑名单为JSON
func (m *BlocklistManager) ExportBlacklist(ctx context.Context, userID int64) (string, error) {
	blacklist, err := m.GetBlacklist(ctx, userID)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(blacklist)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ImportBlacklist 从JSON导入黑名单
func (m *BlocklistManager) ImportBlacklist(ctx context.Context, userID int64, jsonData string) error {
	var blacklist []int64
	if err := json.Unmarshal([]byte(jsonData), &blacklist); err != nil {
		return err
	}

	return m.SyncBlacklistFromDB(ctx, userID, blacklist)
}
