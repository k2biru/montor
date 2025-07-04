package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsgxxDecode(t *testing.T) {
	type args struct {
		header *MsgHeader
		pkt    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msgxx
		wantErr bool
	}{
		{
			name: "case: ok ",
			args: args{
				header: &MsgHeader{
					BodyLength: 11,
				},
				pkt: hex.Str2Byte("19050215293b" +
					"0001030501"),
			},
			want: Msgxx{
				Header: &MsgHeader{
					BodyLength: 11,
				},
				Raw: hex.Str2Byte("19050215293b" +
					"0001030501"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body:   tt.args.pkt,
				Header: tt.args.header,
			}
			msg := Msgxx{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsgxxEncode(t *testing.T) {
	type args struct {
		msg Msgxx
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
				msg: Msgxx{
					Header: &MsgHeader{
						CommandID:  0x82,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Raw: hex.Str2Byte("19050215293b" +
						"0001030501"),
				},
			},
			want: hex.Str2Byte("82FE313233343536373839000000000000000001000b" +
				"19050215293b" +
				"0001030501"),
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
