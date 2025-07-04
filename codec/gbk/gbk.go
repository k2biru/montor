package gbk

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBK UTF-8
func GBK2UTF8(src []byte) ([]byte, error) {
	dst, err := io.ReadAll(transform.NewReader(bytes.NewBuffer(src), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// UTF-8 GBK
func UTF82GBK(src []byte) ([]byte, error) {
	dst, err := io.ReadAll(transform.NewReader(bytes.NewBuffer(src), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return nil, err
	}
	return dst, nil
}
