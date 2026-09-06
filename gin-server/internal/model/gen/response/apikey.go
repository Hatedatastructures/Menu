package response

type GetMyApiKeyRes struct {
	Id uint `json:"id"` // ID
	ApiKey string `json:"apiKey"` // API Key(脱敏)
	BaseUrl string `json:"baseUrl"` // API地址
	Model string `json:"model"` // 模型名称
	Enabled bool `json:"enabled"` // 是否启用
	CreatedAt string `json:"createdAt"` // 创建时间
}

type SaveApiKeyRes struct {
	Id uint `json:"id"` // ID
}

type TestApiKeyRes struct {
	Valid bool `json:"valid"` // 是否有效
	Message string `json:"message"` // 提示信息
}
