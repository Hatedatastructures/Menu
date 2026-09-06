//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysMonitorRouter struct{}

func (s *SysMonitorRouter) InitSysMonitorRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysMonitorRouter := Router.Group("")
	utils.RegisterApi(sysMonitorRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/overview", "获取系统监控概览（首页卡片）-后台使用", sysMonitorApi.GetSysMonitorOverviewHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/detail", "获取系统详细信息-后台使用", sysMonitorApi.GetSystemDetailHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/cpu/usage", "获取实时CPU使用率-后台使用", sysMonitorApi.GetCpuUsageHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/memory/usage", "获取实时内存使用率-后台使用", sysMonitorApi.GetMemoryUsageHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/disk/usage", "获取实时磁盘使用率-后台使用", sysMonitorApi.GetDiskUsageHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/timeline", "获取时间序列监控数据（通用）-后台使用", sysMonitorApi.GetMonitorTimelineHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/network", "获取网络流量信息-后台使用", sysMonitorApi.GetNetworkStatsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/disk/io", "获取磁盘IO信息-后台使用", sysMonitorApi.GetDiskIOHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/process", "获取当前进程信息-后台使用", sysMonitorApi.GetProcessInfoHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sysMonitor/server", "获取服务器基础信息-后台使用", sysMonitorApi.GetServerInfoHandler),
	)
	return sysMonitorRouter
}
