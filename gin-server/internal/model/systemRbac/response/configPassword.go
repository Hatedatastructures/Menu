package response

type GetPasswordConfigRes struct {
	MinLength int `json:"minLength"` // 最小长度
	MaxLength int `json:"maxLength"` // 最大长度
	RequireUppercase bool `json:"requireUppercase"` // 必须包含大写字母
	RequireLowercase bool `json:"requireLowercase"` // 必须包含小写字母
	RequireNumber bool `json:"requireNumber"` // 必须包含数字
	RequireSpecial bool `json:"requireSpecial"` // 必须包含特殊字符
	ExpireDays int `json:"expireDays"` // 密码过期天数
}

type SavePasswordConfigRes struct {
	Version int `json:"version"` // 版本号
}
