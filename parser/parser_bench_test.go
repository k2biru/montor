package parser

import (
	"testing"
)

func BenchmarkConvert_Convert(b *testing.B) {
	conv := NewConvert("speed", uint16(0), 0.1, 0, 2200, false)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		conv.Convert(120)
	}
}

func BenchmarkConvert_Calculate(b *testing.B) {
	conv := NewConvert("speed", uint16(1200), 0.1, 0, 2200, false)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = conv.Calculate()
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	lookup := map[uint8]string{
		1: "working",
		2: "stopped",
	}
	p := NewParser("state", uint8(1), lookup)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Parse("stopped")
	}
}

func BenchmarkParser_String(b *testing.B) {
	p := VehicleStatus()
	p.SetVal(0x01)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.String()
	}
}

func BenchmarkGear_GetGear(b *testing.B) {
	g := Gear()
	g.SetVal(0b1110) // drive
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.GetGear()
	}
}

func BenchmarkSpeed_Calculate(b *testing.B) {
	sp := Speed()
	sp.SetVal(1200) // 120 km/h
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sp.Calculate().AsInt()
	}
}
