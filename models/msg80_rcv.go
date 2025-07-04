package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

func GenerateMsh80Receive(vin string, timeNow time.Time, parameterID []byte) *Msg80Receive {
	return &Msg80Receive{
		Header:       GenerateHeader(vin, 0x80, 0xFE),
		Time:         timeNow,
		ParameterIDs: parameterID,
	}
}

// Parameter Query
type Msg80Receive struct {
	Header       *MsgHeader `json:"header"`
	Time         time.Time  `json:"time"`
	ParameterIDs []byte     `json:"parameterIDs"`
}

func (m *Msg80Receive) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	idLen := hex.ReadByte(pkt, &idx)
	m.ParameterIDs = hex.ReadBytes(pkt, &idx, int(idLen))
	return nil
}

func (m *Msg80Receive) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	pkt = hex.WriteByte(pkt, uint8(len(m.ParameterIDs)))
	pkt = hex.WriteBytes(pkt, m.ParameterIDs)
	return WriteHeader(m, pkt)
}

func (m *Msg80Receive) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg80Receive) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg80Receive) GetMsgSN() string {
	return ""
}
