package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsgHeaderDecode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *MsgHeader
		wantErr bool
	}{
		{
			name: "case1: registrasi ",
			args: args{
				pkt: hex.Str2Byte("01FE3132333435363738393031323334353637010002"),
			},
			want: &MsgHeader{
				CommandID:  0x01,
				Response:   ResponseCommandPackage,
				VIN:        "12345678901234567",
				Encription: 0x01,
				BodyLength: 2,
				Idx:        22,
			},
			wantErr: false,
		},
		{
			name: "case : les then 22 ",
			args: args{
				pkt: hex.Str2Byte("01FE31323334353637383930313233343536370100"),
			},
			want:    &MsgHeader{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := &MsgHeader{}
			err := header.Decode(tt.args.pkt)

			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, header)
		})
	}
}

func TestMsgHeaderEncode(t *testing.T) {
	type args struct {
		header *MsgHeader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case1: registrasi ",
			args: args{
				header: &MsgHeader{
					CommandID:  0x01,
					Response:   ResponseCommandPackage,
					VIN:        "12345678901234567",
					Encription: 0x01,
					BodyLength: 2,
					Idx:        22,
				},
			},
			want:    hex.Str2Byte("01FE3132333435363738393031323334353637010002"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := tt.args.header
			pkt, err := header.Encode()

			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
