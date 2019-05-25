package zinface

import "net"

type IConnection interface {
	Start()

	Stop()
	GetConnID() uint32
	GetTCPConnection() *net.TCPConn
	GetRemoteAddr() net.Addr
	Send(data []byte,cnt int)error
}
type HandleFunc func(*net.TCPConn,[]byte,int) error


