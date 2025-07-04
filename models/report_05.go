package models

import (
	"encoding/json"

	"github.com/k2biru/montor/codec/hex"
)

const locationMin = 9

// type LocationStatus struct {
// 	Valid bool `json:"valid"`
// 	South bool `json:"south"`
// 	West  bool `json:"west"`
// }

// func (m *LocationStatus) Decode(data uint8) {
// 	m.Valid = getBit(data, 0)
// 	m.South = getBit(data, 1)
// 	m.West = getBit(data, 2)

// }
// func (m *LocationStatus) Encode() (num uint8) {
// 	if m.Valid {
// 		num = setBit(num, 0)
// 	}
// 	if m.South {
// 		num = setBit(num, 1)
// 	}
// 	if m.West {
// 		num = setBit(num, 2)
// 	}
// 	return num

// }

type Location struct {
	Status    uint8  `json:"status"`
	Longitude uint32 `json:"longitude"`
	Latidude  uint32 `json:"latidude"`
}

func (m *Location) Decode(pkt []byte, idx *int) error {
	if len(pkt[*idx:]) < locationMin {
		return ErrDecodeMsg
	}
	m.Status = hex.ReadByte(pkt, idx)
	m.Longitude = hex.ReadDoubleWord(pkt, idx)
	m.Latidude = hex.ReadDoubleWord(pkt, idx)
	return nil
}

func (m Location) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.Status)
	pkt = hex.WriteDoubleWord(pkt, m.Longitude)
	pkt = hex.WriteDoubleWord(pkt, m.Latidude)
	return pkt, err
}

func (m Location) GetID() uint8 {
	return 0x05
}

// func setBit(value uint8, bitPos uint) uint8 {
// 	return value | (1 << bitPos)
// }

// func getBit(value uint8, bitPos uint) bool {
// 	return (value & (1 << bitPos)) != 0
// }

func (m *Location) MarshalJSON() ([]byte, error) {
	type Alias Location
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
