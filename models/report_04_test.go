package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestEngineDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Engine
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("0100020003"),
				idx: &[]int{0}[0],
			},
			want: Engine{
				Status:   1,
				Revs:     2,
				FuelRate: 3,
			},
			wantErr: false,
		},
		{
			name: "case: to sort",
			args: args{
				pkt: hex.Str2Byte("01000200"),
				idx: &[]int{0}[0],
			},
			want:    Engine{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Engine{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestEngineDataEncode(t *testing.T) {
	type args struct {
		report Engine
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
				report: Engine{
					Status:   1,
					Revs:     2,
					FuelRate: 3,
				},
			},
			want:    hex.Str2Byte("0100020003"),
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
