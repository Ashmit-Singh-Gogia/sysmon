package main

import (
	"fmt"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/collector"
)

func main() {
	stats := collector.InitStats()
	fmt.Println(stats.CPUDetails.TotalUsage)
}
