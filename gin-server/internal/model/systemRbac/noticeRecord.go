package systemRbac

import "shack/internal/global"

// NoticeRecord 通知阅读记录
type NoticeRecord struct {
	global.GVA_MODEL

	NoticeId   int64 `json:"noticeId" gorm:"column:notice_id;index;not null;comment:通知ID"`
	UserId     int64 `json:"userId" gorm:"column:user_id;index;not null;comment:用户ID"`
	IsRead     bool  `json:"isRead" gorm:"column:is_read;default:false;comment:是否已读"`
	ReadTime   *int64 `json:"readTime" gorm:"column:read_time;comment:阅读时间"`
	Channel    string `json:"channel" gorm:"type:varchar(50);comment:推送渠道"`
	SendStatus int    `json:"sendStatus" gorm:"column:send_status;type:int;default:0;comment:发送状态(0失败1成功)"`
	ErrorMsg   string `json:"errorMsg" gorm:"column:error_msg;type:text;comment:错误信息"`
}

func (NoticeRecord) TableName() string {
	return "sys_notice_records"
}
