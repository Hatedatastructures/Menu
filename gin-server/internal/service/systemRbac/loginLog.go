package systemRbac

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	biz_err "shack/internal/error"
)



// LoginLogService 登录日志服务
type LoginLogService struct{}

// convertToListRes 转换为列表响应(通用)
func (s *LoginLogService) convertToListRes(logs []systemRbac.LoginLog) []res.GetLoginLogListResList {
	list := make([]res.GetLoginLogListResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetLoginLogListResList{
			Id:       int64(log.ID),
			Username: log.Username,
			Ip:       log.IP,
			Location: log.Location,
			Browser:  log.Browser,
			Os:       log.OS,
			Status:   log.Status,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}
	return list
}

// convertToUsernameRes 转换为按用户名查询的响应
func (s *LoginLogService) convertToUsernameRes(logs []systemRbac.LoginLog) []res.GetLoginLogByUsernameResList {
	list := make([]res.GetLoginLogByUsernameResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetLoginLogByUsernameResList{
			Id:       int64(log.ID),
			Username: log.Username,
			Ip:       log.IP,
			Location: log.Location,
			Browser:  log.Browser,
			Os:       log.OS,
			Status:   log.Status,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}
	return list
}

// convertToStatusRes 转换为按状态查询的响应
func (s *LoginLogService) convertToStatusRes(logs []systemRbac.LoginLog) []res.GetLoginLogByStatusResList {
	list := make([]res.GetLoginLogByStatusResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetLoginLogByStatusResList{
			Id:       int64(log.ID),
			Username: log.Username,
			Ip:       log.IP,
			Location: log.Location,
			Browser:  log.Browser,
			Os:       log.OS,
			Status:   log.Status,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}
	return list
}

// convertToTimeRangeRes 转换为按时间范围查询的响应
func (s *LoginLogService) convertToTimeRangeRes(logs []systemRbac.LoginLog) []res.GetLoginLogByTimeRangeResList {
	list := make([]res.GetLoginLogByTimeRangeResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetLoginLogByTimeRangeResList{
			Id:       int64(log.ID),
			Username: log.Username,
			Ip:       log.IP,
			Location: log.Location,
			Browser:  log.Browser,
			Os:       log.OS,
			Status:   log.Status,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}
	return list
}

// convertToRecentRes 转换为最近记录响应
func (s *LoginLogService) convertToRecentRes(logs []systemRbac.LoginLog) []res.GetRecentLoginLogResList {
	list := make([]res.GetRecentLoginLogResList, 0, len(logs))
	for _, log := range logs {
		list = append(list, res.GetRecentLoginLogResList{
			Id:       int64(log.ID),
			Username: log.Username,
			Ip:       log.IP,
			Location: log.Location,
			Browser:  log.Browser,
			Os:       log.OS,
			Status:   log.Status,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}
	return list
}

// GetLoginLogList 获取登录日志列表-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetLoginLogList(
	ctx *gin.Context,
	r req.GetLoginLogListReq,
) (rs res.GetLoginLogListRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var logs []systemRbac.LoginLog

	db := global.GVA_DB.Model(&systemRbac.LoginLog{})

	// 构建查询条件
	if r.Username != "" {
		db = db.Where("username LIKE ?", "%"+r.Username+"%")
	}
	if r.Status != "" {
		db = db.Where("status = ?", r.Status)
	}
	if r.StartTime != "" && r.EndTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", r.StartTime)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", r.EndTime)
			if err == nil {
				db = db.Where("login_time BETWEEN ? AND ?", startTime, endTime)
			}
		}
	}

	db.Count(&total)
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录日志失败")
	}

	rs = res.GetLoginLogListRes{
		Page:      r.Page,
		Size:      r.Size,
		Username:  r.Username,
		Status:    r.Status,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		List:      s.convertToListRes(logs),
		Total:     total,
	}
	return rs, nil
}

// GetLoginLogByUsername 根据用户名查询登录日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetLoginLogByUsername(
	ctx *gin.Context,
	r req.GetLoginLogByUsernameReq,
) (rs res.GetLoginLogByUsernameRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var logs []systemRbac.LoginLog

	db := global.GVA_DB.Model(&systemRbac.LoginLog{}).Where("username = ?", r.Username)
	db.Count(&total)
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录日志失败")
	}

	rs = res.GetLoginLogByUsernameRes{
		Page:  r.Page,
		Size:  r.Size,
		List:  s.convertToUsernameRes(logs),
		Total: total,
	}
	return rs, nil
}

// GetLoginLogByStatus 根据状态查询登录日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetLoginLogByStatus(
	ctx *gin.Context,
	r req.GetLoginLogByStatusReq,
) (rs res.GetLoginLogByStatusRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var logs []systemRbac.LoginLog

	db := global.GVA_DB.Model(&systemRbac.LoginLog{}).Where("status = ?", r.Status)
	db.Count(&total)
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录日志失败")
	}

	rs = res.GetLoginLogByStatusRes{
		Page:  r.Page,
		Size:  r.Size,
		List:  s.convertToStatusRes(logs),
		Total: total,
	}
	return rs, nil
}

// GetLoginLogByTimeRange 按时间范围查询登录日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetLoginLogByTimeRange(
	ctx *gin.Context,
	r req.GetLoginLogByTimeRangeReq,
) (rs res.GetLoginLogByTimeRangeRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var logs []systemRbac.LoginLog

	db := global.GVA_DB.Model(&systemRbac.LoginLog{})

	// 只有当时间都不为空时才添加时间条件
	if r.StartTime != "" && r.EndTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", r.StartTime)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", r.EndTime)
			if err == nil {
				db = db.Where("login_time BETWEEN ? AND ?", startTime, endTime)
			}
		}
	}

	db.Count(&total)
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录日志失败")
	}

	rs = res.GetLoginLogByTimeRangeRes{
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Page:      r.Page,
		Size:      r.Size,
		List:      s.convertToTimeRangeRes(logs),
		Total:     total,
	}
	return rs, nil
}

// GetRecentLoginLog 获取最近N条登录记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetRecentLoginLog(
	ctx *gin.Context,
	r req.GetRecentLoginLogReq,
) (rs res.GetRecentLoginLogRes, err error) {
	if r.Limit <= 0 {
		r.Limit = 10
	}

	var logs []systemRbac.LoginLog
	err = global.GVA_DB.Order("id desc").Limit(r.Limit).Find(&logs).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录日志失败")
	}

	rs = res.GetRecentLoginLogRes{
		Limit: r.Limit,
		List:  s.convertToRecentRes(logs),
	}
	return rs, nil
}

// GetLoginStatistics 获取登录统计信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) GetLoginStatistics(
	ctx *gin.Context,
) (rs res.GetLoginStatisticsRes, err error) {
	// 总数
	var total int64
	global.GVA_DB.Model(&systemRbac.LoginLog{}).Count(&total)

	// 今日数量
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	global.GVA_DB.Model(&systemRbac.LoginLog{}).Where("DATE(login_time) = ?", today).Count(&todayCount)

	// 成功数
	var successCount int64
	global.GVA_DB.Model(&systemRbac.LoginLog{}).Where("status = ?", "成功").Count(&successCount)

	// 失败数
	var failCount int64
	global.GVA_DB.Model(&systemRbac.LoginLog{}).Where("status = ?", "失败").Count(&failCount)

	rs = res.GetLoginStatisticsRes{
		Total:       total,
		TodayCount:  todayCount,
		SuccessCount: successCount,
		FailCount:    failCount,
	}
	return rs, nil
}

// ClearLoginLog 清空登录日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) ClearLoginLog(
	ctx *gin.Context,
) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&[]systemRbac.LoginLog{}).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "清空登录日志失败")
	}
	return nil
}

// DeleteLoginLog 批量删除登录日志-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 19:30:30
func (s *LoginLogService) DeleteLoginLog(
	ctx *gin.Context,
	r req.DeleteLoginLogReq,
) (err error) {
	if r.Ids == "" {
		return biz_err.New(biz_err.PARAM_ERROR, "删除的ID不能为空")
	}

	// 解析ID数组
	idStrs := strings.Split(r.Ids, ",")
	if len(idStrs) == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "删除的ID格式不正确")
	}

	// 转换为int64数组
	var ids []int64
	for _, idStr := range idStrs {
		var id int64
		if _, err := fmt.Sscanf(strings.TrimSpace(idStr), "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "删除的ID格式不正确")
	}

	err = global.GVA_DB.Where("id IN ?", ids).Delete(&systemRbac.LoginLog{}).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除登录日志失败")
	}
	return nil
}