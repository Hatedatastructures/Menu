package systemRbac

import (
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
)

type EmailLogService struct{}

// 获取邮件日志列表-后台使用
func (s *EmailLogService) GetEmailLogList(
	c *gin.Context,
	r req.GetEmailLogListReq,
) (rs res.GetEmailLogListRes, err error) {
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := global.GVA_DB.Model(&systemRbac.EmailLog{})

	if r.ToEmail != "" {
		db = db.Where("to_email LIKE ?", "%"+r.ToEmail+"%")
	}
	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}
	if r.BizType != "" {
		db = db.Where("biz_type = ?", r.BizType)
	}
	if r.StartTime != "" {
		startTime, _ := time.Parse("2006-01-02", r.StartTime)
		db = db.Where("created_at >= ?", startTime)
	}
	if r.EndTime != "" {
		endTime, _ := time.Parse("2006-01-02", r.EndTime)
		endTime = endTime.Add(24 * time.Hour)
		db = db.Where("created_at < ?", endTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var logs []systemRbac.EmailLog
	offset := (page - 1) * pageSize
	err = db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	list := make([]res.GetEmailLogListResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetEmailLogListResList{
			Id:           int64(log.ID),
			ToEmail:      log.ToEmail,
			Subject:      log.Subject,
			TemplateCode: log.TemplateCode,
			Status:       log.Status,
			BizType:      log.BizType,
			SentAt:       log.SentAt.Format("2006-01-02 15:04:05"),
			DeliveredAt:  formatTimePtr(log.DeliveredAt),
			ErrorMsg:     log.ErrorMsg,
			RetryCount:   log.RetryCount,
		})
	}

	rs = res.GetEmailLogListRes{
		Page:      page,
		PageSize:  pageSize,
		ToEmail:   r.ToEmail,
		Status:    r.Status,
		BizType:   r.BizType,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Total:     total,
		List:      list,
	}

	return rs, nil
}

// 获取邮件日志详情-后台使用
func (s *EmailLogService) GetEmailLog(
	c *gin.Context,
	r req.GetEmailLogReq,
) (rs res.GetEmailLogRes, err error) {
	var log systemRbac.EmailLog
	err = global.GVA_DB.First(&log, r.Id).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	rs = res.GetEmailLogRes{
		Id:           int64(log.ID),
		ToEmail:      log.ToEmail,
		Subject:      log.Subject,
		Content:      log.Content,
		TemplateCode: log.TemplateCode,
		Status:       log.Status,
		BizType:      log.BizType,
		SentAt:       formatTimePtr(log.SentAt),
		DeliveredAt:  formatTimePtr(log.DeliveredAt),
		OpenedAt:     formatTimePtr(log.OpenedAt),
		ClickedAt:    formatTimePtr(log.ClickedAt),
		ErrorMsg:     log.ErrorMsg,
		RetryCount:   log.RetryCount,
		Provider:     log.Provider,
	}

	return rs, nil
}

// 重发邮件-后台使用
func (s *EmailLogService) ResendEmail(
	c *gin.Context,
	r req.ResendEmailReq,
) (rs res.ResendEmailRes, err error) {
	var log systemRbac.EmailLog
	err = global.GVA_DB.First(&log, r.Id).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	// 创建新的邮件任务
	newLog := systemRbac.EmailLog{
		ToEmail:      log.ToEmail,
		Subject:      log.Subject,
		Content:      log.Content,
		TemplateCode: log.TemplateCode,
		Status:       "pending",
		RetryCount:   log.RetryCount + 1,
	}
	global.GVA_DB.Create(&newLog)

	rs = res.ResendEmailRes{
		NewTaskId: fmt.Sprintf("email_%d", newLog.ID),
	}

	return rs, nil
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
