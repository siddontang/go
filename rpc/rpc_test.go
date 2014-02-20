package rpc

import (
	"errors"
	"sync"
	"testing"
)

var testServerOnce sync.Once
var testClientOnce sync.Once

var testServer *Server
var testClient *Client

func newTestServer() *Server {
	f := func() {
		testServer = NewServer("127.0.0.1:11182")
		go testServer.Start()
	}

	testServerOnce.Do(f)

	return testServer
}

func newTestClient() *Client {
	f := func() {
		testClient = NewClient("127.0.0.1:11182", 10)
	}

	testClientOnce.Do(f)

	return testClient
}

func OnlineRpc(id int) (int, string, error) {
	return id * 10, "abc", errors.New("hello world")
}

func TestRpc(t *testing.T) {
	defer func() {
		e := recover()
		if s, ok := e.(string); ok {
			println(s)
		}

		if err, ok := e.(error); ok {
			println(err.Error())
		}
	}()
	s := newTestServer()

	s.Register("online_rpc", OnlineRpc)

	c := newTestClient()

	var r func(int) (int, string, error)
	if err := c.MakeRpc("online_rpc", &r); err != nil {
		t.Fatal(err)
	}

	r(10)
}
