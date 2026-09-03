package systemRbac

import (
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailLimitService struct{}

// 获取限流日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailLimitService) GetEmailLimitLog(
	c *gin.Context,
	r req.GetEmailLimitLogReq,
) (rs res.GetEmailLimitLogRes, err error) {
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	db := global.GVA_DB.Model(&systemRbac.EmailLimitLog{})

	if r.Email != "" {
		db = db.Where("email = ?", r.Email)
	}
	if r.Ip != "" {
		db = db.Where("ip = ?", r.Ip)
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

	var logs []systemRbac.EmailLimitLog
	offset := (page - 1) * pageSize
	db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs)

	list := make([]res.GetEmailLimitLogResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetEmailLimitLogResList{
			Id:        int64(log.ID),
			Email:     log.Email,
			Ip:        log.IP,
			LimitType: log.LimitType,
			Action:    log.Action,
			Blocked:   log.Blocked,
			Reason:    log.Reason,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetEmailLimitLogRes{
		Page:      page,
		PageSize:  pageSize,
		Email:     r.Email,
		Ip:        r.Ip,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Total:     total,
		List:      list,
	}

	return rs, nil
}

// 添加黑名单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailLimitService) AddEmailBlacklist(
	c *gin.Context,
	r req.AddEmailBlacklistReq,
) (err error) {
	blacklist := systemRbac.EmailBlacklist{
		Email:  r.Email,
		Reason: r.Reason,
		Type:   r.Type,
		Status: "active",
	}

	err = global.GVA_DB.Create(&blacklist).Error
	if err != nil {
		global.GVA_LOG.Error("添加黑名单失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}

	return nil
}

// 删除黑名单-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailLimitService) DeleteEmailBlacklist(
	c *gin.Context,
	r req.DeleteEmailBlacklistReq,
) (err error) {
	err = global.GVA_DB.Delete(&systemRbac.EmailBlacklist{}, r.Id).Error
	if err != nil {
		global.GVA_LOG.Error("删除黑名单失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR)
	}
	return nil
}

// 获取黑名单列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailLimitService) GetEmailBlacklistList(
	c *gin.Context,
	r req.GetEmailBlacklistListReq,
) (rs res.GetEmailBlacklistListRes, err error) {
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	global.GVA_DB.Model(&systemRbac.EmailBlacklist{}).Count(&total)

	var blacklists []systemRbac.EmailBlacklist
	offset := (page - 1) * pageSize
	global.GVA_DB.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&blacklists)

	list := make([]res.GetEmailBlacklistListResList, 0, len(blacklists))
	for _, item := range blacklists {
		list = append(list, res.GetEmailBlacklistListResList{
			Id:        int64(item.ID),
			Email:     item.Email,
			Reason:    item.Reason,
			Type:      item.Type,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: "",
		})
	}

	rs = res.GetEmailBlacklistListRes{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}

	return rs, nil
}

