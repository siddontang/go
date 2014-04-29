package main

import (
	"flag"
	"github.com/siddontang/golib/log"
)

var logFile = flag.String("logfile", "./logd.log", "file to log")
var net = flag.String("net", "tcp", "server listen protocol, like tcp, udp or unix")
var addr = flag.String("addr", "127.0.0.1:11183", "server listen address")

func main() {
	flag.Parse()

	s, err := log.NewServer(*logFile, *net, *addr)
	if err != nil {
		panic(err)
	}

	s.Run()
}
