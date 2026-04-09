package collector

import (
	"time"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
	"github.com/Ashmit-Singh-Gogia/sysmon/internal/proc"
)

func InitStats() metrics.StatInfo {
	var stats metrics.StatInfo
	var prevSnap metrics.CPUSnapshot
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		currSnap, err := proc.ReadStat()
		if err != nil {
			continue
		}
		if prevSnap.Time.IsZero() {
			prevSnap = currSnap
			continue
		}
		cpuUsage := metrics.ComputeCpuUsage(prevSnap, currSnap)
		stats.CPUDetails = cpuUsage
		stats.BootTime = currSnap.BootTime
		stats.ContextSwitches = currSnap.ContextSwitches
		stats.HardWareInterrupts = currSnap.HardWareInterrupts
		stats.SoftWareInterrupts = currSnap.SoftWareInterrupts
		stats.NumberOfCores = uint64(len(currSnap.Cores))
		stats.Processes = currSnap.Processes
		stats.ProcsBlocked = currSnap.ProcsBlocked
		stats.ProcsRunning = currSnap.ProcsRunning
		return stats
	}
}
