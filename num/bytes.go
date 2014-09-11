package num

import (
	"encoding/binary"
)

//all are bigendian format

func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func Uint16ToBytes(u uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, u)
	return buf
}

func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func Uint32ToBytes(u uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}

func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}
