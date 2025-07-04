package parser

func ExtremeMaxVoltageSingleBatterymV() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "extremeMaxVoltageSingleBattery",
		v:     value,
		unit:  1, // 0.001 /v -> 1 /mV
		min:   0,
		max:   15000,
		round: true,
	}
}

func ExtremeMinVoltageSingleBatterymV() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "extremeMinVoltageSingleBatterymV",
		v:     value,
		unit:  1, // 0.001 /v -> 1 /mV
		min:   0,
		max:   15000,
		round: true,
	}
}

// unit in C
func ExtremeMaxTempProbe() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "extremeMaxTempProbe",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   210,
		round: true,
	}
}

// unit in C
func ExtremeMinTempProbe() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "extremeMinTempProbe",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   210,
		round: true,
	}
}
