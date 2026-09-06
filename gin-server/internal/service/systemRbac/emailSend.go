package systemRbac

import (
	"fmt"
	"html/template"
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

type EmailSendService struct{}

// 发送普通邮件-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailSendService) SendEmail(
	c *gin.Context,
	r req.SendEmailReq,
) (rs res.SendEmailRes, err error) {
	// 1. 验证邮件地址
	if len(r.To) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "收件人邮箱不能为空")
	}
	if r.Subject == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "邮件主题不能为空")
	}
	if r.Content == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "邮件内容不能为空")
	}

	// 2. 获取邮件配置
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "邮件配置未设置")
	}

	// 3. 确定邮件类型
	isHTML := r.Type == "html" || r.Type == ""

	// 4. 创建邮件日志
	taskId := fmt.Sprintf("email_%d", time.Now().Unix())
	emailLog := systemRbac.EmailLog{
		ToEmail: strings.Join(r.To, ","),
		Subject: r.Subject,
		Content: r.Content,
		Status:  "pending",
	}
	global.GVA_DB.Create(&emailLog)

	// 5. 发送邮件（同步）
	mailer := NewMailer(*config)
	err = mailer.SendEmail(r.To, r.Cc, r.Bcc, r.Subject, r.Content, isHTML)

	// 6. 更新日志状态
	if err != nil {
		global.GVA_DB.Model(&emailLog).Updates(map[string]interface{}{
			"status":  "failed",
			"error_msg": err.Error(),
		})
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "邮件发送失败: "+err.Error())
	}

	global.GVA_DB.Model(&emailLog).Update("status", "sent")

	rs = res.SendEmailRes{
		TaskId: taskId,
	}

	global.GVA_LOG.Info("发送普通邮件成功", zap.Strings("to", r.To), zap.String("subject", r.Subject))

	return rs, nil
}

// 发送模板邮件-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailSendService) SendTemplateEmail(
	c *gin.Context,
	r req.SendTemplateEmailReq,
) (err error) {
	// 1. 验证参数
	if len(r.To) == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "收件人邮箱不能为空")
	}

	// 2. 获取模板
	var template systemRbac.EmailTemplate
	err = global.GVA_DB.Where("code = ? AND status = ?", r.TemplateCode, "active").First(&template).Error
	if err != nil {
		global.GVA_LOG.Error("模板不存在或未启用", zap.String("code", r.TemplateCode))
		return biz_err.New(biz_err.EMAIL_NOT_FOUND)
	}

	// 3. 渲染模板
	subject, content, err := s.renderTemplate(template.Subject, template.Content, r.Data)
	if err != nil {
		global.GVA_LOG.Error("渲染模板失败", zap.Error(err))
		return biz_err.New(biz_err.SYSTEM_ERROR, "渲染模板失败: "+err.Error())
	}

	// 4. 获取邮件配置
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return biz_err.New(biz_err.PARAM_ERROR, "邮件配置未设置")
	}

	// 5. 创建邮件日志
	emailLog := systemRbac.EmailLog{
		ToEmail:      strings.Join(r.To, ","),
		Subject:      subject,
		Content:      content,
		TemplateCode: r.TemplateCode,
		Status:       "pending",
	}
	global.GVA_DB.Create(&emailLog)

	// 6. 发送邮件
	mailer := NewMailer(*config)
	err = mailer.SendEmail(r.To, nil, nil, subject, content, true)

	// 7. 更新日志状态
	if err != nil {
		global.GVA_DB.Model(&emailLog).Updates(map[string]interface{}{
			"status":    "failed",
			"error_msg": err.Error(),
		})
		return biz_err.New(biz_err.SYSTEM_ERROR, "邮件发送失败: "+err.Error())
	}

	global.GVA_DB.Model(&emailLog).Update("status", "sent")

	global.GVA_LOG.Info("发送模板邮件成功", zap.Strings("to", r.To), zap.String("template", r.TemplateCode))

	return nil
}

// 批量发送邮件-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailSendService) SendBatchEmail(
	c *gin.Context,
	r req.SendBatchEmailReq,
) (rs res.SendBatchEmailRes, err error) {
	// 1. 验证参数
	if len(r.Recipients) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "收件人列表不能为空")
	}

	// 2. 获取模板
	var template systemRbac.EmailTemplate
	err = global.GVA_DB.Where("code = ? AND status = ?", r.TemplateCode, "active").First(&template).Error
	if err != nil {
		return rs, biz_err.New(biz_err.EMAIL_NOT_FOUND)
	}

	// 3. 获取邮件配置
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "邮件配置未设置")
	}

	// 4. 创建批量任务
	batchId := fmt.Sprintf("batch_%d", time.Now().Unix())
	mailer := NewMailer(*config)

	successCount := 0
	failCount := 0

	// 5. 逐个发送邮件
	for _, recipient := range r.Recipients {
		// 渲染模板
		subject, content, err := s.renderTemplate(template.Subject, template.Content, recipient.Data)
		if err != nil {
			global.GVA_LOG.Error("渲染模板失败", zap.String("email", recipient.Email), zap.Error(err))
			failCount++
			continue
		}

		// 创建邮件日志
		emailLog := systemRbac.EmailLog{
			BatchID:      batchId,
			ToEmail:      recipient.Email,
			Subject:      subject,
			Content:      content,
			TemplateCode: r.TemplateCode,
			Status:       "pending",
		}
		global.GVA_DB.Create(&emailLog)

		// 发送邮件
		err = mailer.SendEmail([]string{recipient.Email}, nil, nil, subject, content, true)
		if err != nil {
			global.GVA_LOG.Error("发送邮件失败", zap.String("email", recipient.Email), zap.Error(err))
			global.GVA_DB.Model(&emailLog).Updates(map[string]interface{}{
				"status":    "failed",
				"error_msg": err.Error(),
			})
			failCount++
		} else {
			global.GVA_DB.Model(&emailLog).Update("status", "sent")
			successCount++
		}
	}

	rs = res.SendBatchEmailRes{
		BatchId:    batchId,
		TotalCount: len(r.Recipients),
	}

	global.GVA_LOG.Info("批量发送邮件完成",
		zap.String("batchId", batchId),
		zap.Int("total", len(r.Recipients)),
		zap.Int("success", successCount),
		zap.Int("fail", failCount),
	)

	return rs, nil
}

// renderTemplate 渲染模板
func (s *EmailSendService) renderTemplate(subjectStr, contentStr string, data map[string]interface{}) (subject, content string, err error) {
	// 渲染主题
	subjectTmpl, err := template.New("subject").Parse(subjectStr)
	if err != nil {
		return "", "", err
	}

	var subjectBuf strings.Builder
	err = subjectTmpl.Execute(&subjectBuf, data)
	if err != nil {
		return "", "", err
	}
	subject = subjectBuf.String()

	// 渲染内容
	contentTmpl, err := template.New("content").Parse(contentStr)
	if err != nil {
		return "", "", err
	}

	var contentBuf strings.Builder
	err = contentTmpl.Execute(&contentBuf, data)
	if err != nil {
		return "", "", err
	}
	content = contentBuf.String()

	return subject, content, nil
}

// 添加一个校验模板 校验传递参数与模板参数是否一致!