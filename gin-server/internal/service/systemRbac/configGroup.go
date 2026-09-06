package systemRbac

import (
	"encoding/json"
	"strconv"
	"strings"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ConfigGroupService struct{}

// CreateConfigGroup 创建配置分组-后台使用
func (s *ConfigGroupService) CreateConfigGroup(
	ctx *gin.Context,
	r req.CreateConfigGroupReq,
) (rs res.CreateConfigGroupRes, err error) {
	// 参数校验
	if r.Code == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "分组编码不能为空")
	}
	if r.Name == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "分组名称不能为空")
	}

	// 检查编码是否已存在
	var count int64
	global.GVA_DB.Model(&systemRbac.ConfigGroup{}).Where("code = ?", r.Code).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "分组编码已存在")
	}

	// 构建Config
	var configJson datatypes.JSON
	if r.Config != "" {
		configJson = datatypes.JSON(r.Config)
	}

	// 创建分组
	group := systemRbac.ConfigGroup{
		Code:   r.Code,
		Name:   r.Name,
		Config: configJson,
		Sort:   r.Sort,
		Status: r.Status,
		Remark: r.Remark,
	}
	err = global.GVA_DB.Create(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建分组失败")
	}

	rs = res.CreateConfigGroupRes{
		Id: strconv.FormatUint(uint64(group.ID), 10),
	}
	return rs, nil
}

// GetConfigGroupList 配置分组列表-后台使用
func (s *ConfigGroupService) GetConfigGroupList(
	ctx *gin.Context,
	r req.GetConfigGroupListReq,
) (rs res.GetConfigGroupListRes, err error) {
	// 设置默认值
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.ConfigGroup{})
	if r.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+r.Keyword+"%")
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 分页查询
	var groups []systemRbac.ConfigGroup
	offset := (r.Page - 1) * r.Size
	err = db.Offset(offset).Limit(r.Size).Order("sort ASC, id DESC").Find(&groups).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询分组列表失败")
	}

	// 构建返回数据
	rs.List = make([]res.GetConfigGroupListResList, 0, len(groups))
	for _, group := range groups {
		rs.List = append(rs.List, res.GetConfigGroupListResList{
			Id:        strconv.FormatUint(uint64(group.ID), 10),
			Code:      group.Code,
			Name:      group.Name,
			Sort:      group.Sort,
			Status:    int(group.Status),
			Remark:    group.Remark,
			CreatedAt: group.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: group.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	rs.Total = total
	rs.Page = r.Page
	rs.Size = r.Size
	rs.Keyword = r.Keyword
	return rs, nil
}

// GetConfigGroupDetail 获取配置分组详情-后台使用
func (s *ConfigGroupService) GetConfigGroupDetail(
	ctx *gin.Context,
	r req.GetConfigGroupDetailReq,
) (rs res.GetConfigGroupDetailRes, err error) {
	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "无效的分组ID")
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "分组不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询分组详情失败")
	}

	rs = res.GetConfigGroupDetailRes{
		Id:        strconv.FormatUint(uint64(group.ID), 10),
		Code:      group.Code,
		Name:      group.Name,
		Config:    string(group.Config),
		Sort:      group.Sort,
		Status:    group.Status,
		Remark:    group.Remark,
		Version:   group.Version,
		CreatedAt: group.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: group.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return rs, nil
}

// UpdateConfigGroup 更新配置分组-后台使用
func (s *ConfigGroupService) UpdateConfigGroup(
	ctx *gin.Context,
	r req.UpdateConfigGroupReq,
) (err error) {
	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return biz_err.New(biz_err.PARAM_ERROR, "无效的分组ID")
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "分组不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询分组失败")
	}

	// 检查编码是否已存在（排除自己）
	if r.Code != "" && r.Code != group.Code {
		var count int64
		global.GVA_DB.Model(&systemRbac.ConfigGroup{}).Where("code = ? AND id != ?", r.Code, id).Count(&count)
		if count > 0 {
			return biz_err.New(biz_err.PARAM_ERROR, "分组编码已存在")
		}
		group.Code = r.Code
	}
	if r.Name != "" {
		group.Name = r.Name
	}
	if r.Config != "" {
		group.Config = datatypes.JSON(r.Config)
	}
	if r.Sort != 0 {
		group.Sort = r.Sort
	}
	if r.Status != 0 {
		group.Status = r.Status
	}
	if r.Remark != "" {
		group.Remark = r.Remark
	}

	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新分组失败")
	}
	return nil
}

// DeleteConfigGroup 删除配置分组-后台使用
func (s *ConfigGroupService) DeleteConfigGroup(
	ctx *gin.Context,
	r req.DeleteConfigGroupReq,
) (err error) {
	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return biz_err.New(biz_err.PARAM_ERROR, "无效的分组ID")
	}

	result := global.GVA_DB.Delete(&systemRbac.ConfigGroup{}, id)
	if result.Error != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除分组失败")
	}
	if result.RowsAffected == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "分组不存在")
	}
	return nil
}

// BatchDeleteConfigGroups 批量删除配置分组-后台使用
func (s *ConfigGroupService) BatchDeleteConfigGroups(
	ctx *gin.Context,
	r req.BatchDeleteConfigGroupsReq,
) (err error) {
	if r.Ids == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "ID列表不能为空")
	}

	// 解析ID数组
	idStrs := strings.Split(r.Ids, ",")
	var ids []uint64
	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "无效的ID列表")
	}

	result := global.GVA_DB.Where("id IN ?", ids).Delete(&systemRbac.ConfigGroup{})
	if result.Error != nil {
		return biz_err.New(biz_err.DB_ERROR, "批量删除分组失败")
	}
	return nil
}

// GetConfigByCode 根据编码获取配置-后台使用
func (s *ConfigGroupService) GetConfigByCode(
	ctx *gin.Context,
	r req.GetConfigByCodeReq,
) (rs res.GetConfigByCodeRes, err error) {
	if r.Code == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "分组编码不能为空")
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", r.Code).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "分组不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询配置失败")
	}

	rs = res.GetConfigByCodeRes{
		Id:      strconv.FormatUint(uint64(group.ID), 10),
		Code:    group.Code,
		Name:    group.Name,
		Config:  string(group.Config),
		Version: group.Version,
	}
	return rs, nil
}

// SaveConfig 保存配置-后台使用
func (s *ConfigGroupService) SaveConfig(
	ctx *gin.Context,
	r req.SaveConfigReq,
) (rs res.SaveConfigRes, err error) {
	if r.Code == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "分组编码不能为空")
	}
	if r.Config == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "配置内容不能为空")
	}

	// 验证JSON格式
	var configJson map[string]interface{}
	err = json.Unmarshal([]byte(r.Config), &configJson)
	if err != nil {
		return rs, biz_err.New(biz_err.JSON_PARSE, "配置JSON格式错误")
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", r.Code).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "分组不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询配置失败")
	}

	group.Config = datatypes.JSON(r.Config)
	group.Version++

	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存配置失败")
	}

	rs = res.SaveConfigRes{
		Version: group.Version,
	}
	return rs, nil
}

// RefreshConfigCache 刷新配置缓存-后台使用
func (s *ConfigGroupService) RefreshConfigCache(
	ctx *gin.Context,
) (rs res.RefreshConfigCacheRes, err error) {
	// TODO: 这里可以添加缓存刷新逻辑，例如清除Redis缓存
	// 假设有一个配置缓存的key前缀
	// global.GVA_REDIS.Del(ctx, "config:*")

	rs = res.RefreshConfigCacheRes{
		Success: true,
	}
	return rs, nil
}

