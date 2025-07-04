package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestBateriTemperaturesDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    BatteriesTemperatures
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("0101000401020304"),
				idx: &[]int{0}[0],
			},
			want: BatteriesTemperatures{
				{
					AssyNo: 1,
					Temperature: []uint8{
						1, 2, 3, 4,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "case: sort",
			args: args{
				pkt: hex.Str2Byte("0100"),
				idx: &[]int{0}[0],
			},
			want:    BatteriesTemperatures{},
			wantErr: true,
		},
		{
			name: "case: less",
			args: args{
				pkt: hex.Str2Byte("0201000401020304020004010203"),
				idx: &[]int{0}[0],
			},
			want: BatteriesTemperatures{
				{
					AssyNo: 1,
					Temperature: []uint8{
						1, 2, 3, 4,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := BatteriesTemperatures{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestBateriTemperaturesEncode(t *testing.T) {
	type args struct {
		data DriveMotor
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				data: DriveMotor{
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
			want:    hex.Str2Byte("010203000400050600070008"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.data.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}

// func TestDriveMotorsDecode(t *testing.T) {
// 	type args struct {
// 		pkt []byte
// 		idx *int
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    DriveMotors
// 		wantErr bool
// 	}{
// 		{
// 			name: "case: ok",
// 			args: args{
// 				pkt: hex.Str2Byte("01010203000400050600070008"),
// 				idx: &[]int{0}[0],
// 			},
// 			want: DriveMotors{
// 				{
// 					SerialNumber:         1,
// 					Status:               2,
// 					ControlerTemperature: 3,
// 					Speed:                4,
// 					Torque:               5,
// 					Temperature:          6,
// 					InputVoltage:         7,
// 					InputCurrent:         8,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "case: sort",
// 			args: args{
// 				pkt: hex.Str2Byte("010203000400050600070008"),
// 				idx: &[]int{0}[0],
// 			},
// 			want:    DriveMotors{},
// 			wantErr: true,
// 		},
// 		{
// 			name: "case: 2 drive",
// 			args: args{
// 				pkt: hex.Str2Byte("02010203000400050600070008" +
// 					"060503000400050600070008"),
// 				idx: &[]int{0}[0],
// 			},
// 			want: DriveMotors{
// 				{
// 					SerialNumber:         1,
// 					Status:               2,
// 					ControlerTemperature: 3,
// 					Speed:                4,
// 					Torque:               5,
// 					Temperature:          6,
// 					InputVoltage:         7,
// 					InputCurrent:         8,
// 				},
// 				{
// 					SerialNumber:         6,
// 					Status:               5,
// 					ControlerTemperature: 3,
// 					Speed:                4,
// 					Torque:               5,
// 					Temperature:          6,
// 					InputVoltage:         7,
// 					InputCurrent:         8,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "case: 2 drive, -1",
// 			args: args{
// 				pkt: hex.Str2Byte("02010203000400050600070008" +
// 					"0605030004000506000700"),
// 				idx: &[]int{0}[0],
// 			},
// 			want:    DriveMotors{},
// 			wantErr: true,
// 		},
// 		{
// 			name: "case: 1+1",
// 			args: args{
// 				pkt: hex.Str2Byte("010102030004000506000700080F"),
// 				idx: &[]int{0}[0],
// 			},
// 			want: DriveMotors{
// 				&DriveMotor{
// 					SerialNumber:         1,
// 					Status:               2,
// 					ControlerTemperature: 3,
// 					Speed:                4,
// 					Torque:               5,
// 					Temperature:          6,
// 					InputVoltage:         7,
// 					InputCurrent:         8,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			vd := DriveMotors{}
// 			err := vd.Decode(tt.args.pkt, tt.args.idx)
// 			require.Equal(t, tt.wantErr, err != nil, err)
// 			require.Equal(t, tt.want, vd)
// 		})
// 	}
// }

// func TestDriveMotorsEncode(t *testing.T) {
// 	type args struct {
// 		data DriveMotors
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []byte
// 		wantErr bool
// 	}{
// 		{
// 			name: "case: ok",
// 			args: args{
// 				data: DriveMotors{
// 					{
// 						SerialNumber:         1,
// 						Status:               2,
// 						ControlerTemperature: 3,
// 						Speed:                4,
// 						Torque:               5,
// 						Temperature:          6,
// 						InputVoltage:         7,
// 						InputCurrent:         8,
// 					},
// 				},
// 			},
// 			want:    hex.Str2Byte("01010203000400050600070008"),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			pkt, err := tt.args.data.Encode()
// 			require.Equal(t, tt.wantErr, err != nil, err)
// 			require.Equal(t, tt.want, pkt)
// 		})
// 	}
// }
