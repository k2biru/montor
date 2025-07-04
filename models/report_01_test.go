package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestVehicleDataDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    VehicleData
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("0102030004000000050006000708090A000B0C0D"),
				idx: &[]int{0}[0],
			},
			want: VehicleData{
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
			wantErr: false,
		},
		{
			name: "case: to sort",
			args: args{
				pkt: hex.Str2Byte("0102030004000000050006000708090A000B0C"),
				idx: &[]int{0}[0],
			},
			want:    VehicleData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := VehicleData{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestVehicleDataEncode(t *testing.T) {
	type args struct {
		vehicle VehicleData
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
				vehicle: VehicleData{
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
			want:    hex.Str2Byte("0102030004000000050006000708090A000B0C0D"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.vehicle.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
