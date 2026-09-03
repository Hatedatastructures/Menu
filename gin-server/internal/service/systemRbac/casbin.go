package systemRbac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"shack/internal/global"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"shack/internal/utils"

	biz_err "shack/internal/error"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// 权限缓存 key 前缀
	cacheKeyPrefix = "casbin:policy:"
	// 缓存过期时间 (1小时，与 casbin 同步)
	cacheExpire = 60 * 60
)

// CasbinService Casbin 权限服务
type CasbinService struct {
	enforcer *casbin.SyncedCachedEnforcer
}

// NewCasbinService 创建 CasbinService 实例
func NewCasbinService() *CasbinService {
	return &CasbinService{
		enforcer: nil,
	}
}

var CasbinServiceApp = NewCasbinService()

// ============ 私有辅助方法 ============

// getEnforcer 获取 casbin enforcer
func (s *CasbinService) getEnforcer() *casbin.SyncedCachedEnforcer {
	if s.enforcer == nil {
		s.enforcer = utils.GetCasbin()
	}
	return s.enforcer
}

// buildCacheKey 构建缓存 key
func (s *CasbinService) buildCacheKey(authorityId uint) string {
	return fmt.Sprintf("%s%d", cacheKeyPrefix, authorityId)
}

// getFromCache 从缓存获取权限列表
func (s *CasbinService) getFromCache(ctx context.Context, authorityId uint) ([]res.GetPolicyPathByAuthorityIdResList, error) {
	if global.GVA_REDIS == nil {
		return nil, nil
	}

	key := s.buildCacheKey(authorityId)
	data, err := global.GVA_REDIS.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, err
	}

	var result []res.GetPolicyPathByAuthorityIdResList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// setCache 设置缓存
func (s *CasbinService) setCache(ctx context.Context, authorityId uint, data []res.GetPolicyPathByAuthorityIdResList) error {
	if global.GVA_REDIS == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := s.buildCacheKey(authorityId)
	return global.GVA_REDIS.Set(ctx, key, jsonData, cacheExpire).Err()
}

// deleteCache 删除缓存
func (s *CasbinService) deleteCache(ctx context.Context, authorityId uint) error {
	if global.GVA_REDIS == nil {
		return nil
	}

	key := s.buildCacheKey(authorityId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

// validateApis 验证 API 是否在权限列表中 (严格模式)
func (s *CasbinService) validateApis(adminAuthorityID uint, casbinInfos []req.UpdateCasbinReqCasbininfo) error {
	if !global.GVA_CONFIG.System.UseStrictAuth {
		return nil
	}

	apis, err := ApiServiceApp.GetAllApis(adminAuthorityID)
	if err != nil {
		return err
	}

	apiMap := make(map[string]bool)
	for _, api := range apis {
		apiMap[api.Path+api.Method] = true
	}

	for _, info := range casbinInfos {
		if !apiMap[info.Path+info.Method] {
			return errors.New("存在API不在权限列表中: " + info.Path + " " + info.Method)
		}
	}
	return nil
}

// deduplicateRules 权限去重
func (s *CasbinService) deduplicateRules(authorityId string, casbinInfos []req.UpdateCasbinReqCasbininfo) [][]string {
	deduplicateMap := make(map[string]bool)
	var rules [][]string

	for _, v := range casbinInfos {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	return rules
}

// convertToResList 转换为响应结构
func (s *CasbinService) convertToResList(pathMaps []req.UpdateCasbinReqCasbininfo) []res.GetPolicyPathByAuthorityIdResList {
	result := make([]res.GetPolicyPathByAuthorityIdResList, 0, len(pathMaps))
	for _, v := range pathMaps {
		result = append(result, res.GetPolicyPathByAuthorityIdResList{
			Path:   v.Path,
			Method: v.Method,
		})
	}
	return result
}

// ============ 公开服务方法 ============

// 更新权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月01日 17:04:13
func (s *CasbinService) UpdateCasbin(
	ctx *gin.Context,
	r req.UpdateCasbinReq,
) (err error) {
	// 参数校验
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "权限ID不能为空")
	}
	if len(r.CasbinInfos) == 0 {
		return nil // 空权限无需处理
	}

	// 验证 API 权限 (严格模式)
	if err := s.validateApis(ctx.GetUint("authorityId"), r.CasbinInfos); err != nil {
		return biz_err.New(biz_err.PARAM_ERROR, err.Error())
	}

	// 清除旧权限
	e := s.getEnforcer()
	if e == nil {
		return biz_err.New(biz_err.DB_ERROR, "Casbin未初始化")
	}
	authorityIdStr := strconv.Itoa(int(r.AuthorityId))
	_, _ = e.RemoveFilteredPolicy(0, authorityIdStr)

	// 权限去重
	rules := s.deduplicateRules(authorityIdStr, r.CasbinInfos)

	if len(rules) == 0 {
		return nil
	}

	// 添加新权限
	success, err := e.AddPolicies(rules)
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "添加权限失败: "+err.Error())
	}
	if !success {
		return biz_err.New(biz_err.DB_ERROR, "存在相同API，添加失败")
	}

	// 删除缓存
	_ = s.deleteCache(ctx.Request.Context(), r.AuthorityId)

	return nil
}

// 获取权限列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月01日 17:04:13
func (s *CasbinService) GetPolicyPathByAuthorityId(
	ctx *gin.Context,
	r req.GetPolicyPathByAuthorityIdReq,
) (rs res.GetPolicyPathByAuthorityIdRes, err error) {
	// 参数校验
	if r.AuthorityId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "权限ID不能为空")
	}

	// 尝试从缓存获取
	cacheCtx := ctx.Request.Context()
	if cached, err := s.getFromCache(cacheCtx, r.AuthorityId); err == nil && cached != nil {
		rs.List = cached
		return rs, nil
	}

	// 从 casbin 获取
	e := s.getEnforcer()
	if e == nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "Casbin未初始化")
	}
	authorityId := strconv.Itoa(int(r.AuthorityId))
	list, err := e.GetFilteredPolicy(0, authorityId)
	if err != nil {
		zap.L().Error("获取权限列表失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR, "获取权限列表失败")
	}

	// 转换为响应结构
	pathMaps := make([]req.UpdateCasbinReqCasbininfo, 0, len(list))
	for _, v := range list {
		pathMaps = append(pathMaps, req.UpdateCasbinReqCasbininfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	rs.List = s.convertToResList(pathMaps)

	// 设置缓存
	_ = s.setCache(cacheCtx, r.AuthorityId, rs.List)

	return rs, nil
}

// 清除权限-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月01日 17:04:13
func (s *CasbinService) ClearCasbin(
	ctx *gin.Context,
	r req.ClearCasbinReq,
) (err error) {
	// 参数校验
	if r.AuthorityId == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "权限ID不能为空")
	}

	// 清除权限
	e := s.getEnforcer()
	if e == nil {
		return biz_err.New(biz_err.DB_ERROR, "Casbin未初始化")
	}
	authorityId := strconv.Itoa(int(r.AuthorityId))
	_, err = e.RemoveFilteredPolicy(0, authorityId)
	if err != nil {
		zap.L().Error("清除权限失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR, "清除权限失败")
	}

	// 删除缓存
	_ = s.deleteCache(ctx.Request.Context(), r.AuthorityId)

	return nil
}

// 刷新Casbin-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月01日 17:04:13
func (s *CasbinService) FreshCasbin(
	ctx *gin.Context,
) (err error) {
	// 刷新 casbin 策略
	e := s.getEnforcer()
	if e == nil {
		return biz_err.New(biz_err.DB_ERROR, "Casbin未初始化")
	}
	if err := e.LoadPolicy(); err != nil {
		zap.L().Error("刷新Casbin失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR, "刷新Casbin失败")
	}

	// 清除所有缓存
	if global.GVA_REDIS != nil {
		keys, err := global.GVA_REDIS.Keys(ctx.Request.Context(), cacheKeyPrefix+"*").Result()
		if err == nil && len(keys) > 0 {
			_ = global.GVA_REDIS.Del(ctx.Request.Context(), keys...).Err()
		}
	}

	return nil
}

// ============ 内部使用的方法 (供其他服务调用) ============

// ClearCasbinDB 使用数据库方式清除权限 (需要 tx)
func (s *CasbinService) ClearCasbinDB(tx *gorm.DB, authorityId string) error {
	return tx.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityId).Error
}

// AddPoliciesDB 批量添加权限 (需要 tx)
func (s *CasbinService) AddPoliciesDB(tx *gorm.DB, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	var casbinRules []gormadapter.CasbinRule
	for _, rule := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rule[0],
			V1:    rule[1],
			V2:    rule[2],
		})
	}
	return tx.Create(&casbinRules).Error
}

// SyncPolicyDB 同步策略到数据库 (需要 tx)
func (s *CasbinService) SyncPolicyDB(tx *gorm.DB, authorityId string, rules [][]string) error {
	if err := s.ClearCasbinDB(tx, authorityId); err != nil {
		return err
	}
	return s.AddPoliciesDB(tx, rules)
}
