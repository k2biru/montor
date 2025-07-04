package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestBateriesDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    BatteriesVoltages
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("0101000200030004000506000700070007000700070007"),
				idx: &[]int{0}[0],
			},
			want: BatteriesVoltages{
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
			wantErr: false,
		},
		{
			name: "case: sort",
			args: args{
				pkt: hex.Str2Byte("01010002000300040005"),
				idx: &[]int{0}[0],
			},
			want:    BatteriesVoltages{},
			wantErr: true,
		},
		{
			name: "case: less",
			args: args{
				pkt: hex.Str2Byte("02010002000300040005060007000700070007000700070200020003000400050600070007000700070007"),
				idx: &[]int{0}[0],
			},
			want: BatteriesVoltages{
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := BatteriesVoltages{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestBateriesEncode(t *testing.T) {
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
