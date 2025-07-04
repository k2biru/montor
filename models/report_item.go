package models

import (
	"encoding/json"
	"strconv"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

type GeneralReport interface {
	Decode(pkt []byte, idx *int) error
	Encode() (pkt []byte, err error)
	GetID() uint8
	MarshalJSON() ([]byte, error)
}
type reportProperties struct {
	generateValue func() GeneralReport
	name          string
}

var reportTable = map[uint8]reportProperties{
	0x01: {
		generateValue: func() GeneralReport { return &VehicleData{} },
		name:          "vehicleData",
	},
	0x02: {
		generateValue: func() GeneralReport { return &DriveMotors{} },
		name:          "DriveMotor",
	},
	0x03: {
		generateValue: func() GeneralReport { return &FuelCell{} },
		name:          "fuelCell",
	},
	0x04: {
		generateValue: func() GeneralReport { return &Engine{} },
		name:          "engine",
	},
	0x05: {
		generateValue: func() GeneralReport { return &Location{} },
		name:          "location",
	},
	0x06: {
		generateValue: func() GeneralReport { return &Extreme{} },
		name:          "extreme",
	},
	0x07: {
		generateValue: func() GeneralReport { return &Alarm{} },
		name:          "alarm",
	},
	0x08: {
		generateValue: func() GeneralReport { return &BatteriesVoltages{} },
		name:          "batteriesVoltages",
	},
	0x09: {
		generateValue: func() GeneralReport { return &BatteriesTemperatures{} },
		name:          "batteriesTemperatures",
	},
}

type Reports map[uint8]GeneralReport

func (m *Reports) reset() {
	*m = make(map[uint8]GeneralReport)
}
func getAs[T GeneralReport](r Reports, id uint8) (T, error) {
	v, ok := r[id]
	var zero T
	if !ok {
		return zero, ErrNotAvailable
	}
	vv, ok := v.(T)
	if !ok {
		return zero, ErrMissmatch
	}
	return vv, nil

}

func (m Reports) As01() (*VehicleData, error) {
	return getAs[*VehicleData](m, 0x01)
}

func (m Reports) As02() (*DriveMotors, error) {
	return getAs[*DriveMotors](m, 0x02)
}

func (m Reports) As03() (*FuelCell, error) {
	return getAs[*FuelCell](m, 0x03)
}

func (m Reports) As04() (*Engine, error) {
	return getAs[*Engine](m, 0x04)
}

func (m Reports) As05() (*Location, error) {
	return getAs[*Location](m, 0x05)
}
func (m Reports) As06() (*Extreme, error) {
	return getAs[*Extreme](m, 0x06)
}
func (m Reports) As07() (*Alarm, error) {
	return getAs[*Alarm](m, 0x07)
}
func (m Reports) As08() (*BatteriesVoltages, error) {
	return getAs[*BatteriesVoltages](m, 0x08)
}
func (m Reports) As09() (*BatteriesTemperatures, error) {
	return getAs[*BatteriesTemperatures](m, 0x09)
}

func (m *Reports) Decode(pkt []byte, idx *int) error {
	//reset
	m.reset()

	for len(pkt[*idx:]) > 1 {
		// get type
		vType := hex.ReadByte(pkt, idx)

		// match lookuptable
		lt, ok := reportTable[vType]
		if !ok {
			return errors.Wrapf(ErrDecodeMsg, "unknow report 0x%02x", vType)
		}
		// generate item
		val := lt.generateValue()
		// decode item
		err := val.Decode(pkt, idx)
		if err != nil {
			return err
		}
		(*m)[vType] = val
	}
	return nil
}

func (m Reports) Encode() (pkt []byte, err error) {
	// sorting as slice
	kv := sortAscending(m)

	// encode every slice
	for _, v := range kv {
		p, err := v.Value.Encode()
		if err != nil {
			return pkt, err
		}
		if v.Key != v.Value.GetID() {
			return pkt, errors.Wrapf(ErrEncodeMsg, "expected id 0x%02x got 0x%02x", v.Key, v.Value.GetID())
		}
		pkt = hex.WriteByte(pkt, v.Key)
		pkt = hex.WriteBytes(pkt, p)
	}
	return pkt, err
}

func (m Reports) MarshalJSON() ([]byte, error) {
	// Initialize a map to hold the JSON structure.
	jsonMap := make(map[string]interface{})

	// Iterate through the map and marshal each report.
	for id, report := range m {

		// Determine the JSON key based on the ID.
		name := strconv.Itoa(int(id))
		if prop, ok := reportTable[id]; ok {
			name = prop.name
		}

		// Marshal the report into JSON.
		data, err := report.MarshalJSON()
		if err != nil {
			return nil, err
		}

		// Unmarshal the data into a generic interface to populate the map.
		var jsonObj interface{}
		if err := json.Unmarshal(data, &jsonObj); err != nil {
			return nil, err
		}

		jsonMap[name] = jsonObj
	}

	// Marshal the entire map into JSON.
	return json.Marshal(jsonMap)
}
