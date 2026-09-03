package request

type GenerateCaptchaRotateReq struct {
	ImageSquareSize int `json:"imageSquareSize" form:"imageSquareSize"` // 主图大小(可选,默认220)
	ThumbImageSquareSize []int `json:"thumbImageSquareSize" form:"thumbImageSquareSize"` // 缩略图大小(可选)
	RangeAnglePosMin int `json:"rangeAnglePosMin" form:"rangeAnglePosMin"` // 随机角度最小值(可选)
	RangeAnglePosMax int `json:"rangeAnglePosMax" form:"rangeAnglePosMax"` // 随机角度最大值(可选)
	ThumbImageAlpha float32 `json:"thumbImageAlpha" form:"thumbImageAlpha"` // 缩略图透明度(可选)
}

type VerifyCaptchaRotateReq struct {
	Token string `json:"token" form:"token"` // 验证码令牌
	Angle float32 `json:"angle" form:"angle"` // 用户旋转的角度
}

type SaveCaptchaRotateConfigReq struct {
	ImageSquareSize int `json:"imageSquareSize" form:"imageSquareSize"` // 主图大小(默认220)
	ThumbImageSquareSize []int `json:"thumbImageSquareSize" form:"thumbImageSquareSize"` // 缩略图大小
	RangeAnglePosMin int `json:"rangeAnglePosMin" form:"rangeAnglePosMin"` // 随机角度最小值
	RangeAnglePosMax int `json:"rangeAnglePosMax" form:"rangeAnglePosMax"` // 随机角度最大值
	ThumbImageAlpha float32 `json:"thumbImageAlpha" form:"thumbImageAlpha"` // 缩略图透明度
}
