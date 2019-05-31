package config

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

//全局配置文件的类
type Globalobj struct {
	Host string //当前监听的IP
	Port int    //当前的监听Port
	Name string //当前zinxserver的名称
	Version          string //当前框架的版本号
	MaxPackageSize   uint32 //每次Read一次的最大长度
	WorkerPoolSize   uint32 //当前服务器要开启多少了worker Goroutine
	MaxWorkerTaskLen uint32 //每个worker的对应的消息对象 允许缓存的最大任务Request数量
	MaxConn          uint32 //当前server的最大链接数量
}

//
var Globalobject *Globalobj


func (g *Globalobj)LoadConfig(){
	data,err:=ioutil.ReadFile("conf/zinx.json")
	if err!=nil{
		fmt.Println("load config error")
		return
	}
	err=json.Unmarshal(data,&Globalobject)
	if err!=nil{
		panic(err)
	}

}

func init(){
	Globalobject=&Globalobj{
		Name:"ZinxServerApp",
		Host:"0.0.0.0",
		Port:8999,
		Version:"V0.4",
		MaxPackageSize:512,

	}
	Globalobject.LoadConfig()
}