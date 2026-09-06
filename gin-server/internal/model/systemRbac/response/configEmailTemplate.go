package response

type GetEmailTemplateConfigRes struct {
	VerifyCode string `json:"verifyCode"` // 验证码邮件模板
	ResetPassword string `json:"resetPassword"` // 重置密码邮件模板
	Welcome string `json:"welcome"` // 欢迎邮件模板
}

type SaveEmailTemplateConfigRes struct {
	Version int `json:"version"` // 版本号
}
