package request

type SaveSystemConfigReq struct {
	SiteName string `json:"siteName" form:"siteName"` // 站点名称
	SiteDescription string `json:"siteDescription" form:"siteDescription"` // 站点描述
	SiteLogo string `json:"siteLogo" form:"siteLogo"` // 站点Logo
	Copyright string `json:"copyright" form:"copyright"` // 版权信息
	Icp string `json:"icp" form:"icp"` // ICP备案号
	WatermarkEnabled bool `json:"watermarkEnabled" form:"watermarkEnabled"` // 启用水印
	WatermarkType string `json:"watermarkType" form:"watermarkType"` // 水印类型
	WatermarkCustomText string `json:"watermarkCustomText" form:"watermarkCustomText"` // 自定义水印文本
	WatermarkOpacity float32 `json:"watermarkOpacity" form:"watermarkOpacity"` // 水印透明度
}
