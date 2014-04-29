package log

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path"
)

//a log server for handling SocketHandler send log

type Server struct {
	closed   bool
	listener net.Listener
	fd       *os.File
}

func NewServer(fileName string, protocol string, addr string) (*Server, error) {
	s := new(Server)

	s.closed = false

	var err error

	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	s.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	s.listener, err = net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true

	s.fd.Close()

	s.listener.Close()
	return nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.onRead(conn)
	}
}

func (s *Server) onRead(c net.Conn) {
	br := bufio.NewReaderSize(c, 1024)

	var bufLen uint32

	for {
		if err := binary.Read(br, binary.BigEndian, &bufLen); err != nil {
			c.Close()
			return
		}

		buf := make([]byte, bufLen, bufLen+1)

		if _, err := io.ReadFull(br, buf); err != nil && err != io.ErrUnexpectedEOF {
			c.Close()
			return
		} else {
			if len(buf) == 0 {
				continue
			}
			if buf[len(buf)-1] != '\n' {
				buf = append(buf, '\n')
			}

			s.fd.Write(buf)
		}

	}
}
