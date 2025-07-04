package parser

func ParamLocalStorageDurrationMs() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "localStorageDurrationMs",
		v:     value,
		unit:  1, // 1/ms
		min:   0,
		max:   60000,
		round: true,
	}
}

func ParamReportDurrationDefaultMs() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "reportDurrationDefaultMs",
		v:     value,
		unit:  1, // 1/ms
		min:   0,
		max:   60000,
		round: true,
	}
}

func ParamReportDurrationAlarmMs() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "reportDurrationAlarmMs",
		v:     value,
		unit:  1, // 1/ms
		min:   0,
		max:   60000,
		round: true,
	}
}

func ParamHeartBeatSec() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "HeartBeatSec",
		v:     value,
		unit:  1, // 1/sec
		min:   0,
		max:   240,
		round: true,
	}
}

func ParamTeminalResponseTimeup() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "terminalResponseTimeup",
		v:     value,
		unit:  1, // 1/sec
		min:   0,
		max:   600,
		round: true,
	}
}

func ParamPlatformResponseTimeup() *Convert[uint16] {
	value := uint16(0)
	value -= 1
	return &Convert[uint16]{
		name:  "platformResponseTimeup",
		v:     value,
		unit:  1, // 1/sec
		min:   0,
		max:   600,
		round: true,
	}
}

func ParamLoginIntervalMin() *Convert[uint8] {
	value := uint8(0)
	value -= 1
	return &Convert[uint8]{
		name:  "loginIntervalMin",
		v:     value,
		unit:  1, // 1/minute
		min:   0,
		max:   240,
		round: true,
	}
}
