package client

import (
	"GoWebAgent/config"
	"GoWebAgent/model"
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	servAddr = flag.String("addr2", config.Get("server_proxy_addr"), "server address")
)


// 将字节转为request
func DecodeRequest(data []byte, reqHost string) (*http.Request, error) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}
	req.Host = reqHost
	scheme := "http"
	req.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI))
	req.RequestURI = ""

	return req, nil
}

func CliStart() {
	flag.Parse()

	ch := make(chan *protocol.Message)

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*servAddr, "")
	xclient := client.NewBidirectionalXClient("RHServer", client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	defer xclient.Close()

	var args model.ReqDTO
	var reply model.ResDTO
	err := xclient.Call(context.Background(), "Mul", &args, &reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	fmt.Printf("客户端已启动，将代理对%v的访问\n",config.Get("client_addr"))
	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
		req,err := DecodeRequest(msg.Payload,config.Get("client_addr"))
		if err != nil {
			log.Fatalf("failed to DecodeRequest: %v", err)
		}
		args.MsgID = msg.Metadata["msgId"]
		args.HttpBody,err = ForwardHandler(req)
		if err != nil {
			log.Fatalf("failed to ForwardHandler: %v", err)
		}
		err = xclient.Call(context.Background(), "CallBack", &args, &reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
	}
}

func ForwardHandler(request *http.Request) ([]byte,error){
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	//设置本地vpn代理
	if config.Get("client_vpn") != "" {
		proxy, err := url.Parse(config.Get("client_vpn"))
		if err != nil {}
		netTransport := &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * time.Duration(5),
		}
		// 请求本地指定的HTTP服务器
		client = &http.Client{
			Timeout: time.Second * 30,
			Transport: netTransport,
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyBytes,nil
}