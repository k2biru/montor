package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

type MsgReply struct {
	Header *MsgHeader `json:"header"`
	Time   time.Time  `json:"time"`
}

func (m *MsgReply) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	return nil
}

func (m *MsgReply) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	return WriteHeader(m, pkt)
}

func (m *MsgReply) GetHeader() *MsgHeader {
	return m.Header
}

func (m *MsgReply) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m MsgReply) GetMsgSN() string {
	return ""
}
