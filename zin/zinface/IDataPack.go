package zinface

type DataPack interface {
	//获取二进制包的头部长度
	GetHeadLen() uint32
	//封包方法——将Message打包成|datalen|dataID|data|的结构
	Pack(msg IMessage) ([]byte,error)
	//拆包方法——将|datalen|dataID|data|拆解到Message结构体
	UnPack([]byte) (IMessage,error)
}