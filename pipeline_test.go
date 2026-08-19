package montor_test

import (
	"bytes"
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/k2biru/montor"
	"github.com/k2biru/montor/codec/hex"
	"github.com/k2biru/montor/models"
	"github.com/stretchr/testify/require"
)

type mockConn struct {
	*bytes.Buffer
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

type mockPipeHooks struct {
	getProcessFn func(id uint8) (*montor.Action, error)
	preProcessFn func(ctx context.Context, msg models.GBT32960Msg) (context.Context, error)
	postDecodeFn func(msg models.GBT32960Msg)
	decryptFn    func(encType uint8, vin string, pkt []byte) ([]byte, error)
	postRecvFn   func(pkt []byte) ([]byte, error)
	postSendFn   func(pkt []byte)
}

func (m *mockPipeHooks) GetProcess(id uint8) (*montor.Action, error) {
	if m.getProcessFn != nil {
		return m.getProcessFn(id)
	}
	return nil, errors.New("unsupported command")
}

func (m *mockPipeHooks) PreProcess(ctx context.Context, msg models.GBT32960Msg) (context.Context, error) {
	if m.preProcessFn != nil {
		return m.preProcessFn(ctx, msg)
	}
	return ctx, nil
}

func (m *mockPipeHooks) PostDecode(msg models.GBT32960Msg) {
	if m.postDecodeFn != nil {
		m.postDecodeFn(msg)
	}
}

func (m *mockPipeHooks) Decrypt(encType uint8, vin string, pkt []byte) ([]byte, error) {
	if m.decryptFn != nil {
		return m.decryptFn(encType, vin, pkt)
	}
	return pkt, nil
}

func (m *mockPipeHooks) PostRecvFrame(in []byte) ([]byte, error) {
	if m.postRecvFn != nil {
		return m.postRecvFn(in)
	}
	return in, nil
}

func (m *mockPipeHooks) PostSendFrame(pkt []byte) {
	if m.postSendFn != nil {
		m.postSendFn(pkt)
	}
}

func TestPipeline_ProcessRead(t *testing.T) {
	sampleMsgBytes := hex.Str2Byte("232307FE3132333435363738393031323334353637000000c8")

	tests := []struct {
		name    string
		stream  []byte
		hooks   *mockPipeHooks
		wantErr bool
	}{
		{
			name:   "case: successful process read and reply send",
			stream: sampleMsgBytes,
			hooks: &mockPipeHooks{
				getProcessFn: func(id uint8) (*montor.Action, error) {
					return &montor.Action{
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
		},
		{
			name:   "case: frame handler recv error (empty stream)",
			stream: []byte{},
			hooks: &mockPipeHooks{},
			wantErr: true,
		},
		{
			name:   "case: packet codec decode error (corrupt checksum)",
			stream: hex.Str2Byte("232307FE313233343536373839303132333435363700000006180101010101FF"),
			hooks:  &mockPipeHooks{},
			wantErr: true,
		},
		{
			name:   "case: processor process error (unsupported command)",
			stream: sampleMsgBytes,
			hooks: &mockPipeHooks{
				getProcessFn: func(id uint8) (*montor.Action, error) {
					return nil, errors.New("unsupported command ID")
				},
			},
			wantErr: true,
		},
		{
			name:   "case: no reply outgoing data (returns nil without send)",
			stream: sampleMsgBytes,
			hooks: &mockPipeHooks{
				getProcessFn: func(id uint8) (*montor.Action, error) {
					return &montor.Action{
						ForceProcess: true,
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: nil,
							}
						},
					}, nil
				},
			},
			wantErr: false,
		},
		{
			name:   "case: postRecvFn error in frame handler",
			stream: sampleMsgBytes,
			hooks: &mockPipeHooks{
				postRecvFn: func(pkt []byte) ([]byte, error) {
					return nil, errors.New("post recv error")
				},
			},
			wantErr: true,
		},
		{
			name:   "case: encode error on outgoing reply",
			stream: sampleMsgBytes,
			hooks: &mockPipeHooks{
				getProcessFn: func(id uint8) (*montor.Action, error) {
					return &montor.Action{
						GenData: func() *models.ProcessData {
							return &models.ProcessData{
								Incoming: &models.Msg07{},
								Outgoing: &models.Msg07{Header: nil},
							}
						},
						Process: func(ctx context.Context, pd *models.ProcessData) error {
							return nil
						},
					}, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConn{Buffer: bytes.NewBuffer(tt.stream)}
			pipe := montor.NewPipeline(conn, tt.hooks)
			err := pipe.ProcessRead(context.Background())
			require.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}

func TestPipeline_ProcessWrite(t *testing.T) {
	tests := []struct {
		name        string
		processData *models.ProcessData
		wantErr     bool
		wantOutput  bool
	}{
		{
			name:        "case: nil process data",
			processData: nil,
			wantErr:     false,
			wantOutput:  false,
		},
		{
			name: "case: nil outgoing msg",
			processData: &models.ProcessData{
				Incoming: &models.Msg07{},
				Outgoing: nil,
			},
			wantErr:    false,
			wantOutput: false,
		},
		{
			name: "case: valid outgoing write",
			processData: &models.ProcessData{
				Outgoing: &models.Msg07{
					Header: &models.MsgHeader{
						CommandID:  0x07,
						Response:   0x01,
						VIN:        "12345678901234567",
						Encription: 0x00,
						BodyLength: 0,
					},
				},
			},
			wantErr:    false,
			wantOutput: true,
		},
		{
			name: "case: encode error due to missing header",
			processData: &models.ProcessData{
				Outgoing: &models.Msg07{
					Header: nil,
				},
			},
			wantErr:    true,
			wantOutput: false,
		},
	}

	hooks := &mockPipeHooks{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConn{Buffer: bytes.NewBuffer(nil)}
			pipe := montor.NewPipeline(conn, hooks)
			err := pipe.ProcessWrite(context.Background(), tt.processData)
			require.Equal(t, tt.wantErr, err != nil, err)
			if tt.wantOutput {
				require.Greater(t, conn.Len(), 0)
			}
		})
	}
}
