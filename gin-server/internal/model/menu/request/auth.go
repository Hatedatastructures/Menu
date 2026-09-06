package request

type RegisterReq struct {
	Email string `json:"email" form:"email" validate:"required,email"` // 用户邮箱
	Password string `json:"password" form:"password" validate:"required,max=128,min=8"` // 密码(8-128位)
	DisplayName string `json:"displayName" form:"displayName" validate:"required,max=64,min=1"` // 显示名称
}

type LoginReq struct {
	Email string `json:"email" form:"email" validate:"required,email"` // 用户邮箱
	Password string `json:"password" form:"password" validate:"required,max=128,min=1"` // 密码
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" form:"refreshToken"` // 刷新令牌
}
