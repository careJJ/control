package net

import (
	"zin/zinface"
	"fmt"
	"net"
	"zin/config"
)

type Server struct {
	IPVersion string
	IP        string
	Port      int
	Name      string
	//路由属性
	//Router zinface.IRouter
	//多路由的消息管理模块
	MsgHandler zinface.IMsgHandler
	//链接管理模块
	connMgr zinface.IConnManager

	//该server创建链接之后自动调用Hook函数
	OnConnStart func(conn zinface.IConnection)
	//该server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn zinface.IConnection)
}

func NewServer(name string) zinface.IServer {
	s := &Server{
		Name:       config.Globalobj.Name,
		IPVersion:  "tcp4",
		IP:         config.Globalobject.Host,
		Port:       config.Globalobject.Port,
		MsgHandler: NewMsgHandler(),
		connMgr:    NewConnManager(),
	}
	return s
}

//启动服务器
//原生socket 服务器编程
func (s *Server) Start() {

	fmt.Printf("[start]Server Linstenner at IP :%s,Port :%d,is starting..\n", s.IP, s.Port)
	//0 启动worker工作池
	s.MsgHandler.StartWorkerPool()
	//1.创建套接字：得到一个TCP的addr
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
	cid = 0

	//3.阻塞等待客户端请求
	go func() {
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//创建一个Connection对象
			//判断当前server链接数量是否已经最大值
			if s.connMgr.Len() >= int(config.Globalobject.MaxConn) {
				//当前链接已经满了
				fmt.Println("---> Too many Connection MAxConn = ", config.Globalobject.MaxConn)
				conn.Close()
				continue
			}

			//此时conn已经和对端客户端链接
			dealConn := NewConnection(s,conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	//服务器停止  应该清空当前全部的链接
	s.connMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start()
	select {} //main函数不退出
}

func (s *Server) AddRouter(msgId uint32, router zinface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router SUCC!!!msgID = ", msgId)
}

func (s *Server) GetConnMgr() zinface.IConnManager {
	return s.connMgr
}

//注册 创建链接之后 调用的 Hook函数 的方法
func (s *Server) AddOnConnStart(hookFunc func(conn zinface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册 销毁链接之前调用的Hook函数 的方法
func (s *Server) AddOnConnStop(hookFunc func(conn zinface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用 创建链接之后的HOOK函数的方法
func (s *Server) CallOnConnStart(conn zinface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用 销毁链接之前调用的HOOk函数的方法
func (s *Server) CallOnConnStop(conn zinface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
