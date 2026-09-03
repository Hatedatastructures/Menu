package request

type RegisterReq struct {
	UserName string `json:"userName" form:"userName"` // 用户登录名
	PassWord string `json:"passWord" form:"passWord"` // 用户登录密码
	NickName string `json:"nickName" form:"nickName"` // 用户昵称
	HeaderImg string `json:"headerImg" form:"headerImg"` // 用户头像
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 用户角色ID
	Enable int `json:"enable" form:"enable"` // 是否启用
	AuthorityIds []uint `json:"authorityIds" form:"authorityIds"` // 多角色ID
	Phone string `json:"phone" form:"phone"` // 用户手机号
	Email string `json:"email" form:"email"` // 用户邮箱
}

type LoginReq struct {
	Username string `json:"username" form:"username"` // 用户名
	Password string `json:"password" form:"password"` // 密码
	Captcha string `json:"captcha" form:"captcha"` // 验证码
	CaptchaId string `json:"captchaId" form:"captchaId"` // 验证码ID
}

type ChangePasswordReq struct {
	Password string `json:"password" form:"password"` // 原密码
	NewPassword string `json:"newPassword" form:"newPassword"` // 新密码
}

type GetUserInfoListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Username string `json:"username" form:"username"`
	NickName string `json:"nickName" form:"nickName"`
	Phone string `json:"phone" form:"phone"`
	Email string `json:"email" form:"email"`
}

type SetUserAuthorityReq struct {
	Id int `json:"id" form:"id"`
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type SetUserAuthoritiesReq struct {
	Id int64 `json:"id" form:"id"` // 用户ID
	AuthorityIds []uint `json:"authorityIds" form:"authorityIds"` // 角色ID数组
}

type DeleteUserReq struct {
	Id int `json:"id" form:"id"`
}

type SetUserInfoReq struct {
	Id int64 `json:"id" form:"id"` // 用户ID
	NickName string `json:"nickName" form:"nickName"` // 用户昵称
	HeaderImg string `json:"headerImg" form:"headerImg"` // 用户头像
	Phone string `json:"phone" form:"phone"` // 用户手机号
	Email string `json:"email" form:"email"` // 用户邮箱
	Enable int `json:"enable" form:"enable"` // 是否启用
}

type SetSelfInfoReq struct {
	Id int64 `json:"id" form:"id"` // 用户ID
	UserName string `json:"userName" form:"userName"` // 用户登录名
	NickName string `json:"nickName" form:"nickName"` // 用户昵称
	HeaderImg string `json:"headerImg" form:"headerImg"` // 用户头像
	Phone string `json:"phone" form:"phone"` // 用户手机号
	Email string `json:"email" form:"email"` // 用户邮箱
}

type SetSelfSettingReq struct {
	OriginSetting map[string]interface{} `json:"originSetting" form:"originSetting"` // 配置JSON
	Uid uint `json:"uid" form:"uid"` // 用户ID
}

type GetUserInfReq struct {
	Uuid string `json:"uuid" form:"uuid"`
}

type FindUserByIdReq struct {
	Id int `json:"id" form:"id"`
}

type FindUserByUuidReq struct {
	Uuid string `json:"uuid" form:"uuid"`
}

type ResetPasswordReq struct {
	Id int64 `json:"id" form:"id"` // 用户ID
	Password string `json:"password" form:"password"` // 新密码
}
