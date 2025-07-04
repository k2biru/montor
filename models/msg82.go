package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

// Control
type Msg82 struct {
	Header  *MsgHeader `json:"header"`
	Time    time.Time  `json:"time"`
	Control Control    `json:"control"`
}

func (m *Msg82) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	return m.Control.Decode(pkt, &idx)
}

func (m *Msg82) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	control, _ := m.Control.Encode()
	// no error return
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "encode control id 0x%02x", m.Control.GetID())
	// }
	pkt = hex.WriteBytes(pkt, control)
	return WriteHeader(m, pkt)
}

func (m *Msg82) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg82) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg82) GetMsgSN() string {
	return ""
}
