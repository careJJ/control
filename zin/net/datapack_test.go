package net

import (
	"testing"
	"fmt"
	"net"
	"io"
	"zin/zinface"
)

// 函数名 Test开头  后面的函数名 自定义
//形参 (t *testing.T)
func TestDataPack(t *testing.T) {
	fmt.Println("test datapack...")

	//模拟写一个server，收到二进制流，进行解包


	listerer,err:=net.Listen("tcp","127.0.0.1:7777")
	if err!=nil{
		fmt.Println("server listener err",err)
		return
	}
	go func() {
		for{
			conn,err:=listerer.Accept()
			if err!=nil {
				fmt.Println("Accept err",err)
				return
			}
			//读写业务
			go func(conn *net.Conn) {
				//读取客户端的请求
				dp:=NewDataPack()
				for{
					headData:=make([]byte,dp.GetHeadLen())
					_,err:=io.ReadFull(*conn,headData)
					if err!=nil {
						fmt.Println("server unpack err",err)
						return
					}
					if msgHead.GetMsgLen()>0{
						//数据区有内容,需要进行第二次读取
						//将msgHead进行向下安装，将IMessage 换换成Message
						msg:=msgHead.(*Message)
					}
				}




			}(&conn)
		}



	}()

}