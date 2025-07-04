package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

const batteryFrameMinimumSize = 10
const batteriesFrameMinimumSize = batteryFrameMinimumSize + 1

type BatteryVoltages struct {
	AssyNo                      uint8    `json:"assamblyNo"`
	Voltage                     uint16   `json:"voltage"`
	Current                     uint16   `json:"current"`
	BatteriesTotalNumber        uint16   `json:"batteriesTotalNumber"`
	BatteryStartNumberFrame     uint16   `json:"batteryStartNumberFrame"`
	SingleBatteryVoltageOnFrame []uint16 `json:"singleBatteryVoltageOnFrame"`
}

func (m *BatteryVoltages) Decode(pkt []byte, idx *int) error {
	if ll := len(pkt[*idx:]); ll < batteryFrameMinimumSize {
		return ErrDecodeMsg
	}
	m.AssyNo = hex.ReadByte(pkt, idx)
	m.Voltage = hex.ReadWord(pkt, idx)
	m.Current = hex.ReadWord(pkt, idx)
	m.BatteriesTotalNumber = hex.ReadWord(pkt, idx)
	m.BatteryStartNumberFrame = hex.ReadWord(pkt, idx)
	bSize := hex.ReadByte(pkt, idx)
	if ll := len(pkt[*idx:]); ll < 2*int(bSize) {
		return ErrDecodeMsg
	}
	m.SingleBatteryVoltageOnFrame = make([]uint16, 0, bSize)
	for i := 0; i < int(bSize); i++ {
		batV := hex.ReadWord(pkt, idx)
		m.SingleBatteryVoltageOnFrame = append(m.SingleBatteryVoltageOnFrame, batV)
	}
	return nil
}

func (m *BatteryVoltages) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.AssyNo)
	pkt = hex.WriteWord(pkt, m.Voltage)
	pkt = hex.WriteWord(pkt, m.Current)
	pkt = hex.WriteWord(pkt, m.BatteriesTotalNumber)
	pkt = hex.WriteWord(pkt, m.BatteryStartNumberFrame)
	bSize := uint8(len(m.SingleBatteryVoltageOnFrame))
	pkt = hex.WriteByte(pkt, bSize)
	for _, v := range m.SingleBatteryVoltageOnFrame {
		pkt = hex.WriteWord(pkt, v)
	}
	return pkt, err
}

type BatteriesVoltages []*BatteryVoltages

func (m BatteriesVoltages) GetID() uint8 {
	return 0x08
}

func (m *BatteriesVoltages) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < batteriesFrameMinimumSize {
		return ErrDecodeMsg
	}
	total := hex.ReadByte(pkt, idx)
	expextedLen := (int(total) * batteryFrameMinimumSize)

	if pLen := len(pkt[*idx:]); pLen < expextedLen {
		return errors.Wrapf(ErrDecodeMsg, "expected len %d got %d", expextedLen, pLen)
	}
	// reset all
	*m = make([]*BatteryVoltages, 0)

	// decode
	for i := 0; i < int(total); i++ {
		dm := BatteryVoltages{}
		err := dm.Decode(pkt, idx)
		if err != nil {
			return err
		}
		*m = append(*m, &dm)
	}

	return nil
}

func (m BatteriesVoltages) Encode() (pkt []byte, err error) {
	total := len(m)
	pkt = hex.WriteByte(pkt, uint8(total))
	for _, v := range m {
		vp, _ := v.Encode()
		// DriveMotor.Encode never return error
		pkt = hex.WriteBytes(pkt, vp)

	}
	return pkt, err
}

func (m *BatteriesVoltages) MarshalJSON() ([]byte, error) {
	type Alias BatteriesVoltages
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
