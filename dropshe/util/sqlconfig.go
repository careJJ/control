package util

import (

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"encoding/json"
	//使用了beego框架的配置文件读取模块
	"github.com/astaxie/beego/config"
)

//初始化数据库

var (
	Server_name  string //项目名称
	Server_addr  string //服务器ip地址
	Server_port  string //服务器端口
	Redis_addr   string //redis ip地址
	Redis_port   string //redis port端口
	Redis_dbnum  string //redis db 编号
	Mysql_addr   string //mysql ip 地址
	Mysql_port   string //mysql 端口
	Mysql_dbname string //mysql db name
	Fastdfs_port string //fastdfs 端口
	Fastdfs_addr string //fastdfs ip
)

func InitConfig() {
	//从配置文件读取配置信息
	appconf, err := config.NewConfig("ini", "/home/itcast/workspace/go/src/dropshe/conf/app.conf")
	if err != nil {
		beego.Debug(err)
		return
	}
	Server_name = appconf.String("appname")
	Server_addr = appconf.String("httpaddr")
	Server_port = appconf.String("httpport")
	Redis_addr = appconf.String("redisaddr")
	Redis_port = appconf.String("redisport")
	Redis_dbnum = appconf.String("redisdbnum")
	Mysql_addr = appconf.String("mysqladdr")
	Mysql_port = appconf.String("mysqlport")
	Mysql_dbname = appconf.String("mysqldbname")
	Fastdfs_port = appconf.String("fastdfsport")
	Fastdfs_addr = appconf.String("fastdfsaddr")
	return
}

//连接redis
func RedisOpen(server_name,redis_addr,redis_port,redis_dbnum string)(bm cache.Cache,err error){
	redis_config_map := map[string]string{
		"key":Server_name,
		"conn":redis_addr+":"+redis_port,
		"dbNum":redis_dbnum,
	}
	//map->json
	redis_config,_:=json.Marshal(redis_config_map)


	//连接redis
	bm,err =cache.NewCache("redis",string(redis_config))
	if err!=nil{
		beego.Error("练接redis错误")

		return nil,err
	}
	return bm,nil
}

func init() {
	InitConfig()
	InitLogs()
}
