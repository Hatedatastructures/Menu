package request

type AuthLoginReq struct {
	Username string `json:"username" form:"username"` // 用户名
	Password string `json:"password" form:"password"` // 密码
	LoginType string `json:"loginType" form:"loginType"` // 登录类型(默认PASSWORD)

	// 验证码类型 (image/slider/click/rotate)
	CaptchaType string `json:"captchaType" form:"captchaType"` // 验证码类型

	// 图片验证码
	CaptchaUuid string `json:"captchaUuid" form:"captchaUuid"` // 验证码UUID(图片验证码时必填)
	CaptchaCode string `json:"captchaCode" form:"captchaCode"` // 验证码(图片验证码时必填)

	// 滑块验证码
	SliderToken string `json:"sliderToken" form:"sliderToken"` // 滑块验证码token
	SliderX int `json:"sliderX" form:"sliderX"` // 滑动X坐标
	SliderY int `json:"sliderY" form:"sliderY"` // 滑动Y坐标

	// 点选验证码
	ClickToken string `json:"clickToken" form:"clickToken"` // 点选验证码token
	ClickPoints []VerifyCaptchaClickReqPoint `json:"clickPoints" form:"clickPoints"` // 用户点击坐标

	// 旋转验证码
	RotateToken string `json:"rotateToken" form:"rotateToken"` // 旋转验证码token
	RotateAngle float32 `json:"rotateAngle" form:"rotateAngle"` // 旋转角度
}

type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"` // 旧密码
	NewPassword string `json:"newPassword" form:"newPassword"` // 新密码
}

type AuthRegisterReq struct {
	Username string `json:"username" form:"username"` // 用户名(4-20位字母数字下划线)
	Password string `json:"password" form:"password"` // 密码
	Nickname string `json:"nickname" form:"nickname"` // 昵称
	Email string `json:"email" form:"email"` // 邮箱(注册验证开启时必填)
	Phone string `json:"phone" form:"phone"` // 手机号(注册验证开启时必填)
	Uuid string `json:"uuid" form:"uuid"` // 验证码UUID(验证码开启时必填)
	Code string `json:"code" form:"code"` // 验证码(验证码开启时必填)
	HeaderImg string `json:"headerImg" form:"headerImg"` // 头像
}
