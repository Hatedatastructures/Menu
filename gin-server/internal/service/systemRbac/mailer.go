package systemRbac

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"

	biz_err "shack/internal/error"

	"gopkg.in/gomail.v2"
	"go.uber.org/zap"
	"text/template"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	SSL      bool
}

// Mailer 邮件发送器
type Mailer struct {
	config EmailConfig
}

// NewMailer 创建邮件发送器
func NewMailer(config EmailConfig) *Mailer {
	return &Mailer{
		config: config,
	}
}

// SendEmail 发送邮件
func (m *Mailer) SendEmail(
	to []string,
	cc []string,
	bcc []string,
	subject string,
	content string,
	isHTML bool,
) error {
	if !m.config.Enabled {
		return biz_err.New(biz_err.PARAM_ERROR, "邮件服务未启用")
	}

	// 创建邮件消息
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.formatAddress())
	msg.SetHeader("To", to...)
	if len(cc) > 0 {
		msg.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		msg.SetHeader("Bcc", bcc...)
	}
	msg.SetHeader("Subject", subject)

	if isHTML {
		msg.SetBody("text/html", content)
	} else {
		msg.SetBody("text/plain", content)
	}

	// 创建SMTP拨号器
	dialer := gomail.NewDialer(
		m.config.Host,
		m.config.Port,
		m.config.Username,
		m.config.Password,
	)

	// SSL/TLS配置
	if m.config.SSL {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         m.config.Host,
		}
	}

	// 发送邮件
	err := dialer.DialAndSend(msg)
	if err != nil {
		global.GVA_LOG.Error("发送邮件失败",
			zap.Strings("to", to),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	global.GVA_LOG.Info("邮件发送成功",
		zap.Strings("to", to),
		zap.String("subject", subject),
	)

	return nil
}

// SendTestEmail 发送测试邮件
func (m *Mailer) SendTestEmail(to string) error {
	subject := "邮件配置测试"
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// HTML格式测试邮件
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; }
        .info { background: white; padding: 15px; margin: 10px 0; border-left: 4px solid #667eea; }
        .success { color: #28a745; font-weight: bold; }
        .footer { text-align: center; margin-top: 20px; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🎉 邮件配置测试</h2>
        </div>
        <div class="content">
            <p class="success">✓ 邮件服务配置成功！</p>

            <div class="info">
                <strong>测试时间：</strong><br>
                %s
            </div>

            <div class="info">
                <strong>邮件服务器信息：</strong><br>
                服务器：%s<br>
                端口：%d<br>
                用户名：%s<br>
                发件人：%s<br>
                加密方式：%s
            </div>

            <div class="info">
                <strong>收件人：</strong><br>
                %s
            </div>

            <p>如果您收到此邮件，说明您的邮件服务配置正确，可以正常发送邮件了！</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
        </div>
    </div>
</body>
</html>`,
		timestamp,
		m.config.Host,
		m.config.Port,
		m.config.Username,
		m.config.FromName,
		map[bool]string{true: "SSL/TLS", false: "STARTTLS"}[m.config.SSL],
		to,
	)

	return m.SendEmail([]string{to}, nil, nil, subject, htmlContent, true)
}

// formatAddress 格式化发件人地址
func (m *Mailer) formatAddress() string {
	if m.config.FromName != "" {
		return fmt.Sprintf("%s <%s>", m.config.FromName, m.config.Username)
	}
	return m.config.Username
}

// GetEmailConfigFromDB 从数据库获取邮件配置
func GetEmailConfigFromDB() (*EmailConfig, error) {
	var configGroup systemRbac.ConfigGroup
	err := global.GVA_DB.Where("code = ?", "email").First(&configGroup).Error
	if err != nil {
		return nil, fmt.Errorf("邮件配置未设置")
	}

	var configMap map[string]interface{}
	if configGroup.Config != nil {
		json.Unmarshal(configGroup.Config, &configMap)
	} else {
		configMap = make(map[string]interface{})
	}

	config := &EmailConfig{
		Enabled:  getBoolFromMap(configMap, "enabled", false),
		Host:     getStringFromMap(configMap, "host", ""),
		Port:     getIntFromMap(configMap, "port", 465),
		Username: getStringFromMap(configMap, "username", ""),
		Password: getStringFromMap(configMap, "password", ""),
		FromName: getStringFromMap(configMap, "fromName", ""),
		SSL:      getBoolFromMap(configMap, "ssl", true),
	}

	return config, nil
}

// RenderTemplate 渲染模板
func RenderTemplate(templateContent string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	return buf.String(), nil
}

// SendEmailWithTemplate 使用模板发送邮件
func (m *Mailer) SendEmailWithTemplate(
	to []string,
	cc []string,
	bcc []string,
	templateContent string,
	templateData map[string]interface{},
	isHTML bool,
) error {
	// 渲染模板
	content, err := RenderTemplate(templateContent, templateData)
	if err != nil {
		return err
	}

	// 从模板数据中获取主题
	subject := "邮件"
	if subj, ok := templateData["subject"].(string); ok {
		subject = subj
	}

	return m.SendEmail(to, cc, bcc, subject, content, isHTML)
}

// SendTestEmailWithTemplate 使用指定模板发送测试邮件
func (m *Mailer) SendTestEmailWithTemplate(to string, templateId uint) error {
	// 从数据库获取模板
	var emailTemplate systemRbac.EmailTemplate
	err := global.GVA_DB.Where("id = ? AND status = ?", templateId, "active").First(&emailTemplate).Error
	if err != nil {
		return fmt.Errorf("模板不存在或未启用: %w", err)
	}

	// 准备模板数据
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	templateData := map[string]interface{}{
		"subject":     emailTemplate.Subject,
		"timestamp":   timestamp,
		"email":       to,
		"host":        m.config.Host,
		"port":        m.config.Port,
		"username":    m.config.Username,
		"fromName":    m.config.FromName,
		"encryption":  map[bool]string{true: "SSL/TLS", false: "STARTTLS"}[m.config.SSL],
	}

	// 渲染模板
	subject, err := RenderTemplate(emailTemplate.Subject, templateData)
	if err != nil {
		return fmt.Errorf("渲染邮件主题失败: %w", err)
	}

	content, err := RenderTemplate(emailTemplate.Content, templateData)
	if err != nil {
		return fmt.Errorf("渲染邮件内容失败: %w", err)
	}

	return m.SendEmail([]string{to}, nil, nil, subject, content, true)
}
