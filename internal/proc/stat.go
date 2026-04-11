package proc

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
)

func ScanCPUAndCore(t string, fields []string, cpuTimes *metrics.CPUTimes, snapShot *metrics.CPUSnapshot) error {
	if len(fields) < 9 {
		return nil
	}
	cpuTimes.ID = fields[0]
	intFields := make([]uint64, len(fields))
	for i := 1; i < len(fields); i++ {
		num, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return err
		}
		intFields[i] = num
	}
	cpuTimes.User = intFields[1]
	cpuTimes.Nice = intFields[2]
	cpuTimes.System = intFields[3]
	cpuTimes.Idle = intFields[4]
	cpuTimes.Iowait = intFields[5]
	cpuTimes.IRQ = intFields[6]
	cpuTimes.SoftIRQ = intFields[7]
	cpuTimes.Steal = intFields[8]
	if t == "cpu" {
		snapShot.Total = *cpuTimes
	} else {
		snapShot.Cores = append(snapShot.Cores, *cpuTimes)
	}
	return nil
}

func ReadStat() (metrics.CPUSnapshot, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return metrics.CPUSnapshot{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	return parseStatFromReader(file, time.Now)
}

func parseStatFromReader(reader io.Reader, now func() time.Time) (metrics.CPUSnapshot, error) {
	var snapShot metrics.CPUSnapshot
	snapShot.Time = now()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		// process the line
		output := strings.Fields(line)

		if strings.HasPrefix(output[0], "cpu") {
			if output[0] == "cpu" {
				var totalSnapShot metrics.CPUTimes
				err := ScanCPUAndCore("cpu", output, &totalSnapShot, &snapShot)
				if err != nil {
					return metrics.CPUSnapshot{}, err
				}
			} else {
				// core
				var coreSnapShot metrics.CPUTimes
				err := ScanCPUAndCore("core", output, &coreSnapShot, &snapShot)
				if err != nil {
					return metrics.CPUSnapshot{}, err
				}
			}
		} else {
			// switch statement for the remaining file details
			if len(output) < 2 {
				continue
			}
			intFields := make([]uint64, len(output))
			for i := 1; i < len(output); i++ {
				num, err := strconv.ParseUint(output[i], 10, 64)
				if err != nil {
					return metrics.CPUSnapshot{}, err
				}
				intFields[i] = num
			}
			switch output[0] {
			case "intr":
				snapShot.HardWareInterrupts = intFields[1]
			case "ctxt":
				snapShot.ContextSwitches = intFields[1]
			case "btime":
				snapShot.BootTime = intFields[1]
			case "processes":
				snapShot.Processes = intFields[1]
			case "procs_running":
				snapShot.ProcsRunning = intFields[1]
			case "procs_blocked":
				snapShot.ProcsBlocked = intFields[1]
			case "softirq":
				snapShot.SoftWareInterrupts = intFields[1]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return metrics.CPUSnapshot{}, err
	}
	return snapShot, nil
}
