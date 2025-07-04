package parser

func EngineStatus() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "engineStatus",
		lookup: map[uint8]string{
			0x01: "startup",
			0x02: "shutdown",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func EngineCrankshaffSpeed() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "engineCrankshaffSpeed",
		v:     value,
		unit:  1, // 1 /r/min
		min:   0,
		max:   60000,
		round: true,
	}
}

func EngineFuelRate() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "fngineFuelRate",
		v:     value,
		unit:  0.01, // 0.01 /L/100km
		min:   0,
		max:   600,
		round: true,
	}
}
