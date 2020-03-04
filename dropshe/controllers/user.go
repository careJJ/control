package controllers

import (
	"dropshe/models"
	"dropshe/util"
	"encoding/base64"
	"net/http"
	"time"

	//"github.com/durango/go-credit-card"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
)

type UserController struct {
	beego.Controller
}
//展示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//面向对象
func RespFunc(this *beego.Controller, resp map[string]interface{}) {
	//1.把容器传递给前段
	this.Data["json"] = resp
	//2.指定传递方式 也就是值前端的ajax
	this.ServeJSON()
}

/*1.创建orm对象
2.创建查询对象
3.给查询条件赋值
4.查询*/

//this.TplName只是重新渲染页面，并不执行任何方法。
// this.Redirect()跳回本页面时执行Get绑定的方法，一般不绑定就执行controller中的Get()方法。

//处理注册业务
func (this *UserController) HandleRegister() {
	//获取数据
	email := this.GetString("Email")
	o := orm.NewOrm()
	var user models.User
	qs:=o.QueryTable("User").Filter("Email",email).Exist()
	if qs==true{
		beego.Error("该邮箱已存在")
	}
	//校验邮箱格式
	//把字符串全部大写
	//邮箱正则   ^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == "" {
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register.html"
		return
	}


	firstname := this.GetString("FirstName")
	lastname := this.GetString("LastName")
	pwd1 := this.GetString("password")
	mdpwd := util.String2md5(pwd1)

	pwd := util.AddSalt2string(mdpwd)

	//校验数据
	if email == "" || firstname == "" || lastname == "" || pwd == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "register.html"
		return
	}

	//处理数据
	//orm插入数据

	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname

	user.Password = pwd
	o.Insert(&user)

	//激活页面
	this.Ctx.SetCookie("mail", user.Email, 60*10)
	//跳转到登录页面
	this.Redirect("/login", 302)

	//返回数据
}
//展示邮箱激活
func(this*UserController)ShowEmail(){
	this.TplName = "register-email.html"
}

//处理邮箱激活业务
func(this*UserController)HandleEmail(){
	//获取数据
	email := this.GetString("email")
	pwd := this.GetString("password")
	rpwd := this.GetString("repassword")
	//校验数据
	if email == "" || pwd == ""|| rpwd == ""{
		beego.Error("输入数据不完整")
		this.Data["errmsg"] = "输入数据不完整"
		this.TplName = "register-email.html"
		return
	}
	//两次密码是否一直
	if pwd != rpwd{
		beego.Error("两次密码输入不一致")
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register-email.html"
		return
	}
	//校验邮箱格式
	//把字符串全部大写
	reg ,_:=regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == ""{
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register-email.html"
		return
	}

	//处理数据
	//发送邮件
	//utils     全局通用接口  工具类  邮箱配置   25,587
	config := `{"username":"dropshesupeng@hotmail.com","password":"Dropshe#888","host":"smtp.office365.com","port":587}`
	emailReg :=utils.NewEMail(config)
	//内容配置
	emailReg.Subject = "dropshe active"
	emailReg.From = "dropshesupeng@hotmail.com"
	emailReg.To = []string{email}
	usermail := this.Ctx.GetCookie("mail")
	emailReg.HTML = `<a href="https://app.dropshe.com/active?usermail=`+usermail+`"> Click to activate the user</a>`

	//发送
	err := emailReg.Send()
	beego.Error(err)

	//插入邮箱   更新邮箱字段
	//o := orm.NewOrm()
	//var user models.User
	//user.Email = email
	//err = o.Read(&user,"Name")
	//if err != nil {
	//	beego.Error("错误处理")
	//	return
	//}
	//user.Email = email
	//o.Update(&user)


	//返回数据
	this.Ctx.WriteString("邮件已发送，请去目标邮箱激活用户！")
}

//激活
func(this*UserController)Active(){
	//获取数据
	email := this.GetString("userName")

	if email == "" {
		beego.Error("用户名错误")
		this.Redirect("/register-email",302)
		return
	}

	//处理数据   本质上是更新active
	o := orm.NewOrm()
	var user models.User
	user.Email = email

	err := o.Read(&user,"Email")
	if err != nil {
		beego.Error("用户名不存在")
		this.Redirect("/register-email",302)
		return
	}
	user.Active = true
	o.Update(&user,"Active")

	//返回数据
	this.Redirect("/login",302)
}


//google创建凭据 -> OAuth 客户端 ID -> 网页应用，之后输入 JavaScript 来源、重定向 URI
//展示登录界面
func (this *UserController) ShowLogin() {
	//获取cookie数据，如果获取查到了，说明上一次记住用户名，不然的话，不记住用户名
	email := this.Ctx.GetCookie("Email")
	//解密
	dec, _ := base64.StdEncoding.DecodeString(email)
	if email != "" {
		this.Data["email"] = string(dec)
		this.Data["checked"] = "checked"
	} else {
		this.Data["email"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}

//处理登录业务
func (this *UserController) HandleLogin() {
	//获取数据
	email := this.GetString("Email")
	pwd1 := this.GetString("Password")
	//校验数据
	if email == "" || pwd1 == "" {
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "login.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var user models.User
	mdpwd := util.String2md5(pwd1)
	pwd := util.AddSalt2string(mdpwd)
	//赋值
	//验证邮箱格式
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result != "" {
		user.Email = email
		err := o.Read(&user, "Email")
		if err != nil {
			this.Data["errmsg"] = "邮箱未注册"
			this.TplName = "login.html"
			return
		}
		if user.Password != pwd {
			this.Data["errmsg"] = "密码错误"
			this.TplName = "login.html"
			return
		}
	}
	//返回数据u  cookie不能存中文  base64   序列化
	m1 := this.GetString("m1")
	if m1 == "2" {
		this.Ctx.SetCookie("LoginName", user.Email, 60*60)
	} else {
		this.Ctx.SetCookie("LoginName", user.Email, -1)
	}

	this.SetSession("email", user.Email)
	this.Redirect("/hub", 302)

}


//使用google账户登录
func (this *UserController)GoogleLogin(){
	http.HandleFunc("/login/oauth", util.HandleGoogleLogin)
	//获取google用户的公开信息json
	GoogleUserJson:=util.HandleGoogleCallback

	//根据获取的json存库

	//写入session

	//跳转？


}

//facebook登录





//展示账户个人信息页
func (this *UserController) ShowAccount() {
	//查询用户名、电话号和默认地址

	var user models.User
	//给查询对象赋值
	email := this.GetSession("email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	user.Email = email.(string)
	//o.Read(&user,"Email")
	firstname:=user.FirstName
	lastname:=user.LastName
	this.Data["Frist name"] = firstname
	this.Data["Last name"] = lastname
	////传地址
	//var addr models.Address
	//qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name",user.Name)
	//qs.Filter("IsDefault",true).One(&addr)
	//this.Data["addr"] = addr


	//this.Layout = "login.html"
	this.TplName = "genereal.html"
}

//更新个人信息的操作（下拉框选择国家、语言、时区）
func (this *UserController) UpdateInfo() {
	//orm插入数据
	o := orm.NewOrm()
	var user models.User
	email:=this.GetSession("email")
	user.Email=email.(string)
	//var countrys *[]models.Country
	//下拉框获取选中的类型
	//Beego中该页面Controller的Post()方法可通过this.input().Get("country")来获得select中country的value值，这些值就是被选中option的value值
	//o.QueryTable("Country").All(&countrys)
	country := this.Input().Get("Country")
	this.Data["Country"] = country
	//更新数据库
	//o.Insert(&country)
	//var timezones *[]models.TimeZone
	//o.QueryTable("TimeZone").All(&timezones)
	timezone := this.Input().Get("TimeZone")
	this.Data["timezone"] = timezone
	//更新数据库
	//o.Insert(&timezone)
	//var languages *[]models.Language
	//o.QueryTable("Language").All(&languages)
	language := this.Input().Get("Language")
	this.Data["language"] = language
	//更新数据库
	//o.Insert(&language)
	//名字可能被更改，再次插入数据库
	firstname:=this.GetString("Frist name")
	lastname:=this.GetString("Last name")
	this.Data["Frist name"] = firstname
	this.Data["Last name"] = lastname
	user.Country=country
	user.TimeZone=timezone
	user.Language=language
	o.Insert(&user)

	//测试看是否已经关联到user表中

}

//更改密码   建议删除current password
func (this *UserController) ChangePassword() {
	o := orm.NewOrm()
	var user models.User
	//获取账号名
	email := this.GetSession("Email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	this.Data["Email"] = email
	//获取旧密码
	//oldpwd := this.GetString("Current Password")
	//this.Data["Current Password"] = oldpwd
	//设置新的密码
	newpwd := this.GetString("New Password")
	this.Data["New Password"] = newpwd
	//确认密码
	repwd := this.GetString("Confirm Password")
	this.Data["Confirm Password"] = repwd
	//校验数据
	if newpwd == "" || repwd == "" {
		beego.Error("更改密码不能为空")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "genereal.html"
		return
	}
	if repwd != newpwd {
		beego.Error("两次输入密码不一致")
		this.Data["errmsg"] = "两次输入密码不一致"
		this.TplName = "genereal.html"
		return
	}
	mdpwd := util.String2md5(newpwd)
	pwd := util.AddSalt2string(mdpwd)
	user.Password = pwd
	o.Update(&user.Password)
}

//展示billing页面
func (this *UserController) ShowBilling() {
	//判定登录
	email := this.GetSession("Email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	//获取数据
	o := orm.NewOrm()
	var card models.CreditCard
	var user models.User
	//给查询对象赋值
	user.Email = email.(string)
	o.Read(&user)
	card.Email = email.(string)
	username := user.LastName + " " + user.FirstName
	this.Data["Username"] = username
	//Card Brand应该修改成???
	cardbrand := card.Company
	this.Data["Card Brand"] = cardbrand
	//卡号的部分隐藏

	before := card.Number[:3]
	after := card.Number[12:]
	cardnumber := before + "********" + after
	this.Data["Card Number"] = cardnumber
	//updatetime

	//数据校验
	if username == "" || cardbrand == "" || cardnumber == "" {
		beego.Error("获取数据失败")
		this.TplName = "genereal.html"
		return
	}
	this.TplName = "billing.html"

	//this.TplName="billing.html"
	//var card models.CreditCard
	//username:=this.GetString("Username")
	//cardbrand:=this.GetString("Card Brand")
	//cardnumber:=this.GetString("Card Number")
	//updatetime:=this.GetString("Update Time")

}

//点击update跳转到添加卡的页面
func (this *UserController) HandleCard() {
	//重定向到hub/setting/setting_billing
	this.Redirect("hub/setting/setting_billing", 302)
}

//添加新卡
func (this *UserController) AddCard() {

	cardnumber := this.GetString("卡号")
	month := this.GetString("月")
	year := this.GetString("年")
	cvc := this.GetString("CVC")
	//创建时间
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t:=(time.Now().In(loc))

	//cardbrand:=this.GetString("Card Brand")
	//校验数据
	if cardnumber == "" || month == "" || year == "" || cvc == "" {
		beego.Error("获取数据错误")
		this.TplName = "setting_billing.html"
		return
	}

	o := orm.NewOrm()
	var card models.CreditCard
	o.Read(&card)
	//给查询对象赋值
	//创建新卡
	//newcard := creditcard.Card{cardnumber,cvc,month,year}
	//判定公司

	//更新数据库
	card.Number = cardnumber
	card.Month = month
	card.Year = year
	card.CVC = cvc
	card.Active = true
	card.Updatetime=t
	//card.Company=cardbrand
	o.Insert(&card)

	this.Redirect("/accout/billing", 302)
}

//展示管理支付账户页面
func (this *UserController) ShowPayment() {
	//判定登录
	email := this.GetSession("Email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	//需要判定是电子钱包还是信用卡！


	//点击EDIT直接进入billing页面

	var card models.CreditCard
	//var user models.User
	//给查询对象赋值
	card.Email = email.(string)

	updatetime := card.Updatetime
	Active := card.Active
	//卡号的加密
	before := card.Number[:3]
	after := card.Number[12:]
	cardnumber := before + "********" + after
	this.Data["Card Number"] = cardnumber


	this.Data["Update Time"] = updatetime
	if Active == true {
		this.Data["Active"] = "Active"
	}
	if Active==false{
		this.Data["Active"] = "Inactive"
	}
	//渲染页面
		this.TplName="payment.html"

}

//管理支付账户
func (this *UserController) HandlePayment() {
	//关系到电子钱包，后续完成
}

//展示交易明细页面
func (this *UserController) ShowInvoies() {
	//关系到电子钱包，后续完成
	this.TplName="invoices.html"
}



//退出登录
func (this *UserController) Logout() {
	this.DelSession("email")
	//跳转页面
	this.Redirect("/login", 302)
}
