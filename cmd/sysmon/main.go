package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/collector"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cpuStream := collector.StreamStats(ctx, 1*time.Second)
	memStream := collector.StreamMemInfo(ctx, 1*time.Second)

	timeout := time.After(10 * time.Second)

	for {
		select {
		case cpu := <-cpuStream:
			fmt.Printf("CPU: %.2f%% | Cores: %d\n",
				cpu.CPUDetails.TotalUsage,
				cpu.NumberOfCores,
			)

		case mem := <-memStream:
			fmt.Printf("MEM: %.2f%% Used | Used: %d MB\n",
				mem.UsagePercent,
				mem.Used/1024/1024,
			)

		case <-timeout:
			fmt.Println("Test complete")
			return
		}
	}
}
