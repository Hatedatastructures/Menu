package response

type GetLoginConfigRes struct {
	CaptchaEnabled bool `json:"captchaEnabled"` // 验证码
	CaptchaType string `json:"captchaType"` // 验证码类型
	MaxRetryCount int `json:"maxRetryCount"` // 最大重试次数
	LockTime int `json:"lockTime"` // 锁定时间(分钟)
	RememberMe bool `json:"rememberMe"` // 记住我
	SingleLogin bool `json:"singleLogin"` // 单点登录
	// 滑块验证码配置
	SliderCaptchaWidth int `json:"sliderCaptchaWidth"` // 滑块宽度
	SliderCaptchaHeight int `json:"sliderCaptchaHeight"` // 滑块高度
	SliderThumbWidth int `json:"sliderThumbWidth"` // 缩略图宽度
	SliderThumbHeight int `json:"sliderThumbHeight"` // 缩略图高度
	SliderVerticalPadding int `json:"sliderVerticalPadding"` // 垂直内边距
	SliderHorizontalPadding int `json:"sliderHorizontalPadding"` // 水平内边距
	SliderShowTheme bool `json:"sliderShowTheme"` // 显示主题
	SliderTitle string `json:"sliderTitle"` // 标题
	SliderButtonText string `json:"sliderButtonText"` // 按钮文字
	SliderIconSize int `json:"sliderIconSize"` // 图标大小
	SliderDotSize int `json:"sliderDotSize"` // 圆点大小
}

type SaveLoginConfigRes struct {
	Version int `json:"version"` // 版本号
}
