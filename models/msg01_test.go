package models

import (
	"testing"
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestMsg01Decode(t *testing.T) {
	type args struct {
		pkt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Msg01
		wantErr bool
	}{
		{
			name: "case: 2024.01.01 02.03.04 UTC8 standard ",
			args: args{
				pkt: hex.Str2Byte("180101020304" +
					"00033839353130303739353138393332333035363235" +
					"010102"),
			},
			want: Msg01{
				Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
				SerialNumber: 3,
				ICCID:        "89510079518932305625",
				EnerygyStorageSystem: EnergyStorageSys{
					Coding: 1,
					Raw:    []byte{0x02},
				},
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
			want: Msg01{
				Time:         time.Time{},
				SerialNumber: 3,
				ICCID:        "89510079518932305625",
				EnerygyStorageSystem: EnergyStorageSys{
					Coding: 1,
					Raw:    []byte{0x02},
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
			msg := Msg01{}
			err := msg.Decode(&pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, msg)
		})
	}
}

func TestMsg01Encode(t *testing.T) {
	type args struct {
		msg Msg01
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
				msg: Msg01{
					Header: &MsgHeader{
						CommandID:  0x01,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, GBT32960Timezone()).UTC(),
					SerialNumber: 3,
					ICCID:        "89510079518932305625",
					EnerygyStorageSystem: EnergyStorageSys{
						Coding: 1,
						Raw:    []byte{0x02},
					},
				},
			},
			want: hex.Str2Byte("01FE313233343536373839000000000000000001001f" +
				"180101020304" +
				"00033839353130303739353138393332333035363235" +
				"010102"),
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
