package request

type CreateEmailTaskReq struct {
	Name         string                 `json:"name" form:"name"`                         // 任务名称
	Description  string                 `json:"description" form:"description"`           // 任务描述
	CronExpr     string                 `json:"cronExpr" form:"cronExpr"`                 // Cron表达式
	TemplateCode string                 `json:"templateCode" form:"templateCode"`         // 使用的模板编码
	TargetType   string                 `json:"targetType" form:"targetType"`             // 目标类型: all_users/specific_users/user_group/custom
	TargetConfig map[string]interface{} `json:"targetConfig" form:"targetConfig"`         // 目标配置
	TemplateData map[string]interface{} `json:"templateData" form:"templateData"`         // 模板数据
	Status       string                 `json:"status" form:"status"`                     // 任务状态: active/inactive
}

type CreateEmailTaskReqRe struct {
	Id int64 `json:"id" form:"id"` // 任务ID
	NextRunTime string `json:"nextRunTime" form:"nextRunTime"` // 下次执行时间
}

type UpdateEmailTaskReq struct {
	Id           int                    `json:"id" form:"id"`
	Name         string                 `json:"name" form:"name"`                       // 任务名称
	Description  string                 `json:"description" form:"description"`         // 任务描述
	CronExpr     string                 `json:"cronExpr" form:"cronExpr"`               // Cron表达式
	TemplateCode string                 `json:"templateCode" form:"templateCode"`       // 使用的模板编码
	TargetType   string                 `json:"targetType" form:"targetType"`           // 目标类型
	TargetConfig map[string]interface{} `json:"targetConfig" form:"targetConfig"`       // 目标配置
	TemplateData map[string]interface{} `json:"templateData" form:"templateData"`       // 模板数据
	Status       string                 `json:"status" form:"status"`                   // 任务状态
}

type DeleteEmailTaskReq struct {
	Id int `json:"id" form:"id"`
}

type GetEmailTaskListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Status string `json:"status" form:"status"`
}

type StartEmailTaskReq struct {
	Id int `json:"id" form:"id"`
}

type StopEmailTaskReq struct {
	Id int `json:"id" form:"id"`
}

type RunEmailTaskReq struct {
	Id int `json:"id" form:"id"`
}
