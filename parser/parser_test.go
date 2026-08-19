package parser

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	lookup := map[uint8]string{
		1: "working",
		2: "stopped",
	}

	tests := []struct {
		name       string
		initialVal uint8
		parseVal   string
		wantString string
		wantVal    uint8
	}{
		{
			name:       "case: lookup match working",
			initialVal: 1,
			parseVal:   "working",
			wantString: "working",
			wantVal:    1,
		},
		{
			name:       "case: lookup match stopped via parse",
			initialVal: 1,
			parseVal:   "stopped",
			wantString: "stopped",
			wantVal:    2,
		},
		{
			name:       "case: parse unknown string keeps original value",
			initialVal: 1,
			parseVal:   "unknown_state",
			wantString: "working",
			wantVal:    1,
		},
		{
			name:       "case: initial unknown value returns 'unknow'",
			initialVal: 99,
			parseVal:   "",
			wantString: "unknow",
			wantVal:    99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser("test_parser", tt.initialVal, lookup)
			require.Equal(t, tt.initialVal, p.GetVal())

			if tt.parseVal != "" {
				p.Parse(tt.parseVal)
			}
			p.SetVal(p.GetVal())
			require.Equal(t, tt.wantVal, p.GetVal())
			require.Equal(t, tt.wantString, p.String())
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name       string
		val        uint16
		unit       float64
		min        int
		max        int
		round      bool
		convertVal int
		wantVal    uint16
		wantInt    int
		wantFloat  float64
		wantErr    error
	}{
		{
			name:       "case: valid convert without round (speed conversion)",
			unit:       0.1,
			min:        0,
			max:        2200,
			round:      false,
			convertVal: 120, // 120 / 0.1 = 1200
			wantVal:    1200,
			wantInt:    120,
			wantFloat:  120.0,
			wantErr:    nil,
		},
		{
			name:       "case: valid convert with round (temperature with offset -40)",
			unit:       1.0,
			min:        -40,
			max:        210,
			round:      true,
			convertVal: 25, // (25 - (-40)) / 1.0 = 65
			wantVal:    65,
			wantInt:    25,
			wantFloat:  25.0,
			wantErr:    nil,
		},
		{
			name:       "case: convert under min underflows value",
			unit:       1.0,
			min:        0,
			max:        100,
			round:      false,
			convertVal: -10,
			wantVal:    65535, // 0xFFFF underflow
			wantInt:    100,
			wantFloat:  100.0,
			wantErr:    ErrInvalid,
		},
		{
			name:       "case: convert over max underflows value",
			unit:       1.0,
			min:        0,
			max:        100,
			round:      false,
			convertVal: 200,
			wantVal:    65535,
			wantInt:    100,
			wantFloat:  100.0,
			wantErr:    ErrInvalid,
		},
		{
			name:       "case: calculate exception value (0xFFFE = 65534)",
			val:        65534, // 0xFFFF - 1
			unit:       1.0,
			min:        0,
			max:        100,
			round:      false,
			convertVal: -999,
			wantVal:    65534,
			wantInt:    0,
			wantFloat:  0,
			wantErr:    ErrException,
		},
		{
			name:       "case: calculate over max size error",
			val:        500,
			unit:       1.0,
			min:        0,
			max:        100,
			round:      false,
			convertVal: -999,
			wantVal:    500,
			wantInt:    100,
			wantFloat:  100,
			wantErr:    ErrOverMaxSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := NewConvert("test_convert", tt.val, tt.unit, tt.min, tt.max, tt.round)
			if tt.convertVal != -999 {
				conv.Convert(tt.convertVal)
			}
			require.Equal(t, tt.wantVal, conv.GetVal())

			calc := conv.Calculate()
			require.Equal(t, tt.wantErr, calc.Error())
			if tt.wantErr == nil {
				require.Equal(t, tt.wantInt, calc.AsInt())
				require.InDelta(t, tt.wantFloat, calc.AsFloat(), 0.001)
			}
		})
	}
}

func TestBitOperations(t *testing.T) {
	t.Run("setBit and getBit", func(t *testing.T) {
		var val uint8 = 0
		val = setBit(val, 0) // 1
		val = setBit(val, 3) // 1 | 8 = 9

		require.True(t, getBit(val, 0))
		require.False(t, getBit(val, 1))
		require.False(t, getBit(val, 2))
		require.True(t, getBit(val, 3))
		require.Equal(t, uint8(9), val)
	})

	t.Run("setByte and getByte", func(t *testing.T) {
		var val uint8 = 0
		val = setByte(val, 5, 2, 3)
		require.Equal(t, uint8(20), val)
		require.Equal(t, uint8(5), getByte(val, 2, 3))
	})
}

func TestCoordinate(t *testing.T) {
	tests := []struct {
		name       string
		val        uint32
		status     uint8
		bitLoc     uint8
		inputCoord float64
		wantCoord  float64
		wantVal    uint32
		wantStatus uint8
	}{
		{
			name:       "case: positive longitude (East)",
			val:        121473611,
			status:     0,
			bitLoc:     0,
			wantCoord:  121.473611,
			wantVal:    121473611,
			wantStatus: 0,
		},
		{
			name:       "case: negative latitude (South)",
			val:        31230416,
			status:     1,
			bitLoc:     0,
			wantCoord:  -31.230416,
			wantVal:    31230416,
			wantStatus: 1,
		},
		{
			name:       "case: parse float coordinate",
			inputCoord: -121.500000,
			status:     0,
			bitLoc:     1,
			wantCoord:  -121.500000,
			wantVal:    121500000,
			wantStatus: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord := &coordinate{bitLoc: tt.bitLoc}
			if tt.inputCoord != 0 {
				coord.Parse(tt.inputCoord)
			} else {
				coord.SetVal(tt.val, tt.status)
			}

			require.InDelta(t, tt.wantCoord, coord.Coordinate(), 0.000001)

			gotVal, gotStatus := coord.GetVal(tt.status)
			require.Equal(t, tt.wantVal, gotVal)
			if math.Signbit(tt.wantCoord) {
				require.True(t, getBit(gotStatus, tt.bitLoc))
			}
		})
	}
}

func TestVehicleDataParsers(t *testing.T) {
	t.Run("Vehicle Status and Enums", func(t *testing.T) {
		require.Equal(t, "start", VehicleStatus().SetVal(0x01).String())
		require.Equal(t, "parkingCharge", Charging().SetVal(0x01).String())
		require.Equal(t, "electric", OperatingMode().SetVal(0x01).String())
		require.Equal(t, "work", DCDCStatus().SetVal(0x01).String())
	})

	t.Run("Vehicle Converters", func(t *testing.T) {
		require.NotNil(t, Speed())
		require.NotNil(t, Odometer())
		require.NotNil(t, TotalVoltage())
		require.NotNil(t, TotalCurrent())
		require.NotNil(t, Insulator())
		require.NotNil(t, Throttle())
		require.NotNil(t, Brake())
	})

	t.Run("Gear Operations", func(t *testing.T) {
		g := Gear()
		g.SetVal(0)
		g.SetGear("drive").SetBraking(true).SetPower(true)

		require.Equal(t, "drive", g.GetGear())
		require.True(t, g.GetBraking())
		require.True(t, g.GetPower())
		require.Equal(t, uint8(0b111110), g.GetVal())
	})
}

func TestDriveMotorParsers(t *testing.T) {
	t.Run("Drive Motor Enums", func(t *testing.T) {
		require.Equal(t, "consumePower", DriveMotorStatus().SetVal(0x01).String())
	})

	t.Run("Drive Motor Converters", func(t *testing.T) {
		require.NotNil(t, DriveMotorControlerTemperature())
		require.NotNil(t, DriveMotorSpeed())
		require.NotNil(t, DriveMotorTorque())
		require.NotNil(t, DriveMotorTemperature())
		require.NotNil(t, DriveMotorInputVoltage())
		require.NotNil(t, DriveMotorInputCurrent())
	})
}

func TestFuelCellParsers(t *testing.T) {
	t.Run("Fuel Cell Enums", func(t *testing.T) {
		require.Equal(t, "working", FuelCellDCDCStatus().SetVal(0x01).String())
	})

	t.Run("Fuel Cell Converters", func(t *testing.T) {
		require.NotNil(t, FuelCellVoltage())
		require.NotNil(t, FuelCellCurrent())
		require.NotNil(t, FuelCellRate())
		require.NotNil(t, FuelCellTemperature())
		require.NotNil(t, FuelCellHydrogenSysTemp())
		require.NotNil(t, FuelCellHydrogenSysConct())
		require.NotNil(t, FuelCellHydrogenSysPress())
	})
}

func TestLocationParsers(t *testing.T) {
	t.Run("Location Coordinates & Azimuth", func(t *testing.T) {
		require.NotNil(t, LocationLatidude())
		require.NotNil(t, LocationLongitude())
		require.NotNil(t, Azimuth())
	})

	t.Run("LocationValid Operations", func(t *testing.T) {
		lv := LocationValid()
		lv.SetValue(1)
		require.Equal(t, uint8(1), lv.GetValue())
		require.False(t, lv.Valid())

		lvFalse := LocationValid()
		lvFalse.SetValue(0)
		require.Equal(t, uint8(0), lvFalse.GetValue())
		require.True(t, lvFalse.Valid())

		lv.Convert(true)
		require.True(t, lv.Valid())
	})
}

func TestParameterParsers(t *testing.T) {
	t.Run("Parameter Converters", func(t *testing.T) {
		require.NotNil(t, ParamLocalStorageDurrationMs())
		require.NotNil(t, ParamReportDurrationDefaultMs())
		require.NotNil(t, ParamReportDurrationAlarmMs())
		require.NotNil(t, ParamHeartBeatSec())
		require.NotNil(t, ParamTeminalResponseTimeup())
		require.NotNil(t, ParamPlatformResponseTimeup())
		require.NotNil(t, ParamLoginIntervalMin())
	})
}

func TestEngineParsers(t *testing.T) {
	t.Run("Engine Enums & Converters", func(t *testing.T) {
		require.Equal(t, "startup", EngineStatus().SetVal(0x01).String())
		require.NotNil(t, EngineCrankshaffSpeed())
		require.NotNil(t, EngineFuelRate())
	})
}

func TestExtremeParsers(t *testing.T) {
	t.Run("Extreme Converters", func(t *testing.T) {
		require.NotNil(t, ExtremeMaxVoltageSingleBatterymV())
		require.NotNil(t, ExtremeMinVoltageSingleBatterymV())
		require.NotNil(t, ExtremeMaxTempProbe())
		require.NotNil(t, ExtremeMinTempProbe())
	})
}

func TestBatteriesParsers(t *testing.T) {
	t.Run("Batteries Voltages & Temperatures Converters", func(t *testing.T) {
		require.NotNil(t, BatteriesTemperaturesTemp())
		require.NotNil(t, BatteriesVoltagesMilli())
		require.NotNil(t, BatteriesVoltagesVolt())
		require.NotNil(t, BatteriesVoltagesCurrent())
	})
}
