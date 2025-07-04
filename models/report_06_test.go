package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestExtremeDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Extreme
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("01020003040500060708090A0B0C"),
				idx: &[]int{0}[0],
			},
			want: Extreme{
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
		{
			name: "case: to sort",
			args: args{
				pkt: hex.Str2Byte("01020003040500060708090A0B"),
				idx: &[]int{0}[0],
			},
			want:    Extreme{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Extreme{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestExtremeDataEncode(t *testing.T) {
	type args struct {
		report Extreme
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
				report: Extreme{
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
			want:    hex.Str2Byte("01020003040500060708090A0B0C"),
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
