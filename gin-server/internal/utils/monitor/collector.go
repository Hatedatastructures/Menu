package monitor

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// CPUInfo CPU信息
type CPUInfo struct {
	Os      string  // 操作系统
	Arch    string  // 系统架构
	CoreNum int     // CPU核心数
	Model   string  // CPU型号
	Load    float64 // 系统负载
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total        float64 // 总内存(MB)
	Used         float64 // 已用内存(MB)
	Free         float64 // 空闲内存(MB)
	Usage        float64 // 使用率(%)
	ProcessUsed float64 // 当前进程占用(MB)
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	GoVersion        string // Go版本
	GoroutineNum     int    // 协程数量
	ProcessStartTime string // 进程启动时间
	Uptime           string // 运行时长
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Name          string // 服务器名称
	Ip            string // 服务器IP
	OsName        string // 操作系统
	OsVersion     string // 系统版本
	Platform      string // 平台信息
	KernelVersion string // 内核版本
	Arch          string // 架构
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Path        string  // 盘符路径
	Total       float64 // 总大小(GB)
	Available   float64 // 可用大小(GB)
	Used        float64 // 已用大小(GB)
	UsageRate   float64 // 使用率(%)
	UsedMB      int64   // 已用大小(MB)
	TotalMB     int64   // 总大小(MB)
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	BytesSent   uint64 // 发送字节数
	BytesRecv   uint64 // 接收字节数
	PacketsSent uint64 // 发送包数
	PacketsRecv uint64 // 接收包数
}

// DiskIOInfo 磁盘IO信息
type DiskIOInfo struct {
	ReadBytes  uint64 // 读取字节数
	WriteBytes uint64 // 写入字节数
	ReadCount  uint64 // 读取次数
	WriteCount uint64 // 写入次数
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	Pid           int     // 进程ID
	CpuPercent    float64 // CPU占用(%)
	MemoryPercent float64 // 内存占用(%)
	MemoryUsed    float64 // 内存使用(MB)
	ThreadCount   int     // 线程数
	GoroutineNum  int     // 协程数
}

// MonitorData 监控数据
type MonitorData struct {
	CPUInfo      CPUInfo
	MemoryInfo   MemoryInfo
	RuntimeInfo  RuntimeInfo
	ServerInfo   ServerInfo
	DiskInfos    []DiskInfo
	NetworkInfo  NetworkInfo
	DiskIOInfo   DiskIOInfo
	ProcessInfo  ProcessInfo
	CpuUsage     float64 // 当前CPU使用率(%)
	MemoryUsage float64 // 当前内存使用率(%)
	DiskUsage   float64 // 磁盘使用率(%)
}

// Collector 系统监控数据收集器
type Collector struct {
	processStartTime time.Time
	prevNetStats     *net.IOCountersStat
	prevDiskStats    *disk.IOCountersStat
}

// New 创建新的监控数据收集器
func New() *Collector {
	return &Collector{
		processStartTime: time.Now(),
	}
}

// GetCPUInfo 获取CPU信息
func (c *Collector) GetCPUInfo() (CPUInfo, error) {
	info := CPUInfo{
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		CoreNum: runtime.NumCPU(),
	}

	// 获取CPU型号 (仅Linux/macOS)
	if cpus, err := cpu.Info(); err == nil && len(cpus) > 0 {
		info.Model = cpus[0].ModelName
	}

	return info, nil
}

// GetMemoryInfo 获取内存信息
func (c *Collector) GetMemoryInfo() (MemoryInfo, error) {
	info := MemoryInfo{}

	// 系统内存
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.Total = float64(memInfo.Total) / MB
		info.Used = float64(memInfo.Used) / MB
		info.Free = float64(memInfo.Free) / MB
		info.Usage = memInfo.UsedPercent
	}

	// 当前进程内存
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err == nil {
		if memInfo, err := proc.MemoryInfo(); err == nil {
			info.ProcessUsed = float64(memInfo.RSS) / MB
		}
	}

	return info, nil
}

// GetRuntimeInfo 获取运行时信息
func (c *Collector) GetRuntimeInfo() RuntimeInfo {
	return RuntimeInfo{
		GoVersion:        runtime.Version(),
		GoroutineNum:      runtime.NumGoroutine(),
		ProcessStartTime: c.processStartTime.Format("2006-01-02 15:04:05"),
		Uptime:           formatUptime(time.Since(c.processStartTime)),
	}
}

// GetServerInfo 获取服务器信息
func (c *Collector) GetServerInfo() (ServerInfo, error) {
	info := ServerInfo{}

	// 主机信息
	hostInfo, err := host.Info()
	if err == nil {
		info.Name = hostInfo.Hostname
		info.OsName = hostInfo.OS
		info.OsVersion = hostInfo.PlatformVersion
		info.Platform = hostInfo.Platform
		info.KernelVersion = hostInfo.KernelVersion
		info.Arch = hostInfo.PlatformFamily
	}

	// 获取第一个非回环IP
	if nics, err := net.Interfaces(); err == nil {
		for _, nic := range nics {
			if nic.Name == "eth0" || nic.Name == "en0" || nic.Name == "ens0" {
				for _, addr := range nic.Addrs {
					if addr.Addr != "127.0.0.1" && addr.Addr != "::1" {
						info.Ip = addr.Addr
						break
					}
				}
			}
		}
		// 如果没找到特定网卡，取第一个非回环IP
		if info.Ip == "" {
			for _, nic := range nics {
				if nic.Name != "lo" && nic.Name != "loopback" {
					for _, addr := range nic.Addrs {
						if addr.Addr != "127.0.0.1" && addr.Addr != "::1" {
							info.Ip = addr.Addr
							break
						}
					}
				}
				if info.Ip != "" {
					break
				}
			}
		}
	}

	return info, nil
}

// GetDiskInfos 获取磁盘信息
func (c *Collector) GetDiskInfos() ([]DiskInfo, error) {
	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	infos := make([]DiskInfo, 0, len(parts))
	for _, part := range parts {
		// 跳过只读和虚拟文件系统
		if part.Fstype == "tmpfs" || part.Fstype == "squashfs" || part.Fstype == "overlay" {
			continue
		}

		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		infos = append(infos, DiskInfo{
			Path:      part.Mountpoint,
			Total:     float64(usage.Total) / GB,
			Available: float64(usage.Free) / GB,
			Used:      float64(usage.Used) / GB,
			UsageRate: usage.UsedPercent,
			UsedMB:    int64(usage.Used) / MB,
			TotalMB:   int64(usage.Total) / MB,
		})
	}

	return infos, nil
}

// GetNetworkInfo 获取网络信息
func (c *Collector) GetNetworkInfo() (NetworkInfo, error) {
	info := NetworkInfo{}

	nics, err := net.IOCounters(true)
	if err != nil {
		return info, err
	}

	for _, nic := range nics {
		if nic.Name == "eth0" || nic.Name == "en0" || nic.Name == "ens0" {
			info.BytesSent += nic.BytesSent
			info.BytesRecv += nic.BytesRecv
			info.PacketsSent += nic.PacketsSent
			info.PacketsRecv += nic.PacketsRecv
		}
	}

	// 如果没有特定网卡，汇总所有非回环网卡
	if info.BytesSent == 0 && info.BytesRecv == 0 {
		for _, nic := range nics {
			if nic.Name != "lo" && nic.Name != "loopback" {
				info.BytesSent += nic.BytesSent
				info.BytesRecv += nic.BytesRecv
				info.PacketsSent += nic.PacketsSent
				info.PacketsRecv += nic.PacketsRecv
			}
		}
	}

	return info, nil
}

// GetDiskIOInfo 获取磁盘IO信息
func (c *Collector) GetDiskIOInfo() (DiskIOInfo, error) {
	info := DiskIOInfo{}

	ioCounters, err := disk.IOCounters()
	if err != nil {
		return info, err
	}

	for _, diskIO := range ioCounters {
		info.ReadBytes += diskIO.ReadBytes
		info.WriteBytes += diskIO.WriteBytes
		info.ReadCount += diskIO.ReadCount
		info.WriteCount += diskIO.WriteCount
	}

	return info, nil
}

// GetProcessInfo 获取进程信息
func (c *Collector) GetProcessInfo() (ProcessInfo, error) {
	info := ProcessInfo{}

	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return info, err
	}

	// PID
	info.Pid = int(proc.Pid)

	// CPU占用
	cpuPercent, _ := proc.CPUPercent()
	info.CpuPercent = cpuPercent

	// 内存占用
	memPercent, _ := proc.MemoryPercent()
	info.MemoryPercent = float64(memPercent)

	// 内存使用
	if memInfo, err := proc.MemoryInfo(); err == nil {
		info.MemoryUsed = float64(memInfo.RSS) / MB
	}

	// 线程数
	threads, _ := proc.Threads()
	info.ThreadCount = len(threads)

	// 协程数
	info.GoroutineNum = runtime.NumGoroutine()

	return info, nil
}

// GetCpuUsage 获取实时CPU使用率
func (c *Collector) GetCpuUsage() (float64, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percent) > 0 {
		return percent[0], nil
	}
	return 0, nil
}

// GetMemoryUsage 获取实时内存使用率
func (c *Collector) GetMemoryUsage() (float64, float64, float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return memInfo.UsedPercent, float64(memInfo.Used) / MB, float64(memInfo.Total) / MB, nil
}

// GetDiskUsage 获取实时磁盘使用率
func (c *Collector) GetDiskUsage() (float64, error) {
	infos, err := c.GetDiskInfos()
	if err != nil || len(infos) == 0 {
		return 0, err
	}

	// 简单取C盘使用率
	for _, info := range infos {
		if len(info.Path) >= 2 && info.Path[1] == ':' {
			return info.UsageRate, nil
		}
	}
	return infos[0].UsageRate, nil
}

// GetAll 收集所有监控数据
func (c *Collector) GetAll() (MonitorData, error) {
	data := MonitorData{}

	// CPU信息
	cpuInfo, _ := c.GetCPUInfo()
	data.CPUInfo = cpuInfo

	// 内存信息
	memoryInfo, _ := c.GetMemoryInfo()
	data.MemoryInfo = memoryInfo

	// 运行时信息
	data.RuntimeInfo = c.GetRuntimeInfo()

	// 服务器信息
	serverInfo, _ := c.GetServerInfo()
	data.ServerInfo = serverInfo

	// 磁盘信息
	diskInfos, _ := c.GetDiskInfos()
	data.DiskInfos = diskInfos

	// 网络信息
	networkInfo, _ := c.GetNetworkInfo()
	data.NetworkInfo = networkInfo

	// 磁盘IO信息
	diskIOInfo, _ := c.GetDiskIOInfo()
	data.DiskIOInfo = diskIOInfo

	// 进程信息
	processInfo, _ := c.GetProcessInfo()
	data.ProcessInfo = processInfo

	// 实时指标
	data.CpuUsage, _ = c.GetCpuUsage()
	memUsage, _, _, _ := c.GetMemoryUsage()
	data.MemoryUsage = memUsage
	data.DiskUsage, _ = c.GetDiskUsage()

	return data, nil
}

// formatUptime 格式化运行时长
func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分%d秒", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}