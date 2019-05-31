package main

import (
	"zin/zinface"
	"fmt"
	"mmo_game_server/core"
	"zin/net"
	"mmo_game_server/apis"
)

//当前客户端建立链接之后触发Hook函数
func OnConnectionAdd(conn zinface.IConnection){
	fmt.Println("conn Add")
	//创建一个玩家 将链接和玩家模块绑定
	p:=core.NewPlayer(conn)
	//创建一个玩家 将链接和玩家模块绑定
	p.ReturnPid()
	//给客户端发送一个msgID:200
	p.ReturnPlayerPosition()

	//上线成功了
	//将玩家对象添加到世界管理器中
	core.WorldMgrObj.AddPlayer(p)

	//给conn添加一个属性 pid属性
	conn.SetProperty("pid", p.Pid)

	fmt.Println("player ID = ",p.Pid,"Online...","Play num = ", len(core.WorldMgrObj.Players))



}

func main() {
	s := net.NewServer("MMO Game Server")

	//注册一些 链接创建/销毁的 Hook钩子函数
	s.AddOnConnStart(OnConnectionAdd)

	//针对MsgID2建立路由业务
	s.AddRouter(2,&apis.WorldChat{})

	//注册一些路由业务
	s.Serve()
}
