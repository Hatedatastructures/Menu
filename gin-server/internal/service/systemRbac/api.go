package systemRbac

import (
	"errors"
	"strconv"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	"github.com/gin-gonic/gin"
	biz_err "shack/internal/error"
	"gorm.io/gorm"
)

type ApiService struct{}

var ApiServiceApp = new(ApiService)

// 有效的HTTP方法
var validMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

// 创建API-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月27日 23:36:39
func (s *ApiService) CreateApi(
	ctx *gin.Context,
	r req.CreateApiReq,
) (rs res.CreateApiRes, err error) {
	// 参数校验
	if r.Name == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "接口名称不能为空")
	}
	if r.Path == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请求路径不能为空")
	}
	if r.Method == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "请求方法不能为空")
	}

	// 校验请求方法
	r.Method = validateMethod(r.Method)
	if r.Method == "" {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "请求方法必须是 GET/POST/PUT/PATCH/DELETE")
	}

	// 检查API是否已存在
	var count int64
	global.GVA_DB.Model(&systemRbac.Api{}).Where("path = ? AND method = ?", r.Path, r.Method).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "该API已存在")
	}

	// 创建API
	api := systemRbac.Api{
		Name:   r.Name,
		Path:   r.Path,
		Method: r.Method,
		Status: r.Status,
	}
	err = global.GVA_DB.Create(&api).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建API失败")
	}

	rs = res.CreateApiRes{
		Id: strconv.FormatUint(uint64(api.ID), 10),
	}
	return rs, nil
}

// 校验并规范化HTTP方法
func validateMethod(method string) string {
	m := method
	// 统一转换为大写
	for i := range m {
		if m[i] >= 'a' && m[i] <= 'z' {
			m = m[:i] + string(m[i]-'a'+'A') + m[i+1:]
			break
		}
	}
	if validMethods[m] {
		return m
	}
	return ""
}

// API列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月27日 23:36:39
func (s *ApiService) GetApiList(
	ctx *gin.Context,
	r req.GetApiListReq,
) (rs res.GetApiListRes, err error) {
	// 设置默认值
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.Api{})
	if r.Keyword != "" {
		db = db.Where("name LIKE ? OR path LIKE ?", "%"+r.Keyword+"%", "%"+r.Keyword+"%")
	}
	if r.Method != "" {
		r.Method = validateMethod(r.Method)
		if r.Method != "" {
			db = db.Where("method = ?", r.Method)
		}
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 分页查询
	var apis []systemRbac.Api
	offset := (r.Page - 1) * r.Size
	err = db.Offset(offset).Limit(r.Size).Order("id desc").Find(&apis).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询API列表失败")
	}

	// 构建返回数据
	rs.List = make([]res.GetApiListResList, 0, len(apis))
	for _, api := range apis {
		rs.List = append(rs.List, res.GetApiListResList{
			Id:        strconv.FormatUint(uint64(api.ID), 10),
			Name:      api.Name,
			Path:      api.Path,
			Method:    api.Method,
			Status:    api.Status,
			CreatedAt: api.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	rs.Keyword = r.Keyword
	rs.Page = r.Page
	rs.Size = r.Size
	rs.Method = r.Method
	rs.Total = total
	return rs, nil
}

// API详情-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月27日 23:36:39
func (s *ApiService) GetApiDetail(
	ctx *gin.Context,
	r req.GetApiDetailReq,
) (rs res.GetApiDetailRes, err error) {
	// 参数校验
	if r.Id == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "API ID不能为空")
	}

	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_FORMAT, "API ID格式错误")
	}

	// 查询API
	var api systemRbac.Api
	err = global.GVA_DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "API不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询API失败")
	}

	// 构建返回数据
	rs = res.GetApiDetailRes{
		Id:        strconv.FormatUint(uint64(api.ID), 10),
		Name:      api.Name,
		Path:      api.Path,
		Method:    api.Method,
		Status:    api.Status,
		CreatedAt: api.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return rs, nil
}

// 更新API-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月27日 23:36:39
func (s *ApiService) UpdateApi(
	ctx *gin.Context,
	r req.UpdateApiReq,
) (err error) {
	// 参数校验
	if r.Id == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "API ID不能为空")
	}

	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return biz_err.New(biz_err.PARAM_FORMAT, "API ID格式错误")
	}

	// 查询API
	var api systemRbac.Api
	err = global.GVA_DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "API不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询API失败")
	}

	// 校验请求方法
	if r.Method != "" {
		r.Method = validateMethod(r.Method)
		if r.Method == "" {
			return biz_err.New(biz_err.PARAM_ERROR, "请求方法必须是 GET/POST/PUT/PATCH/DELETE")
		}
	}

	// 检查API是否被其他API使用
	if r.Path != "" && r.Method != "" {
		if r.Path != api.Path || r.Method != api.Method {
			var count int64
			global.GVA_DB.Model(&systemRbac.Api{}).Where("path = ? AND method = ? AND id != ?", r.Path, r.Method, id).Count(&count)
			if count > 0 {
				return biz_err.New(biz_err.PARAM_ERROR, "该API已存在")
			}
		}
	}

	// 更新API信息
	if r.Name != "" {
		api.Name = r.Name
	}
	if r.Path != "" {
		api.Path = r.Path
	}
	if r.Method != "" {
		api.Method = r.Method
	}
	api.Status = r.Status

	err = global.GVA_DB.Save(&api).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新API失败")
	}
	return nil
}

// 删除API-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月27日 23:36:39
func (s *ApiService) DeleteApi(
	ctx *gin.Context,
	r req.DeleteApiReq,
) (err error) {
	// 参数校验
	if r.Id == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "API ID不能为空")
	}

	id, err := strconv.ParseUint(r.Id, 10, 32)
	if err != nil {
		return biz_err.New(biz_err.PARAM_FORMAT, "API ID格式错误")
	}

	// 检查API是否存在
	var api systemRbac.Api
	err = global.GVA_DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "API不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询API失败")
	}

	// 软删除
	err = global.GVA_DB.Delete(&api).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除API失败")
	}
	return nil
}

// 获取所有API列表-内部使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月01日
func (s *ApiService) GetAllApis(authorityID uint) (apis []systemRbac.Api, err error) {
	err = global.GVA_DB.Order("id desc").Find(&apis).Error
	return apis, err
}

