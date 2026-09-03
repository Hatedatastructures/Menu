package response

type GetSysMonitorOverviewRes struct {
	CpuUsage float64 `json:"cpuUsage"` // CPU使用率(
	MemoryUsage float64 `json:"memoryUsage"` // 内存使用率(
	DiskUsage float64 `json:"diskUsage"` // 磁盘使用率(
	Uptime string `json:"uptime"` // 系统运行时长
	GoroutineNum int `json:"goroutineNum"` // 协程数量
}

type GetSystemDetailRes struct {
	CpuInfo GetSystemDetailResCpuinfo `json:"cpuInfo"`
	MemoryInfo GetSystemDetailResMemoryinfo `json:"memoryInfo"`
	RuntimeInfo GetSystemDetailResRuntimeinfo `json:"runtimeInfo"`
	ServerInfo GetSystemDetailResServerinfo `json:"serverInfo"`
	DiskInfos []GetSystemDetailResDiskinfo `json:"diskInfos"` // []
	NetworkInfo GetSystemDetailResNetworkinfo `json:"networkInfo"`
	ProcessInfo GetSystemDetailResProcessinfo `json:"processInfo"`
	CpuUsage float64 `json:"cpuUsage"` // 当前CPU使用率(
	MemoryUsage float64 `json:"memoryUsage"` // 当前内存使用率(
}

type GetSystemDetailResCpuinfo struct {
	Os string `json:"os"` // 操作系统
	Arch string `json:"arch"` // 系统架构
	CoreNum int `json:"coreNum"` // CPU核心数
	Model string `json:"model"` // CPU型号
	Load float64 `json:"load"` // 系统负载（Linux有效）
}

type GetSystemDetailResMemoryinfo struct {
	Total float64 `json:"total"` // 总内存(MB)
	Used float64 `json:"used"` // 已用内存(MB)
	Free float64 `json:"free"` // 空闲内存(MB)
	Usage float64 `json:"usage"` // 使用率(
	ProcessUsed float64 `json:"processUsed"` // 当前进程占用(MB)
}

type GetSystemDetailResRuntimeinfo struct {
	GoVersion string `json:"goVersion"` // Go版本
	GoroutineNum int `json:"goroutineNum"` // 协程数量
	ProcessStartTime string `json:"processStartTime"` // 进程启动时间
	Uptime string `json:"uptime"` // 运行时长
}

type GetSystemDetailResServerinfo struct {
	Name string `json:"name"` // 服务器名称
	Ip string `json:"ip"` // 服务器IP
	OsName string `json:"osName"` // 操作系统
	OsVersion string `json:"osVersion"` // 系统版本
	Platform string `json:"platform"` // 平台信息
}

type GetSystemDetailResDiskinfo struct {
	Path string `json:"path"` // 盘符路径
	Total float64 `json:"total"` // 总大小(GB)
	Available float64 `json:"available"` // 可用大小(GB)
	Used float64 `json:"used"` // 已用大小(GB)
	UsageRate float64 `json:"usageRate"` // 使用率(
}

type GetSystemDetailResNetworkinfo struct {
	BytesSent uint64 `json:"bytesSent"` // 发送字节数
	BytesRecv uint64 `json:"bytesRecv"` // 接收字节数
}

type GetSystemDetailResProcessinfo struct {
	Pid int `json:"pid"` // 进程ID
	CpuPercent float64 `json:"cpuPercent"` // 进程CPU占用(
	MemoryPercent float64 `json:"memoryPercent"` // 进程内存占用(
}

type GetCpuUsageRes struct {
	Usage float64 `json:"usage"` // CPU使用率(
	Timestamp int64 `json:"timestamp"` // 时间戳
}

type GetMemoryUsageRes struct {
	Usage float64 `json:"usage"` // 内存使用率(
	Used float64 `json:"used"` // 已用内存(MB)
	Total float64 `json:"total"` // 总内存(MB)
	Timestamp int64 `json:"timestamp"` // 时间戳
}

type GetDiskUsageRes struct {
	Usage float64 `json:"usage"` // 磁盘使用率(
	Timestamp int64 `json:"timestamp"` // 时间戳
}

type GetMonitorTimelineRes struct {
	Type string `json:"type"`
	Points int `json:"points"`
	List []GetMonitorTimelineResList `json:"list"` // []
}

type GetMonitorTimelineResList struct {
	Value float64 `json:"value"` // 数值
	Timestamp int64 `json:"timestamp"` // 时间戳
}

type GetNetworkStatsRes struct {
	BytesSent uint64 `json:"bytesSent"` // 发送字节数
	BytesRecv uint64 `json:"bytesRecv"` // 接收字节数
	PacketsSent uint64 `json:"packetsSent"` // 发送包数
	PacketsRecv uint64 `json:"packetsRecv"` // 接收包数
}

type GetDiskIORes struct {
	ReadBytes uint64 `json:"readBytes"` // 读取字节数
	WriteBytes uint64 `json:"writeBytes"` // 写入字节数
	ReadCount uint64 `json:"readCount"` // 读取次数
	WriteCount uint64 `json:"writeCount"` // 写入次数
}

type GetProcessInfoRes struct {
	Pid int `json:"pid"` // 进程ID
	CpuPercent float64 `json:"cpuPercent"` // CPU占用(
	MemoryPercent float64 `json:"memoryPercent"` // 内存占用(
	MemoryUsed float64 `json:"memoryUsed"` // 内存使用(MB)
	ThreadCount int `json:"threadCount"` // 线程数
	GoroutineNum int `json:"goroutineNum"` // 协程数
}

type GetServerInfoRes struct {
	Hostname string `json:"hostname"` // 主机名
	Ip string `json:"ip"` // IP地址
	Os string `json:"os"` // 操作系统
	Platform string `json:"platform"` // 平台
	PlatformVersion string `json:"platformVersion"` // 平台版本
	KernelVersion string `json:"kernelVersion"` // 内核版本
	Arch string `json:"arch"` // 架构
}
