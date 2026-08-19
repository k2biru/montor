package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoCParser(t *testing.T) {
	type args struct {
		val uint8
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// {
		// 	name: "case: 1",
		// 	args: args{
		// 		val: 1,
		// 	},
		// 	want: 1,
		// },
		// {
		// 	name: "case: 50",
		// 	args: args{
		// 		val: 50,
		// 	},
		// 	want: 50,
		// },
		// {
		// 	name: "case: 100",
		// 	args: args{
		// 		val: 100,
		// 	},
		// 	want: 100,
		// },
		// {
		// 	name: "case: 101 convert to max",
		// 	args: args{
		// 		val: 101,
		// 	},
		// 	want: 100,
		// },
		// {
		// 	name: "case: 200 convert to max",
		// 	args: args{
		// 		val: 200,
		// 	},
		// 	want: 100,
		// },
		// {
		// 	name: "case: 0xFE convert to 0",
		// 	args: args{
		// 		val: 0xFE,
		// 	},
		// 	want: 0,
		// },
		{
			name: "case: 0xFF convert to 0",
			args: args{
				val: 0xFF,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SoC().SetVal(tt.args.val).Calculate().AsInt()

			if call := SoC().SetVal(tt.args.val).Calculate(); call.Error() != nil {
				res = call.AsInt()
			}

			require.Equal(t, tt.want, res)
		})
	}
}
