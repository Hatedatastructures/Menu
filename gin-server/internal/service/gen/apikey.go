package gen

import (
	"context"
	"strings"
	"time"

	"shack/internal/global"
	"shack/internal/model/gen"
	req "shack/internal/model/gen/request"
	res "shack/internal/model/gen/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"gorm.io/gorm"
)

type ApikeyService struct{}

// GetMyApiKey 获取我的API Key配置
func (s *ApikeyService) GetMyApiKey(
	ctx *gin.Context,
) (rs res.GetMyApiKeyRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	var apiKey gen.GenApiKey
	err = global.GVA_DB.Where("user_id = ?", userID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "未配置API Key")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询API Key失败")
	}

	rs = res.GetMyApiKeyRes{
		Id:        apiKey.ID,
		ApiKey:    maskApiKey(apiKey.ApiKey),
		BaseUrl:   apiKey.BaseURL,
		Model:     apiKey.Model,
		Enabled:   true,
		CreatedAt: apiKey.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return rs, nil
}

// SaveApiKey 保存/更新API Key
func (s *ApikeyService) SaveApiKey(
	ctx *gin.Context,
	r req.SaveApiKeyReq,
) (rs res.SaveApiKeyRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	if r.ApiKey == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "API Key不能为空")
	}

	// 设置默认值
	baseURL := r.BaseUrl
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	model := r.Model
	if model == "" {
		model = "deepseek-chat"
	}

	var existing gen.GenApiKey
	err = global.GVA_DB.Where("user_id = ?", userID).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询API Key失败")
	}

	if existing.ID > 0 {
		// 更新
		existing.ApiKey = r.ApiKey
		existing.BaseURL = baseURL
		existing.Model = model
		if err := global.GVA_DB.Save(&existing).Error; err != nil {
			return rs, biz_err.New(biz_err.DB_ERROR, "更新API Key失败")
		}
		rs.Id = existing.ID
	} else {
		// 新建
		newKey := gen.GenApiKey{
			UserID:  userID,
			ApiKey:  r.ApiKey,
			BaseURL: baseURL,
			Model:   model,
		}
		if err := global.GVA_DB.Create(&newKey).Error; err != nil {
			return rs, biz_err.New(biz_err.DB_ERROR, "保存API Key失败")
		}
		rs.Id = newKey.ID
	}

	return rs, nil
}

// DeleteApiKey 删除API Key
func (s *ApikeyService) DeleteApiKey(
	ctx *gin.Context,
) (err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR)
	}

	result := global.GVA_DB.Where("user_id = ?", userID).Delete(&gen.GenApiKey{})
	if result.Error != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除API Key失败")
	}
	if result.RowsAffected == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "未配置API Key")
	}
	return nil
}

// TestApiKey 测试API Key是否可用
func (s *ApikeyService) TestApiKey(
	ctx *gin.Context,
	r req.TestApiKeyReq,
) (rs res.TestApiKeyRes, err error) {
	if r.ApiKey == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "API Key不能为空")
	}

	baseURL := r.BaseUrl
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	model := r.Model
	if model == "" {
		model = "deepseek-chat"
	}

	client := openai.NewClient(
		option.WithAPIKey(r.ApiKey),
		option.WithBaseURL(baseURL),
	)

	testCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Chat.Completions.New(testCtx, openai.ChatCompletionNewParams{
		Model: model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("hi"),
		},
		MaxTokens: openai.Int(5),
	})

	if err != nil {
		rs.Valid = false
		errMsg := err.Error()
		if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "invalid") {
			rs.Message = "API Key无效"
		} else if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline") {
			rs.Message = "连接超时，请检查API地址"
		} else {
			rs.Message = "验证失败: " + errMsg
		}
		return rs, nil
	}

	rs.Valid = true
	rs.Message = "API Key有效"
	return rs, nil
}

// maskApiKey 脱敏API Key
func maskApiKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
