# montor

[![Go Reference](https://pkg.go.dev/badge/github.com/k2biru/montor.svg)](https://pkg.go.dev/github.com/k2biru/montor)
[![Go CI](https://github.com/k2biru/montor/actions/workflows/go.yml/badge.svg)](https://github.com/k2biru/montor/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k2biru/montor)](https://goreportcard.com/report/github.com/k2biru/montor)
[![Coverage Status](https://img.shields.io/badge/Coverage-98.6%25-brightgreen.svg)](coverage.html)

`montor` adalah pustaka Go berkinerja tinggi dan mesin pipeline pemrosesan bingkai untuk protokol **GB/T 32960** — Standar Nasional Tiongkok untuk telemetri dan komunikasi data Kendaraan Listrik (Electric Vehicle / EV).

Bahasa / Language / 语言: [English](README.md) | [Bahasa Indonesia](README.id.md) | [简体中文](README.zh-CN.md)

---

## Fitur Utama

- **Kodek Protokol (Codec)**: Enkoding/dekoding biner lengkap untuk header bingkai GB/T 32960, muatan (payload), validasi kode pemeriksaan checksum (BCC), dan bingkai terenkripsi.
- **Mesin Pipeline**: Pipeline pembaca/penulis asinkron untuk penanganan koneksi socket TCP/Net antara terminal T-BOX dan platform telemetri kendaraan.
- **Model Data (Data Models)**: Mendukung seluruh tipe pesan utama GB/T 32960 (`Msg01` Laporan Real-time, `Msg02` Laporan Ulang, `Msg03` Login, `Msg04` Logout, `Msg07` Heartbeat, `Msg80` Kueri Parameter, `Msg81` Pengaturan Parameter, `Msg82` Kontrol).
- **Parser Telemetri Sub-sistem**: Konverter bertipe kuat untuk Status Kendaraan, Motor Penggerak (Drive Motor), Sel Bahan Bakar (Fuel Cell), Mesin, Lokasi/GPS, Nilai Ekstrem, Alarm, serta Tegangan & Suhu Sub-sistem Baterai.
- **Arsitektur Berbasis Hook**: Hook yang dapat diperluas (`PipeHooks`, `FrameHandlerHooks`, `PacketCodecHooks`, `ProcesssHooks`) untuk enkripsi/dekripsi kustom (seperti RSA/AES), otentikasi, dan penanganan proses logika kustom.
- **Cakupan Pengujian Luas**: **98.6% cakupan pernyataan (statement coverage) secara keseluruhan** (100.0% cakupan pernyataan pada paket `parser`).
- **Uji Benchmark**: Jalur kode yang dioptimalkan dengan alokasi memori rendah/nol untuk pemrosesan aliran data berkecepatan tinggi (high throughput).

---

## Instalasi

```bash
go get github.com/k2biru/montor
```

---

## Komponen Utama

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

1. **`FrameHandler`**: Membaca aliran bingkai yang diawali `0x23 0x23` (`##`) dari `io.ReadWriter`.
2. **`PacketCodec`**: Memverifikasi checksum BCC, mendekode header/enkripsi, dan mengenkode bingkai balasan keluar.
3. **`Processor`**: Mengarahkan ID perintah yang masuk ke logika aksi yang sesuai dan menghasilkan bingkai respons.
4. **`parser`**: Mengonversi nilai mentah BCD dan offset biner menjadi satuan data telemetri standar (kecepatan dalam km/h, tegangan dalam Volt, suhu dalam °C).

---

## Contoh Penggunaan

### 1. Inisialisasi Server Pipeline

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
			// Logika penanganan pesan telemetri
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

### 2. Menggunakan Parser Telemetri Sub-sistem

```go
package main

import (
	"fmt"
	"github.com/k2biru/montor/parser"
)

func main() {
	// Membaca Status Kendaraan
	status := parser.VehicleStatus().SetVal(0x01).String()
	fmt.Println("Status Kendaraan:", status) // "working"

	// Konversi Kecepatan (presisi 0.1 km/h, offset 0)
	speed := parser.Speed().SetVal(1200).AsFloat()
	fmt.Println("Kecepatan Kendaraan:", speed, "km/h") // 120.0 km/h

	// Konversi Suhu (presisi 1 °C, offset -40)
	temp := parser.ExtremeMaxTempProbe().SetVal(65).AsInt()
	fmt.Println("Suhu Sensor Maksimum:", temp, "°C") // 25 °C
}
```

---

## Pengujian & Benchmark

### Menjalankan Uji Unit & Laporan Coverage

```bash
# Jalankan seluruh uji unit dengan profil coverage pernyataan
go test -v -coverprofile=coverage.out ./...

# Tampilkan ringkasan coverage per fungsi
go tool cover -func=coverage.out

# Buat laporan HTML visual coverage
go tool cover -html=coverage.out -o coverage.html
```

#### Ringkasan Cakupan Pengujian (Test Coverage)

| Paket | Cakupan Pernyataan (Coverage) | Status |
| :--- | :---: | :---: |
| `github.com/k2biru/montor/parser` | **100.0%** | Sempurna |
| `github.com/k2biru/montor/models` | **98.7%** | Mendekati 100% |
| `github.com/k2biru/montor/codec/hex` | **98.0%** | Mendekati 100% |
| `github.com/k2biru/montor` (Root) | **97.9%** | Mendekati 100% |
| `github.com/k2biru/montor/codec/gbk` | **87.5%** | Tinggi |
| **Total Cakupan Repositori** | **98.6%** | **Sangat Baik** |

---

### Menjalankan Suite Benchmark

```bash
# Jalankan seluruh uji benchmark beserta lacak alokasi memori
go test -bench=. -benchmem ./...
```

---

## Riwayat Versi & Catatan Perubahan

Lihat [`CHANGELOG.md`](CHANGELOG.md) untuk catatan rinci setiap tag rilis (`v0.0.1` hingga `v0.1.4`).

---

## Lisensi

Didistribusikan di bawah Lisensi MIT. Lihat [`LICENSE`](LICENSE) untuk informasi lebih lanjut.
