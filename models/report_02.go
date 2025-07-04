package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

const driveMotorSize = 12
const driveMotorsMinimum = driveMotorSize + 1

type DriveMotor struct {
	SerialNumber         uint8  `json:"serialNumber"`
	Status               uint8  `json:"status"`
	ControlerTemperature uint8  `json:"controlerTemperature"`
	Speed                uint16 `json:"speed"`
	Torque               uint16 `json:"torque"`
	Temperature          uint8  `json:"temperature"`
	InputVoltage         uint16 `json:"inputVoltage"`
	InputCurrent         uint16 `json:"inputCurrent"`
}

func (m *DriveMotor) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < driveMotorSize {
		return ErrDecodeMsg
	}
	m.SerialNumber = hex.ReadByte(pkt, idx)
	m.Status = hex.ReadByte(pkt, idx)
	m.ControlerTemperature = hex.ReadByte(pkt, idx)
	m.Speed = hex.ReadWord(pkt, idx)
	m.Torque = hex.ReadWord(pkt, idx)
	m.Temperature = hex.ReadByte(pkt, idx)
	m.InputVoltage = hex.ReadWord(pkt, idx)
	m.InputCurrent = hex.ReadWord(pkt, idx)
	return nil
}

func (m *DriveMotor) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.SerialNumber)
	pkt = hex.WriteByte(pkt, m.Status)
	pkt = hex.WriteByte(pkt, m.ControlerTemperature)
	pkt = hex.WriteWord(pkt, m.Speed)
	pkt = hex.WriteWord(pkt, m.Torque)
	pkt = hex.WriteByte(pkt, m.Temperature)
	pkt = hex.WriteWord(pkt, m.InputVoltage)
	pkt = hex.WriteWord(pkt, m.InputCurrent)
	return pkt, err
}

type DriveMotors []*DriveMotor

func (m DriveMotors) GetID() uint8 {
	return 0x02
}

func (m *DriveMotors) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < driveMotorsMinimum {
		return ErrDecodeMsg
	}
	total := hex.ReadByte(pkt, idx)
	expextedLen := (int(total) * driveMotorSize)

	if pLen := len(pkt[*idx:]); pLen < expextedLen {
		return errors.Wrapf(ErrDecodeMsg, "expected len %d got %d", expextedLen, pLen)
	}
	// reset all
	*m = make([]*DriveMotor, 0)

	// decode
	for i := 0; i < int(total); i++ {
		dm := DriveMotor{}
		_ = dm.Decode(pkt, idx)
		// DriveMotor.Decode never return error because expextedLen hadbeen check
		// if  err != nil {
		// 	return err
		// }
		*m = append(*m, &dm)
	}

	return nil
}

func (m DriveMotors) Encode() (pkt []byte, err error) {
	total := len(m)
	pkt = hex.WriteByte(pkt, uint8(total))
	for _, v := range m {
		vp, _ := v.Encode()
		// DriveMotor.Encode never return error
		// if err != nil {
		// 	return nil, err
		// }
		pkt = hex.WriteBytes(pkt, vp)

	}
	return pkt, err
}

func (m *DriveMotors) MarshalJSON() ([]byte, error) {
	type Alias DriveMotors
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
