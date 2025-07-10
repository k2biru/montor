package montor

import (
	"time"

	"github.com/k2biru/montor/models"

	"github.com/pkg/errors"
)

var (
	ErrEmptyPacket  = errors.New("Empty packet")
	ErrVerifyFailed = errors.New("Verify failed")
	ErrEncodeType   = errors.New("Error data type")
	ErrActiveClose  = errors.New("Active close")
)

type PacketCodecHooks interface {
	Decrypt(encType uint8, vin string, pkt []byte) ([]byte, error)
}

type PacketCodec interface {
	Decode([]byte) (*models.PacketData, error)
	Encode(models.GBT32960Msg) ([]byte, error)
}

func NewPacketCodec(hooks PacketCodecHooks) PacketCodec {
	return &packetCodec{
		timeNow: time.Now,
		hooks:   hooks,
	}
}

type packetCodec struct {
	timeNow func() time.Time
	hooks   PacketCodecHooks
}

func (m packetCodec) genVerifier(pkt []byte) []byte {
	var code byte
	for _, v := range pkt {
		code ^= v
	}
	pkt = append(pkt, code)
	return pkt
}

func (m packetCodec) verify(pkt []byte) ([]byte, error) {
	n := len(pkt)
	if n == 0 {
		return nil, ErrEmptyPacket
	}

	expected := pkt[n-1]
	// verify using BCC algoritm
	var actual byte = 0
	for _, b := range pkt[:n-1] {
		actual ^= b
	}
	if expected != actual {
		return nil, errors.Wrapf(ErrVerifyFailed, "verify expect=%v, but actual=%v", expected, actual)
	}

	return pkt[:n-1], nil
}

func (m *packetCodec) Decode(payload []byte) (*models.PacketData, error) {
	// verify
	pkt, err := m.verify(payload)
	if err != nil {
		return nil, err
	}
	verifyCode := payload[len(payload)-1]

	// decode
	pd := &models.PacketData{
		Header:     &models.MsgHeader{},
		VerifyCode: verifyCode,
	}

	err = pd.Header.Decode(pkt)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to decode packet")
	}

	pd.Header.TimeCreated = m.timeNow()
	pd.Body = pkt[pd.Header.Idx:]

	// post processor for decrypt
	if pd.Header.Encription != models.EncryptNone {
		if m.hooks == nil {
			return nil, errors.Wrap(models.ErrUnset, "Decrypt not set")
		}
		pd.Body, err = m.hooks.Decrypt(pd.Header.Encription, pd.Header.VIN, pkt[pd.Header.Idx:])
		if err != nil {
			return nil, errors.Wrap(err, "Fail to decrypt packet")
		}
	} else {
		pd.Body = pkt[pd.Header.Idx:]
	}

	pd.Header.Idx = 0 // reset idx
	return pd, nil

}
func (m *packetCodec) Encode(data models.GBT32960Msg) ([]byte, error) {
	pkt, err := data.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to encode msg")
	}
	pkt = m.genVerifier(pkt)
	return pkt, nil
}
