package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

type Msg04 struct {
	Header       *MsgHeader `json:"header"`
	Time         time.Time  `json:"time"`
	SerialNumber uint16     `json:"serialNumber"`
}

func (m *Msg04) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	m.SerialNumber = hex.ReadWord(pkt, &idx)
	return nil
}

func (m *Msg04) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	pkt = hex.WriteWord(pkt, m.SerialNumber)

	return WriteHeader(m, pkt)
}

func (m *Msg04) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg04) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg04) GetMsgSN() string {
	return ""
}

func (m Msg04) GetTime() time.Time {
	return m.Time
}
