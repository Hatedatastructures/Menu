package response

type RegisterRes struct {
	User RegisterResUser `json:"user"`
	AccessToken string `json:"accessToken"` // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
	AccessTokenExpiresInSeconds int `json:"accessTokenExpiresInSeconds"` // 访问令牌过期秒数
}

type RegisterResUser struct {
	Id string `json:"id"` // 用户ID
	Email string `json:"email"` // 邮箱
	DisplayName string `json:"displayName"` // 显示名称
	IsAdmin bool `json:"isAdmin"` // 是否管理员
}

type LoginRes struct {
	User LoginResUser `json:"user"`
	AccessToken string `json:"accessToken"` // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
	AccessTokenExpiresInSeconds int `json:"accessTokenExpiresInSeconds"` // 访问令牌过期秒数
}

type LoginResUser struct {
	Id string `json:"id"` // 用户ID
	Email string `json:"email"` // 邮箱
	DisplayName string `json:"displayName"` // 显示名称
	IsAdmin bool `json:"isAdmin"` // 是否管理员
}

type RefreshTokenRes struct {
	User RefreshTokenResUser `json:"user"`
	AccessToken string `json:"accessToken"` // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
	AccessTokenExpiresInSeconds int `json:"accessTokenExpiresInSeconds"` // 访问令牌过期秒数
}

type RefreshTokenResUser struct {
	Id string `json:"id"` // 用户ID
	Email string `json:"email"` // 邮箱
	DisplayName string `json:"displayName"` // 显示名称
	IsAdmin bool `json:"isAdmin"` // 是否管理员
}
