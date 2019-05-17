package controllers

import "github.com/astaxie/beego"

type CartController struct {
	beego.Controller
}

func (this*CartController)HandleAddCart(){
	id,err := this.GetInt("goodsId")
	num,err2 := this.GetInt("num")
	//返回ajax步骤
	//定义一个map容器
	resp := make(map[string]interface{})

	//封装，集成，多态
	defer RespFunc(&this.Controller,resp)

	//校验数据
	if err != nil || err2 != nil{
		resp["errno"] = 1
		resp["errmsg"] = "输入数据不完整"
		return
	}
	//校验登录状态
	name := this.GetSession("name")
	if name == nil{
		resp["errno"] = 2
		resp["errmsg"] = "当前用户未登录，不能添加购物车"
		return
	}
	beego.Info(id,num)

}