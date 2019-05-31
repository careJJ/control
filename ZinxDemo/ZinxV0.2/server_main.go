package main

import (
	"zin/net"
	"zin/zinface"
	"fmt"
)

type PingRouter struct {
	net.Baserouter
}

func (this *PingRouter)PreHandle(reqeust zinface.IRequest)  {
	fmt.Println("Call Router PreHandler...")
	_,err:=reqeust.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if err!=nil {
		fmt.Println("call ping error")
	}
	}
func (this *PingRouter)Handle(reqeust zinface.IRequest)  {
	fmt.Println("Call Router Handler...")
	_,err:=reqeust.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if err!=nil {
		fmt.Println("call ping ping error")
	}
}
func (this *PingRouter)PostHandle(reqeust zinface.IRequest)  {
	fmt.Println("Call Router PostHandler...")
	_,err:=reqeust.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if err!=nil {
		fmt.Println("call back after ping error")
	}
}

func main(){
	s:=net.NewServer("zinx v0.2")
	s.AddRouter(&PingRouter{})//真正处理核心业务
	s.Serve()
	return
}
