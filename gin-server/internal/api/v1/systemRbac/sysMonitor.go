package systemRbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"
	biz_err "shack/internal/error"
)

type SysMonitorApi struct{}

// GetSysMonitorOverviewHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetSysMonitorOverviewHandler 获取系统监控概览（首页卡片）-后台使用
// @Description GetSysMonitorOverviewHandler 获取系统监控概览（首页卡片）-后台使用
// @Success 200 {object} vo.Result{data=_.GetSysMonitorOverviewRes}
// @Router /sysMonitor/overview [GET]
func (s *SysMonitorApi) GetSysMonitorOverviewHandler(c *gin.Context) {
	data, err := sysMonitorService.GetSysMonitorOverview(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetSystemDetailHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetSystemDetailHandler 获取系统详细信息-后台使用
// @Description GetSystemDetailHandler 获取系统详细信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetSystemDetailRes}
// @Router /sysMonitor/detail [GET]
func (s *SysMonitorApi) GetSystemDetailHandler(c *gin.Context) {
	data, err := sysMonitorService.GetSystemDetail(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetCpuUsageHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetCpuUsageHandler 获取实时CPU使用率-后台使用
// @Description GetCpuUsageHandler 获取实时CPU使用率-后台使用
// @Success 200 {object} vo.Result{data=_.GetCpuUsageRes}
// @Router /sysMonitor/cpu/usage [GET]
func (s *SysMonitorApi) GetCpuUsageHandler(c *gin.Context) {
	data, err := sysMonitorService.GetCpuUsage(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetMemoryUsageHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetMemoryUsageHandler 获取实时内存使用率-后台使用
// @Description GetMemoryUsageHandler 获取实时内存使用率-后台使用
// @Success 200 {object} vo.Result{data=_.GetMemoryUsageRes}
// @Router /sysMonitor/memory/usage [GET]
func (s *SysMonitorApi) GetMemoryUsageHandler(c *gin.Context) {
	data, err := sysMonitorService.GetMemoryUsage(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDiskUsageHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetDiskUsageHandler 获取实时磁盘使用率-后台使用
// @Description GetDiskUsageHandler 获取实时磁盘使用率-后台使用
// @Success 200 {object} vo.Result{data=_.GetDiskUsageRes}
// @Router /sysMonitor/disk/usage [GET]
func (s *SysMonitorApi) GetDiskUsageHandler(c *gin.Context) {
	data, err := sysMonitorService.GetDiskUsage(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetMonitorTimelineHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetMonitorTimelineHandler 获取时间序列监控数据（通用）-后台使用
// @Description GetMonitorTimelineHandler 获取时间序列监控数据（通用）-后台使用
// @Param data body req.GetMonitorTimelineReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetMonitorTimelineRes}
// @Router /sysMonitor/timeline [GET]
func (s *SysMonitorApi) GetMonitorTimelineHandler(c *gin.Context) {
	var req req.GetMonitorTimelineReq
	// query 参数
	{
		val := c.Query("type")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Type = parsed
		}
	}
	{
		val := c.Query("points")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Points = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := sysMonitorService.GetMonitorTimeline(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetNetworkStatsHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetNetworkStatsHandler 获取网络流量信息-后台使用
// @Description GetNetworkStatsHandler 获取网络流量信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetNetworkStatsRes}
// @Router /sysMonitor/network [GET]
func (s *SysMonitorApi) GetNetworkStatsHandler(c *gin.Context) {
	data, err := sysMonitorService.GetNetworkStats(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetDiskIOHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetDiskIOHandler 获取磁盘IO信息-后台使用
// @Description GetDiskIOHandler 获取磁盘IO信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetDiskIORes}
// @Router /sysMonitor/disk/io [GET]
func (s *SysMonitorApi) GetDiskIOHandler(c *gin.Context) {
	data, err := sysMonitorService.GetDiskIO(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetProcessInfoHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetProcessInfoHandler 获取当前进程信息-后台使用
// @Description GetProcessInfoHandler 获取当前进程信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetProcessInfoRes}
// @Router /sysMonitor/process [GET]
func (s *SysMonitorApi) GetProcessInfoHandler(c *gin.Context) {
	data, err := sysMonitorService.GetProcessInfo(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetServerInfoHandler
// @Tags systemRbacsysMonitorApi
// @Summary GetServerInfoHandler 获取服务器基础信息-后台使用
// @Description GetServerInfoHandler 获取服务器基础信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetServerInfoRes}
// @Router /sysMonitor/server [GET]
func (s *SysMonitorApi) GetServerInfoHandler(c *gin.Context) {
	data, err := sysMonitorService.GetServerInfo(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}