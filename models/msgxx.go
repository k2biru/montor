package models

import (
	"github.com/k2biru/montor/codec/hex"
)

type Msgxx struct {
	Header *MsgHeader `json:"header"`
	Raw    []byte     `json:"raw"`
}

func (m *Msgxx) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	m.Raw = hex.ReadBytes(pkt, &idx, int(m.Header.BodyLength))
	return nil
}

func (m *Msgxx) Encode() (pkt []byte, err error) {
	pkt = hex.WriteBytes(pkt, m.Raw)
	return WriteHeader(m, pkt)
}

func (m *Msgxx) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msgxx) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}

func (m Msgxx) GetMsgSN() string {
	return ""
}
