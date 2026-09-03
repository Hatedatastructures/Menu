package response

type GetRegisterConfigRes struct {
	Enabled bool `json:"enabled"` // 开放注册
	VerifyEmail bool `json:"verifyEmail"` // 邮箱验证
	VerifyPhone bool `json:"verifyPhone"` // 手机验证
	DefaultRole string `json:"defaultRole"` // 默认角色
	NeedAudit bool `json:"needAudit"` // 需要审核
	DefaultImage string `json:"defaultImage"` // 默认头像
}

type SaveRegisterConfigRes struct {
	Version int `json:"version"` // 版本号
}
