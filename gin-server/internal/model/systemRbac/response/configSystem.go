package response

type GetSystemConfigRes struct {
	SiteName string `json:"siteName"` // 站点名称
	SiteDescription string `json:"siteDescription"` // 站点描述
	SiteLogo string `json:"siteLogo"` // 站点Logo
	Copyright string `json:"copyright"` // 版权信息
	Icp string `json:"icp"` // ICP备案号
	WatermarkEnabled bool `json:"watermarkEnabled"` // 启用水印
	WatermarkType string `json:"watermarkType"` // 水印类型
	WatermarkCustomText string `json:"watermarkCustomText"` // 自定义水印文本
	WatermarkOpacity float32 `json:"watermarkOpacity"` // 水印透明度
}

type SaveSystemConfigRes struct {
	Version int `json:"version"` // 版本号
}
