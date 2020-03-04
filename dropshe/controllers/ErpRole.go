package controllers

import (
	"dropshe/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ErpController struct {
	beego.Controller
}


//展示注册页面
func (this *ErpController) ShowErpRegister() {
	this.TplName = "erpregister.html"
}


/*1.创建orm对象
2.创建查询对象
3.给查询条件赋值
4.查询*/

//this.TplName只是重新渲染页面，并不执行任何方法。
// this.Redirect()跳回本页面时执行Get绑定的方法，一般不绑定就执行controller中的Get()方法。

//处理注册业务
func (this *ErpController) HandleErpRegister() {
	//获取数据
	name := this.GetString("Name")
	pwd:= this.GetString("Password")
	//校验数据
	if name == "" || pwd == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//orm插入数据
	o := orm.NewOrm()
	var erpuser models.ErpUser
	erpuser.Name=name
	erpuser.Password=pwd
	o.Insert(&erpuser)
	//激活页面
	this.Ctx.SetCookie("name", erpuser.Name, 60*10)
	//跳转到登录页面
	this.Redirect("/login", 302)
	//返回数据
}



//展示登录界面
func (this *ErpController) ShowErpLogin() {

	this.TplName = "login.html"

}

//处理登录业务
func (this *ErpController) HandleErpLogin() {
	//获取数据
	name := this.GetString("name")
	pwd := this.GetString("Password")
	//校验数据
	if name == "" || pwd == "" {
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "login.html"
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var erpuser models.ErpUser
	erpuser.Name=name
	o.Read(&erpuser)


	//赋值
	if erpuser.Password != pwd {
			this.Data["errmsg"] = "密码错误"
			this.TplName = "login.html"
			return
	} else{
		this.SetSession("name", erpuser)
		this.Redirect("/index", 302)
	}


}

//展现erp首页
func (this *ErpController) ShowErpIndex() {
	this.TplName="index.html"
}

//




