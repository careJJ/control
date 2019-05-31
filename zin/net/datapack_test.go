package net

import (
	"testing"
	"fmt"
	"net"
	"io"

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
						fmt.Println("read head error",err)
						break
					}
					msgHead,err:=dp.UnPack(headData)
					if err!=nil {
						fmt.Println("server unpack err",err)
						return
					}
					if msgHead.GetMsgLen()>0{
						//数据区有内容,需要进行第二次读取
						//将msgHead进行向下安装，将IMessage 换换成Message
						msg:=msgHead.(*Message)
						//给msg的Data属性开辟，长度就是数据的长度
						msg.Data=make([]byte,msg.GetMsgLen())
						//根据dataken的长度进行第二次read
						_,err:=io.ReadFull(*conn,msg.Data)
						if err!=nil{
							fmt.Println("server unpack data error",err)
							return
						}
						fmt.Println("Recv MsgID = ",msg.Id," datalen = ", msg.Datalen, "data = ", string(msg.Data))
					}
				}
			}(&conn)
		}
	}()
 conn,err:=net.Dial("tcp","127.0.0.1:7777")
	if err!=nil{
		fmt.Println("client dail err:",err)
		return
	}
	//封包，创建dp拆包，封包的工具
	dp:=NewDataPack()
	//模拟粘包过程法宝
	//封装第一个包
	msg1:=&Message{
		Id:1,
		Datalen:4,
		Data:[]byte{'1','2','3','4'},
	}
	sendData1,err:=dp.Pack(msg1)
	if err!=nil{
		fmt.Println("client send data1 error",err)
		return
	}
	msg2:=&Message{
		Id:2,
		Datalen:5,
		Data:[]byte{'5','6','7','8','9'},
	}
	sendData2,err:=dp.Pack(msg2)
	if err!=nil{
		fmt.Println("client send data2 error",err)
		return
	}
	sendData1=append(sendData1,sendData2...)
	conn.Write(sendData1)



}