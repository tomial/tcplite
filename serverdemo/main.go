package main

import (
	"fmt"

	"github.com/tomial/tcplite/iface"
	"github.com/tomial/tcplite/net"
)

type MyRouter struct {
	net.BaseRouter
}

func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("Called DoConnectionBegin.")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("Connection is lost...")
	fmt.Printf("Connection [%d] is lost.\n", conn.GetConnID())
}

func (r *MyRouter) Handle(request iface.IRequest) {
	fmt.Printf("Received data with msgId: %d, data: %s\n", request.GetMsgID(), request.GetData())
	fmt.Printf("Send reply to client\n")
	request.GetConnection().SendMsg(0, request.GetData())
}

type TestHandlerRouter struct {
	net.BaseRouter
}

func (t *TestHandlerRouter) Handle(request iface.IRequest) {
	fmt.Printf("Received data with msgId: %d, data: %s\n", request.GetMsgID(), request.GetData())
	fmt.Printf("Send reply to client\n")
	request.GetConnection().SendMsg(1, request.GetData())
}

func main() {

	s := net.NewServer("[tcp lite V0.9]")
	s.AddRouter(0, &MyRouter{})
	s.AddRouter(1, &TestHandlerRouter{})

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.Serve()
}
