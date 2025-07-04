package models

type ProcessData struct {
	Incoming GBT32960Msg // received message
	Outgoing GBT32960Msg // message sent
}
