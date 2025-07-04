package hex

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStr2Byte(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: short str",
			args: args{str: "123456789012"},
			want: []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12},
		},
		{
			name: "case2: long str",
			args: args{str: "000000000002080301CD779E0728C032003C0000008F230125145158"},
			want: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x08, 0x03, 0x01, 0xCD, 0x77, 0x9E, 0x07, 0x28, 0xC0, 0x32, 0x00, 0x3C, 0x00, 0x00, 0x00, 0x8F,
				0x23, 0x01, 0x25, 0x14, 0x51, 0x58},
		},
		{
			name: "case3: str of odd length",
			args: args{str: "12345678901"},
			want: []byte{0x12, 0x34, 0x56, 0x78, 0x90},
		}, {
			name: "case4: invalid str",
			args: args{str: "12345Z78"},
			want: []byte{0x12, 0x34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Str2Byte(tt.args.str)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestByte2Str(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1: byte to str",
			args: args{src: []byte{0x12, 0x34, 0x56, 0x78}},
			want: "12345678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Byte2Str(tt.args.src)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReadByte(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "case1: read 0x9a",
			args: args{
				pkt: []byte{0x9a},
				idx: &[]int{0}[0],
			},
			want: 154,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadByte(tt.args.pkt, tt.args.idx)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.args.idx, &[]int{1}[0])
		})
	}
}

func TestWriteByte(t *testing.T) {
	type args struct {
		pkt []byte
		num uint8
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: write 0x9a",
			args: args{
				pkt: []byte{},
				num: 154,
			},
			want: []byte{0x9a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteByte(tt.args.pkt, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadWord(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "case1: read 0x9a 0xff",
			args: args{
				pkt: []byte{0x9a, 0xff},
				idx: &[]int{0}[0],
			},
			want: 0x9aff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadWord(tt.args.pkt, tt.args.idx); got != tt.want {
				t.Errorf("ReadWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteWord(t *testing.T) {
	type args struct {
		pkt []byte
		num uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: write 0x9a 0xff",
			args: args{
				pkt: []byte{},
				num: 0x9aff,
			},
			want: []byte{0x9a, 0xff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteWord(tt.args.pkt, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadDoubleWord(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "case1: read 0x9a 0xff 0xcc 0xdd",
			args: args{
				pkt: []byte{0x9a, 0xff, 0xcc, 0xdd},
				idx: &[]int{0}[0],
			},
			want: 0x9affccdd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadDoubleWord(tt.args.pkt, tt.args.idx); got != tt.want {
				t.Errorf("ReadDoubleWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteDoubleWord(t *testing.T) {
	type args struct {
		pkt []byte
		num uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: write 0x9a 0xff 0xcc 0xdd",
			args: args{
				pkt: []byte{},
				num: 0x9affccdd,
			},
			want: []byte{0x9a, 0xff, 0xcc, 0xdd},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteDoubleWord(tt.args.pkt, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteDoubleWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBytes(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: read 0x9a 0xff 0xcc 0xdd",
			args: args{
				pkt: []byte{0x9a, 0xff, 0xcc, 0xdd},
				idx: &[]int{0}[0],
			},
			want: []byte{0x9a, 0xff, 0xcc, 0xdd},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadBytes(tt.args.pkt, tt.args.idx, len(tt.args.pkt))
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.args.idx, &[]int{4}[0])
		})
	}
}

func TestWriteBytes(t *testing.T) {
	type args struct {
		pkt []byte
		num []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: write 0x9a 0xff 0xcc 0xdd",
			args: args{
				pkt: []byte{},
				num: []byte{0x9a, 0xff, 0xcc, 0xdd},
			},
			want: []byte{0x9a, 0xff, 0xcc, 0xdd},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteBytes(tt.args.pkt, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadString(t *testing.T) {
	type args struct {
		pkt []byte
		idx *int
		n   int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantLen int
	}{
		{
			name: "case: read sugeng dalu",
			args: args{
				pkt: Str2Byte("737567656E672064616C75"),
				idx: &[]int{0}[0],
				n:   11,
			},
			want:    "sugeng dalu",
			wantLen: 11,
		},
		{
			name: "case: null data",
			args: args{
				pkt: Str2Byte("00"),
				idx: &[]int{0}[0],
				n:   11,
			},
			want:    "",
			wantLen: 0,
		},
		{
			name: "case: no data",
			args: args{
				pkt: Str2Byte(""),
				idx: &[]int{0}[0],
				n:   11,
			},
			want:    "",
			wantLen: 0,
		},
		{
			name: "case: read sugeng then null",
			args: args{
				pkt: Str2Byte("737567656E6700"),
				idx: &[]int{0}[0],
				n:   7,
			},
			want:    "sugeng",
			wantLen: 7,
		},
		{
			name: "case: read sugeng then null 2 ",
			args: args{
				pkt: Str2Byte("737567656E670003000200"),
				idx: &[]int{0}[0],
				n:   10,
			},
			want:    "sugeng",
			wantLen: 10,
		},
		{
			name: "case: read too long",
			args: args{
				pkt: Str2Byte("737567656E6700"),
				idx: &[]int{0}[0],
				n:   11,
			},
			want:    "",
			wantLen: 0,
		},
		{
			name: "case: read sugeng 京",
			args: args{
				pkt: Str2Byte("737567656E67BEA900"),
				idx: &[]int{0}[0],
				n:   9,
			},
			want:    "sugeng京",
			wantLen: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadString(tt.args.pkt, tt.args.idx, tt.args.n); got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
			require.Equal(t, tt.wantLen, *tt.args.idx)

		})
	}
}

func TestWriteString(t *testing.T) {
	type args struct {
		pkt []byte
		str string
		n   int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case1: write sugeng dalu",
			args: args{
				str: "sugeng dalu",
				n:   11,
			},
			want: Str2Byte("737567656E672064616C75"),
		},
		{
			name: "case: empty",
			args: args{
				str: "",
				n:   11,
			},
			want: Str2Byte("0000000000000000000000"),
		},
		{
			name: "case1: write sugeng and space",
			args: args{
				str: "sugeng",
				n:   11,
			},
			want: Str2Byte("737567656E670000000000"),
		},
		{
			name: "case1: write number as string",
			args: args{
				str: "77777AA",
				n:   5,
			},
			want: Str2Byte("3737373737"),
		},
		{
			name: "case: write s京",
			args: args{
				str: "s京",
				n:   3,
			},
			want: Str2Byte("73BEA9"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt := WriteString(tt.args.pkt, tt.args.str, tt.args.n)
			require.Equal(t, tt.want, pkt)
			// if got := WriteString(tt.args.pkt, tt.args.str, tt.args.n); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("WriteString() = %v, want %v", got, tt.want)
			// }
		})
	}
}
