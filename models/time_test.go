package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestTimeDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "case1: 2024.01.01 02.03.04 UTC8 ",
			args: args{
				pkt: hex.Str2Byte("180101020304"),
				idx: &[]int{0}[0],
			},
			want:    time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
			wantErr: false,
		},
		{
			name: "case2: invalid year ",
			args: args{
				pkt: hex.Str2Byte("F80F01020304"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "case3: invalid month ",
			args: args{
				pkt: hex.Str2Byte("180F01020304"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "case4: invalid date ",
			args: args{
				pkt: hex.Str2Byte("18012A020304"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "case5: invalid hour ",
			args: args{
				pkt: hex.Str2Byte("180101300304"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "case6: invalid minute ",
			args: args{
				pkt: hex.Str2Byte("18010102A304"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "case7: invalid second ",
			args: args{
				pkt: hex.Str2Byte("1801010203F1"),
				idx: &[]int{0}[0],
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gbtt, err := TimeDecode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, gbtt)
		})
	}
}

func TestTimeEncode(t *testing.T) {
	type args struct {
		time time.Time
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: 2024.01.01 02.03.04 UTC8 ",
			args: args{
				time: time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()),
			},
			want: hex.Str2Byte("180101020304"),
		},
		{
			name: "case1: 2024.01.01 02.03.04 UTC0 ",
			args: args{
				time: time.Date(2024, 1, 1, 2, 3, 4, 0, time.UTC),
			},
			want: hex.Str2Byte("1801010A0304"),
		},
		{
			name: "case1: 2099.12.31 15.59.59 UTC0 ",
			args: args{
				time: time.Date(2099, 12, 31, 15, 59, 59, 0, time.UTC),
			},
			want: hex.Str2Byte("630C1F173B3B"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := TimeEncode(tt.args.time)
			require.Equal(t, tt.want, pkt)
		})
	}
}
