package tests

import (
	"testing"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
	"github.com/Ashmit-Singh-Gogia/sysmon/internal/proc"
)

func TestScanCPUAndCore(t *testing.T) {
	t.Run("maps total cpu fields by position", func(t *testing.T) {
		fields := []string{"cpu", "10", "11", "12", "13", "14", "15", "16", "17", "0", "0"}
		var snap metrics.CPUSnapshot
		var times metrics.CPUTimes

		err := proc.ScanCPUAndCore("cpu", fields, &times, &snap)
		if err != nil {
			t.Fatalf("ScanCPUAndCore returned error: %v", err)
		}

		if snap.Total.ID != "cpu" {
			t.Fatalf("expected ID cpu, got %s", snap.Total.ID)
		}
		if snap.Total.User != 10 || snap.Total.Nice != 11 || snap.Total.System != 12 {
			t.Fatalf("unexpected user/nice/system mapping: %+v", snap.Total)
		}
		if snap.Total.Idle != 13 || snap.Total.Iowait != 14 || snap.Total.IRQ != 15 {
			t.Fatalf("unexpected idle/iowait/irq mapping: %+v", snap.Total)
		}
		if snap.Total.SoftIRQ != 16 || snap.Total.Steal != 17 {
			t.Fatalf("unexpected softirq/steal mapping: %+v", snap.Total)
		}
	})

	t.Run("appends core cpu data", func(t *testing.T) {
		fields := []string{"cpu0", "1", "2", "3", "4", "5", "6", "7", "8", "0", "0"}
		var snap metrics.CPUSnapshot
		var times metrics.CPUTimes

		err := proc.ScanCPUAndCore("core", fields, &times, &snap)
		if err != nil {
			t.Fatalf("ScanCPUAndCore returned error: %v", err)
		}

		if len(snap.Cores) != 1 {
			t.Fatalf("expected 1 core snapshot, got %d", len(snap.Cores))
		}
		if snap.Cores[0].ID != "cpu0" || snap.Cores[0].Iowait != 5 || snap.Cores[0].Steal != 8 {
			t.Fatalf("unexpected core mapping: %+v", snap.Cores[0])
		}
	})

	t.Run("ignores short cpu line", func(t *testing.T) {
		fields := []string{"cpu", "1", "2"}
		var snap metrics.CPUSnapshot
		var times metrics.CPUTimes

		err := proc.ScanCPUAndCore("cpu", fields, &times, &snap)
		if err != nil {
			t.Fatalf("expected nil error for short line, got %v", err)
		}
		if snap.Total.ID != "" {
			t.Fatalf("expected no total data for short line, got %+v", snap.Total)
		}
	})
}

func TestScanCPUAndCore_InvalidNumber(t *testing.T) {
	fields := []string{"cpu", "1", "x", "3", "4", "5", "6", "7", "8", "0", "0"}
	var snap metrics.CPUSnapshot
	var times metrics.CPUTimes

	err := proc.ScanCPUAndCore("cpu", fields, &times, &snap)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestReadStat(t *testing.T) {
	snap, err := proc.ReadStat()
	if err != nil {
		t.Fatalf("ReadStat returned error: %v", err)
	}
	if snap.Time.IsZero() {
		t.Fatal("expected non-zero snapshot time")
	}
	if snap.Total.ID != "cpu" {
		t.Fatalf("expected total cpu ID to be cpu, got %s", snap.Total.ID)
	}
}
