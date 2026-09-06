package systemRbac

import (
	"shack/internal/global"
	"time"
)

// EmailTask 邮件定时任务
type EmailTask struct {
	global.GVA_MODEL
	Name         string          `json:"name" gorm:"type:varchar(100);not null;comment:任务名称"`
	Description  string          `json:"description" gorm:"type:text;comment:任务描述"`
	CronExpr     string          `json:"cronExpr" gorm:"type:varchar(100);not null;comment:Cron表达式"`
	TemplateCode string          `json:"templateCode" gorm:"type:varchar(50);not null;comment:模板编码"`
	TargetType   string          `json:"targetType" gorm:"type:varchar(20);not null;default:all_users;comment:目标类型"`
	TargetConfig string          `json:"targetConfig" gorm:"type:json;comment:目标配置JSON"`
	TemplateData string          `json:"templateData" gorm:"type:json;comment:模板数据JSON"`
	Status       string          `json:"status" gorm:"type:varchar(20);not null;default:active;comment:任务状态"`
	NextRunTime  *time.Time      `json:"nextRunTime" gorm:"column:next_run_time;comment:下次执行时间"`
	LastRunTime  *time.Time      `json:"lastRunTime" gorm:"column:last_run_time;comment:上次执行时间"`
	LastRunStatus string         `json:"lastRunStatus" gorm:"type:varchar(20);column:last_run_status;comment:上次执行状态"`
	TotalRuns    int             `json:"totalRuns" gorm:"column:total_runs;default:0;comment:总执行次数"`
	SuccessRuns  int             `json:"successRuns" gorm:"column:success_runs;default:0;comment:成功次数"`
	FailedRuns   int             `json:"failedRuns" gorm:"column:failed_runs;default:0;comment:失败次数"`
	CreatedBy    uint            `json:"createdBy" gorm:"column:created_by;comment:创建人ID"`
	UpdatedBy    uint            `json:"updatedBy" gorm:"column:updated_by;comment:更新人ID"`
}

// TableName 指定表名
func (EmailTask) TableName() string {
	return "sys_email_tasks"
}
