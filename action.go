package gbt32960

import (
	"context"

	"github.com/k2biru/montor/models"
)

type Action struct {
	IsReply bool
	GenData func() *models.ProcessData
	Process func(context.Context, *models.ProcessData) error
}

// type ProcessOption map[uint8]*Action
type ProcessOps interface {
	GetProcess(id uint8) (*Action, error)
	PreProcess(context.Context, models.GBT32960Msg) (context.Context, error)
	PostDecode(models.GBT32960Msg)
	Decrypt(encType uint8, vin string, pkt []byte) ([]byte, error)
}
