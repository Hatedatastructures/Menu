package response

type GenerateCaptchaClickRes struct {
	Token string `json:"token"` // 验证码令牌
	MasterImage string `json:"masterImage"` // 主图Base64
	ThumbImage string `json:"thumbImage"` // 缩略图Base64
	Width int `json:"width"` // 主图宽度
	Height int `json:"height"` // 主图高度
}

type SaveCaptchaClickConfigRes struct {
	Version int `json:"version"` // 版本号
}

type VerifyCaptchaClickRes struct {
	Success bool `json:"success"` // 是否验证成功
	Message string `json:"message"` // 提示信息
}

type GetCaptchaClickConfigRes struct {
	MasterSizeWidth int `json:"masterSizeWidth"` // 主图宽度
	MasterSizeHeight int `json:"masterSizeHeight"` // 主图高度
	MasterRangeLenMin int `json:"masterRangeLenMin"` // 随机内容长度最小值
	MasterRangeLenMax int `json:"masterRangeLenMax"` // 随机内容长度最大值
	MasterRangeAnglePosMin int `json:"masterRangeAnglePosMin"` // 随机角度最小值
	MasterRangeAnglePosMax int `json:"masterRangeAnglePosMax"` // 随机角度最大值
	MasterRangeSizeMin int `json:"masterRangeSizeMin"` // 随机内容大小最小值
	MasterRangeSizeMax int `json:"masterRangeSizeMax"` // 随机内容大小最大值
	MasterRangeColors []string `json:"masterRangeColors"` // 随机颜色
	MasterDisplayShadow bool `json:"masterDisplayShadow"` // 是否显示阴影
	MasterShadowColor string `json:"masterShadowColor"` // 阴影颜色
	MasterShadowPointX int `json:"masterShadowPointX"` // 阴影X偏移
	MasterShadowPointY int `json:"masterShadowPointY"` // 阴影Y偏移
	MasterImageAlpha float32 `json:"masterImageAlpha"` // 主图透明度
	MasterUseShapeOriginalColor bool `json:"masterUseShapeOriginalColor"` // 是否使用图形原始颜色
	ThumbSizeWidth int `json:"thumbSizeWidth"` // 缩略图宽度
	ThumbSizeHeight int `json:"thumbSizeHeight"` // 缩略图高度
	ThumbRangeVerifyLenMin int `json:"thumbRangeVerifyLenMin"` // 校验内容随机长度最小值
	ThumbRangeVerifyLenMax int `json:"thumbRangeVerifyLenMax"` // 校验内容随机长度最大值
	ThumbDisabledRangeVerifyLen bool `json:"thumbDisabledRangeVerifyLen"` // 禁用校验内容随机长度
	ThumbRangeThumbSizeMin int `json:"thumbRangeThumbSizeMin"` // 缩略内容随机大小最小值
	ThumbRangeThumbSizeMax int `json:"thumbRangeThumbSizeMax"` // 缩略内容随机大小最大值
	ThumbRangeThumbColors []string `json:"thumbRangeThumbColors"` // 缩略随机颜色范围
	ThumbRangeThumbBgColors []string `json:"thumbRangeThumbBgColors"` // 缩略随机背景颜色范围
	ThumbIsThumbNonDeformAbility bool `json:"thumbIsThumbNonDeformAbility"` // 缩略图内容不变形
	ThumbBgDistort int `json:"thumbBgDistort"` // 缩略图背景扭曲
	ThumbBgCirclesNum int `json:"thumbBgCirclesNum"` // 缩略图绘制小圆点数量
	ThumbBgSlimLineNum int `json:"thumbBgSlimLineNum"` // 缩略图绘制线条数量
	Version int `json:"version"` // 版本号
}

type PreviewCaptchaClickRes struct {
	MasterImage string `json:"masterImage"` // 主图Base64
	ThumbImage string `json:"thumbImage"` // 缩略图Base64
	Dots []PreviewCaptchaClickResDot `json:"dots"` // []
}

type PreviewCaptchaClickResDot struct {
	X int `json:"x"` // X坐标
	Y int `json:"y"` // Y坐标
	W int `json:"w"` // 宽度
	H int `json:"h"` // 高度
	Content string `json:"content"` // 内容
}
