package montor_test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/k2biru/montor"
	"github.com/k2biru/montor/codec/hex"
	"github.com/k2biru/montor/models"
)

type benchConn struct {
	readBuf  []byte
	readPos  int
	writeBuf *bytes.Buffer
}

func newBenchConn(raw []byte) *benchConn {
	return &benchConn{
		readBuf:  raw,
		writeBuf: bytes.NewBuffer(make([]byte, 0, 1024)),
	}
}

func (b *benchConn) Read(p []byte) (n int, err error) {
	if b.readPos >= len(b.readBuf) {
		b.readPos = 0 // loop stream for continuous benchmarking
	}
	n = copy(p, b.readBuf[b.readPos:])
	b.readPos += n
	return n, nil
}

func (b *benchConn) Write(p []byte) (n int, err error) {
	b.writeBuf.Reset()
	return b.writeBuf.Write(p)
}

func (b *benchConn) Close() error                       { return nil }
func (b *benchConn) LocalAddr() net.Addr                { return nil }
func (b *benchConn) RemoteAddr() net.Addr               { return nil }
func (b *benchConn) SetDeadline(t time.Time) error      { return nil }
func (b *benchConn) SetReadDeadline(t time.Time) error  { return nil }
func (b *benchConn) SetWriteDeadline(t time.Time) error { return nil }

func BenchmarkPipeline_ProcessRead(b *testing.B) {
	// Valid Msg07 frame with header + checksum: 232307FE3132333435363738393031323334353637000000c8
	sampleFrame := hex.Str2Byte("232307FE3132333435363738393031323334353637000000c8")
	conn := newBenchConn(sampleFrame)
	hooks := &mockPipeHooks{
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
	}
	pipe := montor.NewPipeline(conn, hooks)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pipe.ProcessRead(ctx)
	}
}

func BenchmarkPipeline_ProcessWrite(b *testing.B) {
	conn := newBenchConn(nil)
	hooks := &mockPipeHooks{}
	pipe := montor.NewPipeline(conn, hooks)
	ctx := context.Background()
	processData := &models.ProcessData{
		Outgoing: &models.Msg07{
			Header: &models.MsgHeader{
				CommandID:  0x07,
				Response:   0x01,
				VIN:        "12345678901234567",
				Encription: 0x00,
				BodyLength: 0,
			},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pipe.ProcessWrite(ctx, processData)
	}
}

func BenchmarkPacketCodec_Decode(b *testing.B) {
	payload := hex.Str2Byte("01FE313233343536373839000000000000000001002A130116173738000131323334353637383939383736353433323130300304313233343435363739383730FD")
	pc := montor.NewPacketCodec(nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = pc.Decode(payload)
	}
}

func BenchmarkPacketCodec_Encode(b *testing.B) {
	pc := montor.NewPacketCodec(nil)
	msg := &models.Msg01{
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
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = pc.Encode(msg)
	}
}

func BenchmarkFrameHandler_Recv(b *testing.B) {
	sampleFrame := hex.Str2Byte("232307FE3132333435363738393031323334353637000000c8")
	conn := newBenchConn(sampleFrame)
	fh := montor.NewFrameHandler(conn, &hooks{})
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = fh.Recv(ctx)
	}
}

func BenchmarkFrameHandler_Send(b *testing.B) {
	conn := newBenchConn(nil)
	fh := montor.NewFrameHandler(conn, &hooks{})
	frame := montor.Frame(hex.Str2Byte("010203040506"))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fh.Send(frame)
	}
}

func BenchmarkProcessor_Process(b *testing.B) {
	validHeader := &models.MsgHeader{
		CommandID:  0x07,
		Response:   0xFE,
		VIN:        "12345678901234567",
		Encription: 0x00,
		BodyLength: 0,
	}
	pkt := &models.PacketData{
		Header:     validHeader,
		Body:       []byte{},
		VerifyCode: 0xC8,
	}

	hooks := &mockPipeHooks{
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
	}
	proc := montor.NewProcessor(hooks)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = proc.Process(ctx, pkt)
	}
}
