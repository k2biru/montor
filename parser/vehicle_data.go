package parser

func VehicleStatus() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "vehicleStatus",
		lookup: map[uint8]string{
			0x01: "start",
			0x02: "shutdown",
			0x03: "other",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func Charging() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "charging",
		lookup: map[uint8]string{
			0x01: "parkingCharge",
			0x02: "driveCharge",
			0x03: "notCharged",
			0x04: "completed",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func OperatingMode() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "operatingMode",
		lookup: map[uint8]string{
			0x01: "electric",
			0x02: "hybrid",
			0x03: "fuelOil",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func Speed() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "speed",
		v:     value,
		unit:  0.1, // 0,1 km/h
		min:   0,
		max:   220,
		round: true,
	}
}

func Odometer() *Convert[uint32] {
	value := uint32(0)
	value -= 1
	return &Convert[uint32]{
		name:  "odometer",
		v:     value,
		unit:  0.1, // 0,1 /km
		min:   0,
		max:   9999999,
		round: false,
	}
}

func TotalVoltage() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "totalVoltage",
		v:     value,
		unit:  0.1, // 0,1 /V
		min:   0,
		max:   1000,
		round: true,
	}
}

func TotalCurrent() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "totalCurrent",
		v:     value,
		unit:  0.1, // 0,1 /A
		min:   -1000,
		max:   1000,
		round: true,
	}
}

func SoC() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "SoC",
		v:     value,
		unit:  1, // 1 /%
		min:   0,
		max:   100,
		round: true,
	}
}

func DCDCStatus() *Parser[uint8] {
	return &Parser[uint8]{
		v:    0xff,
		name: "dcDCStatus",
		lookup: map[uint8]string{
			0x01: "work",
			0x02: "break",
			0xfe: "exception",
			0xff: "invalid",
		},
	}
}

func Insulator() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "insulator",
		v:     value,
		unit:  0.1, // 1 /kOhm
		min:   0,
		max:   60000,
		round: true,
	}
}
func Gear() *gear {
	var g gear
	return &g
}

type gear uint8

func gearGear() *Parser[uint8] {
	return &Parser[uint8]{
		name: "gear",
		v:    0,
		lookup: map[uint8]string{
			0b0000: "netral",
			0b0001: "1st",
			0b0010: "2nd",
			0b0011: "3rd",
			0b0100: "4th",
			0b0101: "5th",
			0b0110: "6th",
			0b1101: "reverse",
			0b1110: "drive",
			0b1111: "park",
		},
	}
}

func (m *gear) SetVal(val uint8) *gear {
	*m = gear(val)
	return m
}

func (m *gear) GetVal() uint8 {
	return uint8(*m)
}

func (m gear) GetGear() string {
	// gear at bit 0~3
	return gearGear().SetVal(getByte(uint8(m), 0, 4)).String()
}

func (m gear) GetBraking() bool {
	return getBit(uint8(m), 4) // braking at bit 4
}

func (m gear) GetPower() bool {
	return getBit(uint8(m), 5) // power at bit 5
}

func (m *gear) SetBraking(val bool) *gear {
	if val {
		*m = gear(setBit(uint8(*m), 4)) // braking at bit 4
	}
	return m
}
func (m *gear) SetPower(val bool) *gear {
	if val {
		*m = gear(setBit(uint8(*m), 5)) // power at bit 5
	}
	return m
}
func (m *gear) SetGear(val string) *gear {
	v := gearGear().Parse(val).GetVal()
	// gear at bit 0~3
	*m = gear(setByte(uint8(*m), v, 0, 4))
	return m
}

func Throttle() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "throttle",
		v:     value,
		unit:  1, // 1 /%
		min:   0,
		max:   100,
		round: true,
	}
}

func Brake() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "brake",
		v:     value,
		unit:  1, // 1 /%
		min:   0,
		max:   100,
		round: true,
	}
}
