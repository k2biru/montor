package montor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/k2biru/montor/models"
	"github.com/stretchr/testify/require"
)

type mockProcessHooks struct {
	getProcessFn func(id uint8) (*Action, error)
	preProcessFn func(ctx context.Context, msg models.GBT32960Msg) (context.Context, error)
	postDecodeFn func(msg models.GBT32960Msg)
}

func (m *mockProcessHooks) GetProcess(id uint8) (*Action, error) {
	if m.getProcessFn != nil {
		return m.getProcessFn(id)
	}
	return nil, ErrMsgNotSupportted
}

func (m *mockProcessHooks) PreProcess(ctx context.Context, msg models.GBT32960Msg) (context.Context, error) {
	if m.preProcessFn != nil {
		return m.preProcessFn(ctx, msg)
	}
	return ctx, nil
}

func (m *mockProcessHooks) PostDecode(msg models.GBT32960Msg) {
	if m.postDecodeFn != nil {
		m.postDecodeFn(msg)
	}
}

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
		{
			name: "case: same incoming and outgoing type",
			args: args{
				data: &models.ProcessData{
					Incoming: &models.Msg07{
						Header: models.GenerateHeader("12345678901234567", 0x07, 0xFE),
					},
					Outgoing: &models.Msg07{},
				},
			},
			want: &models.Msg07{
				Header: models.GenerateHeader("12345678901234567", 0x07, 0x01),
			},
			wantErr: false,
		},
		{
			name: "case: not same incoming and outgoing type",
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

func TestProcessor_Process(t *testing.T) {
	validHeader := &models.MsgHeader{
		CommandID:  0x07,
		Response:   0xFE,
		VIN:        "12345678901234567",
		Encription: 0x00,
		BodyLength: 0,
	}

	validPacket := &models.PacketData{
		Header:     validHeader,
		Body:       []byte{},
		VerifyCode: 0xC8,
	}

	tests := []struct {
		name      string
		pkt       *models.PacketData
		hooks     *mockProcessHooks
		wantErr   bool
		checkData func(t *testing.T, res *models.ProcessData)
	}{
		{
			name: "case: unsupported command ID",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return nil, errors.New("unsupported")
				},
			},
			wantErr: true,
		},
		{
			name: "case: nil GenData",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{GenData: nil}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "case: nil Incoming data",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{Incoming: nil}
						},
					}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "case: preProcess error sets response to 0x02",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: &models.Msg07{},
							}
						},
					}, nil
				},
				preProcessFn: func(ctx context.Context, msg models.GBT32960Msg) (context.Context, error) {
					return ctx, errors.New("preprocess error")
				},
			},
			wantErr: false,
			checkData: func(t *testing.T, res *models.ProcessData) {
				require.NotNil(t, res)
				require.Equal(t, uint8(0x02), res.Outgoing.GetHeader().Response)
			},
		},
		{
			name: "case: preProcess error with nil Outgoing",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: nil,
							}
						},
					}, nil
				},
				preProcessFn: func(ctx context.Context, msg models.GBT32960Msg) (context.Context, error) {
					return ctx, errors.New("preprocess error")
				},
			},
			wantErr: false,
			checkData: func(t *testing.T, res *models.ProcessData) {
				require.NotNil(t, res)
				require.Nil(t, res.Outgoing)
			},
		},
		{
			name: "case: invalid response flag without force process",
			pkt: &models.PacketData{
				Header: &models.MsgHeader{
					CommandID: 0x07,
					Response:  0x01, // invalid (expected 0xFE)
				},
			},
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						ForceProcess: false,
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: &models.Msg07{},
							}
						},
					}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "case: process function returns error",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: &models.Msg07{},
							}
						},
						Process: func(ctx context.Context, pd *models.ProcessData) error {
							return errors.New("process failed")
						},
					}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "case: successful process",
			pkt:  validPacket,
			hooks: &mockProcessHooks{
				getProcessFn: func(id uint8) (*Action, error) {
					return &Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: &models.Msg07{},
							}
						},
						Process: func(ctx context.Context, pd *models.ProcessData) error {
							return nil
						},
					}, nil
				},
			},
			wantErr: false,
			checkData: func(t *testing.T, res *models.ProcessData) {
				require.NotNil(t, res)
				require.Equal(t, uint8(0x01), res.Outgoing.GetHeader().Response)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(tt.hooks)
			res, err := p.Process(context.Background(), tt.pkt)
			require.Equal(t, tt.wantErr, err != nil, err)
			if tt.checkData != nil {
				tt.checkData(t, res)
			}
		})
	}
}
