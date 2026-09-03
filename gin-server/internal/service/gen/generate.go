package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"shack/internal/global"
	"shack/internal/model/gen"
	req "shack/internal/model/gen/request"
	res "shack/internal/model/gen/response"
	"shack/internal/model/systemRbac"
	"shack/internal/utils"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	pdfoxide "github.com/yfedoseev/pdf_oxide/go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	templateCode = "gen_question" // 提示词模板编码
)

type GenerateService struct{}

// Generate 上传PDF并生成答案解析(异步)
func (s *GenerateService) Generate(
	ctx *gin.Context,
	r req.GenerateReq,
) (rs res.GenerateRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	// 1. 检查用户是否配置了API Key
	var apiKey gen.GenApiKey
	err = global.GVA_DB.Where("user_id = ?", userID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "请先配置API Key")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询API Key失败")
	}

	// 2. 获取文件信息
	var file systemRbac.SysFile
	err = global.GVA_DB.First(&file, r.FileId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "文件不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询文件失败")
	}

	// 3. 获取提示词模板
	var template systemRbac.Template
	err = global.GVA_DB.Where("code = ? AND status = ?", templateCode, "active").First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "提示词模板不存在，请联系管理员配置")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询提示词模板失败")
	}

	// 4. 创建历史记录(状态: pending)
	history := gen.GenHistory{
		UserID:   userID,
		FileName: r.FileName,
		FileID:   r.FileId,
		Status:   "pending",
	}
	if err := global.GVA_DB.Create(&history).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "创建历史记录失败")
	}

	// 5. 异步执行生成
	go doGenerate(history, file, apiKey, template)

	rs = res.GenerateRes{
		HistoryId: history.ID,
		Status:    "pending",
	}
	return rs, nil
}

// GetGenerateResult 获取生成进度/结果
func (s *GenerateService) GetGenerateResult(
	ctx *gin.Context,
	r req.GetGenerateResultReq,
) (rs res.GetGenerateResultRes, err error) {
	userID := getCurrentUserID(ctx)
	if userID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR)
	}

	var history gen.GenHistory
	err = global.GVA_DB.Where("id = ? AND user_id = ?", r.Id, userID).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "记录不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询记录失败")
	}

	rs = res.GetGenerateResultRes{
		Id:            history.ID,
		Status:        history.Status,
		ErrMsg:        history.ErrMsg,
		QuestionCount: history.QuestionCount,
	}

	// 解析Result JSON
	if len(history.Result) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(history.Result, &result); err == nil {
			rs.Result = result
		}
	}

	return rs, nil
}

// doGenerate 异步执行生成任务
func doGenerate(history gen.GenHistory, file systemRbac.SysFile, apiKey gen.GenApiKey, template systemRbac.Template) {
	db := global.GVA_DB

	// 更新状态为 generating
	db.Model(&gen.GenHistory{}).Where("id = ?", history.ID).Update("status", "generating")

	// 1. 读取PDF文本
	pdfPath := filepath.Join("statics", file.FilePath)
	text, err := extractPDFText(pdfPath)
	if err != nil {
		updateHistoryFailed(db, history.ID, "读取PDF文件失败: "+err.Error())
		return
	}

	if strings.TrimSpace(text) == "" {
		updateHistoryFailed(db, history.ID, "PDF文件内容为空或无法解析")
		return
	}

	// 2. 更新原始内容
	db.Model(&gen.GenHistory{}).Where("id = ?", history.ID).Update("raw_content", text)

	// 3. 构建提示词
	systemPrompt := template.Content
	userPrompt := fmt.Sprintf("以下是题目内容：\n\n%s\n\n请按照要求生成答案和解析。", text)

	// 4. 调用DeepSeek API
	client := openai.NewClient(
		option.WithAPIKey(apiKey.ApiKey),
		option.WithBaseURL(apiKey.BaseURL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: apiKey.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
	})

	if err != nil {
		updateHistoryFailed(db, history.ID, "调用AI接口失败: "+err.Error())
		return
	}

	if len(completion.Choices) == 0 {
		updateHistoryFailed(db, history.ID, "AI返回结果为空")
		return
	}

	// 5. 解析AI返回的JSON
	content := completion.Choices[0].Message.Content
	resultJSON := extractJSON(content)

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(resultJSON), &parsed); err != nil {
		// 如果解析失败，尝试清理后重新解析
		cleaned := cleanJSONString(resultJSON)
		if err2 := json.Unmarshal([]byte(cleaned), &parsed); err2 != nil {
			updateHistoryFailed(db, history.ID, "AI返回的JSON格式无效")
			global.GVA_LOG.Error("AI返回JSON解析失败",
				zap.String("content", content),
				zap.Error(err),
			)
			return
		}
	}

	// 6. 统计题目数量
	questionCount := countQuestions(parsed)

	// 7. 更新历史记录
	resultBytes, _ := json.Marshal(parsed)
	db.Model(&gen.GenHistory{}).Where("id = ?", history.ID).Updates(map[string]interface{}{
		"status":         "done",
		"result":         gen.JSON(resultBytes),
		"question_count": questionCount,
		"err_msg":        "",
	})

	global.GVA_LOG.Info("生成完成",
		zap.Uint("historyId", history.ID),
		zap.Int("questionCount", questionCount),
	)
}

// extractPDFText 提取PDF文本内容
func extractPDFText(filePath string) (string, error) {
	doc, err := pdfoxide.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()

	text, err := doc.ExtractAllText()
	if err != nil {
		return "", fmt.Errorf("读取PDF文本失败: %w", err)
	}
	return text, nil
}

// extractJSON 从AI回复中提取JSON部分
func extractJSON(content string) string {
	content = strings.TrimSpace(content)

	// 尝试从markdown代码块中提取
	if idx := strings.Index(content, "```json"); idx >= 0 {
		start := idx + 7
		end := strings.Index(content[start:], "```")
		if end >= 0 {
			return strings.TrimSpace(content[start : start+end])
		}
	}
	if idx := strings.Index(content, "```"); idx >= 0 {
		start := idx + 3
		// 跳过可能的语言标识行
		if newlineIdx := strings.Index(content[start:], "\n"); newlineIdx >= 0 {
			start = start + newlineIdx + 1
		}
		end := strings.Index(content[start:], "```")
		if end >= 0 {
			return strings.TrimSpace(content[start : start+end])
		}
	}

	// 尝试找JSON对象
	if idx := strings.Index(content, "{"); idx >= 0 {
		// 找最后一个}
		lastIdx := strings.LastIndex(content, "}")
		if lastIdx > idx {
			return content[idx : lastIdx+1]
		}
	}

	return content
}

// cleanJSONString 清理JSON字符串
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)
	// 移除BOM
	s = strings.TrimPrefix(s, "\xef\xbb\xbf")
	// 移除前后非JSON字符
	if idx := strings.Index(s, "{"); idx > 0 {
		s = s[idx:]
	}
	if idx := strings.LastIndex(s, "}"); idx >= 0 && idx < len(s)-1 {
		s = s[:idx+1]
	}
	return s
}

// countQuestions 统计题目数量
func countQuestions(data map[string]interface{}) int {
	count := 0
	sections, ok := data["sections"].([]interface{})
	if !ok {
		return 0
	}
	for _, section := range sections {
		sec, ok := section.(map[string]interface{})
		if !ok {
			continue
		}
		questions, ok := sec["questions"].([]interface{})
		if !ok {
			continue
		}
		count += len(questions)
	}
	return count
}

// updateHistoryFailed 更新历史记录为失败状态
func updateHistoryFailed(db *gorm.DB, historyID uint, errMsg string) {
	db.Model(&gen.GenHistory{}).Where("id = ?", historyID).Updates(map[string]interface{}{
		"status":  "failed",
		"err_msg": errMsg,
	})
	global.GVA_LOG.Error("生成失败", zap.Uint("historyId", historyID), zap.String("errMsg", errMsg))
}

// getCurrentUserID 从上下文获取当前用户ID
func getCurrentUserID(c *gin.Context) uint {
	return utils.GetUserID(c)
}
