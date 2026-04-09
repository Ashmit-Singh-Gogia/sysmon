# sysmon

Linux-only CLI + TUI system monitor built in Go.

It reads kernel-exposed system data from low-level Linux interfaces, computes real-time metrics, and renders a clean terminal dashboard.

## Purpose

This project solves a simple problem: show accurate, live system state in a terminal with minimal overhead.

Why parse Linux system files directly:

- You read raw kernel counters instead of shelling out to external tools.
- You control parsing and calculations end-to-end.
- You can reason about accuracy from first principles.
- You keep dependencies small and behavior predictable.

How this differs from tools like htop:

- htop is feature-rich and general-purpose.
- sysmon is code-first and architecture-first.
- sysmon is designed as an extensible monitoring engine plus a focused TUI.
- The internals are intentionally simple to modify, test, and verify.

## Features

- CPU usage: total and per-core, delta-based.
- Memory usage: parsed from /proc/meminfo.
- Disk usage: filesystem-level stats.
- Network usage: rate-based throughput (bytes/sec) from counter deltas.
- Real-time refresh loop.
- Structured terminal UI with sections and progress bars.

## Architecture

Data flow:

```text
/proc and kernel interfaces
				|
				v
	 proc layer (raw parse)
				|
				v
	snapshot models (point-in-time counters)
				|
				v
 metrics layer (delta + derived values)
				|
				v
			tui layer (render)
```

Layer responsibilities:

- proc layer
	- Reads files like /proc/stat, /proc/meminfo, /proc/net/dev.
	- Parses text into typed snapshots.
	- No UI logic. No presentation concerns.
- metrics layer
	- Converts snapshots into usable metrics.
	- Performs delta calculations across refresh intervals.
	- Encodes formulas for CPU percent, network rates, and utilization values.
- tui layer
	- Renders the latest computed state.
	- Handles layout, bars, sections, and refresh display.
	- Does not parse raw kernel files.

## How Metrics Work

### CPU (delta-based)

/proc/stat exposes cumulative counters since boot.
Absolute values are not instant usage.

You need two snapshots:

- previous counters
- current counters

Compute deltas:

- total_delta = current_total - previous_total
- idle_delta = current_idle - previous_idle
- busy_delta = total_delta - idle_delta

CPU percentage:

- cpu_percent = (busy_delta / total_delta) * 100

The same method applies to total CPU and each core.

### Memory (/proc/meminfo)

Preferred approach uses MemAvailable.

- total = MemTotal
- available = MemAvailable
- used = total - available
- used_percent = (used / total) * 100

Why MemAvailable:

- It better reflects reclaimable memory than just MemFree.
- It avoids over-reporting memory pressure.

### Network (rate-based)

/proc/net/dev exposes cumulative byte counters.

Use two snapshots separated by elapsed time:

- rx_rate = (rx_bytes_now - rx_bytes_prev) / seconds_elapsed
- tx_rate = (tx_bytes_now - tx_bytes_prev) / seconds_elapsed

This yields bytes/sec.
You can aggregate per-interface rates for host totals.

## Project Structure

```text
.
├── cmd/
│   └── sysmon/
│       └── main.go
├── internal/
│   ├── alerts/
│   ├── collector/
│   ├── config/
│   ├── metrics/
│   ├── proc/
│   └── ui/
├── tests/
├── Makefile
├── go.mod
└── README.md
```

Module responsibilities:

- cmd/sysmon
	- Entry point, startup wiring.
- internal/proc
	- Raw data readers and parsers for Linux system files.
- internal/metrics
	- Data models and metric computations.
- internal/collector
	- Refresh loop and snapshot collection orchestration.
- internal/ui
	- Terminal rendering and layout.
- internal/alerts
	- Alerting hooks and threshold logic.
- internal/config
	- Runtime configuration and defaults.
- tests
	- Integration-style and external-package tests.

## Design Principles

- Accuracy over approximation
	- Use kernel counters and mathematically correct deltas.
- Real-time delta-based metrics
	- Never treat cumulative counters as instant rates.
- Clear separation of concerns
	- Parsing, computation, and rendering stay independent.
- Minimal system overhead
	- Efficient parsing, low allocation pressure, and predictable loops.

## Build and Run

Requirements:

- Linux
- Go 1.24+

Common commands:

```bash
# run directly
make run

# build binary
make build

# run all tests
make test

# run all tests with go directly
go test ./... -v

# run only tests package
go test ./tests -v
```

## Testing Strategy

Focus tests on deterministic logic:

- Parser correctness for /proc text formats.
- Delta computations and edge cases.
- Mapping correctness for each counter field.

Recommended test types:

- Unit tests for parser and metric functions.
- Snapshot fixture tests using fixed sample inputs.
- Boundary tests for malformed or partial lines.

## Future Improvements

- Alerts
	- Threshold-based CPU, memory, disk, and network alerts.
- Logging
	- Structured runtime and diagnostic logs.
- Export metrics
	- JSON, CSV, or Prometheus-compatible output.
- Cross-platform support
	- Add abstraction layer for non-Linux providers.

## Extension Guide

When adding a new metric:

1. Add raw parser support in proc layer.
2. Extend snapshot/data models in metrics layer.
3. Add computation logic (including deltas if cumulative).
4. Render it in ui layer.
5. Add unit tests with fixed input snapshots.

If these steps stay isolated per layer, the project remains easy to evolve.