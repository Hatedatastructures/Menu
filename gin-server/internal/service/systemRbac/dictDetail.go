package systemRbac

import (
	"context"
	"encoding/json"
	"errors"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DictDetailService struct{}

// 获取字典详情列表-后台使用
func (s *DictDetailService) GetDictionaryDetailList(
	ctx *gin.Context,
	r req.GetDictionaryDetailListReq,
) (rs res.GetDictionaryDetailListRes, err error) {
	// 分页参数
	limit := r.Size
	if limit <= 0 {
		limit = 10
	}
	offset := limit * (r.Page - 1)
	if offset < 0 {
		offset = 0
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.SysDictionaryDetail{})
	if r.Label != "" {
		db = db.Where("label LIKE ?", "%"+r.Label+"%")
	}
	if r.Value != 0 {
		db = db.Where("value = ?", r.Value)
	}
	if r.SysDictionaryID > 0 {
		db = db.Where("sys_dictionary_id = ?", r.SysDictionaryID)
	}

	// 统计总数
	err = db.Count(&rs.Total).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	// 查询列表
	var details []systemRbac.SysDictionaryDetail
	err = db.Order("sort ASC, id ASC").Limit(limit).Offset(offset).Find(&details).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情列表失败")
	}

	// 转换响应
	rs.List = make([]res.GetDictionaryDetailListResList, 0, len(details))
	for _, d := range details {
		rs.List = append(rs.List, res.GetDictionaryDetailListResList{
			Id:             int(d.ID),
			Label:          d.Label,
			Value:          d.Value,
			Extend:         d.Extend,
			Status:         d.Status != nil && *d.Status,
			Sort:           d.Sort,
			SysDictionaryID: d.SysDictionaryID,
			CreatedAt:      d.CreatedAt,
		})
	}
	return rs, nil
}

// 创建字典详情-后台使用
func (s *DictDetailService) CreateDictionaryDetail(
	ctx *gin.Context,
	r req.CreateDictionaryDetailReq,
) (rs res.CreateDictionaryDetailRes, err error) {
	// 参数校验
	if r.Label == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "展示值不能为空")
	}
	if r.SysDictionaryID <= 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "关联字典ID不能为空")
	}

	// 检查字典是否存在
	var dict systemRbac.SysDictionary
	err = global.GVA_DB.Where("id = ?", r.SysDictionaryID).First(&dict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.DICT_NOT_FOUND, "字典不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "检查字典失败")
	}

	// ✅ 检查字典值是否已存在（在同一字典下）
	var existingDetail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("sys_dictionary_id = ? AND value = ?", r.SysDictionaryID, r.Value).First(&existingDetail).Error
	if err == nil {
		// 找到了已存在的记录，说明value重复
		return rs, biz_err.New(biz_err.DICT_VALUE_DUPLICATE, "字典值已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库查询错误
		return rs, biz_err.New(biz_err.DB_ERROR, "检查字典值失败")
	}

	// 创建字典详情
	detail := systemRbac.SysDictionaryDetail{
		Label:          r.Label,
		Value:          r.Value,
		Extend:         r.Extend,
		Status:         &r.Status,
		Sort:           r.Sort,
		SysDictionaryID: r.SysDictionaryID,
	}
	err = global.GVA_DB.Create(&detail).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建字典详情失败")
	}

	// 清除缓存
	s.deleteDictDetailCache(context.Background(), dict.Type)

	rs = res.CreateDictionaryDetailRes{Id: int(detail.ID)}
	return rs, nil
}

// 字典详情详情-后台使用
func (s *DictDetailService) GetDictionaryDetail(
	ctx *gin.Context,
	r req.GetDictionaryDetailReq,
) (rs res.GetDictionaryDetailRes, err error) {
	var detail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("id = ?", r.Id).First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "字典详情不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	rs = res.GetDictionaryDetailRes{
		Id:             int(detail.ID),
		Label:          detail.Label,
		Value:          detail.Value,
		Extend:         detail.Extend,
		Status:         detail.Status != nil && *detail.Status,
		Sort:           detail.Sort,
		SysDictionaryID: detail.SysDictionaryID,
		CreatedAt:      detail.CreatedAt,
	}
	return rs, nil
}

// 更新字典详情-后台使用
func (s *DictDetailService) UpdateDictionaryDetail(
	ctx *gin.Context,
	r req.UpdateDictionaryDetailReq,
) (err error) {
	var detail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("id = ?", r.Id).First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "字典详情不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	// ✅ 检查字典值是否与其他记录重复（在同一字典下，排除当前记录）
	if r.Value != 0 && r.Value != detail.Value {
		var existingDetail systemRbac.SysDictionaryDetail
		dictID := detail.SysDictionaryID
		if r.SysDictionaryID > 0 {
			dictID = r.SysDictionaryID
		}

		err = global.GVA_DB.Where("sys_dictionary_id = ? AND value = ? AND id != ?", dictID, r.Value, r.Id).First(&existingDetail).Error
		if err == nil {
			// 找到了已存在的记录，说明value重复
			return biz_err.New(biz_err.DICT_VALUE_DUPLICATE, "字典值已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 数据库查询错误
			return biz_err.New(biz_err.DB_ERROR, "检查字典值失败")
		}
	}

	// 更新字段
	if r.Label != "" {
		detail.Label = r.Label
	}
	if r.Extend != "" {
		detail.Extend = r.Extend
	}
	detail.Value = r.Value
	detail.Status = &r.Status
	detail.Sort = r.Sort
	if r.SysDictionaryID > 0 {
		detail.SysDictionaryID = r.SysDictionaryID
	}

	err = global.GVA_DB.Save(&detail).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新字典详情失败")
	}

	// 清除缓存
	var dict systemRbac.SysDictionary
	global.GVA_DB.Where("id = ?", detail.SysDictionaryID).First(&dict)
	s.deleteDictDetailCache(context.Background(), dict.Type)

	return nil
}

// 删除字典详情-后台使用
func (s *DictDetailService) DeleteDictionaryDetail(
	ctx *gin.Context,
	r req.DeleteDictionaryDetailReq,
) (err error) {
	var detail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("id = ?", r.Id).First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "字典详情不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	// 获取字典type用于清除缓存
	var dict systemRbac.SysDictionary
	global.GVA_DB.Where("id = ?", detail.SysDictionaryID).First(&dict)

	// 删除
	err = global.GVA_DB.Delete(&detail).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除字典详情失败")
	}

	// 清除缓存
	s.deleteDictDetailCache(context.Background(), dict.Type)
	return nil
}

// 按字典ID获取字典全部内容-后台使用
func (s *DictDetailService) GetDictionaryListById(
	ctx *gin.Context,
	r req.GetDictionaryListByIdReq,
) (rs res.GetDictionaryListByIdRes, err error) {
	var details []systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("sys_dictionary_id = ? AND status = ?", r.DictionaryId, true).
		Order("sort ASC").Find(&details).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	rs.List = make([]res.GetDictionaryListByIdResList, 0, len(details))
	for _, d := range details {
		rs.List = append(rs.List, res.GetDictionaryListByIdResList{
			Id:     int(d.ID),
			Label:  d.Label,
			Value:  d.Value,
			Extend: d.Extend,
			Status: d.Status != nil && *d.Status,
			Sort:   d.Sort,
		})
	}
	return rs, nil
}

// 按字典type获取字典全部内容-后台使用
func (s *DictDetailService) GetDictionaryListByType(
	ctx *gin.Context,
	r req.GetDictionaryListByTypeReq,
) (rs res.GetDictionaryListByTypeRes, err error) {
	// 尝试从缓存获取
	list, _ := s.getDictDetailFromCache(ctx, dictCachePrefix+r.Type)
	if list != nil {
		rs.List = make([]res.GetDictionaryListByTypeResList, 0, len(list))
		for _, d := range list {
			rs.List = append(rs.List, res.GetDictionaryListByTypeResList{
				Id:     int(d.ID),
				Label:  d.Label,
				Value:  d.Value,
				Extend: d.Extend,
				Status: d.Status != nil && *d.Status,
				Sort:   d.Sort,
			})
		}
		return rs, nil
	}

	// 缓存未命中，从数据库查询
	var details []systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Model(&systemRbac.SysDictionaryDetail{}).
		Joins("JOIN sys_rbac_dictionaries ON sys_rbac_dictionaries.id = sys_rbac_dictionary_details.sys_dictionary_id").
		Where("sys_rbac_dictionaries.type = ? AND sys_rbac_dictionary_details.status = ?", r.Type, true).
		Order("sys_rbac_dictionary_details.sort").Find(&details).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	// 设置缓存
	s.setDictDetailCache(ctx, dictCachePrefix+r.Type, details)

	// 构建响应
	rs.List = make([]res.GetDictionaryListByTypeResList, 0, len(details))
	for _, d := range details {
		rs.List = append(rs.List, res.GetDictionaryListByTypeResList{
			Id:     int(d.ID),
			Label:  d.Label,
			Value:  d.Value,
			Extend: d.Extend,
			Status: d.Status != nil && *d.Status,
			Sort:   d.Sort,
		})
	}
	return rs, nil
}

// 按字典ID和value获取单条字典内容-后台使用
func (s *DictDetailService) GetDictionaryInfoByValue(
	ctx *gin.Context,
	r req.GetDictionaryInfoByValueReq,
) (rs res.GetDictionaryInfoByValueRes, err error) {
	var detail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Where("sys_dictionary_id = ? AND value = ?", r.DictionaryId, r.Value).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "字典详情不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	rs = res.GetDictionaryInfoByValueRes{
		Id:     int(detail.ID),
		Label:  detail.Label,
		Value:  detail.Value,
		Extend: detail.Extend,
		Status: detail.Status != nil && *detail.Status,
		Sort:   detail.Sort,
	}
	return rs, nil
}

// 按字典type和value获取单条字典内容-后台使用
func (s *DictDetailService) GetDictionaryInfoByTypeValue(
	ctx *gin.Context,
	r req.GetDictionaryInfoByTypeValueReq,
) (rs res.GetDictionaryInfoByTypeValueRes, err error) {
	var detail systemRbac.SysDictionaryDetail
	err = global.GVA_DB.Model(&systemRbac.SysDictionaryDetail{}).
		Joins("JOIN sys_rbac_dictionaries ON sys_rbac_dictionaries.id = sys_rbac_dictionary_details.sys_dictionary_id").
		Where("sys_rbac_dictionaries.type = ? AND sys_rbac_dictionary_details.value = ?", r.Type, r.Value).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "字典详情不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询字典详情失败")
	}

	rs = res.GetDictionaryInfoByTypeValueRes{
		Id:     int(detail.ID),
		Label:  detail.Label,
		Value:  detail.Value,
		Extend: detail.Extend,
		Status: detail.Status != nil && *detail.Status,
		Sort:   detail.Sort,
	}
	return rs, nil
}

// 辅助方法：从缓存获取字典详情列表
func (s *DictDetailService) getDictDetailFromCache(ctx context.Context, key string) ([]systemRbac.SysDictionaryDetail, error) {
	if global.GVA_REDIS == nil {
		return nil, nil
	}
	data, err := global.GVA_REDIS.Get(ctx, key).Bytes()
	if err != nil {
		return nil, nil
	}
	var list []systemRbac.SysDictionaryDetail
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// 辅助方法：设置字典详情缓存
func (s *DictDetailService) setDictDetailCache(ctx context.Context, key string, data []systemRbac.SysDictionaryDetail) error {
	if global.GVA_REDIS == nil {
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return global.GVA_REDIS.Set(ctx, key, jsonData, 24*60*60).Err()
}

// 辅助方法：删除字典详情缓存
func (s *DictDetailService) deleteDictDetailCache(ctx context.Context, dictType string) {
	if global.GVA_REDIS != nil && dictType != "" {
		global.GVA_REDIS.Del(ctx, dictCachePrefix+dictType, dictAllCacheKey)
	}
}