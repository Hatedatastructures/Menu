package systemRbac

import "shack/internal/global"

// Notice 通知实体
type Notice struct {
	global.GVA_MODEL

	Title       string `json:"title" gorm:"type:varchar(200);not null;comment:通知标题"`
	Content     string `json:"content" gorm:"type:text;not null;comment:通知内容"`
	NoticeType  int    `json:"noticeType" gorm:"type:int;not null;default:1;comment:通知类型(1通知2公告)"`
	Channels    string `json:"channels" gorm:"type:varchar(500);comment:推送渠道(JSON数组)"`
	TargetType  int    `json:"targetType" gorm:"type:int;not null;default:3;comment:推送对象类型(1指定用户2按部门3全部)"`
	TargetIds   string `json:"targetIds" gorm:"type:text;comment:推送对象ID(JSON数组)"`
	Status      int    `json:"status" gorm:"type:int;not null;default:0;comment:状态(0草稿1已发布2已撤回)"`
	CreateBy    int64  `json:"createBy" gorm:"column:create_by;comment:创建者ID"`
	CreateName  string `json:"createName" gorm:"column:create_name;type:varchar(100);comment:创建者名称"`
	PublishTime *int64 `json:"publishTime" gorm:"column:publish_time;comment:发布时间"`
}

func (Notice) TableName() string {
	return "sys_notices"
}
