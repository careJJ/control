package main

import (
	"fmt"
	"time"
	"net"
)

func main()  {
	fmt.Println("client start..")
	time.Sleep(1*time.Second)
	conn,err:=net.Dial("tcp","127.0.0.1:8999")
	if err!=nil{
		fmt.Println("client start err",err)
		return
	}
	for{
		_,err:=conn.Write([]byte("hello Zinx.."))
		if err!=nil{
			fmt.Println("write conn err",err)
			return
		}
		buf:=make([]byte,512)
		cnt,err:=conn.Read(buf)
		if err!=nil{
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("servar call back : %s, cnt = %d\n",buf,cnt)
	time.Sleep(1 *time.Second)
	}


}