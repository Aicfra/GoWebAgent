package server

import (
	"fmt"
	"github.com/satori/go.uuid"
	rpcx "github.com/smallnest/rpcx/server"
	"net"
	"net/http"
	"net/http/httputil"
)

type httpHandler struct{
	rpcConn *net.Conn
	rpcServer *rpcx.Server
	rpcCallBackMap *map[string][]byte
}

const (
	CLIENT    = "client"
	SERVER    = "server"
	errorHTML = `
<!DOCTYPE html>
<html>
<head>
<title>请求错误</title>
<style>
</style>
</head>
<body>
<h2>无客户端连接</h2>
</body>
</html>`
)

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *h.rpcConn != nil {
		reqBytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Printf("failed to send Request messsage to %s: %v\n", (*h.rpcConn).RemoteAddr().String(), err)
		}
		metadata := make(map[string]string)
		metadata["msgId"] = uuid.NewV4().String()
		err = h.rpcServer.SendMessage(*h.rpcConn, "test_service_path", "test_service_method", metadata, reqBytes)
		if err != nil {
			fmt.Printf("failed to send messsage to %s: %v\n", (*h.rpcConn).RemoteAddr().String(), err)
			h.rpcConn = nil
		}else{
			fmt.Printf("success to send messsage to %s: %v\n", string(reqBytes))
			go func() {
				getMsged := false
				for !getMsged {
					if value, ok := (*h.rpcCallBackMap)[metadata["msgId"]]; ok {
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						w.Header().Set("X-Content-Type-Options", "nosniff")
						w.WriteHeader(http.StatusOK)
						w.Write(value)
						delete(*h.rpcCallBackMap,metadata["msgId"])
						getMsged = true
					}
				}
			}()
		}
	}else{
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorHTML))
	}
}

