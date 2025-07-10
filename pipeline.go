package montor

import (
	"context"
	"net"

	"github.com/k2biru/montor/models"
)

type Pipeline interface {
	ProcessRead(ctx context.Context) error
	ProcessWrite(ctx context.Context, processData *models.ProcessData) error
}

type PipeHooks interface {
	ProcesssHooks
	PacketCodecHooks
	FrameHandlerHooks
}

func NewPipeline(conn net.Conn, po PipeHooks) Pipeline {
	return &pipeline{
		fh: NewFrameHandler(conn, po),
		pc: NewPacketCodec(po),
		mp: NewProcessor(po),
	}
}

type pipeline struct {
	fh FrameHandler
	pc PacketCodec
	mp Processor
}

func (m *pipeline) ProcessWrite(ctx context.Context, processData *models.ProcessData) error {
	if processData == nil || processData.Outgoing == nil { // No need to reply, no follow -up processing
		return nil
	}
	pkt, err := m.pc.Encode(processData.Outgoing)
	if err != nil {
		return err
	}

	return m.fh.Send(pkt)
}
func (m *pipeline) ProcessRead(ctx context.Context) error {
	pkt, err := m.fh.Recv(ctx)
	if err != nil {
		return err
	}

	pd, err := m.pc.Decode(pkt)
	if err != nil {
		return err
	}

	dm, err := m.mp.Process(ctx, pd)
	if err != nil || dm.Outgoing == nil {
		return err // If error or no reply (if no reply err and dm == nil)
	}

	pkt, err = m.pc.Encode(dm.Outgoing)
	if err != nil {
		return err
	}

	return m.fh.Send(pkt)
}
