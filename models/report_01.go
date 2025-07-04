package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const vehicleDataMin = 20

type VehicleData struct {
	Status       uint8  `json:"status"`
	Charging     uint8  `json:"charging"`
	OpMode       uint8  `json:"operatingMode"`
	Speed        uint16 `json:"speed"`
	Odometer     uint32 `json:"odometer"`
	TotalVoltage uint16 `json:"totalVoltage"`
	TotalCurrent uint16 `json:"totalCurrent"`
	SoC          uint8  `json:"stateOfCharge"`
	DCDCStatus   uint8  `json:"dcDcStatus"`
	Gear         uint8  `json:"gear"`
	Insulator    uint16 `json:"insulator"`
	Throttle     uint8  `json:"trottle"`
	Brake        uint8  `json:"brake"`
}

func (m *VehicleData) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < vehicleDataMin {
		return ErrDecodeMsg
	}
	m.Status = hex.ReadByte(pkt, idx)
	m.Charging = hex.ReadByte(pkt, idx)
	m.OpMode = hex.ReadByte(pkt, idx)
	m.Speed = hex.ReadWord(pkt, idx)
	m.Odometer = hex.ReadDoubleWord(pkt, idx)
	m.TotalVoltage = hex.ReadWord(pkt, idx)
	m.TotalCurrent = hex.ReadWord(pkt, idx)
	m.SoC = hex.ReadByte(pkt, idx)
	m.DCDCStatus = hex.ReadByte(pkt, idx)
	m.Gear = hex.ReadByte(pkt, idx)
	m.Insulator = hex.ReadWord(pkt, idx)
	m.Throttle = hex.ReadByte(pkt, idx)
	m.Brake = hex.ReadByte(pkt, idx)
	return nil
}

func (m *VehicleData) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.Status)
	pkt = hex.WriteByte(pkt, m.Charging)
	pkt = hex.WriteByte(pkt, m.OpMode)
	pkt = hex.WriteWord(pkt, m.Speed)
	pkt = hex.WriteDoubleWord(pkt, m.Odometer)
	pkt = hex.WriteWord(pkt, m.TotalVoltage)
	pkt = hex.WriteWord(pkt, m.TotalCurrent)
	pkt = hex.WriteByte(pkt, m.SoC)
	pkt = hex.WriteByte(pkt, m.DCDCStatus)
	pkt = hex.WriteByte(pkt, m.Gear)
	pkt = hex.WriteWord(pkt, m.Insulator)
	pkt = hex.WriteByte(pkt, m.Throttle)
	pkt = hex.WriteByte(pkt, m.Brake)
	return pkt, err
}

func (m VehicleData) GetID() uint8 {
	return 0x01
}

func (m *VehicleData) MarshalJSON() ([]byte, error) {
	type Alias VehicleData
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
