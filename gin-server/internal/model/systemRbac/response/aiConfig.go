package response

type CreateAiConfigRes struct {
	Id uint `json:"id"` // 配置ID
}

type GetAiConfigListRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Name string `json:"name"`
	Enabled bool `json:"enabled"`
	List []GetAiConfigListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetAiConfigListResList struct {
	Id uint `json:"id"` // 配置ID
	Name string `json:"name"` // 配置名称
	Provider string `json:"provider"` // 提供者
	BaseUrl string `json:"baseUrl"` // 接口地址
	ApiKey string `json:"apiKey"` // 密钥(脱敏)
	Model string `json:"model"` // 模型名称
	Type string `json:"type"` // 类型
	Enabled bool `json:"enabled"` // 是否启用
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetAiConfigRes struct {
	Id uint `json:"id"` // 配置ID
	Name string `json:"name"` // 配置名称
	Provider string `json:"provider"` // 提供者
	BaseUrl string `json:"baseUrl"` // 接口地址
	ApiKey string `json:"apiKey"` // 密钥
	Model string `json:"model"` // 模型名称
	Type string `json:"type"` // 类型
	Enabled bool `json:"enabled"` // 是否启用
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type TestAiConfigRes struct {
	Success bool `json:"success"` // 是否成功
	Response string `json:"response"` // AI响应内容
	Error string `json:"error"` // 错误信息
	Duration int64 `json:"duration"` // 耗时(毫秒)
}

type AiChatRes struct {
	Response string `json:"response"` // AI响应内容
	Model string `json:"model"` // 使用的模型
	Duration int64 `json:"duration"` // 耗时(毫秒)
}

type GetEnabledAiConfigsRes struct {
	List []GetEnabledAiConfigsResList `json:"list"` // []
}

type GetEnabledAiConfigsResList struct {
	Id uint `json:"id"` // 配置ID
	Name string `json:"name"` // 配置名称
	Provider string `json:"provider"` // 提供者
	Model string `json:"model"` // 模型名称
	Type string `json:"type"` // 类型
}
