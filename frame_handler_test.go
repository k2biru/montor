package gbt32960_test

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/k2biru/montor/codec/hex"

	"testing"

	gbt32960 "github.com/k2biru/montor"
	"github.com/stretchr/testify/require"
)

type ErrorReadWriter struct {
	writeLimit int
	Err        error
}

// Read implements the io.Reader interface and always returns an error.
func (e *ErrorReadWriter) Read(_ []byte) (n int, err error) {
	return 0, e.Err
}

func (m *ErrorReadWriter) Write(pkt []byte) (int, error) {
	if m.writeLimit <= 0 {
		return 0, m.Err
	}
	pktLen := len(pkt)
	writeN := pktLen
	if pktLen > m.writeLimit {
		writeN = m.writeLimit
		m.writeLimit = 0
	}
	return writeN, nil

}

func TestFrameHandler_recv(t *testing.T) {
	type args struct {
		writer io.ReadWriter
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name: "1 package",
			args: args{
				writer: bytes.NewBuffer(
					hex.Str2Byte("232301FE313233343536373839000000000000000001002A130116173738000131323334353637383939383736353433323130300304313233343435363739383730FD"),
				),
			},
			want: [][]byte{
				hex.Str2Byte("01FE313233343536373839000000000000000001002A130116173738000131323334353637383939383736353433323130300304313233343435363739383730FD"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := gbt32960.NewFrameHandler(tt.args.writer)
			for _, v := range tt.want {
				fr, err := fh.Recv(context.Background())
				require.Equal(t, tt.wantErr, err != nil, err)
				require.Equal(t, gbt32960.Frame(v), fr)
			}
		})
	}
}

func TestFrameHandler_recvErr(t *testing.T) {
	type args struct {
		writer io.ReadWriter
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case : error",
			args: args{
				writer: &ErrorReadWriter{
					Err: errors.New("Mocke error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := gbt32960.NewFrameHandler(tt.args.writer)
			fr, err := fh.Recv(context.Background())
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, gbt32960.Frame(tt.want), fr)
		})
	}
}

func TestFrameHandler_send(t *testing.T) {
	type args struct {
		frame gbt32960.Frame
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case : ok",
			args: args{
				frame: hex.Str2Byte("010203040506"),
			},
			want:    hex.Str2Byte("2323010203040506"),
			wantErr: false,
		},
		{
			name: "case : false",
			args: args{
				frame: hex.Str2Byte(""),
			},
			want:    hex.Str2Byte(""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := bytes.NewBuffer([]byte{})
			fh := gbt32960.NewFrameHandler(writer)
			err := fh.Send(tt.args.frame)
			require.Equal(t, tt.wantErr, err != nil, err)
			require.Equal(t, tt.want, writer.Bytes())
		})
	}
}

func TestFrameHandler_sendErr(t *testing.T) {
	type args struct {
		frame  gbt32960.Frame
		writer io.ReadWriter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case : error",
			args: args{
				frame: hex.Str2Byte("010203040506"),
				writer: &ErrorReadWriter{
					Err: errors.New("Mocke error"),
				},
			},
			wantErr: true,
		},
		{
			name: "case : unfinised",
			args: args{
				frame: hex.Str2Byte("010203040506"),
				writer: &ErrorReadWriter{
					writeLimit: 2,
					Err:        errors.New("Mocke error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := gbt32960.NewFrameHandler(tt.args.writer)
			err := fh.Send(tt.args.frame)
			require.Equal(t, tt.wantErr, err != nil, err)
		})
	}
}
