package request

type GenerateCaptchaClickReq struct {
	MasterSizeWidth int `json:"masterSizeWidth" form:"masterSizeWidth"` // 主图宽度(可选)
	MasterSizeHeight int `json:"masterSizeHeight" form:"masterSizeHeight"` // 主图高度(可选)
	MasterRangeLenMin int `json:"masterRangeLenMin" form:"masterRangeLenMin"` // 随机内容长度最小值(可选)
	MasterRangeLenMax int `json:"masterRangeLenMax" form:"masterRangeLenMax"` // 随机内容长度最大值(可选)
	MasterRangeAnglePosMin int `json:"masterRangeAnglePosMin" form:"masterRangeAnglePosMin"` // 随机角度最小值(可选)
	MasterRangeAnglePosMax int `json:"masterRangeAnglePosMax" form:"masterRangeAnglePosMax"` // 随机角度最大值(可选)
	MasterRangeSizeMin int `json:"masterRangeSizeMin" form:"masterRangeSizeMin"` // 随机内容大小最小值(可选)
	MasterRangeSizeMax int `json:"masterRangeSizeMax" form:"masterRangeSizeMax"` // 随机内容大小最大值(可选)
	MasterRangeColors []string `json:"masterRangeColors" form:"masterRangeColors"` // 随机颜色(可选)
	MasterDisplayShadow bool `json:"masterDisplayShadow" form:"masterDisplayShadow"` // 是否显示阴影(可选)
	MasterShadowColor string `json:"masterShadowColor" form:"masterShadowColor"` // 阴影颜色(可选)
	MasterShadowPointX int `json:"masterShadowPointX" form:"masterShadowPointX"` // 阴影X偏移(可选)
	MasterShadowPointY int `json:"masterShadowPointY" form:"masterShadowPointY"` // 阴影Y偏移(可选)
	MasterImageAlpha float32 `json:"masterImageAlpha" form:"masterImageAlpha"` // 主图透明度(可选)
	MasterUseShapeOriginalColor bool `json:"masterUseShapeOriginalColor" form:"masterUseShapeOriginalColor"` // 是否使用图形原始颜色(可选)
	ThumbSizeWidth int `json:"thumbSizeWidth" form:"thumbSizeWidth"` // 缩略图宽度(可选)
	ThumbSizeHeight int `json:"thumbSizeHeight" form:"thumbSizeHeight"` // 缩略图高度(可选)
	ThumbRangeVerifyLenMin int `json:"thumbRangeVerifyLenMin" form:"thumbRangeVerifyLenMin"` // 校验内容随机长度最小值(可选)
	ThumbRangeVerifyLenMax int `json:"thumbRangeVerifyLenMax" form:"thumbRangeVerifyLenMax"` // 校验内容随机长度最大值(可选)
	ThumbDisabledRangeVerifyLen bool `json:"thumbDisabledRangeVerifyLen" form:"thumbDisabledRangeVerifyLen"` // 禁用校验内容随机长度(可选)
	ThumbRangeThumbSizeMin int `json:"thumbRangeThumbSizeMin" form:"thumbRangeThumbSizeMin"` // 缩略内容随机大小最小值(可选)
	ThumbRangeThumbSizeMax int `json:"thumbRangeThumbSizeMax" form:"thumbRangeThumbSizeMax"` // 缩略内容随机大小最大值(可选)
	ThumbRangeThumbColors []string `json:"thumbRangeThumbColors" form:"thumbRangeThumbColors"` // 缩略随机颜色范围(可选)
	ThumbRangeThumbBgColors []string `json:"thumbRangeThumbBgColors" form:"thumbRangeThumbBgColors"` // 缩略随机背景颜色范围(可选)
	ThumbIsThumbNonDeformAbility bool `json:"thumbIsThumbNonDeformAbility" form:"thumbIsThumbNonDeformAbility"` // 缩略图内容不变形(可选)
	ThumbBgDistort int `json:"thumbBgDistort" form:"thumbBgDistort"` // 缩略图背景扭曲(可选)
	ThumbBgCirclesNum int `json:"thumbBgCirclesNum" form:"thumbBgCirclesNum"` // 缩略图绘制小圆点数量(可选)
	ThumbBgSlimLineNum int `json:"thumbBgSlimLineNum" form:"thumbBgSlimLineNum"` // 缩略图绘制线条数量(可选)
}

type VerifyCaptchaClickReq struct {
	Token string `json:"token" form:"token"` // 验证码令牌
	Points []VerifyCaptchaClickReqPoint `json:"points" form:"points"` // []
}

type VerifyCaptchaClickReqPoint struct {
	X int `json:"x" form:"x"` // X坐标
	Y int `json:"y" form:"y"` // Y坐标
}

type SaveCaptchaClickConfigReq struct {
	MasterSizeWidth int `json:"masterSizeWidth" form:"masterSizeWidth"` // 主图宽度(默认300)
	MasterSizeHeight int `json:"masterSizeHeight" form:"masterSizeHeight"` // 主图高度(默认220)
	MasterRangeLenMin int `json:"masterRangeLenMin" form:"masterRangeLenMin"` // 随机内容长度最小值
	MasterRangeLenMax int `json:"masterRangeLenMax" form:"masterRangeLenMax"` // 随机内容长度最大值
	MasterRangeAnglePosMin int `json:"masterRangeAnglePosMin" form:"masterRangeAnglePosMin"` // 随机角度最小值
	MasterRangeAnglePosMax int `json:"masterRangeAnglePosMax" form:"masterRangeAnglePosMax"` // 随机角度最大值
	MasterRangeSizeMin int `json:"masterRangeSizeMin" form:"masterRangeSizeMin"` // 随机内容大小最小值
	MasterRangeSizeMax int `json:"masterRangeSizeMax" form:"masterRangeSizeMax"` // 随机内容大小最大值
	MasterRangeColors []string `json:"masterRangeColors" form:"masterRangeColors"` // 随机颜色
	MasterDisplayShadow bool `json:"masterDisplayShadow" form:"masterDisplayShadow"` // 是否显示阴影
	MasterShadowColor string `json:"masterShadowColor" form:"masterShadowColor"` // 阴影颜色
	MasterShadowPointX int `json:"masterShadowPointX" form:"masterShadowPointX"` // 阴影X偏移
	MasterShadowPointY int `json:"masterShadowPointY" form:"masterShadowPointY"` // 阴影Y偏移
	MasterImageAlpha float32 `json:"masterImageAlpha" form:"masterImageAlpha"` // 主图透明度
	MasterUseShapeOriginalColor bool `json:"masterUseShapeOriginalColor" form:"masterUseShapeOriginalColor"` // 是否使用图形原始颜色
	ThumbSizeWidth int `json:"thumbSizeWidth" form:"thumbSizeWidth"` // 缩略图宽度(默认150)
	ThumbSizeHeight int `json:"thumbSizeHeight" form:"thumbSizeHeight"` // 缩略图高度(默认40)
	ThumbRangeVerifyLenMin int `json:"thumbRangeVerifyLenMin" form:"thumbRangeVerifyLenMin"` // 校验内容随机长度最小值
	ThumbRangeVerifyLenMax int `json:"thumbRangeVerifyLenMax" form:"thumbRangeVerifyLenMax"` // 校验内容随机长度最大值
	ThumbDisabledRangeVerifyLen bool `json:"thumbDisabledRangeVerifyLen" form:"thumbDisabledRangeVerifyLen"` // 禁用校验内容随机长度
	ThumbRangeThumbSizeMin int `json:"thumbRangeThumbSizeMin" form:"thumbRangeThumbSizeMin"` // 缩略内容随机大小最小值
	ThumbRangeThumbSizeMax int `json:"thumbRangeThumbSizeMax" form:"thumbRangeThumbSizeMax"` // 缩略内容随机大小最大值
	ThumbRangeThumbColors []string `json:"thumbRangeThumbColors" form:"thumbRangeThumbColors"` // 缩略随机颜色范围
	ThumbRangeThumbBgColors []string `json:"thumbRangeThumbBgColors" form:"thumbRangeThumbBgColors"` // 缩略随机背景颜色范围
	ThumbIsThumbNonDeformAbility bool `json:"thumbIsThumbNonDeformAbility" form:"thumbIsThumbNonDeformAbility"` // 缩略图内容不变形
	ThumbBgDistort int `json:"thumbBgDistort" form:"thumbBgDistort"` // 缩略图背景扭曲(1-5)
	ThumbBgCirclesNum int `json:"thumbBgCirclesNum" form:"thumbBgCirclesNum"` // 缩略图绘制小圆点数量
	ThumbBgSlimLineNum int `json:"thumbBgSlimLineNum" form:"thumbBgSlimLineNum"` // 缩略图绘制线条数量
}

type PreviewCaptchaClickReq struct {
	MasterSizeWidth int `json:"masterSizeWidth" form:"masterSizeWidth"` // 主图宽度(可选)
	MasterSizeHeight int `json:"masterSizeHeight" form:"masterSizeHeight"` // 主图高度(可选)
	MasterRangeLenMin int `json:"masterRangeLenMin" form:"masterRangeLenMin"` // 随机内容长度最小值(可选)
	MasterRangeLenMax int `json:"masterRangeLenMax" form:"masterRangeLenMax"` // 随机内容长度最大值(可选)
	MasterRangeAnglePosMin int `json:"masterRangeAnglePosMin" form:"masterRangeAnglePosMin"` // 随机角度最小值(可选)
	MasterRangeAnglePosMax int `json:"masterRangeAnglePosMax" form:"masterRangeAnglePosMax"` // 随机角度最大值(可选)
	MasterRangeSizeMin int `json:"masterRangeSizeMin" form:"masterRangeSizeMin"` // 随机内容大小最小值(可选)
	MasterRangeSizeMax int `json:"masterRangeSizeMax" form:"masterRangeSizeMax"` // 随机内容大小最大值(可选)
	MasterRangeColors []string `json:"masterRangeColors" form:"masterRangeColors"` // 随机颜色(可选)
	MasterDisplayShadow bool `json:"masterDisplayShadow" form:"masterDisplayShadow"` // 是否显示阴影(可选)
	MasterShadowColor string `json:"masterShadowColor" form:"masterShadowColor"` // 阴影颜色(可选)
	MasterShadowPointX int `json:"masterShadowPointX" form:"masterShadowPointX"` // 阴影X偏移(可选)
	MasterShadowPointY int `json:"masterShadowPointY" form:"masterShadowPointY"` // 阴影Y偏移(可选)
	MasterImageAlpha float32 `json:"masterImageAlpha" form:"masterImageAlpha"` // 主图透明度(可选)
	MasterUseShapeOriginalColor bool `json:"masterUseShapeOriginalColor" form:"masterUseShapeOriginalColor"` // 是否使用图形原始颜色(可选)
	ThumbSizeWidth int `json:"thumbSizeWidth" form:"thumbSizeWidth"` // 缩略图宽度(可选)
	ThumbSizeHeight int `json:"thumbSizeHeight" form:"thumbSizeHeight"` // 缩略图高度(可选)
	ThumbRangeVerifyLenMin int `json:"thumbRangeVerifyLenMin" form:"thumbRangeVerifyLenMin"` // 校验内容随机长度最小值(可选)
	ThumbRangeVerifyLenMax int `json:"thumbRangeVerifyLenMax" form:"thumbRangeVerifyLenMax"` // 校验内容随机长度最大值(可选)
	ThumbDisabledRangeVerifyLen bool `json:"thumbDisabledRangeVerifyLen" form:"thumbDisabledRangeVerifyLen"` // 禁用校验内容随机长度(可选)
	ThumbRangeThumbSizeMin int `json:"thumbRangeThumbSizeMin" form:"thumbRangeThumbSizeMin"` // 缩略内容随机大小最小值(可选)
	ThumbRangeThumbSizeMax int `json:"thumbRangeThumbSizeMax" form:"thumbRangeThumbSizeMax"` // 缩略内容随机大小最大值(可选)
	ThumbRangeThumbColors []string `json:"thumbRangeThumbColors" form:"thumbRangeThumbColors"` // 缩略随机颜色范围(可选)
	ThumbRangeThumbBgColors []string `json:"thumbRangeThumbBgColors" form:"thumbRangeThumbBgColors"` // 缩略随机背景颜色范围(可选)
	ThumbIsThumbNonDeformAbility bool `json:"thumbIsThumbNonDeformAbility" form:"thumbIsThumbNonDeformAbility"` // 缩略图内容不变形(可选)
	ThumbBgDistort int `json:"thumbBgDistort" form:"thumbBgDistort"` // 缩略图背景扭曲(可选)
	ThumbBgCirclesNum int `json:"thumbBgCirclesNum" form:"thumbBgCirclesNum"` // 缩略图绘制小圆点数量(可选)
	ThumbBgSlimLineNum int `json:"thumbBgSlimLineNum" form:"thumbBgSlimLineNum"` // 缩略图绘制线条数量(可选)
}
