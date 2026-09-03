package response

type GenerateCaptchaRotateRes struct {
	Token string `json:"token"` // 验证码令牌
	MasterImage string `json:"masterImage"` // 主图Base64(JPEG格式)
	ThumbImage string `json:"thumbImage"` // 缩略图Base64(PNG格式)
	Width int `json:"width"` // 主图宽度
	Height int `json:"height"` // 主图高度
}

type VerifyCaptchaRotateRes struct {
	Success bool `json:"success"` // 是否验证成功
	Message string `json:"message"` // 提示信息
}

type GetCaptchaRotateConfigRes struct {
	ImageSquareSize int `json:"imageSquareSize"` // 主图大小
	ThumbImageSquareSize []int `json:"thumbImageSquareSize"` // 缩略图大小
	RangeAnglePosMin int `json:"rangeAnglePosMin"` // 随机角度最小值
	RangeAnglePosMax int `json:"rangeAnglePosMax"` // 随机角度最大值
	ThumbImageAlpha float32 `json:"thumbImageAlpha"` // 缩略图透明度
	Version int `json:"version"` // 版本号
}

type SaveCaptchaRotateConfigRes struct {
	Version int `json:"version"` // 版本号
}

type PreviewCaptchaRotateRes struct {
	MasterImage string `json:"masterImage"` // 主图Base64
	ThumbImage string `json:"thumbImage"` // 缩略图Base64
	Angle int `json:"angle"` // 目标角度
}
