package models

import (
	"crypto/rsa"
	"time"

	"github.com/k2biru/montor/codec/hex"
)

const (
	ResponseSuccess        uint8 = 0x01
	ResponseFailed         uint8 = 0x02
	ResponseVINDuplicate   uint8 = 0x03
	ResponseCommandPackage uint8 = 0xFE

	headerLength = 22

	EncryptNone uint8 = 0x01
	EncryptRSA  uint8 = 0x02
)

const (
	VINLength = 17
)

func GenerateHeader(vin string, cmdID, response uint8) *MsgHeader {
	return &MsgHeader{
		VIN:        vin,
		CommandID:  cmdID,
		Response:   response,
		Encription: EncryptNone,
	}
}

func GenerateHeaderWithKey(vin string, cmdID, response uint8, rsaPublic *rsa.PublicKey) *MsgHeader {
	return &MsgHeader{
		VIN:        vin,
		CommandID:  cmdID,
		Response:   response,
		Encription: EncryptRSA,
		EncryptKey: rsaPublic,
	}
}

func GenerateHeaderWithKeyNoRSA(vin string, cmdID, response uint8, rsaPublic *rsa.PublicKey) *MsgHeader {
	return &MsgHeader{
		VIN:        vin,
		CommandID:  cmdID,
		Response:   response,
		Encription: EncryptNone,
		EncryptKey: rsaPublic,
	}
}

type MsgHeader struct {
	CommandID         uint8           `json:"commandID"`    // command id
	Response          uint8           `json:"responseFlag"` // command response flag
	VIN               string          `json:"VIN"`          // 17 char string
	Encription        uint8           `json:"encription"`   // metode enkripsi
	BodyLength        uint16          `json:"length"`       // panjang body
	Idx               int             `json:"-"`            // Read the PACKET Header.
	TimeCreated       time.Time       `json:"timeCreated"`  // Time msg create
	EncryptKey        *rsa.PublicKey  `json:"encryptKey"`   // key for encrypt
	PrivateEncryptKey *rsa.PrivateKey `json:"privateKey"`   // key for encrypt
}

// Decoding [] byte into a message head structure
func (m *MsgHeader) Decode(pkt []byte) error {
	if len(pkt) < headerLength {
		return ErrDecodeMsg
	}
	var idx int
	m.CommandID = hex.ReadByte(pkt, &idx)
	m.Response = hex.ReadByte(pkt, &idx)
	m.VIN = hex.ReadString(pkt, &idx, VINLength)
	m.Encription = hex.ReadByte(pkt, &idx)
	m.BodyLength = hex.ReadWord(pkt, &idx)
	m.Idx = idx

	return nil
}

// Code the message head structure into [] byte
func (m *MsgHeader) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.CommandID)
	pkt = hex.WriteByte(pkt, m.Response)
	pkt = hex.WriteString(pkt, m.VIN, VINLength)
	pkt = hex.WriteByte(pkt, m.Encription)
	pkt = hex.WriteWord(pkt, m.BodyLength)
	return pkt, nil
}
