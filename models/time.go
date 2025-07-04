package models

import (
	"time"

	"github.com/k2biru/montor/codec/hex"

	"github.com/pkg/errors"
)

const (
	GBT32960TimeLayout = "060102150405" //"YYMMDDhhmmss"
	timeMinLen         = 6
)

func GBT32960Timezone() *time.Location {
	return time.FixedZone("UTC+8", 8*60*60)
}

func TimeDecode(pkt []byte, idx *int) (time.Time, error) {
	t := time.Time{}
	if l := len(pkt[*idx:]); l < timeMinLen {
		return t, errors.Wrapf(ErrDecodeMsg, "to short, expect %d got %d", timeMinLen, l)
	}
	year := int(hex.ReadByte(pkt, idx))
	month := int(hex.ReadByte(pkt, idx))
	date := int(hex.ReadByte(pkt, idx))
	hour := int(hex.ReadByte(pkt, idx))
	minute := int(hex.ReadByte(pkt, idx))
	second := int(hex.ReadByte(pkt, idx))
	if year < 0 || year > 99 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong year %d", year)
	}
	year += 2000

	if month < 1 || month > 12 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong month %d", month)
		//err
	}
	if date < 1 || date > 31 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong date %d", date)
		//err
	}
	if hour < 0 || hour > 23 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong hour %d", hour)
	}
	if minute < 0 || minute > 59 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong minute %d", minute)
	}

	if second < 0 || second > 59 {
		return t, errors.Wrapf(ErrDecodeMsg, "wrong second %d", second)
	}
	return time.Date(year, time.Month(month),
		date, hour, minute, second, 0, GBT32960Timezone()).UTC(), nil
}

func TimeEncode(t time.Time) (pkt []byte) {
	t = t.In(GBT32960Timezone())
	year := uint8(t.Year() - 2000)
	month := uint8(t.Month())
	date := uint8(t.Day())
	hour := uint8(t.Hour())
	minute := uint8(t.Minute())
	second := uint8(t.Second())
	pkt = hex.WriteByte(pkt, year)
	pkt = hex.WriteByte(pkt, month)
	pkt = hex.WriteByte(pkt, date)
	pkt = hex.WriteByte(pkt, hour)
	pkt = hex.WriteByte(pkt, minute)
	pkt = hex.WriteByte(pkt, second)
	return pkt
}
