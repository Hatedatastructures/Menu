package request

type SendVerifyCodeReq struct {
	Email string `json:"email" form:"email"` // 邮箱地址
	Type string `json:"type" form:"type"` // 验证码类型: register/login/reset_password
	CaptchaToken string `json:"captchaToken" form:"captchaToken"` // 滑块验证token
}

type VerifyCodeReq struct {
	Email string `json:"email" form:"email"` // 邮箱地址
	Code string `json:"code" form:"code"` // 验证码
	Type string `json:"type" form:"type"` // 验证码类型: register/login/reset_password
}

type GetVerifyCodeStatusReq struct {
	Email string `json:"email" form:"email"`
	Type string `json:"type" form:"type"`
}
