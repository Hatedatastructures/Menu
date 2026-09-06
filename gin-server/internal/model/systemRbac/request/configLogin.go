package request

type SaveLoginConfigReq struct {
	CaptchaEnabled bool `json:"captchaEnabled" form:"captchaEnabled"` // 验证码
	CaptchaType string `json:"captchaType" form:"captchaType"` // 验证码类型
	MaxRetryCount int `json:"maxRetryCount" form:"maxRetryCount"` // 最大重试次数
	LockTime int `json:"lockTime" form:"lockTime"` // 锁定时间(分钟)
	RememberMe bool `json:"rememberMe" form:"rememberMe"` // 记住我
	SingleLogin bool `json:"singleLogin" form:"singleLogin"` // 单点登录
	// 滑块验证码配置
	SliderCaptchaWidth int `json:"sliderCaptchaWidth" form:"sliderCaptchaWidth"` // 滑块宽度
	SliderCaptchaHeight int `json:"sliderCaptchaHeight" form:"sliderCaptchaHeight"` // 滑块高度
	SliderThumbWidth int `json:"sliderThumbWidth" form:"sliderThumbWidth"` // 缩略图宽度
	SliderThumbHeight int `json:"sliderThumbHeight" form:"sliderThumbHeight"` // 缩略图高度
	SliderVerticalPadding int `json:"sliderVerticalPadding" form:"sliderVerticalPadding"` // 垂直内边距
	SliderHorizontalPadding int `json:"sliderHorizontalPadding" form:"sliderHorizontalPadding"` // 水平内边距
	SliderShowTheme bool `json:"sliderShowTheme" form:"sliderShowTheme"` // 显示主题
	SliderTitle string `json:"sliderTitle" form:"sliderTitle"` // 标题
	SliderButtonText string `json:"sliderButtonText" form:"sliderButtonText"` // 按钮文字
	SliderIconSize int `json:"sliderIconSize" form:"sliderIconSize"` // 图标大小
	SliderDotSize int `json:"sliderDotSize" form:"sliderDotSize"` // 圆点大小
}
