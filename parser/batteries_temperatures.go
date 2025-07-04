package parser

func BatteriesTemperaturesTemp() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "batteriesTemperaturesTemp",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   210,
		round: true,
	}
}
