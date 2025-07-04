package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg81Decode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg81
		wantErr bool
	}{
		{
			name: "case: ok ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"10" +
					"010001020002030003" +
					"04040534343434" +
					"060006073737373737083838383838" +
					"09090A000A0B000B0C0C" +
					"0D040E656565650F000F1010"),
			},
			want: Msg81{
				Time: time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				Parameters: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x04: uint8(4),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0D: uint8(4),
					0x0E: []uint8(hex.Str2Byte("65656565")),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			wantErr: false,
		},
		{
			name: "case: invalid time ",
			args: args{
				pkt: hex.Str2Byte("18F101020304" +
					"10" +
					"010001020002030003" +
					"04040534343434" +
					"060006073737373737083838383838" +
					"09090A000A0B000B0C0C" +
					"0D040E656565650F000F1010"),
			},
			want: Msg81{
				Time: time.Time{},
				Parameters: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x04: uint8(4),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0D: uint8(4),
					0x0E: []uint8(hex.Str2Byte("65656565")),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := PacketData{
				Body: tt.args.pkt,
			}
			msg := Msg81{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg81Encode(t *testing.T) {
	type args struct {
		msg Msg81
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
				msg: Msg81{
					Header: &MsgHeader{
						CommandID:  0x81,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time: time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					Parameters: Parameters{
						0x01: uint16(1),
						0x02: uint16(2),
						0x03: uint16(3),
						0x04: uint8(4),
						0x05: []uint8(hex.Str2Byte("34343434")),
						0x06: uint16(6),
						0x07: string("77777"),
						0x08: string("88888"),
						0x09: uint8(9),
						0x0A: uint16(10),
						0x0B: uint16(11),
						0x0C: uint8(12),
						0x0D: uint8(4),
						0x0E: []uint8(hex.Str2Byte("65656565")),
						0x0F: uint16(15),
						0x10: uint8(16),
					},
				},
			},
			want: hex.Str2Byte("81FE313233343536373839000000000000000001003c" +
				"180101020304" +
				"10" +
				"010001020002030003" +
				"04040534343434" +
				"060006073737373737083838383838" +
				"09090A000A0B000B0C0C" +
				"0D040E656565650F000F1010"),
			wantErr: false,
		},
		{
			name: "case: unknow ",
			args: args{
				msg: Msg81{
					Header: &MsgHeader{},
					Time:   time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					Parameters: Parameters{
						0x01: uint16(1),
						0xF1: uint16(1),
					},
				},
			},
			want: hex.Str2Byte("180101020304" +
				"00"), // length 0
			wantErr: true,
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
