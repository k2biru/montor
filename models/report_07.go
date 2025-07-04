package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const alarmMin = 14

type Alarm struct {
	AlarmLevel         uint8    `json:"alarmLevel"`
	AlarmBatteryFlag   uint32   `json:"alarmBatteryFlag"`
	AlarmBatteryOthers []uint32 `json:"alarmBatteryOthers"`
	AlarmDriveMotor    []uint32 `json:"alarmDriveMotor"`
	AlarmEngines       []uint32 `json:"alarmEngines"`
	AlarmOthers        []uint32 `json:"alarmOthers"`
}

func (m *Alarm) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < alarmMin {
		return ErrDecodeMsg
	}
	m.AlarmLevel = hex.ReadByte(pkt, idx)
	m.AlarmBatteryFlag = hex.ReadDoubleWord(pkt, idx)
	aLen := hex.ReadByte(pkt, idx)
	if len(pkt[*idx:]) < (int(aLen)*4 + 1) {
		return ErrDecodeMsg
	}
	for i := 0; i < int(aLen); i++ {
		m.AlarmBatteryOthers = append(m.AlarmBatteryOthers, hex.ReadDoubleWord(pkt, idx))
	}
	aLen = hex.ReadByte(pkt, idx)
	if len(pkt[*idx:]) < (int(aLen)*4 + 1) {
		return ErrDecodeMsg
	}
	for i := 0; i < int(aLen); i++ {
		m.AlarmDriveMotor = append(m.AlarmDriveMotor, hex.ReadDoubleWord(pkt, idx))
	}
	aLen = hex.ReadByte(pkt, idx)
	if len(pkt[*idx:]) < (int(aLen)*4 + 1) {
		return ErrDecodeMsg
	}
	for i := 0; i < int(aLen); i++ {
		m.AlarmEngines = append(m.AlarmEngines, hex.ReadDoubleWord(pkt, idx))
	}
	aLen = hex.ReadByte(pkt, idx)
	if len(pkt[*idx:]) < (int(aLen) * 4) {
		return ErrDecodeMsg
	}
	for i := 0; i < int(aLen); i++ {
		m.AlarmOthers = append(m.AlarmOthers, hex.ReadDoubleWord(pkt, idx))
	}
	return nil
}

func (m Alarm) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.AlarmLevel)
	pkt = hex.WriteDoubleWord(pkt, m.AlarmBatteryFlag)
	pkt = hex.WriteByte(pkt, uint8(len(m.AlarmBatteryOthers)))
	for _, v := range m.AlarmBatteryOthers {
		pkt = hex.WriteDoubleWord(pkt, v)
	}

	pkt = hex.WriteByte(pkt, uint8(len(m.AlarmDriveMotor)))
	for _, v := range m.AlarmDriveMotor {
		pkt = hex.WriteDoubleWord(pkt, v)
	}

	pkt = hex.WriteByte(pkt, uint8(len(m.AlarmEngines)))
	for _, v := range m.AlarmEngines {
		pkt = hex.WriteDoubleWord(pkt, v)
	}

	pkt = hex.WriteByte(pkt, uint8(len(m.AlarmOthers)))
	for _, v := range m.AlarmOthers {
		pkt = hex.WriteDoubleWord(pkt, v)
	}
	return pkt, err
}

func (m Alarm) GetID() uint8 {
	return 0x07
}

func (m *Alarm) MarshalJSON() ([]byte, error) {
	type Alias Alarm
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
