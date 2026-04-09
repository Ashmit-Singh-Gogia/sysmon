package metrics

func (c *CPUTimes) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.Iowait + c.IRQ + c.SoftIRQ + c.Steal
}

func (c *CPUTimes) IdleTime() uint64 {
	return c.Idle + c.Iowait
}

func calculteCPUPercentage(prev, curr CPUTimes) float64 {
	prevTotal := prev.Total()
	currTotal := curr.Total()

	prevIdle := prev.IdleTime()
	currIdle := curr.IdleTime()

	totalDiff := currTotal - prevTotal
	idleDiff := currIdle - prevIdle

	return ((float64(totalDiff) - float64(idleDiff)) / float64(totalDiff)) * 100
}

func ComputeCpuUsage(prev, curr CPUSnapshot) CPUUsage {
	totalUsage := calculteCPUPercentage(prev.Total, curr.Total)

	cores := make([]CPUCoreUsage, len(prev.Cores))
	for i := range curr.Cores {
		if i >= len(prev.Cores) {
			continue
		}
		usage := calculteCPUPercentage(prev.Cores[i], curr.Cores[i])
		cores[i] = CPUCoreUsage{
			ID:    curr.Cores[i].ID,
			Usage: usage,
		}
	}
	return CPUUsage{
		TotalUsage: totalUsage,
		Cores:      cores,
	}
}
