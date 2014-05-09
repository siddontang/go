package hack

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestString(t *testing.T) {
	b := []byte("hello world")
	a := String(b)

	if a != "hello world" {
		t.Fatal(a)
	}

	b[0] = 'a'

	if a != "aello world" {
		t.Fatal(a)
	}

	b = append(b, "abc"...)
	if a != "aello world" {
		t.Fatal(a)
	}
}

func TestByte(t *testing.T) {
	a := "hello world"

	b := Slice(a)

	if !bytes.Equal(b, []byte("hello world")) {
		t.Fatal(string(b))
	}
}

func TestInt(t *testing.T) {
	if int64(binary.LittleEndian.Uint64(IntSlice(1))) != 1 {
		t.Fatal("error")
	}

	if int64(binary.LittleEndian.Uint64(IntSlice(-1))) != -1 {
		t.Fatal("error")
	}

	if int64(binary.LittleEndian.Uint64(IntSlice(32768))) != 32768 {
		t.Fatal(1)
	}

	if int64(binary.LittleEndian.Uint64(IntSlice(-32768))) != -32768 {
		t.Fatal(1)
	}
}
