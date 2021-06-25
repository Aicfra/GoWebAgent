package server

import (
	"GoWebAgent/config"
	"GoWebAgent/model"
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/server"
	rpcx "github.com/smallnest/rpcx/server"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type RHServer struct {
	sync.RWMutex
	rpcConnected bool
	rpcConn   *net.Conn
	rpcServer *rpcx.Server
	rpcCallBackMap *map[string][]byte

	httpSvr *http.Server
	httpConnected bool
	handler *httpHandler
}

//开启Http服务
//开启RpcX服务
func (s *RHServer) Start() error{
	s.handler = &httpHandler{}
	s.httpSvr = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.MustGetInt("server_http_port")),
		Handler: s.handler,
	}
	var err error
	err = s.httpSvr.ListenAndServe()
	return err
}

func (s *RHServer) Mul(ctx context.Context, args *model.ReqDTO, reply *model.ResDTO) error {
	c := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	s.handler.rpcServer = s.rpcServer
	s.handler.rpcConn = &c
	s.rpcConn = &c
	reply.HttpBody = "success"
	s.rpcConnected = true
	fmt.Printf("start to send messages to %s\n", (*s.rpcConn).RemoteAddr().String())
	return nil
}

func (s *RHServer) CallBack(ctx context.Context, args *model.ReqDTO, reply *model.ResDTO) error {
	(*s.rpcCallBackMap)[args.MsgID] = args.HttpBody
	reply.HttpBody = "success"
	return nil
}

func SerStart() {
	//开启http服务
	rhs := RHServer{}
	go rhs.Start()
	//开启服务端rpc
	flag.Parse()
	s := server.NewServer()
	rhs.rpcServer = s
	s.Register(&rhs, "")
	s.Serve("tcp", "localhost:9982")
	*rhs.rpcCallBackMap = make(map[string][]byte)
}

