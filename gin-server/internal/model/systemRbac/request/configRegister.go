package request

type SaveRegisterConfigReq struct {
	Enabled bool `json:"enabled" form:"enabled"` // 开放注册
	VerifyEmail bool `json:"verifyEmail" form:"verifyEmail"` // 邮箱验证
	VerifyPhone bool `json:"verifyPhone" form:"verifyPhone"` // 手机验证
	DefaultRole string `json:"defaultRole" form:"defaultRole"` // 默认角色
	NeedAudit bool `json:"needAudit" form:"needAudit"` // 需要审核
	DefaultImage string `json:"defaultImage" form:"defaultImage"` // 默认头像
}
