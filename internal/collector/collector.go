package collector

import (
	"context"
	"time"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
	"github.com/Ashmit-Singh-Gogia/sysmon/internal/proc"
)

func StreamStats(ctx context.Context, interval time.Duration) <-chan metrics.StatInfo {
	ch := make(chan metrics.StatInfo, 1) // buffer = 1 (latest value)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer close(ch)

		var prevSnap metrics.CPUSnapshot

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				currSnap, err := proc.ReadStat()
				if err != nil {
					continue
				}

				// First sample → just initialize baseline
				if prevSnap.Time.IsZero() {
					prevSnap = currSnap
					continue
				}

				stats := metrics.StatInfo{
					CPUDetails:         metrics.ComputeCpuUsage(prevSnap, currSnap),
					BootTime:           currSnap.BootTime,
					ContextSwitches:    currSnap.ContextSwitches,
					HardWareInterrupts: currSnap.HardWareInterrupts,
					SoftWareInterrupts: currSnap.SoftWareInterrupts,
					NumberOfCores:      uint64(len(currSnap.Cores)),
					Processes:          currSnap.Processes,
					ProcsBlocked:       currSnap.ProcsBlocked,
					ProcsRunning:       currSnap.ProcsRunning,
				}

				prevSnap = currSnap

				// Non-blocking send (drop frame if UI is slow)
				select {
				case ch <- stats:
				default:
				}
			}
		}
	}()

	return ch
}

func StreamMemInfo(ctx context.Context, interval time.Duration) <-chan metrics.MemoryComputed {
	ch := make(chan metrics.MemoryComputed, 1) // buffer = 1

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				memSnap, err := proc.ReadMemInfo()
				if err != nil {
					continue
				}

				memStats := metrics.Calculate_Mem_Usage(memSnap)

				select {
				case ch <- memStats:
				default:
				}
			}
		}
	}()

	return ch
}
