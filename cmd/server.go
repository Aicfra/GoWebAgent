package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

var clientConn net.Conn
var connected = false

type Arith int

func (t *Arith) Mul(ctx context.Context, a int) error {
	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	a+=1
	connected = true
	return nil
}

func sStart() {
	flag.Parse()
	ln, _ := net.Listen("tcp", ":9981")
	go http.Serve(ln, nil)
	s := server.NewServer()
	s.Register(new(Arith), "")
	go s.Serve("tcp", *addr)

	for !connected {
		time.Sleep(time.Second)
	}

	fmt.Printf("start to send messages to %s\n", clientConn.RemoteAddr().String())
	for {
		if clientConn != nil {
			err := s.SendMessage(clientConn, "test_service_path", "test_service_method", nil, []byte("abcde"))
			if err != nil {
				fmt.Printf("failed to send messsage to %s: %v\n", clientConn.RemoteAddr().String(), err)
				clientConn = nil
			}
		}
		time.Sleep(time.Second)
	}
}