package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg04Decode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg04
		wantErr bool
	}{
		{
			name: "case: 2024.01.01 02.03.04 UTC8 standard ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"00033839353130303739353138393332333035363235" +
					"010102"),
			},
			want: Msg04{
				Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				SerialNumber: 3,
			},
			wantErr: false,
		},
		{
			name: "case: invalid time ",
			args: args{
				pkt: hex.Str2Byte("18F101020304" +
					"00033839353130303739353138393332333035363235" +
					"010102"),
			},
			want: Msg04{
				Time:         time.Time{},
				SerialNumber: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body: tt.args.pkt,
			}
			msg := Msg04{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg04Encode(t *testing.T) {
	type args struct {
		msg Msg04
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case1: 2024.01.01 02.03.04 UTC8 ",
			args: args{
				msg: Msg04{
					Header: &MsgHeader{
						CommandID:  0x04,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					SerialNumber: 3,
				},
			},
			want: hex.Str2Byte("04FE3132333435363738390000000000000000010008" +
				"180101020304" +
				"0003"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.msg.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
