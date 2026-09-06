package systemRbac

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Redis缓存key前缀
const (
	dictCachePrefix = "dict:"
	dictAllCacheKey = "dict:all"
)

type DictService struct{}

// 辅助方法：从缓存获取字典列表
func (s *DictService) getFromCache(ctx context.Context, key string) ([]systemRbac.SysDictionaryDetail, error) {
	if global.GVA_REDIS == nil {
		return nil, nil
	}
	data, err := global.GVA_REDIS.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, err
	}
	var list []systemRbac.SysDictionaryDetail
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// 辅助方法：设置缓存
func (s *DictService) setCache(ctx context.Context, key string, data interface{}) error {
	if global.GVA_REDIS == nil {
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 缓存24小时
	return global.GVA_REDIS.Set(ctx, key, jsonData, 24*60*60).Err()
}

// 辅助方法：删除缓存
func (s *DictService) deleteCache(ctx context.Context, keys ...string) error {
	if global.GVA_REDIS == nil {
		return nil
	}
	return global.GVA_REDIS.Del(ctx, keys...).Err()
}

// 获取字典列表-后台使用
func (s *DictService) GetDictionaryList(ctx *gin.Context) (rs res.GetDictionaryListRes, err error) {
	var dicts []systemRbac.SysDictionary
	err = global.GVA_DB.Order("id ASC").Find(&dicts).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典列表失败")
	}

	rs.List = make([]res.GetDictionaryListResList, 0, len(dicts))
	for _, d := range dicts {
		rs.List = append(rs.List, res.GetDictionaryListResList{
			Id:        int(d.ID),
			Name:      d.Name,
			Type:      d.Type,
			Status:    d.Status != nil && *d.Status,
			Desc:      d.Desc,
			CreatedAt: d.CreatedAt,
		})
	}
	return rs, nil
}

// 创建字典-后台使用
func (s *DictService) CreateDictionary(ctx *gin.Context, r req.CreateDictionaryReq) (rs res.CreateDictionaryRes, err error) {
	// 参数校验
	if r.Type == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "字典类型不能为空")
	}

	// 检查type是否已存在
	var existing systemRbac.SysDictionary
	err = global.GVA_DB.Where("type = ?", r.Type).First(&existing).Error
	if err == nil {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "字典类型已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return rs, biz_err.New(biz_err.DB_ERROR, "检查字典类型失败")
	}

	// 创建字典
	dict := systemRbac.SysDictionary{
		Name:   r.Name,
		Type:   r.Type,
		Status: &r.Status,
		Desc:   r.Desc,
	}
	err = global.GVA_DB.Create(&dict).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建字典失败")
	}

	// 清除全部字典缓存
	s.deleteCache(context.Background(), dictAllCacheKey)

	rs = res.CreateDictionaryRes{Id: int(dict.ID)}
	return rs, nil
}

// 字典详情-后台使用
func (s *DictService) GetDictionary(ctx *gin.Context, r req.GetDictionaryReq) (rs res.GetDictionaryRes, err error) {
	var dict systemRbac.SysDictionary
	err = global.GVA_DB.Where("id = ?", r.Id).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true).Order("sort")
	}).First(&dict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "字典不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典失败")
	}

	rs = res.GetDictionaryRes{
		Id:                int(dict.ID),
		Name:              dict.Name,
		Type:              dict.Type,
		Status:            dict.Status != nil && *dict.Status,
		Desc:              dict.Desc,
		SysDictionaryDetails: make([]res.GetDictionaryResSysdictionarydetail, 0),
		CreatedAt:         dict.CreatedAt,
	}
	for _, d := range dict.SysDictionaryDetails {
		rs.SysDictionaryDetails = append(rs.SysDictionaryDetails, res.GetDictionaryResSysdictionarydetail{
			Id:    int(d.ID),
			Label: d.Label,
			Value: d.Value,
			Extend: d.Extend,
			Status: d.Status != nil && *d.Status,
			Sort:  d.Sort,
		})
	}
	return rs, nil
}

// 更新字典-后台使用
func (s *DictService) UpdateDictionary(ctx *gin.Context, r req.UpdateDictionaryReq) (err error) {
	var dict systemRbac.SysDictionary
	err = global.GVA_DB.Where("id = ?", r.Id).First(&dict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "字典不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询字典失败")
	}

	// 检查type唯一性
	if r.Type != "" && r.Type != dict.Type {
		var existing systemRbac.SysDictionary
		err = global.GVA_DB.Where("type = ? AND id != ?", r.Type, r.Id).First(&existing).Error
		if err == nil {
			return biz_err.New(biz_err.PARAM_ERROR, "字典类型已存在")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.DB_ERROR, "检查字典类型失败")
		}
		dict.Type = r.Type
	}

	// 更新字段
	if r.Name != "" {
		dict.Name = r.Name
	}
	if r.Desc != "" {
		dict.Desc = r.Desc
	}
	dict.Status = &r.Status

	err = global.GVA_DB.Save(&dict).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新字典失败")
	}

	// 清除缓存
	s.deleteCache(context.Background(), dictAllCacheKey, dictCachePrefix+dict.Type)
	return nil
}

// 删除字典-后台使用
func (s *DictService) DeleteDictionary(ctx *gin.Context, r req.DeleteDictionaryReq) (err error) {
	var dict systemRbac.SysDictionary
	err = global.GVA_DB.Where("id = ?", r.Id).Preload("SysDictionaryDetails").First(&dict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "字典不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询字典失败")
	}

	// 删除字典
	err = global.GVA_DB.Delete(&dict).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除字典失败")
	}

	// 删除关联的字典详情
	if len(dict.SysDictionaryDetails) > 0 {
		global.GVA_DB.Where("sys_dictionary_id = ?", dict.ID).Delete(&[]systemRbac.SysDictionaryDetail{})
	}

	// 清除缓存
	s.deleteCache(context.Background(), dictAllCacheKey, dictCachePrefix+dict.Type)
	return nil
}

// 批量字典查询-前台/后台使用
func (s *DictService) GetDictsByTypes(ctx *gin.Context, r req.GetDictsByTypesReq) (rs res.GetDictsByTypesRes, err error) {
	if r.Types == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "types不能为空")
	}

	types := strings.Split(r.Types, ",")
	for _, t := range types {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}

		// 尝试从缓存获取
		list, _ := s.getFromCache(ctx, dictCachePrefix+t)
		if list == nil {
			// 缓存未命中，从数据库查询
			var details []systemRbac.SysDictionaryDetail
			err = global.GVA_DB.Model(&systemRbac.SysDictionaryDetail{}).
				Joins("JOIN sys_rbac_dictionaries ON sys_rbac_dictionaries.id = sys_rbac_dictionary_details.sys_dictionary_id").
				Where("sys_rbac_dictionaries.type = ? AND sys_rbac_dictionary_details.status = ?", t, true).
				Order("sys_rbac_dictionary_details.sort").Find(&details).Error
			if err != nil {
				continue
			}
			list = details
			// 设置缓存
			s.setCache(ctx, dictCachePrefix+t, list)
		}

		// 根据type设置到对应的字段
		switch t {
		case "status":
			rs.Status = make([]res.GetDictsByTypesResStatu, 0, len(list))
			for _, d := range list {
				rs.Status = append(rs.Status, res.GetDictsByTypesResStatu{Label: d.Label, Value: d.Value})
			}
		case "gender":
			rs.Gender = make([]res.GetDictsByTypesResGender, 0, len(list))
			for _, d := range list {
				rs.Gender = append(rs.Gender, res.GetDictsByTypesResGender{Label: d.Label, Value: d.Value})
			}
		default:
			// 动态字段处理，使用map存储
		}
	}
	return rs, nil
}

// 获取全部字典-前台/后台使用
func (s *DictService) GetAllDicts(ctx *gin.Context) (rs res.GetAllDictsRes, err error) {
	// 从数据库查询所有启用的字典详情和字典信息
	var details []systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Preload("SysDictionary", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true)
	}).Where("sys_rbac_dictionary_details.status = ?", true).Order("sys_dictionary_id, sort").Find(&details).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典失败")
	}

	// 按type分组
	typeGroups := make(map[string][]systemRbac.SysDictionaryDetail)
	for _, d := range details {
		if d.SysDictionary == nil || d.SysDictionary.Type == "" {
			continue
		}
		t := d.SysDictionary.Type
		typeGroups[t] = append(typeGroups[t], d)
	}

	// 设置缓存
	allDetails := make([]systemRbac.SysDictionaryDetail, 0, len(details))
	for _, d := range details {
		allDetails = append(allDetails, d)
	}
	s.setCache(ctx, dictAllCacheKey, allDetails)

	// 构建响应
	for t, list := range typeGroups {
		switch t {
		case "status":
			rs.Status = make([]res.GetAllDictsResStatu, 0, len(list))
			for _, d := range list {
				rs.Status = append(rs.Status, res.GetAllDictsResStatu{Label: d.Label, Value: d.Value})
			}
		case "gender":
			rs.Gender = make([]res.GetAllDictsResGender, 0, len(list))
			for _, d := range list {
				rs.Gender = append(rs.Gender, res.GetAllDictsResGender{Label: d.Label, Value: d.Value})
			}
		}
	}
	return rs, nil
}

// 字典缓存刷新-后台使用
func (s *DictService) RefreshDicts(ctx *gin.Context) (rs res.RefreshDictsRes, err error) {
	// 获取所有字典type
	var dicts []systemRbac.SysDictionary
	err = global.GVA_DB.Where("status = ?", true).Find(&dicts).Error
	if err != nil {
		rs.Success = false
		rs.Message = "查询字典失败"
		return rs, nil
	}

	// 清除所有字典缓存
	keys := []string{dictAllCacheKey}
	for _, d := range dicts {
		keys = append(keys, dictCachePrefix+d.Type)
	}
	s.deleteCache(ctx, keys...)

	rs.Success = true
	rs.Message = "刷新成功"
	return rs, nil
}