package models

import (
	"reflect"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

type Parameters map[uint8]any

func (m *Parameters) Decode(pkt []byte, idx *int) error {
	//reset
	*m = make(map[uint8]any)

	for len(pkt[*idx:]) > 1 {
		// get type
		id := hex.ReadByte(pkt, idx)

		// match lookuptable
		lt, ok := paramLookup[id]
		if !ok {
			return errors.Wrapf(ErrDecodeMsg, "unknow report 0x%02x", id)
		}
		// generate item
		val := lt.generateValue()
		// decode item
		switch val.(type) {
		// case uint32: // no DWORD
		// 	val = hex.ReadDoubleWord(pkt, idx)
		case uint16:
			val = hex.ReadWord(pkt, idx)
		case uint8:
			val = hex.ReadByte(pkt, idx)
		case []uint8:
			length, ok := (*m)[lt.lenID].(uint8)
			if !ok {
				return errors.Wrapf(ErrDecodeMsg, "invalid id 0x%02x ", lt.lenID)
			}
			val = hex.ReadBytes(pkt, idx, int(length))
		case string:
			// read string using fixed length
			val = hex.ReadString(pkt, idx, int(lt.fixedLen))
			// default:
			// 	return errors.Wrapf(ErrDecodeMsg, " unknow type of 0x%02x", id)
		}
		(*m)[id] = val
	}
	return nil
}

func (m *Parameters) encodeItem(id uint8, v any) (pkt []byte, err error) {
	// match lookuptable
	lt, ok := paramLookup[id]
	if !ok {
		return nil, errors.Wrapf(ErrDecodeMsg, "unknow report 0x%02x", id)
	}
	var pktVal []byte
	switch vv := v.(type) {
	case uint16:
		pktVal = hex.WriteWord(pktVal, vv)
	case uint8:
		pktVal = hex.WriteByte(pktVal, vv)
	case []uint8:
		length := uint8(len(vv))
		pkt, _ = m.encodeItem(lt.lenID, length)

		pktVal = hex.WriteBytes(pktVal, vv)
	case string:
		// write string using fixed length
		pktVal = hex.WriteString(pktVal, vv, int(lt.fixedLen))

	default:
		return pkt, errors.Wrapf(ErrDecodeMsg, " unknow type of 0x%02x", id)
	}

	pkt = hex.WriteByte(pkt, id)
	pkt = hex.WriteBytes(pkt, pktVal)
	return pkt, nil
}
func (m *Parameters) deleteUnused() {
	// delete length key
	for _, v := range paramLookup {
		if v.lenID != 0 {
			delete(*m, v.lenID)
		}
	}
}

func (m Parameters) length() int {
	// delete length key
	m.deleteUnused()

	// length current
	length := len(m)
	for key := range m {
		if x, ok := paramLookup[key]; ok && x.lenID != 0x00 {
			length++
		}
	}
	return length
}

func (m Parameters) Encode() (pkt []byte, err error) {
	// delete length key
	m.deleteUnused()

	// sorting as slice
	kv := sortAscending(m)

	// encode every slice
	for _, v := range kv {
		p, err := m.encodeItem(v.Key, v.Value)
		if err != nil {
			return pkt, err
		}
		pkt = hex.WriteBytes(pkt, p)
	}
	return pkt, err
}

func (m Parameters) Add(id uint8, val any) error {
	tb, ok := paramLookup[id]
	if !ok {
		return ErrInvalidID
	}
	inputType := reflect.TypeOf(val)
	generatedType := reflect.TypeOf(tb.generateValue())
	if inputType != generatedType {
		return errors.Wrapf(ErrMissmatch,
			"at id 0x%02x expect %s got %s", id, generatedType, inputType)
	}
	m[id] = val
	return nil
}

func (m Parameters) Delete(id uint8) {
	delete(m, id)
}
func (m Parameters) IsEqual(other Parameters) error {
	if m == nil && other == nil {
		return nil
	}
	if m == nil || other == nil {
		return errors.Wrap(ErrInvalidConvert, "one of the maps is nil")
	}
	if len(m) != len(other) {
		return errors.Wrapf(ErrInvalidLength, "self=%d, other=%d", len(m), len(other))
	}
	for k, v1 := range m {
		v2, ok := other[k]
		if !ok {
			return errors.Wrapf(ErrInvalidID, "key %d is missing in the other map", k)
		}
		if !reflect.DeepEqual(v1, v2) {
			return errors.Wrapf(ErrMissmatch, "values for key %d are not equal: %v != %v", k, v1, v2)
		}
	}
	return nil
}

type GeneralParameter interface {
	Decode(pkt []byte, idx *int) error
	Encode() (pkt []byte, err error)
	GetID() uint8
}

type paramProperties struct {
	generateValue func() any
	lenID         uint8
	fixedLen      uint8
}

var paramLookup = map[uint8]paramProperties{
	0x01: { // waktu penyimpanan lokal
		generateValue: func() any { return uint16(0) },
	},
	0x02: { // waktu pelaporan default
		generateValue: func() any { return uint16(0) },
	},
	0x03: { // waktu pelaporan saat alarm
		generateValue: func() any { return uint16(0) },
	},
	0x04: { // panjang platform domain 0x05
		generateValue: func() any { return uint8(0) },
	},
	0x05: { // domain  platform
		generateValue: func() any { return []byte{} },
		lenID:         0x04,
	},
	0x06: { // port
		generateValue: func() any { return uint16(0) },
	},
	0x07: { // versi hw
		generateValue: func() any { return string("") },
		fixedLen:      5,
	},
	0x08: { // versi sw
		generateValue: func() any { return string("") },
		fixedLen:      5,
	},
	0x09: { // waktu HB
		generateValue: func() any { return uint8(0) },
	},
	0x0A: { // terminal respose timeout
		generateValue: func() any { return uint16(0) },
	},
	0x0B: { // platform response timeout
		generateValue: func() any { return uint16(0) },
	},
	0x0C: { // interval between login
		generateValue: func() any { return uint8(0) },
	},
	0x0D: { // panjang public domain 0x0F
		generateValue: func() any { return uint8(0) },
	},
	0x0E: { // domain public
		generateValue: func() any { return []byte{} },
		lenID:         0x0D,
	},
	0x0F: { // port public
		generateValue: func() any { return uint16(0) },
	},
	0x10: { // monitoring
		generateValue: func() any { return uint8(0) },
	},
	0x80: { // can report
		generateValue: func() any { return uint16(0) },
	},
	0x81: { // Whether the full CAN is uploaded, 0x01 indicates yes, 0x02 indicates no,
		generateValue: func() any { return uint8(0) },
	},
	0x86: { // tspID
		generateValue: func() any { return string("") },
		fixedLen:      32,
	},
}
