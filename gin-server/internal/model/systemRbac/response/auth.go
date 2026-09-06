package response

type GetCaptchaRes struct {
	Uuid string `json:"uuid"` // 验证码UUID
	Img string `json:"img"` // 验证码图片(base64)
}

type AuthLoginRes struct {
	ID uint `json:"ID"` // 用户ID
	Uuid string `json:"uuid"` // 用户UUID
	UserName string `json:"userName"` // 用户名
	NickName string `json:"nickName"` // 昵称
	HeaderImg string `json:"headerImg"` // 头像链接
	AuthorityId uint `json:"authorityId"` // 角色ID
	Authority AuthLoginResAuthority `json:"authority"`
	Authorities []AuthLoginResAuthority `json:"authorities"` // []
	Phone string `json:"phone"` // 手机号
	Email string `json:"email"` // 邮箱
	Enable int `json:"enable"` // 是否启用
	Token string `json:"token,omitempty"` // JWT token
	ExpiresAt int64 `json:"expiresAt,omitempty"` // 过期时间（毫秒时间戳）
}


type AuthLoginResAuthority struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认路由
}

type GetPublicConfigRes struct {
	System GetPublicConfigResSystem `json:"system"`
	Login GetPublicConfigResLogin `json:"login"`
	Register GetPublicConfigResRegister `json:"register"`
	Password GetPublicConfigResPassword `json:"password"`
	Storage GetPublicConfigResStorage `json:"storage"`
	Security GetPublicConfigResSecurity `json:"security"`
}

type GetPublicConfigResSystem struct {
	SiteName string `json:"siteName"` // 站点名称
	SiteDescription string `json:"siteDescription"` // 站点描述
	SiteLogo string `json:"siteLogo"` // 站点Logo
	Copyright string `json:"copyright"` // 版权信息
	Icp string `json:"icp"` // ICP备案号
	WatermarkEnabled bool `json:"watermarkEnabled"` // 是否启用水印
	WatermarkType string `json:"watermarkType"` // 水印类型
	WatermarkCustomText string `json:"watermarkCustomText"` // 自定义水印文字
	WatermarkOpacity float32 `json:"watermarkOpacity"` // 水印透明度
}

type GetPublicConfigResLogin struct {
	CaptchaEnabled bool `json:"captchaEnabled"` // 是否启用验证码
	CaptchaType string `json:"captchaType"` // 验证码类型(image/slider/sms)
	MaxRetryCount int `json:"maxRetryCount"` // 最大重试次数
	RememberMe bool `json:"rememberMe"` // 是否启用记住我
}

type GetPublicConfigResRegister struct {
	Enabled bool `json:"enabled"` // 是否开放注册
	VerifyEmail bool `json:"verifyEmail"` // 是否验证邮箱
	VerifyPhone bool `json:"verifyPhone"` // 是否验证手机号
	NeedAudit bool `json:"needAudit"` // 注册是否需要审核
}

type GetPublicConfigResPassword struct {
	MinLength int `json:"minLength"` // 密码最小长度
	MaxLength int `json:"maxLength"` // 密码最大长度
	RequireUppercase bool `json:"requireUppercase"` // 必须包含大写字母
	RequireLowercase bool `json:"requireLowercase"` // 必须包含小写字母
	RequireNumber bool `json:"requireNumber"` // 必须包含数字
	RequireSpecial bool `json:"requireSpecial"` // 必须包含特殊字符
}

type GetPublicConfigResStorage struct {
	MaxSize int `json:"maxSize"` // 最大文件大小(MB)
	AllowTypes []string `json:"allowTypes"` // 允许的文件类型
}

type GetPublicConfigResSecurity struct {
	DisableDevtool bool `json:"disableDevtool"` // 是否禁用开发者工具
}
