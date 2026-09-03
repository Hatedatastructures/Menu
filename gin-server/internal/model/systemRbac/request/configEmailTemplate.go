package request

type SaveEmailTemplateConfigReq struct {
	VerifyCode string `json:"verifyCode" form:"verifyCode"` // 验证码邮件模板
	ResetPassword string `json:"resetPassword" form:"resetPassword"` // 重置密码邮件模板
	Welcome string `json:"welcome" form:"welcome"` // 欢迎邮件模板
}
