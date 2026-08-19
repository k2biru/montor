# montor

[![Go Reference](https://pkg.go.dev/badge/github.com/k2biru/montor.svg)](https://pkg.go.dev/github.com/k2biru/montor)
[![Go CI](https://github.com/k2biru/montor/actions/workflows/go.yml/badge.svg)](https://github.com/k2biru/montor/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k2biru/montor)](https://goreportcard.com/report/github.com/k2biru/montor)
[![Coverage Status](https://img.shields.io/badge/Coverage-98.6%25-brightgreen.svg)](coverage.html)

`montor` 是一个高性能的 Go 语言 GB/T 32960 国家标准协议解析与数据处理流水线框架（电动汽车远程服务与管理系统通讯协议）。

语言 / Language / Bahasa: [English](README.md) | [简体中文](README.zh-CN.md) | [Bahasa Indonesia](README.id.md)

---

## 核心特性

- **协议编解码器 (Codec)**：完整支持 GB/T 32960 报文头、数据载荷、BCC 校验码生成与校验，以及加密报文处理。
- **流水线引擎 (Pipeline)**：提供基于 TCP/Net Socket 的异步数据流读写流水线，适用于车载终端 (T-BOX) 与企业/国家监管平台之间的数据通信。
- **数据模型 (Data Models)**：完整支持 GB/T 32960 核心命令字（`Msg01` 实时信息上报、`Msg02` 补发信息上报、`Msg03` 车辆登入、`Msg04` 车辆登出、`Msg07` 心跳、`Msg80` 参数查询、`Msg81` 参数设置、`Msg82` 终端控制）。
- **子系统解析器 (Parsers)**：强类型化转换器，支持整车数据、驱动电机数据、燃料电池数据、发动机数据、车辆位置/GPS数据、极值数据、报警数据、可充电储能子系统电压与温度数据。
- **Hook 挂钩架构**：扩展挂钩（`PipeHooks`、`FrameHandlerHooks`、`PacketCodecHooks`、`ProcesssHooks`），方便集成自定义加密/解密算法（如 RSA/AES）、鉴权以及自定义逻辑处理。
- **超高测试覆盖率**：**全项目语句覆盖率达 98.6%**（`parser` 包达到 100.0% 完整覆盖）。
- **性能基准测试 (Benchmark)**：经过高并发零/低内存分配路径优化，满足高吞吐量数据流处理需求。

---

## 安装

```bash
go get github.com/k2biru/montor
```

---

## 核心架构组件

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

1. **`FrameHandler`**：从 `io.ReadWriter` 流中读取以 `0x23 0x23` (`##`) 开头的起始帧。
2. **`PacketCodec`**：校验 BCC 异或校验码，解码报文头与加密载荷，并编码输出应答报文。
3. **`Processor`**：将接收到的命令 ID 分发给相对应的业务逻辑 action，并生成响应数据帧。
4. **`parser`**：将原始 BCD 码和带偏移量的二进制数据转换为标准物理量单位（如速度 km/h、电压 V、温度 °C）。

---

## 使用示例

### 1. 初始化 Pipeline 流水线服务端

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
			// 处理遥测消息业务逻辑
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

### 2. 使用子系统解析器 (Parsers)

```go
package main

import (
	"fmt"
	"github.com/k2biru/montor/parser"
)

func main() {
	// 解析车辆状态
	status := parser.VehicleStatus().SetVal(0x01).String()
	fmt.Println("车辆状态:", status) // "working"

	// 速度转换 (精度 0.1 km/h，偏移量 0)
	speed := parser.Speed().SetVal(1200).AsFloat()
	fmt.Println("车速:", speed, "km/h") // 120.0 km/h

	// 温度转换 (精度 1 °C，偏移量 -40)
	temp := parser.ExtremeMaxTempProbe().SetVal(65).AsInt()
	fmt.Println("最高温度探针:", temp, "°C") // 25 °C
}
```

---

## 测试与性能基准

### 运行单元测试与覆盖率报告

```bash
# 运行所有单元测试并生成覆盖率文件
go test -v -coverprofile=coverage.out ./...

# 查看函数级覆盖率统计
go tool cover -func=coverage.out

# 生成可视化 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

#### 测试覆盖率汇总

| 包名 | 语句覆盖率 | 状态 |
| :--- | :---: | :---: |
| `github.com/k2biru/montor/parser` | **100.0%** | 完整覆盖 |
| `github.com/k2biru/montor/models` | **98.7%** | 接近 100% |
| `github.com/k2biru/montor/codec/hex` | **98.0%** | 接近 100% |
| `github.com/k2biru/montor` (Root) | **97.9%** | 接近 100% |
| `github.com/k2biru/montor/codec/gbk` | **87.5%** | 优秀 |
| **全项目总覆盖率** | **98.6%** | **卓越** |

---

### 运行性能基准测试

```bash
# 运行所有基准测试并统计内存分配情况
go test -bench=. -benchmem ./...
```

---

## 版本历史与更新日志

详细更新日志与版本历史标签 (`v0.0.1` 至 `v0.1.4`) 请参见 [`CHANGELOG.md`](CHANGELOG.md)。

---

## 开源协议

本项目采用 MIT 协议开源。详见 [`LICENSE`](LICENSE)。
