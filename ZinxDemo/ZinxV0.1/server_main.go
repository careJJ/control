package main

import "zin/net"

func main(){
	s:=net.NewServer("zinx v0.1")
	s.Serve()
	return
}
