package tests

import (
	"testing"

	"github.com/Ashmit-Singh-Gogia/sysmon/internal/metrics"
	"github.com/Ashmit-Singh-Gogia/sysmon/internal/proc"
)

func TestCalculateMemUsage(t *testing.T) {
	t.Run("calculates memory and swap usage", func(t *testing.T) {
		input := metrics.MemoryStats{
			Total:     8000,
			Available: 3000,
			Buffers:   100,
			Cached:    200,
			SwapTotal: 4000,
			SwapFree:  1000,
		}

		got := metrics.Calculate_Mem_Usage(input)

		if got.Used != 5000 {
			t.Fatalf("expected used memory 5000, got %d", got.Used)
		}
		if got.UsagePercent != 62.5 {
			t.Fatalf("expected usage percent 62.5, got %v", got.UsagePercent)
		}
		if got.SwapUsed != 3000 {
			t.Fatalf("expected swap used 3000, got %d", got.SwapUsed)
		}
		if got.SwapUsagePercent != 75 {
			t.Fatalf("expected swap usage percent 75, got %v", got.SwapUsagePercent)
		}
		if got.Cache != 300 {
			t.Fatalf("expected cache 300, got %d", got.Cache)
		}
	})

	t.Run("returns zero swap usage when swap total is zero", func(t *testing.T) {
		input := metrics.MemoryStats{
			Total:     1024,
			Available: 512,
			SwapTotal: 0,
			SwapFree:  0,
		}

		got := metrics.Calculate_Mem_Usage(input)

		if got.SwapUsagePercent != 0 {
			t.Fatalf("expected zero swap usage percent, got %v", got.SwapUsagePercent)
		}
		if got.SwapUsed != 0 {
			t.Fatalf("expected zero swap used, got %d", got.SwapUsed)
		}
	})
}

func TestReadMemInfo(t *testing.T) {
	stats, err := proc.ReadMemInfo()
	if err != nil {
		t.Fatalf("ReadMemInfo returned error: %v", err)
	}

	if stats.Total == 0 {
		t.Fatal("expected total memory to be populated")
	}
	if stats.Available == 0 {
		t.Fatal("expected available memory to be populated")
	}
}
