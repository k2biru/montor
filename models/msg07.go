package models

type Msg07 struct {
	Header *MsgHeader `json:"header"`
}

func (m *Msg07) Decode(packet *PacketData) error {
	m.Header = packet.Header
	return nil
}

func (m *Msg07) Encode() (pkt []byte, err error) {
	return WriteHeader(m, pkt)
}

func (m *Msg07) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg07) Copy() GBT32960Msg {
	header := *m.Header
	cp := *m
	cp.Header = &header
	return &cp
}
func (m Msg07) GetMsgSN() string {
	return ""
}
