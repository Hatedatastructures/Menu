package systemRbac

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"shack/internal/global"
	model "shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// Redis key 前缀
	RedisKeyEmailCode       = "email:code:"        // 验证码
	RedisKeyEmailCodeExpire = "email:code:expire:" // 验证码过期时间
	RedisKeyEmailCodeRetry  = "email:code:retry:"  // 验证码重试次数
	RedisKeyEmailCodeUsed   = "email:code:used:"   // 验证码已使用
	RedisKeyEmailLimit      = "limit:email:"       // 邮箱限流

	// 验证码配置
	CodeExpire      = 5 * 60 // 5分钟过期
	CodeMaxRetry    = 5      // 最多验证5次
	CodeResendSec   = 60     // 60秒内不能重发
	EmailLimitSec   = 60     // 邮箱限流60秒
	EmailLimitDaily = 10     // 每天最多10次
)

type VerifyCodeService struct{}

// 发送验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *VerifyCodeService) SendVerifyCode(
	c *gin.Context,
	r req.SendVerifyCodeReq,
) (rs res.SendVerifyCodeRes, err error) {
	email := r.Email
	codeType := r.Type

	// 1. 检查邮箱限流
	limitKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailLimit, email, codeType)
	count, _ := global.GVA_REDIS.Exists(c, limitKey).Result()
	if count > 0 {
		ttl, _ := global.GVA_REDIS.TTL(c, limitKey).Result()
		return rs, biz_err.New(biz_err.EMAIL_LIMIT_EXCEEDED, fmt.Sprintf("请%d秒后再试", ttl))
	}

	// 2. 检查每日发送次数
	dailyKey := fmt.Sprintf("email:daily:%s", email)
	dailyCount, _ := global.GVA_REDIS.Get(c, dailyKey).Int()
	if dailyCount >= EmailLimitDaily {
		return rs, biz_err.New(biz_err.EMAIL_LIMIT_EXCEEDED, "今日发送次数已达上限")
	}

	// 3. 生成6位验证码
	code, err := s.generateCode()
	if err != nil {
		global.GVA_LOG.Error("生成验证码失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR)
	}

	// 4. 存储验证码到Redis
	codeKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailCode, email, codeType)
	err = global.GVA_REDIS.Set(c, codeKey, code, time.Duration(CodeExpire)*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("存储验证码失败", zap.Error(err))
		return rs, biz_err.New(biz_err.CACHE_ERROR)
	}

	// 5. 设置限流
	err = global.GVA_REDIS.Set(c, limitKey, "1", time.Duration(EmailLimitSec)*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("设置限流失败", zap.Error(err))
	}

	// 6. 更新每日发送次数
	err = global.GVA_REDIS.Incr(c, dailyKey).Err()
	if err == nil {
		global.GVA_REDIS.Expire(c, dailyKey, 24*3600*time.Second)
	}

	// 7. 发送邮件 - 使用邮件模板服务
	// 7.1 获取邮件配置
	emailConfig, err := GetEmailConfigFromDB()
	if err != nil {
		global.GVA_LOG.Error("获取邮件配置失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "邮件服务未配置")
	}

	if !emailConfig.Enabled {
		global.GVA_LOG.Warn("邮件服务未启用", zap.String("email", email))
		// 邮件服务未启用时，仅记录日志，不返回错误（用于开发测试）
		global.GVA_LOG.Info("验证码（未发送邮件）", zap.String("email", email), zap.String("code", code), zap.String("type", codeType))
		rs = res.SendVerifyCodeRes{
			ExpireIn:    CodeExpire,
			CanResendIn: EmailLimitSec,
		}
		return rs, nil
	}

	// 7.2 获取邮件模板
	var template model.EmailTemplate
	templateCode := "register_code" // 默认注册验证码模板
 if codeType == "password_reset_code" {
		templateCode = "reset_password_code"
	}

	err = global.GVA_DB.Where("code = ? AND status = ?", templateCode, "active").First(&template).Error
	if err != nil {
		global.GVA_LOG.Error("获取邮件模板失败", zap.String("templateCode", templateCode), zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "邮件模板不存在")
	}

	// 7.3 准备模板数据
	typeDesc := "注册"
	if codeType == "password_reset_code" {
		typeDesc = "重置密码"
	} 

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	templateData := map[string]interface{}{
		"Code":      code,
		"TypeDesc":  typeDesc,
		"Email":     email,
		"Timestamp": timestamp,
		"ExpireIn":  CodeExpire / 60, // 转换为分钟
	}

	// 7.4 渲染主题
	subject, err := RenderTemplate(template.Subject, templateData)
	if err != nil {
		global.GVA_LOG.Error("渲染邮件主题失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "渲染邮件主题失败")
	}

	// 7.5 渲染内容
	content, err := RenderTemplate(template.Content, templateData)
	if err != nil {
		global.GVA_LOG.Error("渲染邮件内容失败", zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "渲染邮件内容失败")
	}

	// 7.6 创建邮件发送器并发送
	mailer := NewMailer(*emailConfig)
	err = mailer.SendEmail([]string{email}, nil, nil, subject, content, true)
	if err != nil {
		global.GVA_LOG.Error("发送邮件失败", zap.String("email", email), zap.Error(err))
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "发送邮件失败")
	}

	global.GVA_LOG.Info("发送验证码邮件成功", zap.String("email", email), zap.String("type", codeType))

	rs = res.SendVerifyCodeRes{
		ExpireIn:    CodeExpire,
		CanResendIn: EmailLimitSec,
	}

	return rs, nil
}

// 验证验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *VerifyCodeService) VerifyCode(
	c *gin.Context,
	r req.VerifyCodeReq,
) (rs res.VerifyCodeRes, err error) {
	email := r.Email
	code := r.Code
	codeType := r.Type

	// 1. 检查验证码是否存在
	codeKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailCode, email, codeType)
	savedCode, err := global.GVA_REDIS.Get(c, codeKey).Result()
	if err != nil {
		return rs, biz_err.New(biz_err.VERIFY_CODE_EXPIRED)
	}

	// 2. 检查是否已使用
	usedKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailCodeUsed, email, codeType)
	isUsed, _ := global.GVA_REDIS.Exists(c, usedKey).Result()
	if isUsed > 0 {
		return rs, biz_err.New(biz_err.VERIFY_CODE_EXPIRED, "验证码已使用")
	}

	// 3. 验证码比对
	if savedCode != code {
		// 增加重试次数
		retryKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailCodeRetry, email, codeType)
		retryCount, _ := global.GVA_REDIS.Incr(c, retryKey).Result()
		global.GVA_REDIS.Expire(c, retryKey, time.Duration(CodeExpire)*time.Second)

		if retryCount >= CodeMaxRetry {
			// 超过次数，删除验证码
			global.GVA_REDIS.Del(c, codeKey)
			return rs, biz_err.New(biz_err.VERIFY_CODE_ERROR, "验证码错误次数过多")
		}
		return rs, biz_err.New(biz_err.VERIFY_CODE_ERROR)
	}

	// 4. 标记为已使用
	global.GVA_REDIS.Set(c, usedKey, "1", time.Duration(CodeExpire)*time.Second)

	// 5. 删除验证码
	global.GVA_REDIS.Del(c, codeKey)

	rs = res.VerifyCodeRes{
		Valid: true,
	}

	return rs, nil
}

// 获取验证码状态-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *VerifyCodeService) GetVerifyCodeStatus(
	c *gin.Context,
	r req.GetVerifyCodeStatusReq,
) (rs res.GetVerifyCodeStatusRes, err error) {
	email := r.Email
	codeType := r.Type

	// 检查限流
	limitKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailLimit, email, codeType)
	ttl, _ := global.GVA_REDIS.TTL(c, limitKey).Result()

	rs = res.GetVerifyCodeStatusRes{
		CanSend:  ttl <= 0,
		ResendIn: int(ttl),
	}

	// 如果验证码存在，获取剩余时间
	codeKey := fmt.Sprintf("%s%s:%s", RedisKeyEmailCode, email, codeType)
	codeTTL, _ := global.GVA_REDIS.TTL(c, codeKey).Result()
	if codeTTL > 0 {
		rs.ExpireIn = int(codeTTL)
	}

	return rs, nil
}

// generateCode 生成6位数字验证码
func (s *VerifyCodeService) generateCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}
	return code, nil
}
