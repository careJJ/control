package net

import (
	"zin/zinface"
	"fmt"
	"net"
)

type Server struct {
	IPVersion string
	IP        string
	Port      int
	Name      string
	//路由属性
	Router zinface.IRouter
}

func NewServer(name string) zinface.IServer {
	s := &Server{
		Name:      config.Globa10bject.Name,
		IPVersion: "tcp4",
		IP:        Globa10bject.Host,
		Port:      Globa10bject.Port,
		Router:nil,
	}
	return s
}

//启动服务器
//原生socket 服务器编程
func (s *Server) Start() {
	//1.创建套接字：得到一个TCP的addr
	fmt.Printf("[start]Server Linstenner at IP :%s,Port :%d,is starting..\n", s.IP, s.Port)
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	//2.监听服务器地址
	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IPVersion, "err", err)
		return
	}
	//生成id的累加器
	var cid uint32
	cid=0

	//3.阻塞等待客户端请求
	go func() {
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//此时conn已经和对端客户端链接
			dealConn:=NewConnection(conn,cid,s.Router)
			cid++
			go dealConn.Start()

				}
			}()
		}
	}()
}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.Start()
	select {

	}//main函数不退出
}
