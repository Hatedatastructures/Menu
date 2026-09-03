package response

type GetProfileRes struct {
	Uuid string `json:"uuid"` // 用户UUID
	Username string `json:"username"` // 用户登录名
	NickName string `json:"nickName"` // 用户昵称
	HeaderImg string `json:"headerImg"` // 用户头像
	Phone string `json:"phone"` // 用户手机号
	Email string `json:"email"` // 用户邮箱
	AuthorityId uint `json:"authorityId"` // 用户角色ID
	Enable int `json:"enable"` // 用户是否被冻结 1正常 2冻结
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type UpdateProfileRes struct {
	Id uint `json:"id"` // 用户ID
}

type UpdateAvatarRes struct {
	HeaderImg string `json:"headerImg"` // 更新后的头像URL
}

type UploadAvatarRes struct {
	Url string `json:"url"` // 头像URL
	Filename string `json:"filename"` // 文件名
	Size int64 `json:"size"` // 文件大小
}
