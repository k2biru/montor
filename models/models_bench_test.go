package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"
)

func BenchmarkHeader_Decode(b *testing.B) {
	raw := hex.Str2Byte("01FE313233343536373839000000000000000001002A")
	hdr := &MsgHeader{}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = hdr.Decode(raw)
	}
}

func BenchmarkHeader_Encode(b *testing.B) {
	hdr := GenerateHeader("12345678901234567", 0x01, 0xFE)
	hdr.BodyLength = 42

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hdr.Encode()
	}
}

func BenchmarkTime_Decode(b *testing.B) {
	raw := []byte{0x18, 0x01, 0x01, 0x0C, 0x1E, 0x2D} // 24-01-01 12:30:45

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		idx := 0
		_, _ = TimeDecode(raw, &idx)
	}
}

func BenchmarkTime_Encode(b *testing.B) {
	now := time.Date(2024, 1, 1, 12, 30, 45, 0, GBT32960Timezone()).UTC()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TimeEncode(now)
	}
}

func BenchmarkMsg01_Encode(b *testing.B) {
	msg := &Msg01{
		Header: &MsgHeader{
			CommandID:  0x01,
			Response:   0xFE,
			VIN:        "123456789",
			Encription: 0x01,
		},
		Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
		SerialNumber: 3,
		ICCID:        "89510079518932305625",
		EnerygyStorageSystem: EnergyStorageSys{
			Coding: 1,
			Raw:    []byte{0x02},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = msg.Encode()
	}
}

func BenchmarkMsg01_Decode(b *testing.B) {
	header := GenerateHeader("123456789", 0x01, 0xFE)
	body := hex.Str2Byte("18010102030400033839353130303739353138393332333035363235010102")
	pkt := &PacketData{
		Header: header,
		Body:   body,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		msg := &Msg01{}
		_ = msg.Decode(pkt)
	}
}
