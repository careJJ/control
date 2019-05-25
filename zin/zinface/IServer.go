package zinface

type IServer interface{
	Start()
	Stop()
	Serve()
	//添加路由方法  暴露给开发者的
	AddRouter(router IRouter)
}
