package metrics

type MemoryStats struct {
	Total     uint64
	Free      uint64
	Used      uint64  // To be computed
	UsedPer   float64 // To be computed
	Available uint64
	Buffers   uint64
	Cached    uint64

	SwapTotal uint64
	SwapFree  uint64

	Active   uint64
	Inactive uint64
	Dirty    uint64
	Slab     uint64
}

type MemoryComputed struct {
	Used             uint64
	UsagePercent     float64
	SwapUsed         uint64
	SwapUsagePercent float64
	Cache            uint64
}
