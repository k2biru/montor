package models

import (
	"encoding/json"
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

type reportItemMock struct {
	byte  uint8
	Word  uint16
	error error
}

func (m *reportItemMock) Decode(pkt []byte, idx *int) error {
	m.byte = hex.ReadByte(pkt, idx)
	m.Word = hex.ReadWord(pkt, idx)
	return nil
}

func (m *reportItemMock) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.byte)
	pkt = hex.WriteWord(pkt, m.Word)
	return pkt, m.error
}

func (m reportItemMock) GetID() uint8 {
	return 0xF9
}

func (m *reportItemMock) MarshalJSON() ([]byte, error) {
	type Alias reportItemMock
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func TestReportItemDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Reports
		wantErr bool
	}{
		{
			name: "case: 1 ok",
			args: args{
				pkt: hex.Str2Byte("010102030004000000050006000708090A000B0C0D"),
				idx: &[]int{0}[0],
			},
			want: Reports{
				0x01: &VehicleData{
					Status:       1,
					Charging:     2,
					OpMode:       3,
					Speed:        4,
					Odometer:     5,
					TotalVoltage: 6,
					TotalCurrent: 7,
					SoC:          8,
					DCDCStatus:   9,
					Gear:         10,
					Insulator:    11,
					Throttle:     12,
					Brake:        13,
				},
			},
			wantErr: false,
		},
		{
			name: "case: ok 9 data",
			args: args{
				pkt: hex.Str2Byte("010102030004000000050006000708090A000B0C0D" + // 1
					"0201010203000400050600070008" + // 2
					"03000100020003040404040400050600070800090A0B" + // 3
					"040100020003" + // 4
					"05050000000100000002" + // 5
					"0601020003040500060708090A0B0C" + // 6
					"07010000000203000000030000000300000003" + // 7
					"0400000004000000040000000400000004" +
					"050000000500000005000000050000000500000005" +
					"06000000060000000600000006000000060000000600000006" +
					"080101000200030004000506000700070007000700070007" +
					"090101000401020304" +
					""),
				idx: &[]int{0}[0],
			},
			want: Reports{
				0x01: &VehicleData{
					Status:       1,
					Charging:     2,
					OpMode:       3,
					Speed:        4,
					Odometer:     5,
					TotalVoltage: 6,
					TotalCurrent: 7,
					SoC:          8,
					DCDCStatus:   9,
					Gear:         10,
					Insulator:    11,
					Throttle:     12,
					Brake:        13,
				},
				0x02: &DriveMotors{
					{
						SerialNumber:         1,
						Status:               2,
						ControlerTemperature: 3,
						Speed:                4,
						Torque:               5,
						Temperature:          6,
						InputVoltage:         7,
						InputCurrent:         8,
					},
				},
				0x03: &FuelCell{
					BatVoltage:                     1,
					BatCurrent:                     2,
					FuelRate:                       3,
					Temperatures:                   []uint8{4, 4, 4, 4},
					HydrogenSysMaxTemp:             5,
					HydrogenSysMaxTempNo:           6,
					HydrogenSysMaxConcentrations:   7,
					HydrogenSysMaxConcentrationsNo: 8,
					HydrogenSysMaxPressure:         9,
					HydrogenSysMaxPressureNo:       10,
					DCStatus:                       11,
				},
				0x04: &Engine{
					Status:   1,
					Revs:     2,
					FuelRate: 3,
				},
				0x05: &Location{
					Status:    5,
					Longitude: 1,
					Latidude:  2,
				},
				0x06: &Extreme{
					MaxVoltageBatAssyNo:      1,
					MaxVoltageSingleBatNo:    2,
					MaxVoltageSingleBatValue: 3,
					MinVoltageBatAssyNo:      4,
					MinVoltageSingleBatNo:    5,
					MinVoltageSingleBatValue: 6,
					MaxTempBatProbeNo:        7,
					MaxTempBatAssyNo:         8,
					MaxTempBatProbeValue:     9,
					MinTempBatAssyNo:         10,
					MinTempBatProbeNo:        11,
					MinTempBatProbeValue:     12,
				},
				0x07: &Alarm{
					AlarmLevel:         1,
					AlarmBatteryFlag:   2,
					AlarmBatteryOthers: []uint32{3, 3, 3},
					AlarmDriveMotor:    []uint32{4, 4, 4, 4},
					AlarmEngines:       []uint32{5, 5, 5, 5, 5},
					AlarmOthers:        []uint32{6, 6, 6, 6, 6, 6},
				},
				0x08: &BatteriesVoltages{
					{
						AssyNo:                  1,
						Voltage:                 2,
						Current:                 3,
						BatteriesTotalNumber:    4,
						BatteryStartNumberFrame: 5,
						SingleBatteryVoltageOnFrame: []uint16{
							7, 7, 7, 7, 7, 7,
						},
					},
				},
				0x09: &BatteriesTemperatures{
					{
						AssyNo: 1,
						Temperature: []uint8{
							1, 2, 3, 4,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "case: unknow",
			args: args{
				pkt: hex.Str2Byte("F3000200030005"),
				idx: &[]int{0}[0],
			},
			want:    Reports{},
			wantErr: true,
		},
		{
			name: "case: 1",
			args: args{
				pkt: hex.Str2Byte("0201010203000400050600070008"),
				idx: &[]int{0}[0],
			},
			want: Reports{
				0x02: &DriveMotors{
					{
						SerialNumber:         1,
						Status:               2,
						ControlerTemperature: 3,
						Speed:                4,
						Torque:               5,
						Temperature:          6,
						InputVoltage:         7,
						InputCurrent:         8,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "case: 1 -1",
			args: args{
				pkt: hex.Str2Byte("02010102030004000506000700"),
				idx: &[]int{0}[0],
			},
			want:    Reports{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Reports{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestReportItemEncode(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case: 1 ok",
			args: args{
				report: Reports{
					0x01: &VehicleData{
						Status:       1,
						Charging:     2,
						OpMode:       3,
						Speed:        4,
						Odometer:     5,
						TotalVoltage: 6,
						TotalCurrent: 7,
						SoC:          8,
						DCDCStatus:   9,
						Gear:         10,
						Insulator:    11,
						Throttle:     12,
						Brake:        13,
					},
				},
			},
			want:    hex.Str2Byte("010102030004000000050006000708090A000B0C0D"),
			wantErr: false,
		},
		{
			name: "case: 2 ok",
			args: args{
				report: Reports{
					0x01: &VehicleData{
						Status:       1,
						Charging:     2,
						OpMode:       3,
						Speed:        4,
						Odometer:     5,
						TotalVoltage: 6,
						TotalCurrent: 7,
						SoC:          8,
						DCDCStatus:   9,
						Gear:         10,
						Insulator:    11,
						Throttle:     12,
						Brake:        13,
					},
					0x02: &DriveMotors{
						{
							SerialNumber:         1,
							Status:               2,
							ControlerTemperature: 3,
							Speed:                4,
							Torque:               5,
							Temperature:          6,
							InputVoltage:         7,
							InputCurrent:         8,
						},
					},
				},
			},
			want: hex.Str2Byte("010102030004000000050006000708090A000B0C0D" +
				"0201010203000400050600070008"),
			wantErr: false,
		},
		{
			name: "case: missmatch id",
			args: args{
				report: Reports{
					0x01: &reportItemMock{
						byte: 1,
						Word: 2,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "case: error encode item",
			args: args{
				report: Reports{
					0xF9: &reportItemMock{
						byte:  1,
						Word:  2,
						error: ErrEncodeMsg,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.report.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}

func TestReportID(t *testing.T) {
	type args struct {
		report GeneralReport
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "case: 01",
			args: args{
				report: &VehicleData{},
			},
			want: 0x01,
		},
		{
			name: "case: 02",
			args: args{
				report: &DriveMotors{},
			},
			want: 0x02,
		},
		{
			name: "case: 03",
			args: args{
				report: &FuelCell{},
			},
			want: 0x03,
		},
		{
			name: "case: 04",
			args: args{
				report: &Engine{},
			},
			want: 0x04,
		},
		{
			name: "case: 05",
			args: args{
				report: &Location{},
			},
			want: 0x05,
		},
		{
			name: "case: 06",
			args: args{
				report: &Extreme{},
			},
			want: 0x06,
		},
		{
			name: "case: 07",
			args: args{
				report: &Alarm{},
			},
			want: 0x07,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.args.report.GetID()
			require.Equal(t, tt.want, id)
		})
	}
}

func TestReportAs01(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x01: &VehicleData{
						Status: 10,
					},
				},
			},
			want: &VehicleData{
				Status: 10,
			},
			wantErr: false,
		},
		{
			name: "case: NONE",
			args: args{
				report: Reports{},
			},
			want:    (*VehicleData)(nil),
			wantErr: true,
		},

		{
			name: "case: missmatch",
			args: args{
				report: Reports{
					0x01: &reportItemMock{},
				},
			},
			want:    (*VehicleData)(nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As01()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs02(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x02: &DriveMotors{
						{
							SerialNumber: 1,
						},
					},
				},
			},
			want: &DriveMotors{
				{
					SerialNumber: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As02()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs03(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x03: &FuelCell{
						BatVoltage:                     1,
						BatCurrent:                     2,
						FuelRate:                       3,
						Temperatures:                   []uint8{4, 4, 4, 4},
						HydrogenSysMaxTemp:             5,
						HydrogenSysMaxTempNo:           6,
						HydrogenSysMaxConcentrations:   7,
						HydrogenSysMaxConcentrationsNo: 8,
						HydrogenSysMaxPressure:         9,
						HydrogenSysMaxPressureNo:       10,
						DCStatus:                       11,
					},
				},
			},
			want: &FuelCell{
				BatVoltage:                     1,
				BatCurrent:                     2,
				FuelRate:                       3,
				Temperatures:                   []uint8{4, 4, 4, 4},
				HydrogenSysMaxTemp:             5,
				HydrogenSysMaxTempNo:           6,
				HydrogenSysMaxConcentrations:   7,
				HydrogenSysMaxConcentrationsNo: 8,
				HydrogenSysMaxPressure:         9,
				HydrogenSysMaxPressureNo:       10,
				DCStatus:                       11,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As03()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs04(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x04: &Engine{
						Status:   1,
						Revs:     2,
						FuelRate: 3,
					},
				},
			},
			want: &Engine{
				Status:   1,
				Revs:     2,
				FuelRate: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As04()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs05(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x05: &Location{
						Status:    5,
						Longitude: 1,
						Latidude:  2,
					},
				},
			},
			want: &Location{
				Status:    5,
				Longitude: 1,
				Latidude:  2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As05()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs06(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x06: &Extreme{
						MaxVoltageBatAssyNo:      1,
						MaxVoltageSingleBatNo:    2,
						MaxVoltageSingleBatValue: 3,
						MinVoltageBatAssyNo:      4,
						MinVoltageSingleBatNo:    5,
						MinVoltageSingleBatValue: 6,
						MaxTempBatProbeNo:        7,
						MaxTempBatAssyNo:         8,
						MaxTempBatProbeValue:     9,
						MinTempBatAssyNo:         10,
						MinTempBatProbeNo:        11,
						MinTempBatProbeValue:     12,
					},
				},
			},
			want: &Extreme{
				MaxVoltageBatAssyNo:      1,
				MaxVoltageSingleBatNo:    2,
				MaxVoltageSingleBatValue: 3,
				MinVoltageBatAssyNo:      4,
				MinVoltageSingleBatNo:    5,
				MinVoltageSingleBatValue: 6,
				MaxTempBatProbeNo:        7,
				MaxTempBatAssyNo:         8,
				MaxTempBatProbeValue:     9,
				MinTempBatAssyNo:         10,
				MinTempBatProbeNo:        11,
				MinTempBatProbeValue:     12,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As06()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}

func TestReportAs07(t *testing.T) {
	type args struct {
		report Reports
	}
	tests := []struct {
		name    string
		args    args
		want    GeneralReport
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Reports{
					0x07: &Alarm{
						AlarmLevel:         1,
						AlarmBatteryFlag:   2,
						AlarmBatteryOthers: []uint32{3, 3, 3},
						AlarmDriveMotor:    []uint32{4, 4, 4, 4},
						AlarmEngines:       []uint32{5, 5, 5, 5, 5},
						AlarmOthers:        []uint32{6, 6, 6, 6, 6, 6},
					},
				},
			},
			want: &Alarm{
				AlarmLevel:         1,
				AlarmBatteryFlag:   2,
				AlarmBatteryOthers: []uint32{3, 3, 3},
				AlarmDriveMotor:    []uint32{4, 4, 4, 4},
				AlarmEngines:       []uint32{5, 5, 5, 5, 5},
				AlarmOthers:        []uint32{6, 6, 6, 6, 6, 6},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.args.report.As07()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, d)
		})
	}
}
