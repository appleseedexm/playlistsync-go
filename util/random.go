package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func RandomBytesInBase64(count int) string {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func RandomBytesInHex(count int) string {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}
