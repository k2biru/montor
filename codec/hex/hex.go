package hex

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/k2biru/montor/codec/gbk"
)

const (
	doubeWordLen = 4
)

// BYTE
func ReadByte(pkt []byte, idx *int) uint8 {
	ans := pkt[*idx]
	*idx++
	return ans
}

func WriteByte(pkt []byte, num uint8) []byte {
	return append(pkt, num)
}

// WORD
func ReadWord(pkt []byte, idx *int) uint16 {
	ans := binary.BigEndian.Uint16(pkt[*idx : *idx+2])
	*idx += 2
	return ans
}

func WriteWord(pkt []byte, num uint16) []byte {
	numPkt := make([]byte, 2)
	binary.BigEndian.PutUint16(numPkt, num)
	return append(pkt, numPkt...)
}

// DWORD
func ReadDoubleWord(pkt []byte, idx *int) uint32 {
	ans := binary.BigEndian.Uint32(pkt[*idx : *idx+4])
	*idx += doubeWordLen
	return ans
}

func WriteDoubleWord(pkt []byte, num uint32) []byte {
	numPkt := make([]byte, doubeWordLen)
	binary.BigEndian.PutUint32(numPkt, num)
	return append(pkt, numPkt...)
}

// BYTES
func ReadBytes(pkt []byte, idx *int, n int) []byte {
	ans := pkt[*idx : *idx+n]
	*idx += n
	return ans
}

func WriteBytes(pkt, arr []byte) []byte {
	return append(pkt, arr...)
}

// STRING using GBK encoding
func ReadString(pkt []byte, idx *int, n int) string {
	if l := len(pkt); l == 0 || l < n {
		return ""
	}
	var data []byte
	for i := 0; i < n; i++ {
		v := ReadByte(pkt, idx)
		data = append(data, v)
	}
	data, err := gbk.GBK2UTF8(data)
	if err != nil {
		// should never get here
		return ""
	}
	stopIndex := strings.IndexRune(string(data), '\x00')
	if stopIndex != -1 {
		data = data[:stopIndex]
	}

	return string(data)
}

// STRING
func WriteString(pkt []byte, str string, n int) []byte {
	arr, err := gbk.UTF82GBK([]byte(str))
	if err != nil {
		// should never get here
		return pkt
	}
	strLen := len(arr)
	if dif := n - strLen; dif > 0 {
		// apend
		for i := 0; i < dif; i++ {
			arr = append(arr, 0x00)
		}
	} else if dif < 0 {
		// crop arr
		arr = arr[:n]
	}
	return WriteBytes(pkt, arr)
}

// tools hex str
func Str2Byte(src string) []byte {
	dst, err := hex.DecodeString(src)
	if err != nil {
		if errors.Is(err, hex.ErrLength) {
			log.Printf("Source str invalid \"%s\", will ignore extra byte", src)
		} else {
			log.Print("Fail to transform hex str to byte array")
		}
	}
	return dst
}

func Byte2Str(src []byte) string {
	return hex.EncodeToString(src)
}
