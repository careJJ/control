package main

import "zin/net"

func main(){
	s:=net.NewServer("zinx v0.2")

	s.Serve()
	return
}
