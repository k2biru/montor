package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

type Msg02 struct {
	Header  *MsgHeader `json:"header"`
	Time    time.Time  `json:"time"`
	Reports Reports    `json:"reports"`
}

func (m *Msg02) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	return m.Reports.Decode(pkt, &idx)
}

func (m *Msg02) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	report, _ := m.Reports.Encode()
	// report encode never error
	// if err != nil {
	// 	return nil, err
	// }
	pkt = hex.WriteBytes(pkt, report)

	return WriteHeader(m, pkt)
}

func (m *Msg02) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg02) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg02) GetMsgSN() string {
	return ""
}

func (m Msg02) GetTime() time.Time {
	return m.Time
}
