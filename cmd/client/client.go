package client

import (
	"GoWebAgent/model"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

var (
	servAddr = flag.String("addr2", "localhost:9982", "server address")
)

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

	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
		args.MsgID = msg.Metadata["msgId"]
		args.HttpBody = []byte(args.MsgID+"---"+"12341234")
		err := xclient.Call(context.Background(), "CallBack", &args, &reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
	}
}