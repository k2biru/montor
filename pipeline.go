package montor

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/k2biru/montor/models"
)

type Component interface {
	GetProcessOption() ProcessOps
}

type Pipeline interface {
	ProcessRead(ctx context.Context) error
	ProcessWrite(ctx context.Context, processData *models.ProcessData) error
}

func NewPipeline(conn net.Conn, c Component) Pipeline {
	return &pipeline{
		fh: NewFrameHandler(conn),
		pc: NewPacketCodec(c.GetProcessOption()),
		mp: NewProcessor(c.GetProcessOption()),
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

	if err := m.fh.Send(pkt); err != nil {
		// log.Warn().Err(err).
		// 	Msg("Failed to send dashboard frame")
	}
	return nil
}
func (m *pipeline) ProcessRead(ctx context.Context) error {
	pkt, err := m.fh.Recv(ctx)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil // EOF, stop read buffer
		}
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

	if err := m.fh.Send(pkt); err != nil {
		// log.Error().Err(err).
		// 	Msg("Failed to send dashboard frame")
	}
	return nil
}
