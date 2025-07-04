package models

import (
	"testing"

	"github.com/k2biru/montor/codec/hex"

	"github.com/stretchr/testify/require"
)

func TestParametersDecode(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name    string
		args    args
		want    Parameters
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				pkt: hex.Str2Byte("010001020002030003" +
					"04040534343434" +
					"060006073737373737083838383838" +
					"09090A000A0B000B0C0C" +
					"0D040E656565650F000F1010"),
				idx: &[]int{0}[0],
			},
			want: Parameters{
				0x01: uint16(1),
				0x02: uint16(2),
				0x03: uint16(3),
				0x04: uint8(4),
				0x05: []uint8(hex.Str2Byte("34343434")),
				0x06: uint16(6),
				0x07: string("77777"),
				0x08: string("88888"),
				0x09: uint8(9),
				0x0A: uint16(10),
				0x0B: uint16(11),
				0x0C: uint8(12),
				0x0D: uint8(4),
				0x0E: []uint8(hex.Str2Byte("65656565")),
				0x0F: uint16(15),
				0x10: uint8(16),
			},
			wantErr: false,
		},
		{
			name: "case: unknow",
			args: args{
				pkt: hex.Str2Byte("F10003"),
				idx: &[]int{0}[0],
			},
			want:    Parameters{},
			wantErr: true,
		},
		{
			name: "case: unknowlength",
			args: args{
				pkt: hex.Str2Byte("0534343434"),
				idx: &[]int{0}[0],
			},
			want:    Parameters{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := Parameters{}
			err := vd.Decode(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, vd)
		})
	}
}

func TestParametersEncode(t *testing.T) {
	type args struct {
		item Parameters
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
				item: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x04: uint8(4),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0D: uint8(4),
					0x0E: []uint8(hex.Str2Byte("65656565")),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			want: hex.Str2Byte("010001020002030003" +
				"04040534343434" +
				"060006073737373737083838383838" +
				"09090A000A0B000B0C0C" +
				"0D040E656565650F000F1010"),
			wantErr: false,
		},
		{
			name: "case: unknow",
			args: args{
				item: Parameters{
					0xF1: uint16(1),
				},
			},
			wantErr: true,
		},
		{
			name: "case: invalid type",
			args: args{
				item: Parameters{
					0x01: uint32(1),
				},
			},
			wantErr: true,
		},
		{
			name: "case: cropted string",
			args: args{
				item: Parameters{
					0x07: string("77777AA"),
					0x08: string("8888855555"),
				},
			},
			want:    hex.Str2Byte("073737373737083838383838"),
			wantErr: false,
		},
		{
			name: "case: empty string",
			args: args{
				item: Parameters{
					0x07: string(""),
				},
			},
			want:    hex.Str2Byte("070000000000"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := tt.args.item.Encode()
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, pkt)
		})
	}
}

func TestParametersLength(t *testing.T) {
	type args struct {
		item Parameters
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case: ok, 16",
			args: args{
				item: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0E: []uint8(hex.Str2Byte("65656565")),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			want: 16,
		},
		{
			name: "case: ok, 14",
			args: args{
				item: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			want: 12,
		},
		{
			name: "case: ok, 4 correct",
			args: args{
				item: Parameters{
					0x04: uint8(4),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x0D: uint8(4),
					0x0E: []uint8(hex.Str2Byte("65656565")),
				},
			},
			want: 4,
		},
		{
			name: "case: ok, 4 incorrect",
			args: args{
				item: Parameters{
					0x04: uint8(6),
					0x05: []uint8(hex.Str2Byte("34343434")),
					0x0D: uint8(8),
					0x0E: []uint8(hex.Str2Byte("65656565")),
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.args.item.length()
			require.Equal(t, tt.want, l)
		})
	}
}

func TestParametersAdd(t *testing.T) {
	type args struct {
		param Parameters
		id    uint8
		val   any
	}
	tests := []struct {
		name    string
		args    args
		want    Parameters
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				param: Parameters{
					0x01: uint16(1),
				},
				id:  0x02,
				val: uint16(0x02),
			},
			want: Parameters{
				0x01: uint16(1),
				0x02: uint16(2),
			},
			wantErr: false,
		},
		{
			name: "case: invalid id",
			args: args{
				param: Parameters{
					0x01: uint16(1),
				},
				id:  0xF6,
				val: uint16(0x02),
			},
			want: Parameters{
				0x01: uint16(1),
			},
			wantErr: true,
		},
		{
			name: "case: invalid type",
			args: args{
				param: Parameters{
					0x01: uint16(1),
				},
				id:  0x02,
				val: "wring type",
			},
			want: Parameters{
				0x01: uint16(1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.param.Add(tt.args.id, tt.args.val)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, tt.args.param)
		})
	}
}

func TestParametersDelete(t *testing.T) {
	type args struct {
		param Parameters
		id    uint8
	}
	tests := []struct {
		name string
		args args
		want Parameters
	}{
		{
			name: "case: ok",
			args: args{
				param: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				id: 0x02,
			},
			want: Parameters{
				0x01: uint16(1),
				0x03: uint16(3),
				0x06: uint16(6),
				0x07: string("77777"),
				0x08: string("88888"),
				0x09: uint8(9),
				0x0A: uint16(10),
				0x0B: uint16(11),
				0x0C: uint8(12),
				0x0F: uint16(15),
				0x10: uint8(16),
			},
		},

		{
			name: "case: notfound",
			args: args{
				param: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				id: 0x32,
			},
			want: Parameters{
				0x01: uint16(1),
				0x02: uint16(2),
				0x03: uint16(3),
				0x06: uint16(6),
				0x07: string("77777"),
				0x08: string("88888"),
				0x09: uint8(9),
				0x0A: uint16(10),
				0x0B: uint16(11),
				0x0C: uint8(12),
				0x0F: uint16(15),
				0x10: uint8(16),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.param.Delete(tt.args.id)
			require.Equal(t, tt.want, tt.args.param)
		})
	}
}

func TestParametersIsEqual(t *testing.T) {
	type args struct {
		param1 Parameters
		param2 Parameters
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case: ok",
			args: args{
				param1: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				param2: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			wantErr: false,
		},
		{
			name:    "case: ok nil",
			args:    args{},
			wantErr: false,
		},
		{
			name: "case: err",
			args: args{
				param1: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			wantErr: true,
		},

		{
			name: "case: err",
			args: args{
				param1: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				param2: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
				},
			},
			wantErr: true,
		},

		{
			name: "case: err key",
			args: args{
				param1: Parameters{
					0x01: uint16(1),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				param2: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
				},
			},
			wantErr: true,
		},

		{
			name: "case: err value",
			args: args{
				param1: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("77777"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
				param2: Parameters{
					0x01: uint16(1),
					0x02: uint16(2),
					0x03: uint16(3),
					0x06: uint16(6),
					0x07: string("false"),
					0x08: string("88888"),
					0x09: uint8(9),
					0x0A: uint16(10),
					0x0B: uint16(11),
					0x0C: uint8(12),
					0x0F: uint16(15),
					0x10: uint8(16),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.param1.IsEqual(tt.args.param2)
			require.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}
