package response

type GenerateCaptchaSliderRes struct {
	Token string `json:"token"`       // 验证码令牌
	MasterImage string `json:"masterImage"` // 主图Base64(JPEG格式)
	TileImage string `json:"tileImage"`     // 拼图块Base64(PNG格式)
	X int `json:"x"` // 目标X坐标（调试用）
	Y int `json:"y"` // 目标Y坐标（调试用）
	TileWidth int `json:"tileWidth"` // 拼图块实际宽度
	TileHeight int `json:"tileHeight"` // 拼图块实际高度
	Width int `json:"width"` // 主图宽度
	Height int `json:"height"` // 主图高度
}

type VerifyCaptchaSliderRes struct {
	Success bool   `json:"success"` // 是否验证成功
	Message string `json:"message"` // 提示信息
}

type SaveCaptchaSliderConfigRes struct {
	Version int `json:"version"` // 版本号
}

type GetCaptchaSliderConfigRes struct {
	MasterSizeWidth int `json:"masterSizeWidth"` // 主图宽度
	MasterSizeHeight int `json:"masterSizeHeight"` // 主图高度
	MasterImageAlpha float32 `json:"masterImageAlpha"` // 主图透明度
	RangeGraphSizeMin int `json:"rangeGraphSizeMin"` // 拼图块随机尺寸最小值
	RangeGraphSizeMax int `json:"rangeGraphSizeMax"` // 拼图块随机尺寸最大值
	RangeGraphAnglePosMin int `json:"rangeGraphAnglePosMin"` // 拼图块随机角度最小值
	RangeGraphAnglePosMax int `json:"rangeGraphAnglePosMax"` // 拼图块随机角度最大值
	GenGraphNumber int `json:"genGraphNumber"` // 拼图块个数
	EnableGraphVerticalRandom bool `json:"enableGraphVerticalRandom"` // 拼图块水平方向是否随机排序
	RangeDeadZoneDirections []string `json:"rangeDeadZoneDirections"` // 贴图盲区方向
	Version int `json:"version"` // 版本号
}

type PreviewCaptchaSliderRes struct {
	MasterImage string `json:"masterImage"` // 主图Base64
	TileImage string `json:"tileImage"` // 拼图块Base64
	X int `json:"x"` // 目标X坐标
	Y int `json:"y"` // 目标Y坐标
	TileWidth int `json:"tileWidth"` // 拼图块实际宽度
	TileHeight int `json:"tileHeight"` // 拼图块实际高度
}
