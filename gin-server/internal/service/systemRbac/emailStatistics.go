package systemRbac

import (
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	"github.com/gin-gonic/gin"
)

type EmailStatisticsService struct{}

// 获取邮件统计-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailStatisticsService) GetEmailStatistics(
	c *gin.Context,
	r req.GetEmailStatisticsReq,
) (rs res.GetEmailStatisticsRes, err error) {
	db := global.GVA_DB.Model(&systemRbac.EmailLog{})

	if r.StartTime != "" {
		startTime, _ := time.Parse("2006-01-02", r.StartTime)
		db = db.Where("created_at >= ?", startTime)
	}
	if r.EndTime != "" {
		endTime, _ := time.Parse("2006-01-02", r.EndTime)
		endTime = endTime.Add(24 * time.Hour)
		db = db.Where("created_at < ?", endTime)
	}

	var totalSent, totalSuccess, totalFailed int64
	db.Count(&totalSent)

	db.Where("status = ?", "success").Count(&totalSuccess)
	db.Where("status = ?", "failed").Count(&totalFailed)

	successRate := "0%"
	if totalSent > 0 {
		rate := float64(totalSuccess) / float64(totalSent) * 100
		successRate = fmt.Sprintf("%.1f%%", rate)
	}

	// 按类型统计
	var byTypeResults []struct {
		BizType string
		Count   int64
	}
	global.GVA_DB.Model(&systemRbac.EmailLog{}).
		Select("biz_type, count(*) as count").
		Group("biz_type").
		Scan(&byTypeResults)

	byType := make(map[string]int)
	for _, result := range byTypeResults {
		byType[result.BizType] = int(result.Count)
	}

	rs = res.GetEmailStatisticsRes{
		StartTime:    r.StartTime,
		EndTime:      r.EndTime,
		TotalSent:    int(totalSent),
		TotalSuccess: int(totalSuccess),
		TotalFailed:  int(totalFailed),
		SuccessRate:  successRate,
		ByType:       byType,
	}

	return rs, nil
}

// 获取队列状态-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 23:24:16
func (s *EmailStatisticsService) GetEmailQueueStatus(
	c *gin.Context,
) (rs res.GetEmailQueueStatusRes, err error) {
	// TODO: 从Redis获取队列状态
	rs = res.GetEmailQueueStatusRes{
		Pending:        0,
		Processing:     0,
		Workers:        5,
		AvgProcessTime: "2s",
	}
	return rs, nil
}
