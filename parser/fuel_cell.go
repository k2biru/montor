package parser

func FuelCellVoltage() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellVoltage",
		v:     value,
		unit:  0.1, // 0.1 /V
		min:   0,
		max:   2000,
		round: true,
	}
}

func FuelCellCurrent() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellCurrent",
		v:     value,
		unit:  0.1, // 0.1 /A
		min:   0,
		max:   2000,
		round: true,
	}
}

func FuelCellRate() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellRate",
		v:     value,
		unit:  0.01, // 0.01 /kg/100km
		min:   0,
		max:   600,
		round: true,
	}
}
func FuelCellTemperature() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "fuelCellTemperature",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   200,
		round: true,
	}
}

func FuelCellHydrogenSysTemp() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellHydrogenSysTemp",
		v:     value,
		unit:  0.1, // 0.1 /C
		min:   -40,
		max:   200,
		round: true,
	}
}

func FuelCellHydrogenSysConct() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellHydrogenSysConct",
		v:     value,
		unit:  1, // 1 /mg/kg
		min:   0,
		max:   60000,
		round: true,
	}
}

func FuelCellHydrogenSysPress() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fuelCellHydrogenSysPress",
		v:     value,
		unit:  0.1, // 0.1 /MPa
		min:   0,
		max:   100,
		round: true,
	}
}

func FuelCellDCDCStatus() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "fuelCellDCDCStatus",
		lookup: map[uint8]string{
			0x01: "working",
			0x02: "disconnect",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}
