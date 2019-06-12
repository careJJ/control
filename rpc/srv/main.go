package main

import (
	"net/rpc"
	"net"
	"fmt"
	"net/http"
)

type Panda struct{

}

func (this *Panda) GetAdd(In int,Out *int)error{
	*Out = In +10086
	return nil
}


func main(){
	//	1结构体实例化
	pd:=new(Panda)
	rpc.Register(pd)
	rpc.HandleHTTP()

	//	4监听网络
	ln,err:=net.Listen("tcp","127.0.0.1:10086")
	if err!=nil{
		fmt.Println("网络错误",err)
	}

	http.Serve(ln,nil)


}