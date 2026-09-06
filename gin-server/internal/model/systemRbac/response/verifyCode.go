package response

type SendVerifyCodeRes struct {
	ExpireIn int `json:"expireIn"` // 验证码有效期(秒)
	CanResendIn int `json:"canResendIn"` // 可重新发送时间(秒)
}

type VerifyCodeRes struct {
	Valid bool `json:"valid"` // 是否验证成功
}

type GetVerifyCodeStatusRes struct {
	Email string `json:"email"`
	Type string `json:"type"`
	CanSend bool `json:"canSend"` // 是否可以发送
	ResendIn int `json:"resendIn"` // 距离可重发时间(秒)
	ExpireIn int `json:"expireIn"` // 验证码剩余有效期(秒)
}
