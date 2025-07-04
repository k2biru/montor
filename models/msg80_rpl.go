package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

// Parameter Reply
type Msg80Reply struct {
	Header     *MsgHeader `json:"header"`
	Time       time.Time  `json:"time"`
	Parameters Parameters `json:"parameters"`
}

func (m *Msg80Reply) Decode(packet *PacketData) error {
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

func (m *Msg80Reply) Encode() (pkt []byte, err error) {
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

func (m *Msg80Reply) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg80Reply) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg80Reply) GetMsgSN() string {
	return ""
}

func (m Msg80Reply) GetTime() time.Time {
	return m.Time
}
