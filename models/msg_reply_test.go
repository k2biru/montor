package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMsgReply_EncodeDecode(t *testing.T) {
	tests := []struct {
		name    string
		msg     *MsgReply
		wantErr bool
	}{
		{
			name: "case: valid reply msg",
			msg: &MsgReply{
				Header: GenerateHeader("12345678901234567", 0x02, 0x01),
				Time:   time.Date(2024, 1, 1, 12, 30, 45, 0, GBT32960Timezone()).UTC(),
			},
			wantErr: false,
		},
		{
			name: "case: reply msg with zero time",
			msg: &MsgReply{
				Header: GenerateHeader("12345678901234567", 0x01, 0x01),
				Time:   time.Time{},
			},
			wantErr: false,
		},
		{
			name: "case: missing header",
			msg: &MsgReply{
				Header: nil,
				Time:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := tt.msg.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)

			if !tt.wantErr {
				// Encoded packet has header (22 bytes) + 6 bytes time body
				body := encoded[headerLength:]
				packet := &PacketData{
					Header: tt.msg.Header,
					Body:   body,
				}
				decoded := &MsgReply{}
				err = decoded.Decode(packet)
				require.NoError(t, err)
				require.Equal(t, tt.msg.Header, decoded.GetHeader())
				require.Equal(t, "", decoded.GetMsgSN())

				cp := tt.msg.Copy()
				require.Equal(t, tt.msg.Header, cp.GetHeader())
			}
		})
	}
}
