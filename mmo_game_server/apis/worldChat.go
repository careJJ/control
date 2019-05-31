package apis

import (
	"zin/net"
	"zin/zinface"
	"github.com/golang/protobuf/proto"
	"mmo_game_server/core"
	"fmt"

	"mmo_game_server/pb"
)


type WorldChat struct{
	net.Baserouter
}
func (wc *WorldChat) Handle(request zinface.IRequest) {
	//1 解析客户端传递进来的protobuf数据
	proto_msg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetMsg().GetMsgData(), proto_msg);err != nil {
		fmt.Println("Talk message unmarshal error ", err)
		return
	}

	//通过获取链接属性，得到当前的玩家ID
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get Pid error ", err)
		return
	}

	//通过pid 来得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 当前的聊天数据广播给全部的在线玩家
	//当前玩家的windows客户端发送过来的消息
	player.SendTalkMsgToAll(proto_msg.GetContent())
}
