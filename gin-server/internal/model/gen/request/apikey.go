package request

type SaveApiKeyReq struct {
	ApiKey string `json:"apiKey" form:"apiKey"` // DeepSeek API Key
	BaseUrl string `json:"baseUrl" form:"baseUrl"` // API地址(不可选,默认https:
	Model string `json:"model" form:"model"` // 模型名称(不可选,默认deepseek-chat)
}

type TestApiKeyReq struct {
	ApiKey string `json:"apiKey" form:"apiKey"` // DeepSeek API Key
	BaseUrl string `json:"baseUrl" form:"baseUrl"` // API地址
	Model string `json:"model" form:"model"` // 模型名称
}
