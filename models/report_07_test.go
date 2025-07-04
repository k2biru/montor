package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestAlarmDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Alarm
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("010000000203000000030000000300000003" +
					"0400000004000000040000000400000004" +
					"050000000500000005000000050000000500000005" +
					"06000000060000000600000006000000060000000600000006"),
				idx: &[]int{0}[0],
			},
			want: Alarm{
				AlarmLevel:         1,
				AlarmBatteryFlag:   2,
				AlarmBatteryOthers: []uint32{3, 3, 3},
				AlarmDriveMotor:    []uint32{4, 4, 4, 4},
				AlarmEngines:       []uint32{5, 5, 5, 5, 5},
				AlarmOthers:        []uint32{6, 6, 6, 6, 6, 6},
			},
			wantErr: false,
		},
		{
			name: "case: fail 1",
			args: args{
				pkt: hex.Str2Byte("010000000203000000030000000300000003"),
				idx: &[]int{0}[0],
			},
			want: Alarm{
				AlarmLevel:       1,
				AlarmBatteryFlag: 2,
			},
			wantErr: true,
		}, {
			name: "case: fail 2",
			args: args{
				pkt: hex.Str2Byte("010000000203000000030000000300000003" +
					"0400000004000000040000000400000004"),
				idx: &[]int{0}[0],
			},
			want: Alarm{
				AlarmLevel:         1,
				AlarmBatteryFlag:   2,
				AlarmBatteryOthers: []uint32{3, 3, 3},
			},
			wantErr: true,
		},
		{
			name: "case: fail 3",
			args: args{
				pkt: hex.Str2Byte("010000000203000000030000000300000003" +
					"0400000004000000040000000400000004" +
					"050000000500000005000000050000000500000005"),
				idx: &[]int{0}[0],
			},
			want: Alarm{
				AlarmLevel:         1,
				AlarmBatteryFlag:   2,
				AlarmBatteryOthers: []uint32{3, 3, 3},
				AlarmDriveMotor:    []uint32{4, 4, 4, 4},
			},
			wantErr: true,
		},

		{
			name: "case: fail 4",
			args: args{
				pkt: hex.Str2Byte("010000000203000000030000000300000003" +
					"0400000004000000040000000400000004" +
					"050000000500000005000000050000000500000005" +
					"060000000600000006000000060000000600000006000000"),
				idx: &[]int{0}[0],
			},
			want: Alarm{
				AlarmLevel:         1,
				AlarmBatteryFlag:   2,
				AlarmBatteryOthers: []uint32{3, 3, 3},
				AlarmDriveMotor:    []uint32{4, 4, 4, 4},
				AlarmEngines:       []uint32{5, 5, 5, 5, 5},
			},
			wantErr: true,
		},

		{
			name: "case: to sort",
			args: args{
				pkt: hex.Str2Byte("01020003040500060708090A0B"),
				idx: &[]int{0}[0],
			},
			want:    Alarm{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Alarm{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestAlarmDataEncode(t *testing.T) {
	type args struct {
		report Alarm
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				report: Alarm{
					AlarmLevel:         1,
					AlarmBatteryFlag:   2,
					AlarmBatteryOthers: []uint32{3, 3, 3},
					AlarmDriveMotor:    []uint32{4, 4, 4, 4},
					AlarmEngines:       []uint32{5, 5, 5, 5, 5},
					AlarmOthers:        []uint32{6, 6, 6, 6, 6, 6},
				},
			},
			want: hex.Str2Byte("010000000203000000030000000300000003" +
				"0400000004000000040000000400000004" +
				"050000000500000005000000050000000500000005" +
				"06000000060000000600000006000000060000000600000006"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.report.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}
