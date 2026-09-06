package systemRbac

import (
	"time"

	"github.com/gin-gonic/gin"
	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	biz_err "shack/internal/error"
)

var OprationServiceApp = new(OprationService)

type OprationService struct{}

// 清空所有操作记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) ClearOperationRecord(
	ctx *gin.Context,
) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&[]systemRbac.SysOperationRecord{}).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "清空操作记录失败")
	}
	return nil
}

// 获取操作统计信息-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetOperationStatistics(
	ctx *gin.Context,
) (rs res.GetOperationStatisticsRes, err error) {
	// 总数
	var total int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Count(&total)

	// 今日数量
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("DATE(created_at) = ?", today).Count(&todayCount)

	// 异常数量
	var errorCount int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("status != ?", 200).Count(&errorCount)

	// 平均延迟
	var avgLatency float64
	var records []systemRbac.SysOperationRecord
	global.GVA_DB.Select("latency").Find(&records)
	if len(records) > 0 {
		var totalLatency int64
		for _, r := range records {
			totalLatency += int64(r.Latency)
		}
		avgLatency = float64(totalLatency) / float64(len(records)) / 1e6 // 转换为ms
	}

	// 按请求方法统计
	var methodStats []res.GetOperationStatisticsResMethodstat
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Select("method, COUNT(*) as count").Group("method").Scan(&methodStats)

	// 按请求路径统计(取前10)
	var pathStats []res.GetOperationStatisticsResPathstat
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Select("path, COUNT(*) as count").Group("path").Order("count DESC").Limit(10).Scan(&pathStats)

	rs = res.GetOperationStatisticsRes{
		Total:       total,
		TodayCount:  todayCount,
		ErrorCount:  errorCount,
		AvgLatency:  avgLatency,
		MethodStats: methodStats,
		PathStats:   pathStats,
	}
	return rs, nil
}

// 根据用户ID获取操作记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetOperationRecordsByUserId(
	ctx *gin.Context,
	r req.GetOperationRecordsByUserIdReq,
) (rs res.GetOperationRecordsByUserIdRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var records []systemRbac.SysOperationRecord

	db := global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("user_id = ?", r.UserId)
	db.Count(&total)
	db = db.Order("id desc").Limit(limit).Offset(offset).Preload("User")
	err = db.Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetOperationRecordsByUserIdResList
	for _, record := range records {
		list = append(list, res.GetOperationRecordsByUserIdResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			Latency:     int64(record.Latency),
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetOperationRecordsByUserIdRes{
		Page:  r.Page,
		Size:  r.Size,
		List:  list,
		Total: total,
	}
	return rs, nil
}

// 按时间范围查询操作记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetOperationRecordsByTimeRange(
	ctx *gin.Context,
	r req.GetOperationRecordsByTimeRangeReq,
) (rs res.GetOperationRecordsByTimeRangeRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var records []systemRbac.SysOperationRecord

	db := global.GVA_DB.Model(&systemRbac.SysOperationRecord{})

	// 只有当时间都不为空时才添加时间条件
	if r.StartTime != "" && r.EndTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", r.StartTime)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", r.EndTime)
			if err == nil {
				db = db.Where("created_at BETWEEN ? AND ?", startTime, endTime)
			}
		}
	}

	db.Count(&total)
	db = db.Order("id desc").Limit(limit).Offset(offset).Preload("User")
	err = db.Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetOperationRecordsByTimeRangeResList
	for _, record := range records {
		list = append(list, res.GetOperationRecordsByTimeRangeResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			Latency:     int64(record.Latency),
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetOperationRecordsByTimeRangeRes{
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Page:      r.Page,
		Size:      r.Size,
		List:      list,
		Total:     total,
	}
	return rs, nil
}

// 获取最近N条记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetRecentOperationRecords(
	ctx *gin.Context,
	r req.GetRecentOperationRecordsReq,
) (rs res.GetRecentOperationRecordsRes, err error) {
	if r.Limit <= 0 {
		r.Limit = 10
	}

	var records []systemRbac.SysOperationRecord
	err = global.GVA_DB.Order("id desc").Limit(r.Limit).Preload("User").Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetRecentOperationRecordsResList
	for _, record := range records {
		list = append(list, res.GetRecentOperationRecordsResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			Latency:     int64(record.Latency),
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetRecentOperationRecordsRes{
		Limit: r.Limit,
		List:  list,
	}
	return rs, nil
}

// 获取异常请求记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetErrorRecords(
	ctx *gin.Context,
	r req.GetErrorRecordsReq,
) (rs res.GetErrorRecordsRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var records []systemRbac.SysOperationRecord

	db := global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("status != ?", 200)
	db.Count(&total)
	db = db.Order("id desc").Limit(limit).Offset(offset).Preload("User")
	err = db.Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetErrorRecordsResList
	for _, record := range records {
		list = append(list, res.GetErrorRecordsResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetErrorRecordsRes{
		Page:  r.Page,
		Size:  r.Size,
		List:  list,
		Total: total,
	}
	return rs, nil
}

// 批量删除过期记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) DeleteExpiredRecords(
	ctx *gin.Context,
	r req.DeleteExpiredRecordsReq,
) (err error) {
	if r.Days <= 0 {
		r.Days = 30 // 默认保留30天
	}

	expiredTime := time.Now().AddDate(0, 0, -r.Days)
	err = global.GVA_DB.Where("created_at < ?", expiredTime).Delete(&systemRbac.SysOperationRecord{}).Error
	if err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除过期记录失败")
	}
	return nil
}

// 按请求方法查询记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetOperationRecordsByMethod(
	ctx *gin.Context,
	r req.GetOperationRecordsByMethodReq,
) (rs res.GetOperationRecordsByMethodRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var records []systemRbac.SysOperationRecord

	db := global.GVA_DB.Model(&systemRbac.SysOperationRecord{})
	// 只有当方法不为空时才添加方法条件
	if r.Method != "" {
		db = db.Where("method = ?", r.Method)
	}
	db.Count(&total)
	db = db.Order("id desc").Limit(limit).Offset(offset).Preload("User")
	err = db.Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetOperationRecordsByMethodResList
	for _, record := range records {
		list = append(list, res.GetOperationRecordsByMethodResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			Latency:     int64(record.Latency),
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetOperationRecordsByMethodRes{
		Page:  r.Page,
		Size:  r.Size,
		List:  list,
		Total: total,
	}
	return rs, nil
}

// 按请求路径查询记录-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) GetOperationRecordsByPath(
	ctx *gin.Context,
	r req.GetOperationRecordsByPathReq,
) (rs res.GetOperationRecordsByPathRes, err error) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	limit := r.Size
	offset := r.Size * (r.Page - 1)

	var total int64
	var records []systemRbac.SysOperationRecord

	db := global.GVA_DB.Model(&systemRbac.SysOperationRecord{})
	// 只有当路径不为空时才添加路径模糊搜索条件
	if r.Path != "" {
		db = db.Where("path LIKE ?", "%"+r.Path+"%")
	}
	db.Count(&total)
	db = db.Order("id desc").Limit(limit).Offset(offset).Preload("User")
	err = db.Find(&records).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR)
	}

	var list []res.GetOperationRecordsByPathResList
	for _, record := range records {
		list = append(list, res.GetOperationRecordsByPathResList{
			Id:          int64(record.ID),
			UserId:      record.UserID,
			Ip:          record.Ip,
			Method:      record.Method,
			Path:        record.Path,
			Status:      record.Status,
			Latency:     int64(record.Latency),
			ErrorMessage: record.ErrorMessage,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	rs = res.GetOperationRecordsByPathRes{
		Path:  r.Path,
		Page:  r.Page,
		Size:  r.Size,
		List:  list,
		Total: total,
	}
	return rs, nil
}

// 统计今日操作数-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年03月28日 15:27:20
func (s *OprationService) CountTodayOperations(
	ctx *gin.Context,
) (rs res.CountTodayOperationsRes, err error) {
	today := time.Now().Format("2006-01-02")

	var count int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("DATE(created_at) = ?", today).Count(&count)

	var successCount int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("DATE(created_at) = ? AND status = ?", today, 200).Count(&successCount)

	var errorCount int64
	global.GVA_DB.Model(&systemRbac.SysOperationRecord{}).Where("DATE(created_at) = ? AND status != ?", today, 200).Count(&errorCount)

	rs = res.CountTodayOperationsRes{
		Count:        count,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
	}
	return rs, nil
}

