package websocket

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func calcAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func calcKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}
