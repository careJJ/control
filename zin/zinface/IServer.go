package zinface

type IServer interface{
	Start()
	Stop()
	Serve()

	//提供一个得到链接管理模块的方法
	GetConnMgr() IConnManager

	//添加路由方法  暴露给开发者的
	AddRouter(router IRouter)

	AddOnConnStart(hookFunc func(conn IConnection))
	AddOnConnStop(hookFunc func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}

