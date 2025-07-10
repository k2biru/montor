package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"
)

const (
	ICCIDLength = 20
)

// Login
type Msg01 struct {
	Header               *MsgHeader       `json:"header"`
	Time                 time.Time        `json:"time"`
	SerialNumber         uint16           `json:"serialNumber"`
	ICCID                string           `json:"iccid"`
	EnerygyStorageSystem EnergyStorageSys `json:"energyStorageSystem"`
}

func (m *Msg01) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	t, err := TimeDecode(pkt, &idx)
	if err != nil {
		t = time.Time{}
	}
	m.Time = t.UTC()
	m.SerialNumber = hex.ReadWord(pkt, &idx)
	m.ICCID = hex.ReadString(pkt, &idx, ICCIDLength)
	return m.EnerygyStorageSystem.Decode(pkt, &idx)
}

func (m *Msg01) Encode() (pkt []byte, err error) {
	t := TimeEncode(m.Time)
	pkt = hex.WriteBytes(pkt, t)
	pkt = hex.WriteWord(pkt, m.SerialNumber)
	pkt = hex.WriteString(pkt, m.ICCID, ICCIDLength)
	pkt = hex.WriteBytes(pkt, m.EnerygyStorageSystem.Encode())

	return WriteHeader(m, pkt)
}

func (m *Msg01) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg01) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}

func (m Msg01) GetMsgSN() string {
	return ""
}

func (m Msg01) GetTime() time.Time {
	return m.Time
}

type EnergyStorageSys struct {
	Coding byte

	// todo : create interface for eatch coding
	Raw []byte
}

func (m *EnergyStorageSys) Decode(pkt []byte, idx *int) error {
	if len(pkt) < *idx {
		return nil
	}
	size := hex.ReadByte(pkt, idx)
	m.Coding = hex.ReadByte(pkt, idx)
	codingSize := int(size * m.Coding)
	m.Raw = hex.ReadBytes(pkt, idx, codingSize)
	return nil
}

func (m *EnergyStorageSys) Encode() (pkt []byte) {
	if m.Coding == 0 {
		return pkt
	}
	size := uint8(len(m.Raw) / int(m.Coding))
	pkt = hex.WriteByte(pkt, size)
	pkt = hex.WriteByte(pkt, m.Coding)
	pkt = hex.WriteBytes(pkt, m.Raw)
	return pkt
}
