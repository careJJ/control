package handler

import (
	"context"
"fmt"
	//"github.com/micro/go-log"
	"sss/IhomeWeb/utils"
	example "sss/PostRet/proto/example"
	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"sss/IhomeWeb/model"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostRet(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 注册服务  PostRet  /api/v1.0/users")


	/*1 初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*2 连接redis*/
	bm ,err :=utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,utils.G_redis_dbnum)
	if err!=nil{
		fmt.Println("redis连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//3获取短信验证码
	value:=bm.Get(req.Mobile)
	value_string,_:=redis.String(value,nil)
	//4.验证短信验证码是否正确
	if value_string!=req.SmsCode{
		fmt.Println("短信验证码错误",value_string,req.SmsCode)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//加密
	user:=models.User{}
	user.Password_hash = utils.Getmd5string(req.Password)


	//插入数据
	o:= orm.NewOrm()
	id  ,err :=o.Insert(&user)
	if err!=nil{
		fmt.Println("用户数据注册插入失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*7  生成sessionid*/
	sessionid :=utils.Getmd5string(req.Mobile+req.Password+strconv.Itoa(int(time.Now().UnixNano())))

	rsp.Sessionid = sessionid

	/*8  通过sessionid 将数据存入redis*/
	bm.Put(sessionid+"user_id",id ,time.Second*600)
	bm.Put(sessionid+"mobile",user.Mobile ,time.Second*600)
	bm.Put(sessionid+"name",user.Name,time.Second*600)


	return nil



}


