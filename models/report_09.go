package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

const batteryTemperatureMinimumSize = 2
const batteryTemperaturesMinimumSize = batteryTemperatureMinimumSize + 1

type BatteryTemperature struct {
	AssyNo      uint8   `json:"assamblyNo                  "`
	Temperature []uint8 `json:"temperature                 "`
}

func (m *BatteryTemperature) Decode(pkt []byte, idx *int) error {
	if ll := len(pkt[*idx:]); ll < batteryTemperatureMinimumSize {
		return ErrDecodeMsg
	}
	m.AssyNo = hex.ReadByte(pkt, idx)
	tSize := hex.ReadWord(pkt, idx)
	if ll := len(pkt[*idx:]); ll < 1*int(tSize) {
		return ErrDecodeMsg
	}
	m.Temperature = make([]uint8, 0, tSize)
	for i := 0; i < int(tSize); i++ {
		temp := hex.ReadByte(pkt, idx)
		m.Temperature = append(m.Temperature, temp)
	}
	return nil
}

func (m *BatteryTemperature) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.AssyNo)
	tSize := uint16(len(m.Temperature))
	pkt = hex.WriteWord(pkt, tSize)
	for _, v := range m.Temperature {
		pkt = hex.WriteByte(pkt, v)

	}
	return pkt, err
}

type BatteriesTemperatures []*BatteryTemperature

func (m BatteriesTemperatures) GetID() uint8 {
	return 0x09
}

func (m *BatteriesTemperatures) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < batteryTemperaturesMinimumSize {
		return ErrDecodeMsg
	}
	total := hex.ReadByte(pkt, idx)
	expextedLen := (int(total) * batteryTemperatureMinimumSize)

	if pLen := len(pkt[*idx:]); pLen < expextedLen {
		return errors.Wrapf(ErrDecodeMsg, "expected len %d got %d", expextedLen, pLen)
	}
	// reset all
	*m = make([]*BatteryTemperature, 0)

	// decode
	for i := 0; i < int(total); i++ {
		dm := BatteryTemperature{}
		err := dm.Decode(pkt, idx)
		if err != nil {
			return err
		}
		*m = append(*m, &dm)
	}

	return nil
}

func (m BatteriesTemperatures) Encode() (pkt []byte, err error) {
	total := len(m)
	pkt = hex.WriteByte(pkt, uint8(total))
	for _, v := range m {
		vp, _ := v.Encode()
		// DriveMotor.Encode never return error
		pkt = hex.WriteBytes(pkt, vp)

	}
	return pkt, err
}

func (m *BatteriesTemperatures) MarshalJSON() ([]byte, error) {
	type Alias BatteriesTemperatures
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
