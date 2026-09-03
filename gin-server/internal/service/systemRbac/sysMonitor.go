package systemRbac

import (
	"time"

	"shack/internal/model/systemRbac/request"
	"shack/internal/model/systemRbac/response"
	"shack/internal/utils/monitor"

	"github.com/gin-gonic/gin"
)

type SysMonitorService struct {
	collector *monitor.Collector
}

var sysMonitorService *SysMonitorService

// InitSysMonitorService 初始化系统监控服务
func InitSysMonitorService() {
	sysMonitorService = &SysMonitorService{
		collector: monitor.New(),
	}
	// 启动后台时间序列采集
	monitor.InitGlobalTimeline(60, 5)
}

// GetSysMonitorService 获取系统监控服务实例
func GetSysMonitorService() *SysMonitorService {
	if sysMonitorService == nil {
		InitSysMonitorService()
	}
	return sysMonitorService
}

// GetSysMonitorOverview 获取系统监控概览（首页卡片）
func (s *SysMonitorService) GetSysMonitorOverview(
	ctx *gin.Context,
) (response.GetSysMonitorOverviewRes, error) {
	data, err := s.collector.GetAll()
	if err != nil {
		return response.GetSysMonitorOverviewRes{}, err
	}

	return response.GetSysMonitorOverviewRes{
		CpuUsage:     data.CpuUsage,
		MemoryUsage:  data.MemoryUsage,
		DiskUsage:    data.DiskUsage,
		Uptime:       data.RuntimeInfo.Uptime,
		GoroutineNum: data.RuntimeInfo.GoroutineNum,
	}, nil
}

// GetSystemDetail 获取系统详细信息
func (s *SysMonitorService) GetSystemDetail(
	ctx *gin.Context,
) (response.GetSystemDetailRes, error) {
	data, err := s.collector.GetAll()
	if err != nil {
		return response.GetSystemDetailRes{}, err
	}

	diskInfos := make([]response.GetSystemDetailResDiskinfo, 0, len(data.DiskInfos))
	for _, d := range data.DiskInfos {
		diskInfos = append(diskInfos, response.GetSystemDetailResDiskinfo{
			Path:      d.Path,
			Total:     d.Total,
			Available: d.Available,
			Used:      d.Used,
			UsageRate: d.UsageRate,
		})
	}

	return response.GetSystemDetailRes{
		CpuInfo: response.GetSystemDetailResCpuinfo{
			Os:      data.CPUInfo.Os,
			Arch:    data.CPUInfo.Arch,
			CoreNum: data.CPUInfo.CoreNum,
			Model:   data.CPUInfo.Model,
			Load:    data.CPUInfo.Load,
		},
		MemoryInfo: response.GetSystemDetailResMemoryinfo{
			Total:       data.MemoryInfo.Total,
			Used:        data.MemoryInfo.Used,
			Free:        data.MemoryInfo.Free,
			Usage:       data.MemoryInfo.Usage,
			ProcessUsed: data.MemoryInfo.ProcessUsed,
		},
		RuntimeInfo: response.GetSystemDetailResRuntimeinfo{
			GoVersion:        data.RuntimeInfo.GoVersion,
			GoroutineNum:     data.RuntimeInfo.GoroutineNum,
			ProcessStartTime: data.RuntimeInfo.ProcessStartTime,
			Uptime:           data.RuntimeInfo.Uptime,
		},
		ServerInfo: response.GetSystemDetailResServerinfo{
			Name:      data.ServerInfo.Name,
			Ip:        data.ServerInfo.Ip,
			OsName:    data.ServerInfo.OsName,
			OsVersion: data.ServerInfo.OsVersion,
			Platform:  data.ServerInfo.Platform,
		},
		DiskInfos: diskInfos,
		NetworkInfo: response.GetSystemDetailResNetworkinfo{
			BytesSent: data.NetworkInfo.BytesSent,
			BytesRecv: data.NetworkInfo.BytesRecv,
		},
		ProcessInfo: response.GetSystemDetailResProcessinfo{
			Pid:           data.ProcessInfo.Pid,
			CpuPercent:    data.ProcessInfo.CpuPercent,
			MemoryPercent: data.ProcessInfo.MemoryPercent,
		},
		CpuUsage:    data.CpuUsage,
		MemoryUsage: data.MemoryUsage,
	}, nil
}

// GetCpuUsage 获取实时CPU使用率
func (s *SysMonitorService) GetCpuUsage(
	ctx *gin.Context,
) (response.GetCpuUsageRes, error) {
	usage, err := s.collector.GetCpuUsage()
	if err != nil {
		return response.GetCpuUsageRes{}, err
	}

	return response.GetCpuUsageRes{
		Usage:     usage,
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

// GetMemoryUsage 获取实时内存使用率
func (s *SysMonitorService) GetMemoryUsage(
	ctx *gin.Context,
) (response.GetMemoryUsageRes, error) {
	usage, used, total, err := s.collector.GetMemoryUsage()
	if err != nil {
		return response.GetMemoryUsageRes{}, err
	}

	return response.GetMemoryUsageRes{
		Usage:     usage,
		Used:      used,
		Total:     total,
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

// GetDiskUsage 获取实时磁盘使用率
func (s *SysMonitorService) GetDiskUsage(
	ctx *gin.Context,
) (response.GetDiskUsageRes, error) {
	usage, err := s.collector.GetDiskUsage()
	if err != nil {
		return response.GetDiskUsageRes{}, err
	}

	return response.GetDiskUsageRes{
		Usage:     usage,
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

// GetMonitorTimeline 获取时间序列监控数据
func (s *SysMonitorService) GetMonitorTimeline(
	ctx *gin.Context,
	r request.GetMonitorTimelineReq,
) (response.GetMonitorTimelineRes, error) {
	timeline := monitor.GetGlobalTimeline()
	points := r.Points
	if points <= 0 {
		points = 60
	}

	var list []monitor.TimelinePoint
	switch r.Type {
	case "cpu":
		list = timeline.GetCPUPoints()
	case "memory":
		list = timeline.GetMemoryPoints()
	case "disk":
		list = timeline.GetDiskPoints()
	default:
		list = timeline.GetCPUPoints()
	}

	if len(list) > points {
		list = list[len(list)-points:]
	}

	resList := make([]response.GetMonitorTimelineResList, 0, len(list))
	for _, p := range list {
		resList = append(resList, response.GetMonitorTimelineResList{
			Value:     p.Value,
			Timestamp: p.Timestamp,
		})
	}

	return response.GetMonitorTimelineRes{
		Type:   r.Type,
		Points: points,
		List:   resList,
	}, nil
}

// GetNetworkStats 获取网络流量信息
func (s *SysMonitorService) GetNetworkStats(
	ctx *gin.Context,
) (response.GetNetworkStatsRes, error) {
	info, err := s.collector.GetNetworkInfo()
	if err != nil {
		return response.GetNetworkStatsRes{}, err
	}

	return response.GetNetworkStatsRes{
		BytesSent:   info.BytesSent,
		BytesRecv:   info.BytesRecv,
		PacketsSent: info.PacketsSent,
		PacketsRecv: info.PacketsRecv,
	}, nil
}

// GetDiskIO 获取磁盘IO信息
func (s *SysMonitorService) GetDiskIO(
	ctx *gin.Context,
) (response.GetDiskIORes, error) {
	info, err := s.collector.GetDiskIOInfo()
	if err != nil {
		return response.GetDiskIORes{}, err
	}

	return response.GetDiskIORes{
		ReadBytes:  info.ReadBytes,
		WriteBytes: info.WriteBytes,
		ReadCount:  info.ReadCount,
		WriteCount: info.WriteCount,
	}, nil
}

// GetProcessInfo 获取当前进程信息
func (s *SysMonitorService) GetProcessInfo(
	ctx *gin.Context,
) (response.GetProcessInfoRes, error) {
	info, err := s.collector.GetProcessInfo()
	if err != nil {
		return response.GetProcessInfoRes{}, err
	}

	return response.GetProcessInfoRes{
		Pid:           info.Pid,
		CpuPercent:    info.CpuPercent,
		MemoryPercent: info.MemoryPercent,
		MemoryUsed:    info.MemoryUsed,
		ThreadCount:   info.ThreadCount,
		GoroutineNum:  info.GoroutineNum,
	}, nil
}

// GetServerInfo 获取服务器基础信息
func (s *SysMonitorService) GetServerInfo(
	ctx *gin.Context,
) (response.GetServerInfoRes, error) {
	info, err := s.collector.GetServerInfo()
	if err != nil {
		return response.GetServerInfoRes{}, err
	}

	return response.GetServerInfoRes{
		Hostname:        info.Name,
		Ip:              info.Ip,
		Os:              info.OsName,
		Platform:        info.Platform,
		PlatformVersion: info.OsVersion,
		KernelVersion:   info.KernelVersion,
		Arch:            info.Arch,
	}, nil
}