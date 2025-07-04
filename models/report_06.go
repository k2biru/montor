package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const extremeMin = 14

type Extreme struct {
	MaxVoltageBatAssyNo      uint8  `json:"maxVoltageBatAssyNo"`
	MaxVoltageSingleBatNo    uint8  `json:"maxVoltageSingleBatNo"`
	MaxVoltageSingleBatValue uint16 `json:"maxVoltageSingleBatValue"`
	MinVoltageBatAssyNo      uint8  `json:"minVoltageBatAssyNo"`
	MinVoltageSingleBatNo    uint8  `json:"minVoltageSingleBatNo"`
	MinVoltageSingleBatValue uint16 `json:"minVoltageSingleBatValue"`
	MaxTempBatProbeNo        uint8  `json:"maxTempBatProbeNo"`
	MaxTempBatAssyNo         uint8  `json:"maxTempBatAssyNo"`
	MaxTempBatProbeValue     uint8  `json:"maxTempBatProbeValue"`
	MinTempBatAssyNo         uint8  `json:"minTempBatAssyNo"`
	MinTempBatProbeNo        uint8  `json:"minTempBatProbeNo"`
	MinTempBatProbeValue     uint8  `json:"minTempBatProbeValue"`
}

func (m *Extreme) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < extremeMin {
		return ErrDecodeMsg
	}
	m.MaxVoltageBatAssyNo = hex.ReadByte(pkt, idx)
	m.MaxVoltageSingleBatNo = hex.ReadByte(pkt, idx)
	m.MaxVoltageSingleBatValue = hex.ReadWord(pkt, idx)
	m.MinVoltageBatAssyNo = hex.ReadByte(pkt, idx)
	m.MinVoltageSingleBatNo = hex.ReadByte(pkt, idx)
	m.MinVoltageSingleBatValue = hex.ReadWord(pkt, idx)
	m.MaxTempBatProbeNo = hex.ReadByte(pkt, idx)
	m.MaxTempBatAssyNo = hex.ReadByte(pkt, idx)
	m.MaxTempBatProbeValue = hex.ReadByte(pkt, idx)
	m.MinTempBatAssyNo = hex.ReadByte(pkt, idx)
	m.MinTempBatProbeNo = hex.ReadByte(pkt, idx)
	m.MinTempBatProbeValue = hex.ReadByte(pkt, idx)
	return nil
}

func (m Extreme) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.MaxVoltageBatAssyNo)
	pkt = hex.WriteByte(pkt, m.MaxVoltageSingleBatNo)
	pkt = hex.WriteWord(pkt, m.MaxVoltageSingleBatValue)
	pkt = hex.WriteByte(pkt, m.MinVoltageBatAssyNo)
	pkt = hex.WriteByte(pkt, m.MinVoltageSingleBatNo)
	pkt = hex.WriteWord(pkt, m.MinVoltageSingleBatValue)
	pkt = hex.WriteByte(pkt, m.MaxTempBatProbeNo)
	pkt = hex.WriteByte(pkt, m.MaxTempBatAssyNo)
	pkt = hex.WriteByte(pkt, m.MaxTempBatProbeValue)
	pkt = hex.WriteByte(pkt, m.MinTempBatAssyNo)
	pkt = hex.WriteByte(pkt, m.MinTempBatProbeNo)
	pkt = hex.WriteByte(pkt, m.MinTempBatProbeValue)
	return pkt, err
}

func (m Extreme) GetID() uint8 {
	return 0x06
}

func (m *Extreme) MarshalJSON() ([]byte, error) {
	type Alias Extreme
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
