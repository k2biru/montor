package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestFuelCellDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    FuelCell
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("000100020003040404040400050600070800090A0B"),
				idx: &[]int{0}[0],
			},
			want: FuelCell{
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
		{
			name: "case: to sort",
			args: args{
				pkt: hex.Str2Byte("000100020003040404040400050600070800090A"),
				idx: &[]int{0}[0],
			},
			want:    FuelCell{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := FuelCell{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestFuelCellDataEncode(t *testing.T) {
	type args struct {
		report FuelCell
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
				report: FuelCell{
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
			want:    hex.Str2Byte("000100020003040404040400050600070800090A0B"),
			wantErr: false,
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
