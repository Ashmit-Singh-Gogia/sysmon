package metrics

import "time"

type CPUTimes struct {
	ID      string
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
}

type CPUSnapshot struct {
	Time               time.Time
	Total              CPUTimes
	Cores              []CPUTimes
	HardWareInterrupts uint64
	ContextSwitches    uint64
	BootTime           uint64
	Processes          uint64
	ProcsRunning       uint64
	ProcsBlocked       uint64
	SoftWareInterrupts uint64
}

type CPUCoreUsage struct {
	ID    string
	Usage float64
}

type CPUUsage struct {
	TotalUsage float64
	Cores      []CPUCoreUsage
}

type StatInfo struct {
	CPUDetails         CPUUsage
	HardWareInterrupts uint64
	ContextSwitches    uint64
	BootTime           uint64
	Processes          uint64
	ProcsRunning       uint64
	ProcsBlocked       uint64
	SoftWareInterrupts uint64
	NumberOfCores      uint64 // might be used
}
