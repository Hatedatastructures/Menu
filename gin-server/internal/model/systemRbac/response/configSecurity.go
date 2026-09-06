package response

type GetSecurityConfigRes struct {
	EncryptEnabled bool `json:"encryptEnabled"` // 接口加密
	EncryptScope string `json:"encryptScope"` // 加密范围
	EncryptPublicKey string `json:"encryptPublicKey"` // RSA公钥
	EncryptPrivateKey string `json:"encryptPrivateKey"` // RSA私钥
	DisableDevtool bool `json:"disableDevtool"` // 禁止前端调试
	TokenName string `json:"tokenName"` // Token名称
	TokenTimeout int `json:"tokenTimeout"` // Token有效期(秒)
	TokenActiveTimeout int `json:"tokenActiveTimeout"` // 活跃超时时间(秒)
	TokenIsConcurrent bool `json:"tokenIsConcurrent"` // 允许多端登录
	TokenIsShare bool `json:"tokenIsShare"` // 共用Token
	TokenStyle string `json:"tokenStyle"` // Token风格
	TokenIsReadHeader bool `json:"tokenIsReadHeader"` // 从Header读取
	TokenIsLog bool `json:"tokenIsLog"` // 输出操作日志
	TokenIsPrint bool `json:"tokenIsPrint"` // 打印版本信息
}

type SaveSecurityConfigRes struct {
	Version int `json:"version"` // 版本号
}

type GenerateRSAKeysRes struct {
	PublicKey string `json:"publicKey"` // 公钥
	PrivateKey string `json:"privateKey"` // 私钥
}
