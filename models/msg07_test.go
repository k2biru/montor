package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsg07_EncodeDecode(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Msg07
		wantErr bool
	}{
		{
			name: "case: valid heartbeat msg07",
			msg: &Msg07{
				Header: GenerateHeader("12345678901234567", 0x07, 0xFE),
			},
			wantErr: false,
		},
		{
			name: "case: missing header",
			msg: &Msg07{
				Header: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := tt.msg.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)

			if !tt.wantErr {
				packet := &PacketData{
					Header: tt.msg.Header,
					Body:   []byte{},
				}
				decoded := &Msg07{}
				err = decoded.Decode(packet)
				require.NoError(t, err)
				require.Equal(t, tt.msg.Header, decoded.Header)
				require.Equal(t, "", decoded.GetMsgSN())

				require.NotNil(t, encoded)
				cp := tt.msg.Copy()
				require.Equal(t, tt.msg.Header, cp.GetHeader())
			}
		})
	}
}
