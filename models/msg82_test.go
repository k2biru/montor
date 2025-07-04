package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg82Decode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg82
		wantErr bool
	}{
		{
			name: "case: ok ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"01"),
			},
			want: Msg82{
				Time:    time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				Control: 0x01,
			},
			wantErr: false,
		},
		{
			name: "case: invalid time ",
			args: args{
				pkt: hex.Str2Byte("18F101020304" +
					"01"),
			},
			want: Msg82{
				Time:    time.Time{},
				Control: 0x01,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body: tt.args.pkt,
			}
			msg := Msg82{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg82Encode(t *testing.T) {
	type args struct {
		msg Msg82
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case1: ok ",
			args: args{
				msg: Msg82{
					Header: &MsgHeader{
						CommandID:  0x82,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time:    time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					Control: 0x01,
				},
			},
			want: hex.Str2Byte("82FE3132333435363738390000000000000000010007" +
				"180101020304" +
				"01"),
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
