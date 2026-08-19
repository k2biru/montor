package models

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHeaderGenerators(t *testing.T) {
	key := &rsa.PublicKey{}

	h1 := GenerateHeader("12345678901234567", 0x01, 0xFE)
	require.Equal(t, uint8(0x01), h1.Encription)
	require.Nil(t, h1.EncryptKey)

	h2 := GenerateHeaderWithKey("12345678901234567", 0x01, 0xFE, key)
	require.Equal(t, uint8(0x02), h2.Encription)
	require.Equal(t, key, h2.EncryptKey)

	h3 := GenerateHeaderWithKeyNoRSA("12345678901234567", 0x01, 0xFE, key)
	require.Equal(t, uint8(0x01), h3.Encription)
	require.Equal(t, key, h3.EncryptKey)
}

func TestMsgInterfacesAndCopy(t *testing.T) {
	hdr := GenerateHeader("12345678901234567", 0x01, 0xFE)

	messages := []GBT32960Msg{
		&Msg01{Header: hdr, SerialNumber: 1},
		&Msg02{Header: hdr},
		&Msg03{Header: hdr},
		&Msg04{Header: hdr, SerialNumber: 2},
		&Msg07{Header: hdr},
		&Msg80Receive{Header: hdr},
		&Msg80Reply{Header: hdr},
		&Msg81{Header: hdr},
		&Msg82{Header: hdr},
		&Msgxx{Header: hdr},
	}

	for _, msg := range messages {
		t.Run(msg.GetHeader().VIN, func(t *testing.T) {
			require.Equal(t, hdr, msg.GetHeader())
			_ = msg.GetMsgSN()

			cp := msg.Copy()
			require.NotNil(t, cp)
			require.Equal(t, msg.GetHeader().VIN, cp.GetHeader().VIN)
		})
	}
}

func TestReportStructures(t *testing.T) {
	t.Run("Report08 Encode and JSON", func(t *testing.T) {
		bv := &BatteryVoltages{
			AssyNo:                      1,
			Voltage:                     3000,
			Current:                     500,
			BatteriesTotalNumber:        2,
			BatteryStartNumberFrame:     1,
			SingleBatteryVoltageOnFrame: []uint16{3200, 3210},
		}

		encSingle, err := bv.Encode()
		require.NoError(t, err)
		require.NotEmpty(t, encSingle)

		bvs := BatteriesVoltages{bv}
		encAll, err := bvs.Encode()
		require.NoError(t, err)
		require.NotEmpty(t, encAll)

		jsonBytes, err := bvs.MarshalJSON()
		require.NoError(t, err)
		require.Contains(t, string(jsonBytes), "assamblyNo")

		// Decode error branches
		idx := 0
		var errBv BatteryVoltages
		require.ErrorIs(t, errBv.Decode([]byte{0x01}, &idx), ErrDecodeMsg)

		idx = 0
		var errBvs BatteriesVoltages
		require.ErrorIs(t, errBvs.Decode([]byte{0x01}, &idx), ErrDecodeMsg)
	})

	t.Run("Report09 Encode and JSON", func(t *testing.T) {
		bt := &BatteryTemperature{
			AssyNo:      1,
			Temperature: []uint8{50, 52},
		}

		encSingle, err := bt.Encode()
		require.NoError(t, err)
		require.NotEmpty(t, encSingle)

		bts := BatteriesTemperatures{bt}
		encAll, err := bts.Encode()
		require.NoError(t, err)
		require.NotEmpty(t, encAll)

		jsonBytes, err := bts.MarshalJSON()
		require.NoError(t, err)
		require.Contains(t, string(jsonBytes), "assamblyNo")

		// Decode error branches
		idx := 0
		var errBt BatteryTemperature
		require.ErrorIs(t, errBt.Decode([]byte{0x01}, &idx), ErrDecodeMsg)

		idx = 0
		var errBts BatteriesTemperatures
		require.ErrorIs(t, errBts.Decode([]byte{0x01}, &idx), ErrDecodeMsg)
	})

	t.Run("Report Item Helper Methods & Errors", func(t *testing.T) {
		reportsList := []GeneralReport{
			&VehicleData{},
			&DriveMotors{},
			&FuelCell{},
			&Engine{},
			&Location{},
			&Extreme{},
			&Alarm{},
		}

		for _, r := range reportsList {
			jsonBytes, err := r.MarshalJSON()
			require.NoError(t, err)
			require.NotEmpty(t, jsonBytes)
		}

		vd := VehicleData{}
		require.Equal(t, uint8(0x01), vd.GetID())

		reps := Reports{
			0x08: &BatteriesVoltages{},
			0x09: &BatteriesTemperatures{},
		}

		as08, err := reps.As08()
		require.NoError(t, err)
		require.NotNil(t, as08)

		as09, err := reps.As09()
		require.NoError(t, err)
		require.NotNil(t, as09)

		jsonBytes, err := reps.MarshalJSON()
		require.NoError(t, err)
		require.NotEmpty(t, jsonBytes)

		// Test invalid item ID decode
		idx := 0
		var rDecode Reports
		invalidPayload := []byte{0xFF, 0x01}
		require.Error(t, rDecode.Decode(invalidPayload, &idx))
	})
}

func TestMsgGetTimeAndGenerators(t *testing.T) {
	now := time.Now().UTC()

	m1 := &Msg01{Time: now}
	m2 := &Msg02{Time: now}
	m3 := &Msg03{Time: now}
	m4 := &Msg04{Time: now}
	m80Rpl := &Msg80Reply{Time: now}
	m81 := &Msg81{Time: now}

	require.Equal(t, now, m1.GetTime())
	require.Equal(t, now, m2.GetTime())
	require.Equal(t, now, m3.GetTime())
	require.Equal(t, now, m4.GetTime())
	require.Equal(t, now, m80Rpl.GetTime())
	require.Equal(t, now, m81.GetTime())

	m80Recv := GenerateMsh80Receive("12345678901234567", now, []byte{0x01})
	require.Equal(t, "12345678901234567", m80Recv.Header.VIN)
	require.Equal(t, now, m80Recv.Time)

	lookup := map[uint8]ParamProperties{
		0x01: {GenerateValue: func() any { return uint16(0) }},
	}
	SetParameterPropertiesLookup(lookup)

	bv := BatteriesVoltages{}
	require.Equal(t, uint8(0x08), bv.GetID())

	bt := BatteriesTemperatures{}
	require.Equal(t, uint8(0x09), bt.GetID())
}
