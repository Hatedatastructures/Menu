package response

type GetEmailConfigRes struct {
	Enabled bool `json:"enabled"` // 启用邮件
	Host string `json:"host"` // SMTP服务器
	Port int `json:"port"` // 端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	FromName string `json:"fromName"` // 发件人名称
	Ssl bool `json:"ssl"` // SSL加密
}

type SaveEmailConfigRes struct {
	Version int `json:"version"` // 版本号
}

type TestEmailRes struct {
	Success bool `json:"success"` // 是否成功
	Message string `json:"message"` // 结果信息
}
