package zinface

//抽象IRequest 一次行请求的数据封装

type IRequest interface {
	//得到当前的请求的链接
	GetConnection() IConnection
	//得到请求的消息
	GetMsg() IMessage
}