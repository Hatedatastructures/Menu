package systemRbac

import (
	"context"
	"encoding/json"
	"shack/internal/global"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	SecurityConfigCacheKey = "system:security:config"
	SecurityConfigTTL      = 5 * time.Minute // 缓存5分钟
)

type ConfigEncryptService struct{}

// GetSecurityConfig 获取安全配置(带Redis缓存)
func (s *ConfigEncryptService) GetSecurityConfig(ctx *gin.Context) (res.GetSecurityConfigRes, error) {
	// 尝试从Redis获取
	cacheKey := SecurityConfigCacheKey
	rdb := global.GVA_REDIS

	if rdb == nil {
		// Redis未初始化，直接从数据库获取
		return s.getSecurityConfigFromDB(ctx)
	}

	// 先从缓存获取
	cachedData, err := rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedData != "" {
		var config res.GetSecurityConfigRes
		if json.Unmarshal([]byte(cachedData), &config) == nil {
			return config, nil
		}
	}

	// 缓存未命中，从数据库获取
	config, err := s.getSecurityConfigFromDB(ctx)
	if err != nil {
		return res.GetSecurityConfigRes{}, err
	}

	// 存入Redis缓存
	if jsonBytes, err := json.Marshal(config); err == nil {
		rdb.Set(context.Background(), cacheKey, jsonBytes, SecurityConfigTTL)
	}

	return config, nil
}

// getSecurityConfigFromDB 从数据库获取安全配置
func (s *ConfigEncryptService) getSecurityConfigFromDB(ctx *gin.Context) (res.GetSecurityConfigRes, error) {
	// 直接调用 ConfigSecurityService 的方法
	css := &ConfigSecurityService{}
	return css.GetSecurityConfig(ctx)
}

// SetSecurityConfig 保存安全配置并更新缓存
func (s *ConfigEncryptService) SetSecurityConfig(ctx *gin.Context, req req.SaveSecurityConfigReq) (res.SaveSecurityConfigRes, error) {
	// 直接调用 ConfigSecurityService 的方法保存到数据库
	css := &ConfigSecurityService{}
	result, err := css.SaveSecurityConfig(ctx, req)
	if err != nil {
		return result, err
	}

	// 更新Redis缓存
	s.refreshCache(ctx)

	return result, nil
}

// refreshCache 刷新缓存
func (s *ConfigEncryptService) refreshCache(ctx *gin.Context) {
	cacheKey := SecurityConfigCacheKey
	rdb := global.GVA_REDIS

	if rdb == nil {
		return
	}

	// 先删除旧缓存
	rdb.Del(context.Background(), cacheKey)

	// 重新从数据库获取并缓存
	config, err := s.getSecurityConfigFromDB(ctx)
	if err != nil {
		return
	}
	if jsonBytes, err := json.Marshal(config); err == nil {
		rdb.Set(context.Background(), cacheKey, jsonBytes, SecurityConfigTTL)
	}
}

// GetAESKey 获取AES加密密钥
func (s *ConfigEncryptService) GetAESKey(ctx *gin.Context) (string, error) {
	config, err := s.GetSecurityConfig(ctx)
	if err != nil {
		return "", err
	}
	// 从配置中获取AES密钥(这里使用RSA私钥作为AES密钥来源，实际应单独配置AES密钥)
	return config.EncryptPrivateKey, nil
}

// IsEncryptEnabled 检查是否启用了加密
func (s *ConfigEncryptService) IsEncryptEnabled(ctx *gin.Context) bool {
	config, err := s.GetSecurityConfig(ctx)
	if err != nil {
		return false
	}
	return config.EncryptEnabled
}

// GetEncryptScope 获取加密范围
func (s *ConfigEncryptService) GetEncryptScope(ctx *gin.Context) string {
	config, err := s.GetSecurityConfig(ctx)
	if err != nil {
		return "none"
	}
	return config.EncryptScope
}

// 便捷函数: 直接从Redis获取配置(不经过gin context)
func GetSecurityConfigFromCache() (res.GetSecurityConfigRes, error) {
	cacheKey := SecurityConfigCacheKey
	rdb := global.GVA_REDIS
	if rdb == nil {
		return res.GetSecurityConfigRes{
			EncryptEnabled: false,
			EncryptScope:   "none",
		}, nil
	}

	cachedData, err := rdb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		// 缓存不存在，返回默认配置
		return res.GetSecurityConfigRes{
			EncryptEnabled: false,
			EncryptScope:   "none",
		}, nil
	} else if err != nil {
		return res.GetSecurityConfigRes{}, err
	}

	var config res.GetSecurityConfigRes
	if err := json.Unmarshal([]byte(cachedData), &config); err != nil {
		return res.GetSecurityConfigRes{}, err
	}
	return config, nil
}

// ClearSecurityConfigCache 清除安全配置缓存(仅供管理后台使用)
func ClearSecurityConfigCache() error {
	cacheKey := SecurityConfigCacheKey
	rdb := global.GVA_REDIS
	return rdb.Del(context.Background(), cacheKey).Err()
}