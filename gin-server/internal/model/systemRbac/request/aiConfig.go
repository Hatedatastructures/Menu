package request

type CreateAiConfigReq struct {
	Name string `json:"name" form:"name"` // 配置名称
	Provider string `json:"provider" form:"provider"` // 提供者
	BaseUrl string `json:"baseUrl" form:"baseUrl"` // 接口地址
	ApiKey string `json:"apiKey" form:"apiKey"` // 密钥
	Model string `json:"model" form:"model"` // 模型名称
	Type string `json:"type" form:"type"` // 类型(openai/兼容openai)
	Enabled bool `json:"enabled" form:"enabled"` // 是否启用
}

type GetAiConfigListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Name string `json:"name" form:"name"`
	Enabled bool `json:"enabled" form:"enabled"`
}

type GetAiConfigReq struct {
	Id uint `json:"id" form:"id"`
}

type UpdateAiConfigReq struct {
	Id uint `json:"id" form:"id"`
	Name string `json:"name" form:"name"` // 配置名称
	Provider string `json:"provider" form:"provider"` // 提供者
	BaseUrl string `json:"baseUrl" form:"baseUrl"` // 接口地址
	ApiKey string `json:"apiKey" form:"apiKey"` // 密钥
	Model string `json:"model" form:"model"` // 模型名称
	Type string `json:"type" form:"type"` // 类型(openai/兼容openai)
	Enabled bool `json:"enabled" form:"enabled"` // 是否启用
}

type DeleteAiConfigReq struct {
	Id uint `json:"id" form:"id"`
}

type BatchDeleteAiConfigReq struct {
	Ids []uint `json:"ids" form:"ids"` // 配置ID数组
}

type TestAiConfigReq struct {
	Id uint `json:"id" form:"id"` // 配置ID
	Prompt string `json:"prompt" form:"prompt"` // 测试提示词
}

type AiChatReq struct {
	ConfigId uint `json:"configId" form:"configId"` // 配置ID
	Prompt string `json:"prompt" form:"prompt"` // 用户输入
}

type ToggleAiConfigStatusReq struct {
	Id uint `json:"id" form:"id"`
	Enabled bool `json:"enabled" form:"enabled"` // 启用状态
}
