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
		val := lt.GenerateValue()
		// decode item
		switch val.(type) {
		// case uint32: // no DWORD
		// 	val = hex.ReadDoubleWord(pkt, idx)
		case uint16:
			val = hex.ReadWord(pkt, idx)
		case uint8:
			val = hex.ReadByte(pkt, idx)
		case []uint8:
			length, ok := (*m)[lt.LenID].(uint8)
			if !ok {
				return errors.Wrapf(ErrDecodeMsg, "invalid id 0x%02x ", lt.LenID)
			}
			val = hex.ReadBytes(pkt, idx, int(length))
		case string:
			// read string using fixed length
			val = hex.ReadString(pkt, idx, int(lt.FixedLen))
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
		pkt, _ = m.encodeItem(lt.LenID, length)

		pktVal = hex.WriteBytes(pktVal, vv)
	case string:
		// write string using fixed length
		pktVal = hex.WriteString(pktVal, vv, int(lt.FixedLen))

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
		if v.LenID != 0 {
			delete(*m, v.LenID)
		}
	}
}

func (m Parameters) length() int {
	// delete length key
	m.deleteUnused()

	// length current
	length := len(m)
	for key := range m {
		if x, ok := paramLookup[key]; ok && x.LenID != 0x00 {
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
	generatedType := reflect.TypeOf(tb.GenerateValue())
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

type ParamProperties struct {
	GenerateValue func() any // value generator return []byte{}, string(""), uint8, or uint16
	LenID         uint8      // length of value from another id for string / []byte value
	FixedLen      uint8      // fixed length of value for string / []byte value
}

var paramLookup = map[uint8]ParamProperties{
	0x01: { // waktu penyimpanan lokal
		GenerateValue: func() any { return uint16(0) },
	},
	0x02: { // waktu pelaporan default
		GenerateValue: func() any { return uint16(0) },
	},
	0x03: { // waktu pelaporan saat alarm
		GenerateValue: func() any { return uint16(0) },
	},
	0x04: { // panjang platform domain 0x05
		GenerateValue: func() any { return uint8(0) },
	},
	0x05: { // domain  platform
		GenerateValue: func() any { return []byte{} },
		LenID:         0x04,
	},
	0x06: { // port
		GenerateValue: func() any { return uint16(0) },
	},
	0x07: { // versi hw
		GenerateValue: func() any { return string("") },
		FixedLen:      5,
	},
	0x08: { // versi sw
		GenerateValue: func() any { return string("") },
		FixedLen:      5,
	},
	0x09: { // waktu HB
		GenerateValue: func() any { return uint8(0) },
	},
	0x0A: { // terminal respose timeout
		GenerateValue: func() any { return uint16(0) },
	},
	0x0B: { // platform response timeout
		GenerateValue: func() any { return uint16(0) },
	},
	0x0C: { // interval between login
		GenerateValue: func() any { return uint8(0) },
	},
	0x0D: { // panjang public domain 0x0F
		GenerateValue: func() any { return uint8(0) },
	},
	0x0E: { // domain public
		GenerateValue: func() any { return []byte{} },
		LenID:         0x0D,
	},
	0x0F: { // port public
		GenerateValue: func() any { return uint16(0) },
	},
	0x10: { // monitoring
		GenerateValue: func() any { return uint8(0) },
	},
}

// set custom parameter lookup for param decoder
func SetParameterPropertiesLookup(id uint8, properties ParamProperties) {
	paramLookup[id] = properties
}
