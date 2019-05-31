package zinface

import "net"

type IConnection interface {
	Start()

	Stop()
	GetConnID() uint32
	GetTCPConnection() *net.TCPConn
	GetRemoteAddr() net.Addr
	Send(msgId uint32,msgData []byte)error
	//设置属性
	SetProperty(key string, value interface{})
	//获取属性
	GetProperty(key string)(interface{}, error)
	//删除属性
	RemoveProperty(key string)
}
//业务处理方法 抽象定义
type HandleFunc func(request IRequest) error


