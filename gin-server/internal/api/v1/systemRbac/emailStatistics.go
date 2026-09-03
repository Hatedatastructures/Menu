package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailStatisticsApi struct{}

// GetEmailStatisticsHandler
// @Tags systemRbacemailStatisticsApi
// @Summary GetEmailStatisticsHandler 获取邮件统计-后台使用
// @Description GetEmailStatisticsHandler 获取邮件统计-后台使用
// @Param data body req.GetEmailStatisticsReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailStatisticsRes}
// @Router /email/statistics [GET]
func (s *EmailStatisticsApi) GetEmailStatisticsHandler(c *gin.Context) {
	var req req.GetEmailStatisticsReq
	// query 参数
	{
		val := c.Query("startTime")
		if val != "" {
			req.StartTime = val

		}
	}
	{
		val := c.Query("endTime")
		if val != "" {
			req.EndTime = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailStatisticsService.GetEmailStatistics(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetEmailQueueStatusHandler
// @Tags systemRbacemailStatisticsApi
// @Summary GetEmailQueueStatusHandler 获取队列状态-后台使用
// @Description GetEmailQueueStatusHandler 获取队列状态-后台使用
// @Success 200 {object} vo.Result{data=_.GetEmailQueueStatusRes}
// @Router /email/queue/status [GET]
func (s *EmailStatisticsApi) GetEmailQueueStatusHandler(c *gin.Context) {
	data, err := emailStatisticsService.GetEmailQueueStatus(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
