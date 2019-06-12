package main

import (
	"net/rpc"
	"fmt"
)

func main(){
	cli,err:=rpc.DialHTTP("tcp","127.0.0.1:10086")
	if err!=nil{
		fmt.Println("链接服务器失败",err)
		return
	}
	var size int
	err = cli.Call("Panda.GetAdd",12306,&size)
	if err!=nil{
		fmt.Println("失败",err)
	return
		}

		fmt.Println("计算结果为",size)


}
