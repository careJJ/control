package zinface

type IConnManager interface {
	Add(conn IConnection)
	Remove(connID uint32)
//根据链接ID得到链接
	Get(connID uint32) (IConnection,error)
	//得到目前服务器的链接总个数
	Len() int
	//清空全部链接的方法
	ClearConn()
}
