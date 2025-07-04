package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const fuelCellMin = 21

type FuelCell struct {
	BatVoltage                     uint16  `json:"batVoltage"`
	BatCurrent                     uint16  `json:"batCurrent"`
	FuelRate                       uint16  `json:"fuelRate"`
	Temperatures                   []uint8 `json:"temperatures"`
	HydrogenSysMaxTemp             uint16  `json:"hydrogenSysMaxTemp"`
	HydrogenSysMaxTempNo           uint8   `json:"hydrogenSysMaxTempNo"`
	HydrogenSysMaxConcentrations   uint16  `json:"hydrogenSysMaxConcentrations"`
	HydrogenSysMaxConcentrationsNo uint8   `json:"hydrogenSysMaxConcentrationsNo"`
	HydrogenSysMaxPressure         uint16  `json:"hydrogenSysMaxPressure"`
	HydrogenSysMaxPressureNo       uint8   `json:"hydrogenSysMaxPressureNo"`
	DCStatus                       uint8   `json:"dcStatus"`
}

func (m *FuelCell) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < fuelCellMin {
		return ErrDecodeMsg
	}
	m.BatVoltage = hex.ReadWord(pkt, idx)
	m.BatCurrent = hex.ReadWord(pkt, idx)
	m.FuelRate = hex.ReadWord(pkt, idx)
	tempLen := hex.ReadByte(pkt, idx)
	m.Temperatures = hex.ReadBytes(pkt, idx, int(tempLen))
	m.HydrogenSysMaxTemp = hex.ReadWord(pkt, idx)
	m.HydrogenSysMaxTempNo = hex.ReadByte(pkt, idx)
	m.HydrogenSysMaxConcentrations = hex.ReadWord(pkt, idx)
	m.HydrogenSysMaxConcentrationsNo = hex.ReadByte(pkt, idx)
	m.HydrogenSysMaxPressure = hex.ReadWord(pkt, idx)
	m.HydrogenSysMaxPressureNo = hex.ReadByte(pkt, idx)
	m.DCStatus = hex.ReadByte(pkt, idx)
	return nil
}

func (m *FuelCell) Encode() (pkt []byte, err error) {
	pkt = hex.WriteWord(pkt, m.BatVoltage)
	pkt = hex.WriteWord(pkt, m.BatCurrent)
	pkt = hex.WriteWord(pkt, m.FuelRate)
	pkt = hex.WriteByte(pkt, uint8(len(m.Temperatures)))
	pkt = hex.WriteBytes(pkt, m.Temperatures)
	pkt = hex.WriteWord(pkt, m.HydrogenSysMaxTemp)
	pkt = hex.WriteByte(pkt, m.HydrogenSysMaxTempNo)
	pkt = hex.WriteWord(pkt, m.HydrogenSysMaxConcentrations)
	pkt = hex.WriteByte(pkt, m.HydrogenSysMaxConcentrationsNo)
	pkt = hex.WriteWord(pkt, m.HydrogenSysMaxPressure)
	pkt = hex.WriteByte(pkt, m.HydrogenSysMaxPressureNo)
	pkt = hex.WriteByte(pkt, m.DCStatus)
	return pkt, err
}

func (m FuelCell) GetID() uint8 {
	return 0x03
}

func (m *FuelCell) MarshalJSON() ([]byte, error) {
	type Alias FuelCell
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
