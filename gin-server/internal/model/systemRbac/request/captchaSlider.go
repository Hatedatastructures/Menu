package request

type GenerateCaptchaSliderReq struct {
	MasterSizeWidth int `json:"masterSizeWidth" form:"masterSizeWidth"` // 主图宽度(可选)
	MasterSizeHeight int `json:"masterSizeHeight" form:"masterSizeHeight"` // 主图高度(可选)
	MasterImageAlpha float32 `json:"masterImageAlpha" form:"masterImageAlpha"` // 主图透明度(可选)
	RangeGraphSizeMin int `json:"rangeGraphSizeMin" form:"rangeGraphSizeMin"` // 拼图块随机尺寸最小值(可选)
	RangeGraphSizeMax int `json:"rangeGraphSizeMax" form:"rangeGraphSizeMax"` // 拼图块随机尺寸最大值(可选)
	RangeGraphAnglePosMin int `json:"rangeGraphAnglePosMin" form:"rangeGraphAnglePosMin"` // 拼图块随机角度最小值(可选)
	RangeGraphAnglePosMax int `json:"rangeGraphAnglePosMax" form:"rangeGraphAnglePosMax"` // 拼图块随机角度最大值(可选)
	GenGraphNumber int `json:"genGraphNumber" form:"genGraphNumber"` // 拼图块个数(可选)
	EnableGraphVerticalRandom bool `json:"enableGraphVerticalRandom" form:"enableGraphVerticalRandom"` // 拼图块水平方向是否随机排序(可选)
}

type VerifyCaptchaSliderReq struct {
	Token string `json:"token" form:"token" binding:"required"` // 验证码令牌
	X     int    `json:"x" form:"x" binding:"required"`        // 滑动终点X坐标
	Y     int    `json:"y" form:"y" binding:"required"`        // 滑动终点Y坐标
}

type SaveCaptchaSliderConfigReq struct {
	MasterSizeWidth             int      `json:"masterSizeWidth" form:"masterSizeWidth"`               // 主图宽度(默认300)
	MasterSizeHeight            int      `json:"masterSizeHeight" form:"masterSizeHeight"`             // 主图高度(默认220)
	MasterImageAlpha           float32  `json:"masterImageAlpha" form:"masterImageAlpha"`             // 主图透明度
	RangeGraphSizeMin           int      `json:"rangeGraphSizeMin" form:"rangeGraphSizeMin"`           // 拼图块随机尺寸最小值
	RangeGraphSizeMax           int      `json:"rangeGraphSizeMax" form:"rangeGraphSizeMax"`           // 拼图块随机尺寸最大值
	RangeGraphAnglePosMin       int      `json:"rangeGraphAnglePosMin" form:"rangeGraphAnglePosMin"`   // 拼图块随机角度最小值
	RangeGraphAnglePosMax       int      `json:"rangeGraphAnglePosMax" form:"rangeGraphAnglePosMax"`   // 拼图块随机角度最大值
	GenGraphNumber              int      `json:"genGraphNumber" form:"genGraphNumber"`               // 拼图块个数
	EnableGraphVerticalRandom   bool     `json:"enableGraphVerticalRandom" form:"enableGraphVerticalRandom"` // 拼图块水平方向是否随机排序
	RangeDeadZoneDirections     []string `json:"rangeDeadZoneDirections" form:"rangeDeadZoneDirections"` // 贴图盲区方向(可选)
}

type PreviewCaptchaSliderReq struct {
	MasterSizeWidth             int      `json:"masterSizeWidth" form:"masterSizeWidth"`               // 主图宽度(可选)
	MasterSizeHeight            int      `json:"masterSizeHeight" form:"masterSizeHeight"`             // 主图高度(可选)
	MasterImageAlpha           float32  `json:"masterImageAlpha" form:"masterImageAlpha"`             // 主图透明度(可选)
	RangeGraphSizeMin           int      `json:"rangeGraphSizeMin" form:"rangeGraphSizeMin"`           // 拼图块随机尺寸最小值(可选)
	RangeGraphSizeMax           int      `json:"rangeGraphSizeMax" form:"rangeGraphSizeMax"`           // 拼图块随机尺寸最大值(可选)
	RangeGraphAnglePosMin       int      `json:"rangeGraphAnglePosMin" form:"rangeGraphAnglePosMin"`   // 拼图块随机角度最小值(可选)
	RangeGraphAnglePosMax       int      `json:"rangeGraphAnglePosMax" form:"rangeGraphAnglePosMax"`   // 拼图块随机角度最大值(可选)
	GenGraphNumber              int      `json:"genGraphNumber" form:"genGraphNumber"`               // 拼图块个数(可选)
	EnableGraphVerticalRandom   bool     `json:"enableGraphVerticalRandom" form:"enableGraphVerticalRandom"` // 拼图块水平方向是否随机排序(可选)
	RangeDeadZoneDirections     []string `json:"rangeDeadZoneDirections" form:"rangeDeadZoneDirections"` // 贴图盲区方向(可选)
}
