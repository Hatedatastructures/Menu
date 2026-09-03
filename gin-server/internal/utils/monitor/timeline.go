package monitor

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	DefaultMaxPoints = 60 // 默认保留60个数据点
	DefaultInterval = 5  // 默认采样间隔5秒
)

// TimelinePoint 时间序列数据点
type TimelinePoint struct {
	Value     float64 // 数值
	Timestamp int64   // 时间戳(毫秒)
}

// TimelineCollector 时间序列数据收集器
type TimelineCollector struct {
	cpuPoints     []TimelinePoint
	memoryPoints  []TimelinePoint
	diskPoints    []TimelinePoint
	maxPoints     int
	interval      time.Duration
	mu            sync.RWMutex
	stopCh        chan struct{}
}

// NewTimelineCollector 创建时间序列收集器
func NewTimelineCollector(maxPoints int, intervalSec int) *TimelineCollector {
	if maxPoints <= 0 {
		maxPoints = DefaultMaxPoints
	}
	if intervalSec <= 0 {
		intervalSec = DefaultInterval
	}

	tc := &TimelineCollector{
		cpuPoints:  make([]TimelinePoint, 0, maxPoints),
		memoryPoints: make([]TimelinePoint, 0, maxPoints),
		diskPoints:  make([]TimelinePoint, 0, maxPoints),
		maxPoints:  maxPoints,
		interval:   time.Duration(intervalSec) * time.Second,
		stopCh:     make(chan struct{}),
	}

	// 启动后台采集
	go tc.startCollecting()

	return tc
}

// startCollecting 启动后台采集
func (tc *TimelineCollector) startCollecting() {
	ticker := time.NewTicker(tc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tc.collect()
		case <-tc.stopCh:
			return
		}
	}
}

// collect 采集当前数据
func (tc *TimelineCollector) collect() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	now := time.Now().UnixMilli()

	// CPU使用率
	cpuUsage, _ := GetCpuUsage()
	tc.cpuPoints = tc.appendPoint(tc.cpuPoints, cpuUsage, now)

	// 内存使用率
	memUsage, _, _, _ := GetMemoryUsagePercent()
	tc.memoryPoints = tc.appendPoint(tc.memoryPoints, memUsage, now)

	// 磁盘使用率
	diskUsage, _ := GetDiskUsage()
	tc.diskPoints = tc.appendPoint(tc.diskPoints, diskUsage, now)
}

// appendPoint 添加数据点并保持最大数量
func (tc *TimelineCollector) appendPoint(points []TimelinePoint, value float64, timestamp int64) []TimelinePoint {
	points = append(points, TimelinePoint{
		Value:     value,
		Timestamp: timestamp,
	})

	// 保持最大数据点数量
	if len(points) > tc.maxPoints {
		points = points[len(points)-tc.maxPoints:]
	}

	return points
}

// GetCPUPoints 获取CPU历史数据
func (tc *TimelineCollector) GetCPUPoints() []TimelinePoint {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	result := make([]TimelinePoint, len(tc.cpuPoints))
	copy(result, tc.cpuPoints)
	return result
}

// GetMemoryPoints 获取内存历史数据
func (tc *TimelineCollector) GetMemoryPoints() []TimelinePoint {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	result := make([]TimelinePoint, len(tc.memoryPoints))
	copy(result, tc.memoryPoints)
	return result
}

// GetDiskPoints 获取磁盘历史数据
func (tc *TimelineCollector) GetDiskPoints() []TimelinePoint {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	result := make([]TimelinePoint, len(tc.diskPoints))
	copy(result, tc.diskPoints)
	return result
}

// GetPointsByType 根据类型获取历史数据
func (tc *TimelineCollector) GetPointsByType(metricType string) []TimelinePoint {
	switch metricType {
	case "cpu":
		return tc.GetCPUPoints()
	case "memory":
		return tc.GetMemoryPoints()
	case "disk":
		return tc.GetDiskPoints()
	default:
		return nil
	}
}

// Stop 停止采集
func (tc *TimelineCollector) Stop() {
	close(tc.stopCh)
}

// 全局时间序列收集器
var globalTimeline *TimelineCollector

// InitGlobalTimeline 初始化全局时间序列收集器
func InitGlobalTimeline(maxPoints int, intervalSec int) {
	globalTimeline = NewTimelineCollector(maxPoints, intervalSec)
}

// GetGlobalTimeline 获取全局时间序列收集器
func GetGlobalTimeline() *TimelineCollector {
	if globalTimeline == nil {
		globalTimeline = NewTimelineCollector(DefaultMaxPoints, DefaultInterval)
	}
	return globalTimeline
}

// GetCpuUsage 获取CPU使用率(包级别函数,用于后台采集)
func GetCpuUsage() (float64, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percent) > 0 {
		return percent[0], nil
	}
	return 0, nil
}

// GetMemoryUsage 获取内存使用率(包级别函数,用于后台采集)
func GetMemoryUsagePercent() (float64, float64, float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return memInfo.UsedPercent, float64(memInfo.Used) / MB, float64(memInfo.Total) / MB, nil
}

// GetDiskUsage 获取磁盘使用率(包级别函数,用于后台采集)
func GetDiskUsage() (float64, error) {
	parts, err := disk.Partitions(true)
	if err != nil {
		return 0, err
	}

	for _, part := range parts {
		if part.Fstype == "tmpfs" || part.Fstype == "squashfs" || part.Fstype == "overlay" {
			continue
		}
		// 优先获取C盘或根分区
		if len(part.Mountpoint) >= 2 && part.Mountpoint[1] == ':' {
			usage, err := disk.Usage(part.Mountpoint)
			if err == nil {
				return usage.UsedPercent, nil
			}
		}
	}

	// 取第一个有效分区
	if len(parts) > 0 {
		usage, err := disk.Usage(parts[0].Mountpoint)
		if err == nil {
			return usage.UsedPercent, nil
		}
	}

	return 0, nil
}