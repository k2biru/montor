package gbt32960

import (
	"testing"
	"time"

	"github.com/k2biru/montor/models"
	"github.com/stretchr/testify/require"
)

func TestPrepareReply(t *testing.T) {
	type args struct {
		data *models.ProcessData
	}
	tests := []struct {
		name    string
		args    args
		want    models.GBT32960Msg
		wantErr bool
	}{
		// {
		// 	name: "case: same",
		// 	args: args{
		// 		data: &models.ProcessData{
		// 			Incoming: &models.Msg07{
		// 				Header: models.GenerateHeader("12345678901234567", 0x07, 0xFE),
		// 			},
		// 			Outgoing: &models.Msg07{},
		// 		},
		// 	},
		// 	want: &models.Msg07{
		// 		Header: models.GenerateHeader("12345678901234567", 0x07, 0x01),
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "case: not same",
			args: args{
				data: &models.ProcessData{
					Incoming: &models.Msg02{
						Header: models.GenerateHeader("12345678901234567", 0x02, 0xFE),
						Time:   time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
					},
					Outgoing: &models.MsgReply{},
				},
			},
			want: &models.MsgReply{
				Header: models.GenerateHeader("12345678901234567", 0x02, 0x01),
				Time:   time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareReply(tt.args.data)
			require.Equal(t, tt.want, tt.args.data.Outgoing)
		})
	}
}
