package gbk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGBK2UTF8(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "case1: 京",
			args:    args{s: []byte{0xbe, 0xa9}},
			want:    []byte("京"),
			wantErr: false,
		},
		{
			name:    "case2: invalid GBK encoding",
			args:    args{s: []byte{0xff}}, // Invalid GBK sequence
			want:    []byte("�"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GBK2UTF8(tt.args.s)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUTF82GBK(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "case1: 京",
			args:    args{s: []byte("s京")},
			want:    []byte{0x73, 0xbe, 0xa9},
			wantErr: false,
		},
		{
			name:    "case2: error encoding",
			args:    args{s: []byte{230, 181}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UTF82GBK(tt.args.s)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, got)
		})
	}
}
