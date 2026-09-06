package systemRbac

import (
	"regexp"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TemplateService struct{}

// 创建模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年05月22日 10:53:03
func (s *TemplateService) CreateTemplate(
	ctx *gin.Context,
	r req.CreateTemplateReq,
) (rs res.CreateTemplateRes, err error) {
	// 参数校验
	if r.Code == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "模板编码不能为空")
	}
	if r.Name == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "模板名称不能为空")
	}
	if r.Tpl == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "模板内容不能为空")
	}

	// 检查编码是否已存在
	var count int64
	err = global.GVA_DB.Model(&systemRbac.Template{}).Where("code = ?", r.Code).Count(&count).Error
	if err != nil {
		global.GVA_LOG.Error("检查模板编码失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "模板编码已存在")
	}

	// 自动解析模板变量
	variables := parseTemplateVariables(r.Tpl)

	// 状态处理, 默认 active
	status := "active"
	if r.Status != "" {
		status = r.Status
	}

	// 引擎处理, 默认 html/template
	engine := "html/template"
	if r.Engine != "" {
		engine = r.Engine
	}

	// 创建模板
	template := systemRbac.Template{
		Code:        r.Code,
		Name:        r.Name,
		Content:     r.Tpl,
		Status:      status,
		Description: r.Description,
		Variables:   datatypes.NewJSONType(variables),
		Engine:      engine,
	}

	err = global.GVA_DB.Create(&template).Error
	if err != nil {
		global.GVA_LOG.Error("创建模板失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	rs = res.CreateTemplateRes{
		ID: template.ID,
	}

	return rs, nil
}

// 删除模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年05月22日 10:53:03
func (s *TemplateService) DeleteTemplate(
	ctx *gin.Context,
	r req.DeleteTemplateReq,
) (err error) {
	// 检查模板是否存在
	var template systemRbac.Template
	err = global.GVA_DB.First(&template, r.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "模板不存在")
		}
		global.GVA_LOG.Error("查询模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	err = global.GVA_DB.Delete(&template).Error
	if err != nil {
		global.GVA_LOG.Error("删除模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 更新模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年05月22日 10:53:03
func (s *TemplateService) UpdateTemplate(
	ctx *gin.Context,
	r req.UpdateTemplateReq,
) (err error) {
	// 检查模板是否存在
	var template systemRbac.Template
	err = global.GVA_DB.First(&template, r.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "模板不存在")
		}
		global.GVA_LOG.Error("查询模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	// 参数校验
	if r.Code == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "模板编码不能为空")
	}
	if r.Name == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "模板名称不能为空")
	}
	if r.Tpl == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "模板内容不能为空")
	}

	// 检查编码是否被其他模板使用
	var count int64
	err = global.GVA_DB.Model(&systemRbac.Template{}).
		Where("code = ? AND id != ?", r.Code, r.Id).
		Count(&count).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}
	if count > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "模板编码已被使用")
	}

	// 自动解析模板变量
	variables := parseTemplateVariables(r.Tpl)

	// 更新模板
	err = global.GVA_DB.Model(&template).Updates(map[string]interface{}{
		"code":        r.Code,
		"name":        r.Name,
		"content":     r.Tpl,
		"status":      r.Status,
		"description": r.Description,
		"engine":      r.Engine,
		"variables":   datatypes.NewJSONType(variables),
	}).Error

	if err != nil {
		global.GVA_LOG.Error("更新模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 获取模板详情-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年05月22日 10:53:03
func (s *TemplateService) GetTemplateDetail(
	ctx *gin.Context,
	r req.GetTemplateDetailReq,
) (rs res.GetTemplateDetailRes, err error) {
	var template systemRbac.Template
	err = global.GVA_DB.First(&template, r.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "模板不存在")
		}
		global.GVA_LOG.Error("查询模板失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	rs = res.GetTemplateDetailRes{
		ID:          template.ID,
		Code:        template.Code,
		Name:        template.Name,
		Content:     template.Content,
		Status:      template.Status,
		Description: template.Description,
		Variables:   template.Variables.Data(),
		Engine:      template.Engine,
		CreatedAt:   template.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// 获取模板分页列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年05月22日 10:53:03
func (s *TemplateService) GetTemplateList(
	ctx *gin.Context,
	r req.GetTemplateListReq,
) (rs res.GetTemplateListRes, err error) {
	// 分页参数
	page := r.Page
	if page < 1 {
		page = 1
	}
	size := r.Size
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.Template{})

	// 条件过滤
	if r.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+r.Keyword+"%", "%"+r.Keyword+"%")
	}
	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	// 分页查询
	var templates []systemRbac.Template
	offset := (page - 1) * size
	err = db.Order("created_at DESC").Limit(size).Offset(offset).Find(&templates).Error
	if err != nil {
		global.GVA_LOG.Error("查询模板列表失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	// 构建响应
	list := make([]res.GetTemplateListResList, 0, len(templates))
	for _, tmpl := range templates {
		list = append(list, res.GetTemplateListResList{
			ID:          tmpl.ID,
			Code:        tmpl.Code,
			Name:        tmpl.Name,
			Content:     tmpl.Content,
			Status:      tmpl.Status,
			Description: tmpl.Description,
			Engine:      tmpl.Engine,
			Variables:   tmpl.Variables.Data(),
			CreatedAt:   tmpl.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tmpl.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetTemplateListRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}

	return rs, nil
}

// ============ 辅助函数 ============

// parseTemplateVariables 从模板内容中解析变量
// 支持 {{.VariableName}} 格式
func parseTemplateVariables(tpl string) []string {
	re := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)
	matches := re.FindAllStringSubmatch(tpl, -1)

	// 去重
	seen := make(map[string]bool)
	var variables []string
	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			if !seen[name] {
				seen[name] = true
				variables = append(variables, name)
			}
		}
	}

	return variables
}

