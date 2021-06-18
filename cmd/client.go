package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

var (
	servAddr = flag.String("addr", "localhost:8972", "server address")
)

func cStart() {
	flag.Parse()

	ch := make(chan *protocol.Message)

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*servAddr, "")
	xclient := client.NewBidirectionalXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	defer xclient.Close()

	args := 1

	reply := 1
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d , %d", args, reply)

	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
	}
}