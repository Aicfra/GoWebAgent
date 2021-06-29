package main

import (
	"GoWebAgent/cmd/client"
	"GoWebAgent/cmd/server"
	"testing"
)

func TestServer(t *testing.T){
	server.SerStart()
}

func TestClient(t *testing.T){
	client.CliStart()
}