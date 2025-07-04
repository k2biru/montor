package parser

func DriveMotorStatus() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "driveMotorStatus",
		lookup: map[uint8]string{
			0x01: "consumePower",
			0x02: "generatePower",
			0x03: "off",
			0x04: "ready",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func DriveMotorControlerTemperature() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "driveMotorControlerTemperature",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   210,
		round: true,
	}
}

func DriveMotorSpeed() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "driveMotorControlerTemperature",
		v:     value,
		unit:  1, // 1 /rpm/min
		min:   -20000,
		max:   45531,
		round: true,
	}
}

func DriveMotorTorque() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "driveMotorTorque",
		v:     value,
		unit:  0.1, // 0.1 /Nm
		min:   -2000,
		max:   4553,
		round: true,
	}
}

func DriveMotorTemperature() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "driveMotorTemperature",
		v:     value,
		unit:  1, // 1 /C
		min:   -40,
		max:   210,
		round: true,
	}
}
func DriveMotorInputVoltage() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "driveMotorInputVoltage",
		v:     value,
		unit:  0.1, // 0.1 /V
		min:   0,
		max:   30000,
		round: true,
	}
}

func DriveMotorInputCurrent() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "driveMotorInputCurrent",
		v:     value,
		unit:  0.1, // 0.1 /A
		min:   -1000,
		max:   1000,
		round: true,
	}
}
