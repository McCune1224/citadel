package services

import (
	"context"
	"sync"
	"time"

	"citadel/highway17/internal/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

type SystemStatsService struct {
	log        *zap.Logger
	mu         sync.RWMutex
	cached     *models.SystemStats
	cacheTTL   time.Duration
	lastUpdate time.Time
}

func NewSystemStatsService(log *zap.Logger, cacheTTL time.Duration) *SystemStatsService {
	return &SystemStatsService{
		log:      log,
		cacheTTL: cacheTTL,
	}
}

// GetStats fetches current system statistics
func (ss *SystemStatsService) GetStats(ctx context.Context) (*models.SystemStats, error) {
	ss.mu.RLock()
	if ss.cached != nil && time.Since(ss.lastUpdate) < ss.cacheTTL {
		defer ss.mu.RUnlock()
		return ss.cached, nil
	}
	ss.mu.RUnlock()

	stats := &models.SystemStats{
		LastUpdated: time.Now(),
	}

	// CPU usage (per-core average)
	cpuPercents, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil {
		ss.log.Sugar().Warnw("failed to get CPU percentage", "error", err)
	} else if len(cpuPercents) > 0 {
		stats.CPUPercent = cpuPercents[0]
	}

	// Memory usage
	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		ss.log.Sugar().Warnw("failed to get memory info", "error", err)
	} else {
		stats.MemoryPercent = vmem.UsedPercent
		stats.MemoryUsedGB = float64(vmem.Used) / (1024 * 1024 * 1024)
		stats.MemoryTotalGB = float64(vmem.Total) / (1024 * 1024 * 1024)
	}

	// Disk usage (root partition)
	diskUsage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		ss.log.Sugar().Warnw("failed to get disk info", "error", err)
	} else {
		stats.DiskPercent = diskUsage.UsedPercent
		stats.DiskUsedGB = float64(diskUsage.Used) / (1024 * 1024 * 1024)
		stats.DiskTotalGB = float64(diskUsage.Total) / (1024 * 1024 * 1024)
	}

	// Uptime
	uptime, err := host.UptimeWithContext(ctx)
	if err != nil {
		ss.log.Sugar().Warnw("failed to get uptime", "error", err)
	} else {
		stats.UptimeSeconds = uptime
	}

	// Process count
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		ss.log.Sugar().Warnw("failed to get process count", "error", err)
	} else {
		stats.ProcessCount = len(processes)
	}

	// Load average (platform specific - only works on Unix-like systems)
	// This is handled by host.SensorsTemperaturesWithContext but we'll skip for now
	stats.LoadAverage[0] = 0
	stats.LoadAverage[1] = 0
	stats.LoadAverage[2] = 0

	// Cache the result
	ss.mu.Lock()
	ss.cached = stats
	ss.lastUpdate = time.Now()
	ss.mu.Unlock()

	return stats, nil
}

// ClearCache clears the system stats cache
func (ss *SystemStatsService) ClearCache() {
	ss.mu.Lock()
	ss.cached = nil
	ss.mu.Unlock()
}
