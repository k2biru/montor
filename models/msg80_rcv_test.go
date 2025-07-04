package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg80ReceiveDecode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg80Receive
		wantErr bool
	}{
		{
			name: "case: ok ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"06010203040506"),
			},
			want: Msg80Receive{
				Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				ParameterIDs: hex.Str2Byte("010203040506"),
			},
			wantErr: false,
		},
		{
			name: "case: invalid time ",
			args: args{
				pkt: hex.Str2Byte("18F101020304" +
					"06010203040506"),
			},
			want: Msg80Receive{
				Time:         time.Time{},
				ParameterIDs: hex.Str2Byte("010203040506"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body: tt.args.pkt,
			}
			msg := Msg80Receive{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg80ReceiveEncode(t *testing.T) {
	type args struct {
		msg Msg80Receive
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
				msg: Msg80Receive{
					Header: &MsgHeader{
						CommandID:  0x80,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					ParameterIDs: hex.Str2Byte("010203040506"),
				},
			},
			want: hex.Str2Byte("80FE313233343536373839000000000000000001000d" +
				"180101020304" +
				"06010203040506"),
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
