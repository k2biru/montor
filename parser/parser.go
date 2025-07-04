package parser

import (
	"math"
)

type ParsableInteger interface {
	uint8 | uint16 | uint32
}

func NewParser[T ParsableInteger](name string, val T, lookup map[T]string) *Parser[T] {
	return &Parser[T]{
		name:   name,
		v:      val,
		lookup: lookup,
	}
}

type Parser[T ParsableInteger] struct {
	name   string
	v      T
	lookup map[T]string
}

func (m *Parser[T]) SetVal(v T) *Parser[T] {
	m.v = v
	return m
}

func (m *Parser[T]) GetVal() T {
	return m.v
}

func (m Parser[T]) String() string {
	if r, ok := m.lookup[m.v]; ok {
		return r
	}
	return "unknow"
}

func (m *Parser[T]) Parse(val string) *Parser[T] {
	for i, v := range m.lookup {
		if v == val {
			m.v = i
			return m
		}
	}
	// log.Error().Err(models.ErrInvalidParse).Str("name", m.name).Str("val", val).Msg("Invalid parse")
	return m
}
func NewConvert[T ParsableInteger](name string, val T, unit float64, min, max int, round bool) *Convert[T] {
	return &Convert[T]{
		name:  name,
		v:     val,
		unit:  unit,
		min:   min,
		max:   max,
		round: round,
	}
}

type Convert[T ParsableInteger] struct {
	name  string
	v     T
	unit  float64
	min   int
	max   int
	round bool
}

func (m *Convert[T]) SetVal(v T) *Convert[T] {
	m.v = v
	return m
}

func (m *Convert[T]) GetVal() T {
	return m.v
}

func (m Convert[T]) Calculate() int {
	fError := func(cause string) {
		// log.Warn().Err(models.ErrInvalidConvert).Str("cause", cause).Str("name", m.name).Any("val", m.v).Msg("Failed convert")
	}
	var exception, invalid T
	exception = 0
	invalid = 0
	exception -= 2
	invalid -= 1
	if m.v == exception {
		fError("exception")
		return 0
	} else if m.v == invalid {
		fError("invalid")
		return 0
	}
	temp := float64(m.v)*m.unit + float64(m.min)
	result := int(temp)
	if m.round {
		result = int(math.Round(temp))
	}
	if result > m.max {
		// log.Warn().Err(models.ErrInvalidConvert).Str("name", m.name).
		// 	Any("val", m.v).Int("result", result).Any("max", m.max).Msg("Value maxed")
		return result
	}

	return result
}

func (m *Convert[T]) Convert(val int) *Convert[T] {
	m.v = 0
	m.v -= 1 // underflow -> 0xFF or 0xFFFF or 0xFFFFFFFF
	if val > m.max || val < m.min {
		// log.Error().Err(models.ErrInvalidConvert).Str("name", m.name).
		// 	Any("input", val).Int("min", m.min).Any("max", m.max).Msg("Invalid convert")
		return m
	}
	temp := float64(val) - float64(m.min)
	temp /= m.unit
	m.v = T(temp)
	if m.round {
		m.v = T(math.Round(temp))
	}
	return m
}

func (m Convert[T]) Calculate2() *calculateValue {
	fError := func(cause string) {
		// log.Warn().Err(models.ErrInvalidConvert).Str("cause", cause).Str("name", m.name).Any("val", m.v).Msg("Failed convert")
	}
	value := &calculateValue{
		val:   0,
		round: m.round,
	}
	var exception, invalid T
	exception = 0
	invalid = 0
	exception -= 2
	invalid -= 1
	switch m.v {
	case exception:
		fError("exception")
		return value
	case invalid:
		fError("invalid")
		return value

	}
	temp := float64(m.v)*m.unit + float64(m.min)
	value.val = temp
	if result := int(temp); result > m.max {
		// log.Warn().Err(models.ErrInvalidConvert).Str("name", m.name).
		// 	Any("val", m.v).Int("result", result).Any("max", m.max).Msg("Value maxed")
		value.val = float64(m.max)
	}
	return value
}

type calculateValue struct {
	val   float64
	round bool
}

func (m calculateValue) AsInt() int {
	result := int(m.val)
	if m.round {
		result = int(math.Round(m.val))
	}
	return result
}

func (m calculateValue) AsFloat() float64 {
	return math.Round(m.val*1000) / 1000
}

func setBit(value uint8, bitPos uint8) uint8 {
	return value | (1 << bitPos)
}

func getBit(value uint8, bitPos uint8) bool {
	return (value & (1 << bitPos)) != 0
}

func setByte(last uint8, value uint8, bitPos uint8, size uint8) uint8 {
	mask := ((1 << size) - 1) << bitPos
	clearedLast := last & ^uint8(mask)
	return clearedLast | ((value & uint8((1<<size)-1)) << bitPos)
}

func getByte(value uint8, bitPos uint8, size uint8) uint8 {
	mask := (1 << size) - 1
	return (value >> bitPos) & uint8(mask)
}

// / location related
type coordinate struct {
	v      float64
	bitLoc uint8
}

func (m *coordinate) SetVal(val uint32, status uint8) *coordinate {
	m.v = float64(val) / 1000000
	if getBit(status, m.bitLoc) {
		m.v *= -1
	}
	return m
}

func (m coordinate) Coordinate() float64 {
	return m.v
}

func (m *coordinate) Parse(coordinate float64) *coordinate {
	m.v = coordinate
	return m
}

func (m *coordinate) GetVal(status uint8) (uint32, uint8) {
	if m.v < 0 {
		status = setBit(status, m.bitLoc)
		m.v *= -1
	}
	result := m.v * 1000000
	return uint32(result), status
}
