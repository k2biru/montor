package gbt32960

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type Frame []byte
type FrameHandler interface {
	Recv(ctx context.Context) (Frame, error)
	Send(Frame) error
}

func NewFrameHandler(buffer io.ReadWriter) FrameHandler {
	return &frameHandler{
		rbuf:   bufio.NewReader(buffer),
		writer: buffer,
	}
}

type frameHandler struct {
	rbuf   *bufio.Reader
	writer io.Writer
}

func (m *frameHandler) Recv(ctx context.Context) (Frame, error) {
	rawBuf := make([]byte, 0, 225)
	buffer := make([]byte, 1)

	// read header

	index := 0

	// lookupCommand := generalAction()
	lengthBuff := make([]byte, 0, 2)
	expectedLength := 0
	startByte := 2
	/////////////
	for {
		_, err := m.rbuf.Read(buffer)
		if err != nil {
			return nil, errors.Wrap(err, "Fail to read stream to framePayload")
		}

		// read 2 byte "##"
		if buffer[0] == 0x23 && startByte > 0 {
			startByte--
			continue
		}
		index++
		rawBuf = append(rawBuf, buffer[0])
		switch index {
		case 21:
			// get expexted length
			lengthBuff = lengthBuff[:0]
			lengthBuff = append(lengthBuff, buffer...)
		case 22:
			lengthBuff = append(lengthBuff, buffer...)
			expectedLength = int(binary.BigEndian.Uint16(lengthBuff)) + 1
		}
		if index > 22 {
			expectedLength--
			if expectedLength <= 0 {
				// reset buffer
				// log.Debug().Int("frame_len", len(rawBuf)).
				// 	Hex("frame", rawBuf).
				// 	Msg("Recv frame")
				return rawBuf, nil
			}
		}
	}
}

func (m *frameHandler) Send(frame Frame) error {
	if len(frame) == 0 {
		// log.Warn().Msg("Frame error, with zero len")
		return nil
	}
	// add start byte
	pkt := []byte{
		0x23, 0x23,
	}
	pkt = append(pkt, frame...)
	for {
		n, err := m.writer.Write(pkt)
		if err != nil {
			return errors.Wrap(err, "Failed to send payload")
		} else if n >= len(pkt) {
			break
		}

		pkt = pkt[n:] // Haven't finished all the data, write again
	}
	// log.Debug().
	// 	Int("frame_len", len(frame)).
	// 	Hex("frame_payload", frame).
	// 	Msg("Sent frame.")
	return nil
}
