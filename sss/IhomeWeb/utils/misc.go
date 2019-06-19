package utils

import (
	//"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/astaxie/beego/cache"
	"fmt"
	"encoding/json"
	_"github.com/astaxie/beego/cache/redis"
	_"github.com/garyburd/redigo/redis"
	_"github.com/gomodule/redigo/redis"
	"crypto/md5"
	"encoding/hex"
)

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求
func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url

	return domain_url
}

func RedisOpen(server_name,redis_addr,redis_port,redis_dbnum string)(bm cache.Cache,err error){
	redis_config_map := map[string]string{
		"key":server_name,
		"conn":redis_addr+":"+redis_port,
		"dbNum":redis_dbnum,
	}
	//map->json
	redis_config,_:=json.Marshal(redis_config_map)


	//链接redis
	bm,err =cache.NewCache("redis",string(redis_config))
	if err!=nil{
		fmt.Println("连接redis错误",err)

		return nil,err
	}
	return bm,nil
}

func Getmd5string(s string)string{
	m :=md5.New()
	return  hex.EncodeToString(m.Sum([]byte(s)))
}