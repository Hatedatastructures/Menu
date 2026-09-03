package ai

import (
	"context"
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"

	biz_err "shack/internal/error"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"gorm.io/gorm"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求参数
type ChatRequest struct {
	ConfigId uint     // AI配置ID, 为0时使用第一个启用的配置
	Messages []ChatMessage
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content  string `json:"content"`
	Model    string `json:"model"`
	Duration int64  `json:"duration"` // 耗时(毫秒)
}

// GetEnabledConfig 获取一个启用的AI配置
func GetEnabledConfig(configId uint) (systemRbac.AiConfig, error) {
	var config systemRbac.AiConfig
	db := global.GVA_DB

	if configId > 0 {
		err := db.Where("id = ? AND enabled = ?", configId, true).First(&config).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return config, biz_err.New(biz_err.PARAM_ERROR, "AI配置不存在或未启用")
			}
			return config, biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
		}
	} else {
		err := db.Where("enabled = ?", true).Order("id asc").First(&config).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return config, biz_err.New(biz_err.PARAM_ERROR, "没有可用的AI配置,请先在后台添加")
			}
			return config, biz_err.New(biz_err.DB_ERROR, "查询AI配置失败")
		}
	}

	return config, nil
}

// ChatCompletion 调用AI聊天补全
func ChatCompletion(req ChatRequest) (ChatResponse, error) {
	var rs ChatResponse

	// 获取AI配置
	config, err := GetEnabledConfig(req.ConfigId)
	if err != nil {
		return rs, err
	}

	// 构建消息
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			messages = append(messages, openai.SystemMessage(msg.Content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(msg.Content))
		default:
			messages = append(messages, openai.UserMessage(msg.Content))
		}
	}

	// 创建客户端
	client := openai.NewClient(
		option.WithAPIKey(config.ApiKey),
		option.WithBaseURL(config.BaseURL),
	)

	// 调用API
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    config.Model,
		Messages: messages,
	})
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		return rs, fmt.Errorf("调用AI接口失败: %w", err)
	}

	if len(completion.Choices) == 0 {
		return rs, fmt.Errorf("AI返回结果为空")
	}

	rs = ChatResponse{
		Content:  completion.Choices[0].Message.Content,
		Model:    config.Model,
		Duration: duration,
	}

	return rs, nil
}

// ChatCompletionWithSystem 带系统提示的快捷调用
func ChatCompletionWithSystem(systemPrompt, userPrompt string) (ChatResponse, error) {
	return ChatCompletion(ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	})
}

// SimpleChat 简单对话(单条用户消息)
func SimpleChat(prompt string) (ChatResponse, error) {
	return ChatCompletion(ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
	})
}
