package models

import "github.com/pkg/errors"

var (
	ErrInvalidParse   = errors.New("Fail to parse")
	ErrInvalidConvert = errors.New("Fail to convert msg")
	ErrDecodeMsg      = errors.New("Fail to decode msg")
	ErrEncodeMsg      = errors.New("Fail to encode msg")
	ErrGenOutgoingMsg = errors.New("Fail to generate outgoing msg")
	ErrNotAvailable   = errors.New("Not available")
	ErrMissmatch      = errors.New("Missmatch")
	ErrInvalidID      = errors.New("Invalid ID")
	ErrInvalidLength  = errors.New("Invalid length")
	ErrPreprocess     = errors.New("Invalid preprocess")
	ErrNotFound       = errors.New("Not found")
	ErrUnset          = errors.New("Unset")
	ErrDecryptMsg     = errors.New("Fail to decrypt msg")
	ErrEncryptMsg     = errors.New("Fail to encrypt msg")
)

type MsgComponent interface {
	Decode(pkt []byte, idx *int) error
	Encode() (pkt []byte)
}
