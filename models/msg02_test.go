package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg02Decode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg02
		wantErr bool
	}{
		{
			name: "case: 2024.01.01 02.03.04 UTC8 standard ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"010102030004000000050006000708090A000B0C0D" +
					"0201010203000400050600070008"),
			},

			want: Msg02{
				Time: time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				Reports: Reports{
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
			wantErr: false,
		},
		{
			name: "case: unknow ",
			args: args{
				pkt: hex.Str2Byte("18F101020304" +
					"F3000200030005"),
			},
			want: Msg02{
				Time:    time.Time{},
				Reports: make(Reports),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body: tt.args.pkt,
			}
			msg := Msg02{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg02Encode(t *testing.T) {
	type args struct {
		msg Msg02
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case1: 2024.01.01 02.03.04 UTC8 ",
			args: args{
				msg: Msg02{
					Header: &MsgHeader{
						CommandID:  0x02,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time: time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					Reports: Reports{
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
			},
			want: hex.Str2Byte("02FE3132333435363738390000000000000000010029" +
				"180101020304" +
				"010102030004000000050006000708090A000B0C0D" +
				"0201010203000400050600070008"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.msg.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
