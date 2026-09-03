package request

type UpdateProfileReq struct {
	NickName string `json:"nickName" form:"nickName"` // 用户昵称
	Phone string `json:"phone" form:"phone"` // 用户手机号
	Email string `json:"email" form:"email"` // 用户邮箱
}

type UpdatePassword1Req struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"` // 旧密码
	NewPassword string `json:"newPassword" form:"newPassword"` // 新密码
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword"` // 确认密码
}

type UpdateAvatarReq struct {
	HeaderImg string `json:"headerImg" form:"headerImg"` // 头像URL
}

type UploadAvatarReq struct {
	File string `json:"file" form:"file"` // 头像文件
}
