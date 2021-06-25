package test

import (
	"GoWebAgent/cmd/client"
	"GoWebAgent/cmd/server"
	"fmt"
	"strconv"
	"testing"
)

func TestServer(t *testing.T){
	server.SerStart()
}

func TestClient(t *testing.T){
	client.CliStart()
}

type People interface {
	GetName() string
	GetAge() int
}

type Chinese struct {
	Name string
	Age int
}

type XXX struct {
	chin *Chinese
}

type YYY struct {
	xxx *XXX
}

func (c Chinese) GetName() string{
	return c.Name
}

func (c Chinese) GetAge() int{
	return c.Age
}

func TestT3(t *testing.T){
	ch := Chinese{Name: "xm",Age: 1}
	var peo People
	peo = ch
	fmt.Printf(peo.GetName())
	fmt.Printf(strconv.Itoa(peo.GetAge()))
}

func TestT4(t *testing.T){
	xx := &XXX{chin: &Chinese{Name: "aaa"}}
	yy := &YYY{xxx: xx}
	ch2 := &Chinese{Name: "1111111111111",Age: 1}
	xx.chin = ch2
	fmt.Println(yy.xxx.chin.Name)
	fmt.Println(xx.chin.Name)
}