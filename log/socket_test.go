package log

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	fileName := "./test_server.log"

	os.Remove(fileName)

	s, err := NewServer(fileName, "tcp", "127.0.0.1:11183")
	if err != nil {
		t.Fatal(err)
	}
	go s.Run()
	defer s.Close()

	var h *SocketHandler
	h, err = NewSocketHandler("tcp", "127.0.0.1:11183")

	_, err = h.Write([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	s.Close()

	var f *os.File
	f, err = os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	buf := make([]byte, 64)
	var n int
	n, err = f.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	buf = buf[0:n]

	if string(buf) != "hello world\n" {
		t.Fatal(string(buf))
	}

	os.Remove(fileName)
}
