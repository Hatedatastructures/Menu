package request

type SaveEmailConfigReq struct {
	Enabled bool `json:"enabled" form:"enabled"` // 启用邮件
	Host string `json:"host" form:"host"` // SMTP服务器
	Port int `json:"port" form:"port"` // 端口
	Username string `json:"username" form:"username"` // 用户名
	Password string `json:"password" form:"password"` // 密码
	FromName string `json:"fromName" form:"fromName"` // 发件人名称
	Ssl bool `json:"ssl" form:"ssl"` // SSL加密
}

type TestEmailReq struct {
	Address string `json:"address" form:"address"` // 收件人邮箱
	TemplateId uint `json:"templateId" form:"templateId"` // 模板ID，可选，不传则使用默认测试模板
}
