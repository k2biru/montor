package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

// Parameter Set
type Msg81 struct {
	Header     *MsgHeader `json:"header"`
	Time       time.Time  `json:"time"`
	Parameters Parameters `json:"parameters"`
}

func (m *Msg81) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	_ = hex.ReadByte(pkt, &idx)
	return m.Parameters.Decode(pkt, &idx)
}

func (m *Msg81) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	param, err := m.Parameters.Encode()
	if err != nil {
		pkt = hex.WriteByte(pkt, 0)
		return pkt, err
	}
	pkt = hex.WriteByte(pkt, uint8(m.Parameters.length()))
	pkt = hex.WriteBytes(pkt, param)
	return WriteHeader(m, pkt)
}

func (m *Msg81) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg81) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg81) GetMsgSN() string {
	return ""
}

func (m Msg81) GetTime() time.Time {
	return m.Time
}
