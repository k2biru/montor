package montor

import (
	"context"
	"reflect"

	"github.com/k2biru/montor/models"

	"github.com/pkg/errors"
)

var (
	ErrMsgInvalidResponseFlag = errors.New("msg response flag invalid")
	ErrMsgNotSupportted       = errors.New("Msg is not supportted")
	ErrNotAuthorized          = errors.New("Not authorized")
	// ErrActiveClose = errors.New("Active close")
)

type Processor interface {
	Process(ctx context.Context, pkt *models.PacketData) (*models.ProcessData, error)
}

func NewProcessor(option ProcessOps) Processor {
	return &processor{
		procOpt: option,
	}
}

type processor struct {
	procOpt ProcessOps
}

func (m *processor) Process(ctx context.Context, pkt *models.PacketData) (*models.ProcessData, error) {
	//is support?
	cmdID := pkt.Header.CommandID
	act, err := m.procOpt.GetProcess(cmdID)
	if err != nil {
		return nil, errors.Wrapf(ErrMsgNotSupportted, "commandID 0x%02x", cmdID)
	}
	// get act
	genDataFn := act.GenData
	if genDataFn == nil {
		return nil, errors.Wrapf(ErrMsgNotSupportted, "no action for commandID 0x%02x", cmdID)
	}
	// prepare msg
	data := genDataFn()
	if data.Incoming == nil {
		return nil, errors.Wrapf(ErrMsgNotSupportted, "no parser for commandID 0x%02x", cmdID)
	}

	// prepare incoming msg
	in := data.Incoming

	// decode incoming msg
	err = in.Decode(pkt)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to decode packet to gbt32960")
	}
	// action in post decode
	m.procOpt.PostDecode(in)

	if r := in.GetHeader().Response; !act.IsReply && r != 0xFE {
		return nil, errors.Wrapf(ErrMsgInvalidResponseFlag, "expect 0xFE got 0x%02x", r)
	}

	// prepare outgoing msg
	if !act.IsReply {
		prepareReply(data)
	}
	// process
	process := act.Process
	if process == nil {
		return data, nil
	}
	ctx, err = m.procOpt.PreProcess(ctx, data.Incoming)
	if err != nil {
		// log.Error().Err(err).Str("vin", data.Incoming.GetHeader().VIN).Msg("Failed preprocess")
		data.Outgoing.GetHeader().Response = 0x02
		return data, nil
	}
	err = process(ctx, data)
	if err != nil {
		return data, errors.Wrap(err, "Fail to process data")
	}

	return data, nil
}

func prepareReply(data *models.ProcessData) {
	if reflect.TypeOf(data.Incoming) == reflect.TypeOf(data.Outgoing) {
		data.Outgoing = data.Incoming.Copy()
		data.Outgoing.GetHeader().Response = 0x01 // respose ok
		return
	}
	if reply, ok := data.Outgoing.(*models.MsgReply); ok {
		// copy the header
		inHeader := *data.Incoming.GetHeader()
		reply.Header = &inHeader
		reply.Header.Response = 0x01 // response ok
		if at, ok := data.Incoming.(models.GBT32960MsgTime); ok {
			reply.Time = at.GetTime()
		}
	}
}

type ChannelGw interface {
	SendDeviceChannel(id, src string, opt1, opt2 uint16, val any) error
}

type BroadcastCh interface {
	SendBroadcastCh(id, src string, opt1, opt2 uint8, val models.GBT32960Msg) error
}
