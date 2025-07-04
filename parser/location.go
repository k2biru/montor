package parser

func LocationLatidude() *coordinate {
	return &coordinate{
		v:      0,
		bitLoc: 1,
	}
}
func LocationLongitude() *coordinate {
	return &coordinate{
		v:      0,
		bitLoc: 2,
	}
}

func LocationValid() *locValid {
	var v locValid
	return &v
}

type locValid bool

func (m *locValid) SetValue(val uint8) *locValid {
	*m = locValid(getBit(val, 0))
	return m
}
func (m locValid) GetValue() uint8 {
	if m {
		return setBit(0, 0)
	}
	return 0
}
func (m locValid) Valid() bool {
	return !bool(m)
}

func (m *locValid) Convert(val bool) *locValid {
	*m = !locValid(val)
	return m
}
func Azimuth() *Convert[uint32] {
	value := uint32(0)
	value -= 1
	return &Convert[uint32]{
		name:  "azimuth",
		v:     value,
		unit:  0.1, // 0,1°
		min:   0,
		max:   3600,
		round: false,
	}
}
