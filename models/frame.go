package models

import (
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type PacketData struct {
	Header     *MsgHeader
	Body       []byte
	VerifyCode uint8
}

func WriteHeader(m GBT32960Msg, pkt []byte) ([]byte, error) {
	header := m.GetHeader()
	if header == nil {
		return nil, errors.Wrapf(ErrEncodeMsg, "header not set of %s", reflect.TypeOf(m))
	}
	var err error
	body := pkt
	header.BodyLength = uint16(len(body))
	headerPkt, err := header.Encode()
	if err != nil {
		return nil, err
	}

	return append(headerPkt, body...), nil
}

type GBT32960Msg interface {
	Decode(*PacketData) error        // Packet -> GBT32960Msg
	Encode() (pkt []byte, err error) // GBT32960Msg -> Packet
	GetHeader() *MsgHeader           // Header
	Copy() GBT32960Msg
	GetMsgSN() string
	// GenOutgoing(incoming GBT32960Msg) error
}
type GBT32960MsgTime interface {
	GetTime() time.Time
}
