package request

type SaveSecurityConfigReq struct {
	EncryptEnabled bool `json:"encryptEnabled" form:"encryptEnabled"` // 接口加密
	EncryptScope string `json:"encryptScope" form:"encryptScope"` // 加密范围
	EncryptPublicKey string `json:"encryptPublicKey" form:"encryptPublicKey"` // RSA公钥
	EncryptPrivateKey string `json:"encryptPrivateKey" form:"encryptPrivateKey"` // RSA私钥
	DisableDevtool bool `json:"disableDevtool" form:"disableDevtool"` // 禁止前端调试
	TokenName string `json:"tokenName" form:"tokenName"` // Token名称
	TokenTimeout int `json:"tokenTimeout" form:"tokenTimeout"` // Token有效期(秒)
	TokenActiveTimeout int `json:"tokenActiveTimeout" form:"tokenActiveTimeout"` // 活跃超时时间(秒)
	TokenIsConcurrent bool `json:"tokenIsConcurrent" form:"tokenIsConcurrent"` // 允许多端登录
	TokenIsShare bool `json:"tokenIsShare" form:"tokenIsShare"` // 共用Token
	TokenStyle string `json:"tokenStyle" form:"tokenStyle"` // Token风格
	TokenIsReadHeader bool `json:"tokenIsReadHeader" form:"tokenIsReadHeader"` // 从Header读取
	TokenIsLog bool `json:"tokenIsLog" form:"tokenIsLog"` // 输出操作日志
	TokenIsPrint bool `json:"tokenIsPrint" form:"tokenIsPrint"` // 打印版本信息
}
