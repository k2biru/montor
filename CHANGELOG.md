# Changelog

All notable changes to the `montor` GB/T 32960 Go library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.1.4] - 2026-08-19

### Added
- **Table-Driven Unit Test Suite**:
  - Implemented comprehensive table-driven unit tests across all packages (`montor`, `codec`, `models`, `parser`).
  - Added dedicated interface and helper test suite (`models/models_coverage_test.go`, `parser/parser_test.go`, `pipeline_test.go`, `packet_codec_test.go`, `processor_test.go`).
  - Achieved **98.6% overall repository statement coverage** (**100.0% coverage in `parser`**).
- **Performance Benchmarks**:
  - Implemented benchmark suites measuring latency and memory allocation (`benchmark_test.go`, `models/models_bench_test.go`, `parser/parser_bench_test.go`).
- **Documentation**:
  - Updated `README.md` with complete architecture diagrams, core component explanations, socket server pipeline usage, telemetry parser examples, and test coverage metrics.


---

## [v0.1.3] - 2026-04-07

### Fixed
- Fixed parameter map deep copy in `Msg81.Copy()` and `Parameters` manipulation.

---

## [v0.1.2] - 2025-09-01

### Fixed
- Fixed parsing overflow handling when converter values exceed max limits (`fix parse over max`).
- Refactored and cleaned up `Processor.Process` flow.

---

## [v0.1.1] - 2025-07-10

### Fixed
- Fixed crash caused by decoding empty `EnergyStorageSys` data in `Msg01`.

---

## [v0.1.0] - 2025-07-10

### Added
- Added custom parameter lookup registration via `SetParameterPropertiesLookup`.
- Added customizable pipeline hooks (`PipeHooks`, `FrameHandlerHooks`, `PacketCodecHooks`, `ProcesssHooks`).

### Fixed
- Fixed stream EOF/EOL error handling in `FrameHandler` that caused high CPU usage under disconnected socket conditions.

---

## [v0.0.1] - 2025-07-04

### Added
- Initial release of GB/T 32960 EV protocol parser and processing engine.
