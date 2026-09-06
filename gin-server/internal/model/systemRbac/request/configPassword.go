package request

type SavePasswordConfigReq struct {
	MinLength int `json:"minLength" form:"minLength"` // 最小长度
	MaxLength int `json:"maxLength" form:"maxLength"` // 最大长度
	RequireUppercase bool `json:"requireUppercase" form:"requireUppercase"` // 必须包含大写字母
	RequireLowercase bool `json:"requireLowercase" form:"requireLowercase"` // 必须包含小写字母
	RequireNumber bool `json:"requireNumber" form:"requireNumber"` // 必须包含数字
	RequireSpecial bool `json:"requireSpecial" form:"requireSpecial"` // 必须包含特殊字符
	ExpireDays int `json:"expireDays" form:"expireDays"` // 密码过期天数
}
