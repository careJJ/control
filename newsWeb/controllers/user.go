package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) HandleRegister() {
	//获取
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	if userName == "" || pwd == "" {
		beego.Error("注册信息不完整")
		this.TplName = "register.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.Pwd = pwd
	id, err := o.Insert(&user)
	if err != nil {
		beego.Error("注册失败")
		this.TplName = "register.html"
		return
	}
	beego.Info(id)
	this.Redirect("/login",302)
}
func (this *UserController) ShowLogin() {
	this.TplName="login.html"
}

func(this *UserController) HandleLogin(){
	userName := this.GetString("userName")
		pwd := this.GetString("password")
		if userName == "" || pwd == "" {
		beego.Error("登陆信息不完整")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user,"Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.TplName = "login.html"
		return
	}
	if user.Pwd != pwd {
		beego.Error("密码错误")
		this.TplName = "login.html"
		return
	}
	this.Redirect("/index",302)


}
