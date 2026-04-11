package metrics

func Calculate_Mem_Usage(memStats MemoryStats) MemoryComputed {
	// Calculate used and used per
	var memoryComputed MemoryComputed
	used := memStats.Total - memStats.Available
	cache := memStats.Cached + memStats.Buffers
	usedPer := float64(used) * 100 / float64(memStats.Total)
	swapUsed := memStats.SwapTotal - memStats.SwapFree
	if memStats.SwapTotal > 0 {
		swapUsedPer := float64(swapUsed) * 100 / float64(memStats.SwapTotal)
		memoryComputed.SwapUsagePercent = swapUsedPer
	} else {
		memoryComputed.SwapUsagePercent = 0
	}
	memoryComputed.Used = used
	memoryComputed.UsagePercent = usedPer
	memoryComputed.SwapUsed = swapUsed
	memoryComputed.Cache = cache
	return memoryComputed
}
