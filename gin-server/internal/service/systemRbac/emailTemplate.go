package systemRbac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailTemplateService struct{}

// 创建邮件模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) CreateEmailTemplate(
	c *gin.Context,
	r req.CreateEmailTemplateReq,
) (rs res.CreateEmailTemplateRes, err error) {
	// 1. 检查编码是否已存在
	var count int64
	err = global.GVA_DB.Model(&systemRbac.EmailTemplate{}).Where("code = ?", r.Code).Count(&count).Error
	if err != nil {
		global.GVA_LOG.Error("检查模板编码失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "模板编码已存在")
	}

	// 2. 序列化变量
	variablesJSON, err := json.Marshal(r.Variables)
	if err != nil {
		global.GVA_LOG.Error("序列化变量失败", zap.Error(err))
		return rs, biz_err.New(biz_err.JSON_PARSE)
	}

	// 3. 创建模板
	template := systemRbac.EmailTemplate{
		Code:        r.Code,
		Name:        r.Name,
		Subject:     r.Subject,
		Content:     r.Content,
		Type:        r.Type,
		Description: r.Description,
		Variables:   string(variablesJSON),
		Status:      "active",
	}

	err = global.GVA_DB.Create(&template).Error
	if err != nil {
		global.GVA_LOG.Error("创建模板失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	rs = res.CreateEmailTemplateRes{
		Id:        int64(template.ID),
		Code:      template.Code,
		CreatedAt: template.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// 更新邮件模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) UpdateEmailTemplate(
	c *gin.Context,
	r req.UpdateEmailTemplateReq,
) (err error) {
	// 1. 检查模板是否存在
	var template systemRbac.EmailTemplate
	err = global.GVA_DB.First(&template, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("模板不存在", zap.Error(err))
		return biz_err.New(biz_err.EMAIL_NOT_FOUND)
	}

	// 2. 检查编码是否被其他模板使用
	var count int64
	err = global.GVA_DB.Model(&systemRbac.EmailTemplate{}).
		Where("code = ? AND id != ?", r.Code, r.Id).
		Count(&count).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}
	if count > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "模板编码已被使用")
	}

	// 3. 序列化变量
	variablesJSON, err := json.Marshal(r.Variables)
	if err != nil {
		return biz_err.New(biz_err.JSON_PARSE)
	}

	// 4. 更新模板
	err = global.GVA_DB.Model(&template).Updates(map[string]interface{}{
		"code":        r.Code,
		"name":        r.Name,
		"subject":     r.Subject,
		"content":     r.Content,
		"type":        r.Type,
		"description": r.Description,
		"variables":   string(variablesJSON),
	}).Error

	if err != nil {
		global.GVA_LOG.Error("更新模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 删除邮件模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) DeleteEmailTemplate(
	c *gin.Context,
	r req.DeleteEmailTemplateReq,
) (err error) {
	err = global.GVA_DB.Delete(&systemRbac.EmailTemplate{}, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("删除模板失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}
	return nil
}

// 获取邮件模板列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) GetEmailTemplateList(
	c *gin.Context,
	r req.GetEmailTemplateListReq,
) (rs res.GetEmailTemplateListRes, err error) {
	// 分页参数
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.EmailTemplate{})

	// 条件过滤
	if r.Type != "" {
		db = db.Where("type = ?", r.Type)
	}
	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}
	if r.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+r.Keyword+"%", "%"+r.Keyword+"%")
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	// 分页查询
	var templates []systemRbac.EmailTemplate
	offset := (page - 1) * pageSize
	err = db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&templates).Error
	if err != nil {
		global.GVA_LOG.Error("查询模板列表失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	// 构建响应
	list := make([]res.GetEmailTemplateListResList, 0, len(templates))
	for _, tmpl := range templates {
		list = append(list, res.GetEmailTemplateListResList{
			Id:        int64(tmpl.ID),
			Code:      tmpl.Code,
			Name:      tmpl.Name,
			Type:      tmpl.Type,
			Status:    tmpl.Status,
			CreatedAt: tmpl.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: tmpl.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetEmailTemplateListRes{
		Page:     page,
		PageSize: pageSize,
		Type:     r.Type,
		Keyword:  r.Keyword,
		Status:   r.Status,
		Total:    total,
		List:     list,
	}

	return rs, nil
}

// 获取邮件模板详情-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) GetEmailTemplate(
	c *gin.Context,
	r req.GetEmailTemplateReq,
) (rs res.GetEmailTemplateRes, err error) {
	var template systemRbac.EmailTemplate
	err = global.GVA_DB.First(&template, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("查询模板失败", zap.Error(err))
		return rs, biz_err.New(biz_err.EMAIL_NOT_FOUND)
	}

	// 反序列化变量
	var variables []res.GetEmailTemplateResVariable
	if template.Variables != "" {
		err = json.Unmarshal([]byte(template.Variables), &variables)
		if err != nil {
			global.GVA_LOG.Error("反序列化变量失败", zap.Error(err))
		}
	}

	rs = res.GetEmailTemplateRes{
		Id:          int64(template.ID),
		Code:        template.Code,
		Name:        template.Name,
		Subject:     template.Subject,
		Content:     template.Content,
		Type:        template.Type,
		Status:      template.Status,
		Description: template.Description,
		Variables:   variables,
		CreatedAt:   template.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return rs, nil
}

// 预览邮件模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) PreviewEmailTemplate(
	c *gin.Context,
	r req.PreviewEmailTemplateReq,
) (err error) {
	// 获取模板
	var template systemRbac.EmailTemplate
	err = global.GVA_DB.Where("code = ?", r.TemplateId).First(&template).Error
	if err != nil {
		return biz_err.New(biz_err.EMAIL_NOT_FOUND)
	}

	// TODO: 渲染模板并返回
	// 这里应该调用模板渲染引擎

	return nil
}

// AI生成邮件模板-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailTemplateService) AiGenerateEmailTemplate(
	c *gin.Context,
	r req.AiGenerateEmailTemplateReq,
) (rs res.AiGenerateEmailTemplateRes, err error) {
	// 1. 获取启用的AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("enabled = ?", true).First(&aiConfig).Error
	if err != nil {
		global.GVA_LOG.Error("未找到启用的AI配置", zap.Error(err))
		return rs, biz_err.New(biz_err.PARAM_ERROR, "未找到启用的AI配置，请先在AI配置中添加并启用配置")
	}

	// 2. 构建AI提示词
	prompt := s.buildEmailTemplatePrompt(r.Prompt, r.Style, r.Language)

	// 3. 调用AI服务生成模板
	aiResponse, err := s.callAIService(aiConfig, prompt)
	if err != nil {
		global.GVA_LOG.Error("调用AI服务失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "调用AI服务失败: "+err.Error())
	}

	// 4. 解析AI返回的JSON结果
	err = s.parseAIResponse(aiResponse, &rs)
	if err != nil {
		global.GVA_LOG.Error("解析AI响应失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "解析AI响应失败: "+err.Error())
	}

	return rs, nil
}

// buildEmailTemplatePrompt 构建邮件模板生成的提示词
func (s *EmailTemplateService) buildEmailTemplatePrompt(userPrompt, style, language string) string {
	styleMap := map[string]string{
		"modern":  "现代简约风格，使用清爽的颜色和简洁的布局",
		"classic": "商务正式风格，使用稳重的颜色和专业的布局",
		"minimal": "极简主义风格，使用黑白灰色调和极简布局",
	}

	styleDesc := styleMap[style]
	if styleDesc == "" {
		styleDesc = "现代简约风格"
	}

	languageDesc := "中文"
	if language == "en-US" {
		languageDesc = "英文"
	}

	prompt := fmt.Sprintf(`你是一个专业的邮件模板设计师。请根据以下要求生成一个邮件模板：

用户需求：%s

设计风格：%s
语言：%s

请按照以下JSON格式返回结果（必须是纯JSON，不要有其他文字）：
{
  "subject": "邮件主题",
  "html": "<html>完整的HTML邮件内容，包含内联CSS样式，适配各种邮件客户端</html>",
  "text": "纯文本版本的邮件内容",
  "variables": ["var1", "var2"]  // 识别出的变量名列表
}

要求：
1. HTML必须是完整的、可直接使用的邮件模板
2. 使用内联CSS样式，确保在各大邮件客户端中正常显示
3. 主题要简洁明了
4. 纯文本版本要去除HTML标签，保留核心内容
5. 自动识别模板中需要替换的变量（如用户名、验证码、链接等）
6. 设计要美观大方，符合%s的特点
7. 使用%s语言
`, userPrompt, styleDesc, languageDesc, styleDesc, languageDesc)

	return prompt
}

// callAIService 调用AI服务
func (s *EmailTemplateService) callAIService(config systemRbac.AiConfig, prompt string) (string, error) {
	// 构建请求URL
	url := config.BaseURL + "/v1/chat/completions"

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个专业的邮件模板设计师，擅长设计美观实用的HTML邮件模板。"},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构建请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI接口返回错误: %s", string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取返回内容
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("响应格式错误: 缺少choices字段")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: choices[0]格式错误")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: message格式错误")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("响应格式错误: content格式错误")
	}

	return content, nil
}

// parseAIResponse 解析AI返回的JSON响应
func (s *EmailTemplateService) parseAIResponse(aiResponse string, rs *res.AiGenerateEmailTemplateRes) error {
	// 尝试提取JSON（AI可能在JSON前后添加了说明文字）
	jsonStart := strings.Index(aiResponse, "{")
	jsonEnd := strings.LastIndex(aiResponse, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return fmt.Errorf("未找到有效的JSON内容")
	}

	jsonStr := aiResponse[jsonStart : jsonEnd+1]

	// 解析JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	// 提取字段
	if subject, ok := result["subject"].(string); ok {
		rs.Subject = subject
	}
	if html, ok := result["html"].(string); ok {
		rs.Html = html
	}
	if text, ok := result["text"].(string); ok {
		rs.Text = text
	}
	if variables, ok := result["variables"].([]interface{}); ok {
		rs.Variables = make([]string, 0, len(variables))
		for _, v := range variables {
			if varStr, ok := v.(string); ok {
				rs.Variables = append(rs.Variables, varStr)
			}
		}
	}

	return nil
}

