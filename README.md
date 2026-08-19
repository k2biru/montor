# montor

[![Go Reference](https://pkg.go.dev/badge/github.com/k2biru/montor.svg)](https://pkg.go.dev/github.com/k2biru/montor)
[![Go CI](https://github.com/k2biru/montor/actions/workflows/go.yml/badge.svg)](https://github.com/k2biru/montor/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k2biru/montor)](https://goreportcard.com/report/github.com/k2biru/montor)
[![Coverage Status](https://img.shields.io/badge/Coverage-98.6%25-brightgreen.svg)](coverage.html)

`montor` is a high-performance Go library and frame processing pipeline for the **GB/T 32960** protocol — the Chinese National Standard for Electric Vehicle (EV) telemetry and data communication (Technical Specifications for Technical Service of Electric Vehicle Telemetry).
Language / 语言 / Bahasa: [English](README.md) | [简体中文](README.zh-CN.md) | [Bahasa Indonesia](README.id.md)

---

## Key Features

- **Protocol Codec**: Full binary encoding/decoding for GB/T 32960 frame headers, payloads, checksums (BCC validation), and encrypted data frames.
- **Pipeline Engine**: Asynchronous reader/writer pipeline for TCP/Net stream socket handling between T-BOX terminals and vehicle telemetry platforms.
- **Data Models**: Supports all core GB/T 32960 message types (`Msg01` Real-time Report, `Msg02` Re-report, `Msg03` Login, `Msg04` Logout, `Msg07` Heartbeat, `Msg80` Parameter Query, `Msg81` Parameter Set, `Msg82` Control).
- **Sub-system Telemetry Parsers**: Strongly typed converters for Vehicle Status, Drive Motors, Fuel Cell, Engine, Location/GPS, Extremes, Alarms, Battery Subsystem Voltages, and Temperatures.
- **Hook-based Architecture**: Extensible hooks (`PipeHooks`, `FrameHandlerHooks`, `PacketCodecHooks`, `ProcesssHooks`) for custom encryption/decryption (e.g. RSA/AES), authentication, and custom process handlers.
- **Extensive Test Coverage**: **98.6% overall statement coverage** (100% statement coverage in `parser`).
- **Benchmarked**: Optimized zero-allocation/low-allocation code paths for high-throughput stream processing.

---

## Installation

```bash
go get github.com/k2biru/montor
```

---

## Core Components

```
                +------------------------------------+
                |              TCP Net               |
                +------------------------------------+
                                  |
                                  v
+------------------------------------------------------------------+
|                            Pipeline                              |
|                                                                  |
|  +------------------+   +------------------+   +---------------+ |
|  |   FrameHandler   |-->|   PacketCodec    |-->|   Processor   | |
|  |  (Stream Reader) |   | (Header / Check) |   | (Action Exec) | |
|  +------------------+   +------------------+   +---------------+ |
+------------------------------------------------------------------+
```

1. **`FrameHandler`**: Reads `0x23 0x23` (`##`) framed streams from `io.ReadWriter`.
2. **`PacketCodec`**: Verifies BCC checksums, decodes headers/encryption, and encodes outgoing responses.
3. **`Processor`**: Routes incoming command IDs to corresponding action logic and generates response frames.
4. **`parser`**: Converts raw BCD and binary offset values into standardized telemetry data types (speed in km/h, voltages in Volts, temperatures in °C).

---

## Usage Examples

### 1. Initializing Pipeline Server

```go
package main

import (
	"context"
	"net"
	"github.com/k2biru/montor"
	"github.com/k2biru/montor/models"
)

type MyHooks struct{}

func (h *MyHooks) GetProcess(id uint8) (*montor.Action, error) {
	return &montor.Action{
		GenData: func() *models.ProcessData {
			return &models.ProcessData{
				Incoming: &models.Msg07{},
				Outgoing: &models.Msg07{},
			}
		},
		Process: func(ctx context.Context, pd *models.ProcessData) error {
			// Handle telemetry message logic
			return nil
		},
	}, nil
}

func (h *MyHooks) PreProcess(ctx context.Context, msg models.GBT32960Msg) (context.Context, error) { return ctx, nil }
func (h *MyHooks) PostDecode(msg models.GBT32960Msg) {}
func (h *MyHooks) Decrypt(encType uint8, vin string, pkt []byte) ([]byte, error) { return pkt, nil }
func (h *MyHooks) PostRecvFrame(in []byte) ([]byte, error) { return in, nil }
func (h *MyHooks) PostSendFrame(pkt []byte) {}

func handleConnection(conn net.Conn) {
	hooks := &MyHooks{}
	pipe := montor.NewPipeline(conn, hooks)

	ctx := context.Background()
	for {
		if err := pipe.ProcessRead(ctx); err != nil {
			break
		}
	}
}
```

### 2. Using Sub-system Telemetry Parsers

```go
package main

import (
	"fmt"
	"github.com/k2biru/montor/parser"
)

func main() {
	// Parse Vehicle Status
	status := parser.VehicleStatus().SetVal(0x01).String()
	fmt.Println("Vehicle Status:", status) // "working"

	// Speed conversion (0.1 km/h per bit, offset 0)
	speed := parser.Speed().SetVal(1200).AsFloat()
	fmt.Println("Vehicle Speed:", speed, "km/h") // 120.0 km/h

	// Temperature conversion (1 °C per bit, offset -40)
	temp := parser.ExtremeMaxTempProbe().SetVal(65).AsInt()
	fmt.Println("Max Probe Temp:", temp, "°C") // 25 °C
}
```

---

## Testing & Benchmarks

### Running Unit Tests & Coverage Report

```bash
# Run all unit tests with statement coverage profile
go test -v -coverprofile=coverage.out ./...

# View function-level coverage summary
go tool cover -func=coverage.out

# Render visual HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

#### Test Coverage Summary

| Package | Statement Coverage | Status |
| :--- | :---: | :---: |
| `github.com/k2biru/montor/parser` | **100.0%** | Complete |
| `github.com/k2biru/montor/models` | **98.7%** | Near 100% |
| `github.com/k2biru/montor/codec/hex` | **98.0%** | Near 100% |
| `github.com/k2biru/montor` (Root) | **97.9%** | Near 100% |
| `github.com/k2biru/montor/codec/gbk` | **87.5%** | High |
| **Total Repository Coverage** | **98.6%** | **High Excellence** |

---

### Running Benchmark Suite

```bash
# Run all benchmarks with memory allocations
go test -bench=. -benchmem ./...
```

---

## Version History & Changelog

See [`CHANGELOG.md`](CHANGELOG.md) for detailed notes on each release tag (`v0.0.1` through `v0.1.4`).

---

## License

Distributed under the MIT License. See [`LICENSE`](LICENSE) for details.