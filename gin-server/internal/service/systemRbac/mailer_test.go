package systemRbac

import (
	"testing"
)

// TestEmailConfig 测试邮件配置结构
func TestEmailConfig(t *testing.T) {
	config := EmailConfig{
		Enabled:  true,
		Host:     "smtp.example.com",
		Port:     465,
		Username: "test@example.com",
		Password: "password",
		FromName: "Test Sender",
		SSL:      true,
	}

	mailer := NewMailer(config)

	if mailer == nil {
		t.Error("创建邮件发送器失败")
	}

	if mailer.config.Host != config.Host {
		t.Errorf("邮件配置不正确，期望 %s，实际 %s", config.Host, mailer.config.Host)
	}
}

// TestFormatAddress 测试发件人地址格式化
func TestFormatAddress(t *testing.T) {
	tests := []struct {
		name     string
		fromName string
		username string
		expected string
	}{
		{
			name:     "有发件人名称",
			fromName: "Test Sender",
			username: "test@example.com",
			expected: "Test Sender <test@example.com>",
		},
		{
			name:     "无发件人名称",
			fromName: "",
			username: "test@example.com",
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := EmailConfig{
				FromName: tt.fromName,
				Username: tt.username,
			}
			mailer := NewMailer(config)
			result := mailer.formatAddress()

			if result != tt.expected {
				t.Errorf("formatAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNewMailer 测试创建邮件发送器
func TestNewMailer(t *testing.T) {
	config := EmailConfig{
		Enabled:  true,
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "user@gmail.com",
		Password: "app-password",
		FromName: "My App",
		SSL:      false,
	}

	mailer := NewMailer(config)

	if mailer == nil {
		t.Fatal("NewMailer() 返回 nil")
	}

	if mailer.config.Host != config.Host {
		t.Errorf("Host 配置不正确，期望 %s，实际 %s", config.Host, mailer.config.Host)
	}

	if mailer.config.Port != config.Port {
		t.Errorf("Port 配置不正确，期望 %d，实际 %d", config.Port, mailer.config.Port)
	}

	if mailer.config.Username != config.Username {
		t.Errorf("Username 配置不正确，期望 %s，实际 %s", config.Username, mailer.config.Username)
	}

	if mailer.config.SSL != config.SSL {
		t.Errorf("SSL 配置不正确，期望 %v，实际 %v", config.SSL, mailer.config.SSL)
	}
}

// BenchmarkFormatAddress 性能测试
func BenchmarkFormatAddress(b *testing.B) {
	config := EmailConfig{
		FromName: "Test Sender",
		Username: "test@example.com",
	}
	mailer := NewMailer(config)

	for i := 0; i < b.N; i++ {
		mailer.formatAddress()
	}
}
