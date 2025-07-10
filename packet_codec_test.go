package montor

import (
	"testing"
	"time"

	"github.com/k2biru/montor/models"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestPacketCodec_Decode(t *testing.T) {
	type args struct {
		time    time.Time
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *models.PacketData
		wantErr bool
	}{
		{
			name: "case : ok",
			args: args{
				time:    time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
				payload: hex.Str2Byte("01FE313233343536373839000000000000000001002A130116173738000131323334353637383939383736353433323130300304313233343435363739383730FD"),
			},
			want: &models.PacketData{
				Header: &models.MsgHeader{
					CommandID:   0x01,
					Response:    0xFE,
					VIN:         "123456789",
					Encription:  0x01,
					BodyLength:  0x002A,
					TimeCreated: time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
				},
				Body:       hex.Str2Byte("130116173738000131323334353637383939383736353433323130300304313233343435363739383730"),
				VerifyCode: 0xFD,
			},
			wantErr: false,
		},

		{
			name: "case : invalid verifyCode",
			args: args{
				payload: hex.Str2Byte("01FE313233343536373839000000000000000001002A130116173738000131323334353637383939383736353433323130300304313233343435363739383730FF"),
			},
			wantErr: true,
		},
		{
			name: "case : empty",
			args: args{
				payload: hex.Str2Byte(""),
			},
			wantErr: true,
		},

		{
			name: "case : invalid tooshort",
			args: args{
				payload: hex.Str2Byte("01FE313233343536373839000000000000000001CF"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &packetCodec{
				timeNow: func() time.Time {
					return tt.args.time
				},
			}
			got, err := pc.Decode(tt.args.payload)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestPacketCodec_Encode(t *testing.T) {
	type args struct {
		pkt []byte
		msg models.GBT32960Msg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []byte
	}{
		{
			name: "case : ok",
			args: args{
				pkt: hex.Str2Byte("00010203040506"),
				msg: &models.Msg01{
					Header: &models.MsgHeader{
						CommandID:  0x01,
						Response:   0xFE,
						VIN:        "123456789",
						Encription: 0x01,
					},
					Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, models.GBT32960Timezone()).UTC(),
					SerialNumber: 3,
					ICCID:        "89510079518932305625",
					EnerygyStorageSystem: models.EnergyStorageSys{
						Coding: 1,
						Raw:    []byte{0x02},
					},
				}},
			want: hex.Str2Byte("01FE313233343536373839000000000000000001001f" +
				"180101020304" +
				"00033839353130303739353138393332333035363235" +
				"010102" +
				"c4"),
			wantErr: false,
		},
		{
			name: "case : no header",
			args: args{
				pkt: hex.Str2Byte("00010203040506"),
				msg: &models.Msg01{
					Time:         time.Date(2024, 1, 1, 2, 3, 4, 0, models.GBT32960Timezone()).UTC(),
					SerialNumber: 3,
					ICCID:        "89510079518932305625",
					EnerygyStorageSystem: models.EnergyStorageSys{
						Coding: 1,
						Raw:    []byte{0x02},
					},
				}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := NewPacketCodec(nil)
			got, err := pc.Encode(tt.args.msg)

			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// func TestPacketCodec1_Decode(t *testing.T) {
// 	type args struct {
// 		time    time.Time
// 		payload []byte
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *models.PacketData
// 		wantErr bool
// 	}{
// 		{
// 			name: "case : ok",
// 			args: args{
// 				time:    time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
// 				payload: hex.Str2Byte("02FE4C4D454C424C315039525243363035303001013E180C0A111C0A02010104574E204E205A10EA27100700000000000000000008010110DF274C006700016710611062106310621063106110601061106210611062106210621062106310621065106210631063106410631063106310631063106410641061106110611062106210611063106210611062106310621063106310621061105F105F1060106010611060106010601060106110611061106110601062106410641063106210621063106210631064106310641064106410631063106210611062106010601063106110621062106110611061106310611063106210641062106210631064106310621063106110631064106310640901010013444443444444444444444444444544444444450601111065012D105F010E4501034301010301000000001A2A10DF274C60010F10B800000502065DB4AA005E7AEF61"),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			pc := &packetCodec{
// 				timeNow: func() time.Time {
// 					return tt.args.time
// 				},
// 			}
// 			got, _ := pc.Decode(tt.args.payload)
// 			pp := NewProcessor()
// 			// pd, err := pp.Process(context.Background(), got)
// 			// fmt.Printf(">>> %s code %v \n", err.Error(), pd)
// 			// require.Equal(t, tt.wantErr, err != nil, err)
// 			// require.Equal(t, tt.want, got)
// 		})
// 	}
// }
