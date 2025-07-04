package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const engineMin = 5

type Engine struct {
	Status   uint8  `json:"status"`
	Revs     uint16 `json:"revolution"`
	FuelRate uint16 `json:"fuelRate"`
}

func (m *Engine) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < engineMin {
		return ErrDecodeMsg
	}
	m.Status = hex.ReadByte(pkt, idx)
	m.Revs = hex.ReadWord(pkt, idx)
	m.FuelRate = hex.ReadWord(pkt, idx)
	return nil
}

func (m *Engine) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.Status)
	pkt = hex.WriteWord(pkt, m.Revs)
	pkt = hex.WriteWord(pkt, m.FuelRate)
	return pkt, err
}

func (m Engine) GetID() uint8 {
	return 0x04
}

func (m *Engine) MarshalJSON() ([]byte, error) {
	type Alias Engine
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
