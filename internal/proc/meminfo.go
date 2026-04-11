package proc

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
)

func ReadMemInfo() (metrics.MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return metrics.MemoryStats{}, err
	}
	defer func() {
		_ = file.Close()
	}()
	return parseMemoryFile(file)
}
func parseMemoryFile(reader io.Reader) (metrics.MemoryStats, error) {
	var MemStats metrics.MemoryStats
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		// Process the line
		output := strings.SplitN(line, ":", 2)
		if len(output) != 2 {
			continue
		}
		key := output[0]

		fields := strings.Fields(output[1])
		if len(fields) == 0 {
			continue
		}
		value, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return metrics.MemoryStats{}, err
		}

		switch key {
		case "MemFree":
			MemStats.Free = value
		case "MemAvailable":
			MemStats.Available = value
		case "Buffers":
			MemStats.Buffers = value
		case "Cached":
			MemStats.Cached = value
		case "SwapFree":
			MemStats.SwapFree = value
		case "SwapTotal":
			MemStats.SwapTotal = value
		case "Dirty":
			MemStats.Dirty = value
		case "Inactive":
			MemStats.Inactive = value
		case "Slab":
			MemStats.Slab = value
		case "MemTotal":
			MemStats.Total = value
		case "Active":
			MemStats.Active = value
		}
	}
	return MemStats, nil
}
