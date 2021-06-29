package main

import (
	"GoWebAgent/cmd"
	"GoWebAgent/cmd/client"
	"GoWebAgent/cmd/server"
	"GoWebAgent/config"
	"flag"
	"github.com/pkg/errors"
)

var (
	mode  = flag.String("mode", config.Get("mode"), "启动模式，可选为client、server")
)

func main() {
	flag.Parse()
	switch *mode {
	case cmd.SERVER:
		server.SerStart()
	case cmd.CLIENT:
		client.CliStart()
	default:
		panic(errors.New("mode error"))
	}
}