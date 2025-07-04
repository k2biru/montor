package parser

func BatteriesVoltagesMilli() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "batteriesVoltagesMilli",
		v:     value,
		unit:  1, // 1 /mV
		min:   0,
		max:   60000,
		round: true,
	}
}

func BatteriesVoltagesVolt() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "batteriesVoltagesVolt",
		v:     value,
		unit:  0.1, // 1 /0.01V
		min:   0,
		max:   1000,
		round: true,
	}
}

func BatteriesVoltagesCurrent() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "batteriesVoltagesCurrent",
		v:     value,
		unit:  0.1, // 1 /0.01A
		min:   -1000,
		max:   1000,
		round: true,
	}
}
