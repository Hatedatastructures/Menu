package systemRbac

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"
	"gorm.io/gorm"
)

type AiConfigService struct{}

var AiConfigServiceApp = new(AiConfigService)

// CreateAiConfig 创建AI配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) CreateAiConfig(
	ctx context.Context,
	r req.CreateAiConfigReq,
) (rs res.CreateAiConfigRes, err error) {
	// 参数校验
	if r.Name == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "配置名称不能为空")
	}
	if r.BaseUrl == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "接口地址不能为空")
	}
	if r.ApiKey == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "密钥不能为空")
	}
	if r.Model == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "模型名称不能为空")
	}
	if r.Type == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "类型不能为空")
	}

	// 检查配置名称是否已存在
	var count int64
	global.GVA_DB.Model(&systemRbac.AiConfig{}).Where("name = ?", r.Name).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "配置名称已存在")
	}

	// 创建AI配置
	aiConfig := systemRbac.AiConfig{
		Name:     r.Name,
		Provider: r.Provider,
		BaseURL:  r.BaseUrl,
		ApiKey:   r.ApiKey,
		Model:    r.Model,
		Type:     r.Type,
		Enabled:  r.Enabled,
	}
	err = global.GVA_DB.Create(&aiConfig).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建AI配置失败")
	}

	rs = res.CreateAiConfigRes{
		Id: aiConfig.ID,
	}
	return rs, nil
}

// GetAiConfigList 获取AI配置列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) GetAiConfigList(
	ctx context.Context,
	r req.GetAiConfigListReq,
) (rs res.GetAiConfigListRes, err error) {
	// 设置默认值
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.AiConfig{})
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}

	// 查询总数
	var total int64
	db.Count(&total)

	// 分页查询
	var aiConfigs []systemRbac.AiConfig
	offset := (r.Page - 1) * r.Size
	err = db.Offset(offset).Limit(r.Size).Order("id desc").Find(&aiConfigs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询AI配置列表失败")
	}

	// 构建返回数据（脱敏ApiKey）
	rs.List = make([]res.GetAiConfigListResList, 0, len(aiConfigs))
	for _, config := range aiConfigs {
		maskedApiKey := maskApiKey(config.ApiKey)
		rs.List = append(rs.List, res.GetAiConfigListResList{
			Id:        config.ID,
			Name:      config.Name,
			Provider:  config.Provider,
			BaseUrl:   config.BaseURL,
			ApiKey:    maskedApiKey,
			Model:     config.Model,
			Type:      config.Type,
			Enabled:   config.Enabled,
			CreatedAt: config.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: config.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs.Page = r.Page
	rs.Size = r.Size
	rs.Name = r.Name
	rs.Enabled = r.Enabled
	rs.Total = total
	return rs, nil
}

// GetAiConfig 获取AI配置详情-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) GetAiConfig(
	ctx context.Context,
	r req.GetAiConfigReq,
) (rs res.GetAiConfigRes, err error) {
	// 参数校验
	if r.Id == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}

	// 查询AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.Id).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 构建返回数据
	rs = res.GetAiConfigRes{
		Id:        aiConfig.ID,
		Name:      aiConfig.Name,
		Provider:  aiConfig.Provider,
		BaseUrl:   aiConfig.BaseURL,
		ApiKey:    aiConfig.ApiKey,
		Model:     aiConfig.Model,
		Type:      aiConfig.Type,
		Enabled:   aiConfig.Enabled,
		CreatedAt: aiConfig.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: aiConfig.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return rs, nil
}

// UpdateAiConfig 更新AI配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) UpdateAiConfig(
	ctx context.Context,
	r req.UpdateAiConfigReq,
) (err error) {
	// 参数校验
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}

	// 查询AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.Id).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 检查配置名称是否被其他配置使用
	if r.Name != "" && r.Name != aiConfig.Name {
		var count int64
		global.GVA_DB.Model(&systemRbac.AiConfig{}).Where("name = ? AND id != ?", r.Name, r.Id).Count(&count)
		if count > 0 {
			return biz_err.New(biz_err.PARAM_ERROR, "配置名称已存在")
		}
	}

	// 更新AI配置信息
	aiConfig.Name = r.Name
	aiConfig.Provider = r.Provider
	aiConfig.BaseURL = r.BaseUrl
	aiConfig.ApiKey = r.ApiKey
	aiConfig.Model = r.Model
	aiConfig.Type = r.Type
	aiConfig.Enabled = r.Enabled

	err = global.GVA_DB.Save(&aiConfig).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新AI配置失败")
	}
	return nil
}

// DeleteAiConfig 删除AI配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) DeleteAiConfig(
	ctx context.Context,
	r req.DeleteAiConfigReq,
) (err error) {
	// 参数校验
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}

	// 检查AI配置是否存在
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.Id).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 软删除
	err = global.GVA_DB.Delete(&aiConfig).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除AI配置失败")
	}
	return nil
}

// BatchDeleteAiConfig 批量删除AI配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) BatchDeleteAiConfig(
	ctx context.Context,
	r req.BatchDeleteAiConfigReq,
) (err error) {
	// 参数校验
	if len(r.Ids) == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "配置ID数组不能为空")
	}

	// 批量软删除
	err = global.GVA_DB.Delete(&[]systemRbac.AiConfig{}, "id in (?)", r.Ids).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "批量删除AI配置失败")
	}
	return nil
}

// TestAiConfig 测试AI接口-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) TestAiConfig(
	ctx context.Context,
	r req.TestAiConfigReq,
) (rs res.TestAiConfigRes, err error) {
	// 参数校验
	if r.Id == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}
	if r.Prompt == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "测试提示词不能为空")
	}

	// 查询AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.Id).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 调用AI接口
	startTime := time.Now()
	response, callErr := callAI(aiConfig, r.Prompt)
	duration := time.Since(startTime).Milliseconds()

	if callErr != nil {
		rs = res.TestAiConfigRes{
			Success:  false,
			Error:    callErr.Error(),
			Duration: duration,
		}
		return rs, nil
	}

	rs = res.TestAiConfigRes{
		Success:  true,
		Response: response,
		Duration: duration,
	}
	return rs, nil
}

// AiChat AI对话调用-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) AiChat(
	ctx context.Context,
	r req.AiChatReq,
) (rs res.AiChatRes, err error) {
	// 参数校验
	if r.ConfigId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}
	if r.Prompt == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "用户输入不能为空")
	}

	// 查询AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.ConfigId).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 检查配置是否启用
	if !aiConfig.Enabled {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "该AI配置未启用")
	}

	// 调用AI接口
	startTime := time.Now()
	response, callErr := callAI(aiConfig, r.Prompt)
	duration := time.Since(startTime).Milliseconds()

	if callErr != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "调用AI接口失败: "+callErr.Error())
	}

	rs = res.AiChatRes{
		Response: response,
		Model:    aiConfig.Model,
		Duration: duration,
	}
	return rs, nil
}

// ToggleAiConfigStatus 切换AI配置启用状态-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) ToggleAiConfigStatus(
	ctx context.Context,
	r req.ToggleAiConfigStatusReq,
) (err error) {
	// 参数校验
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "配置ID不能为空")
	}

	// 查询AI配置
	var aiConfig systemRbac.AiConfig
	err = global.GVA_DB.Where("id = ?", r.Id).First(&aiConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
	}

	// 更新启用状态
	aiConfig.Enabled = r.Enabled
	err = global.GVA_DB.Save(&aiConfig).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新AI配置状态失败")
	}
	return nil
}

// GetEnabledAiConfigs 获取启用的AI配置列表-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 22:50:55
func (s *AiConfigService) GetEnabledAiConfigs(
	ctx context.Context,
) (rs res.GetEnabledAiConfigsRes, err error) {
	// 查询所有启用的AI配置
	var aiConfigs []systemRbac.AiConfig
	err = global.GVA_DB.Where("enabled = ?", true).Order("id desc").Find(&aiConfigs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询AI配置列表失败")
	}

	// 构建返回数据
	rs.List = make([]res.GetEnabledAiConfigsResList, 0, len(aiConfigs))
	for _, config := range aiConfigs {
		rs.List = append(rs.List, res.GetEnabledAiConfigsResList{
			Id:       config.ID,
			Name:     config.Name,
			Provider: config.Provider,
			Model:    config.Model,
			Type:     config.Type,
		})
	}
	return rs, nil
}

// ========== 私有方法 ==========

// callAI 调用AI接口（OpenAI兼容格式）
func callAI(config systemRbac.AiConfig, prompt string) (string, error) {
	// 构建请求URL
	url := config.BaseURL + "/v1/chat/completions"

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": config.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
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

	// 设置请求头（使用完整的ApiKey）
	apiKey := config.ApiKey
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 调试日志：打印API Key长度（不打印完整key保护安全）
	fmt.Printf("DEBUG: API Key length = %d, first 8 chars = %s\n", len(apiKey), apiKey[:8])

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

// maskApiKey 脱敏ApiKey
func maskApiKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "********"
	}
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}
