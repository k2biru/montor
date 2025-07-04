package models

import "github.com/k2biru/montor/codec/hex"

type Control uint8

func (m *Control) Decode(pkt []byte, idx *int) error {
	*m = Control(hex.ReadByte(pkt, idx))
	// todo : support control argument
	return nil
}
func (m *Control) Encode() (pkt []byte, err error) {
	pkt = hex.WriteByte(pkt, m.GetID())
	// todo : support control argument
	// please edit Msg82.Encode if new error added
	return pkt, err
}
func (m Control) GetID() uint8 {
	return uint8(m)
}

// func (m *Control) GetParameter() map[string]any {
// 	return nil
// }
