package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestLocationDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Location
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("010000000100000002"),
				idx: &[]int{0}[0],
			},
			want: Location{
				Status:    1,
				Longitude: 1,
				Latidude:  2,
			},
			wantErr: false,
		},
		{
			name: "case: ok 2",
			args: args{
				pkt: hex.Str2Byte("020000000600000007"),
				idx: &[]int{0}[0],
			},
			want: Location{
				Status:    2,
				Longitude: 6,
				Latidude:  7,
			},
			wantErr: false,
		}, {
			name: "case:err",
			args: args{
				pkt: hex.Str2Byte("0200000006000000"),
				idx: &[]int{0}[0],
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Location{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestLocationEncode(t *testing.T) {
	type args struct {
		reportItem Location
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
				reportItem: Location{
					Status:    5,
					Longitude: 1,
					Latidude:  2,
				},
			},
			want:    hex.Str2Byte("050000000100000002"),
			wantErr: false,
		},
		{
			name: "case: ok2",
			args: args{
				reportItem: Location{
					Status:    2,
					Longitude: 6,
					Latidude:  7,
				},
			},
			want:    hex.Str2Byte("020000000600000007"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.reportItem.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
