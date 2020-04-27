package controllers

import (
	"dropshe/models"
	"dropshe/util"
	"math/rand"
	//"net/http"
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
	email := this.GetString("email")
	beego.Alert(email)
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	o := orm.NewOrm()
	qs := o.QueryTable("User").Filter("Email", email).Exist()
	if qs == true {
		beego.Error("该邮箱已存在")
		resp["errno"] = 1
		resp["errmsg"] = "The mailbox has been registered"
		this.TplName = "register.html"
		return
	}
	//校验邮箱格式
	//把字符串全部大写
	//邮箱正则   ^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$

	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == "" {
		beego.Error("邮箱格式错误")
		resp["errno"] = 2
		resp["errmsg"] = "The mailbox format is incorrect"
		this.TplName = "register.html"
		return
	}
	firstname := this.GetString("first_name")
	lastname := this.GetString("last_name")
	pwd1 := this.GetString("password")
	storename := this.GetString("Shopify_name")
	phone:=this.GetString("Phone_number")
	//mdpwd := util.String2md5(pwd1)
	//pwd := util.AddSalt2string(mdpwd)
	//校验数据
	if email == "" || firstname == "" || lastname == "" || pwd1 == "" || storename == "" ||phone==""{
		beego.Error("获取数据错误")
		//this.Data["errmsg"] = "获取数据错误"
		resp["errno"] = 3
		resp["errmsg"] = "Please complete the registration form"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//orm插入数据
	var user models.User
	//var store models.ShopifyStore
	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname
	user.Password = pwd1
	//user.ShopifyStore.Name=storename
	user.StoreName = storename
	user.Number=phone
	//设置随机数种子
	beego.Alert(phone)
	rand.Seed(time.Now().UnixNano())
	//获取随机数
	//num:=time.Now().Format("20060102150405")
	//num1:=rand.Intn(99999999)
	//user.Id=num1
	_,err:=o.Insert(&user)
	if err!=nil{
		beego.Error(err)
	}
	//beego.Alert(user.Email)
	//激活页面
	//this.Ctx.SetCookie("email", user.Email, 60*10)
	//跳转到登录页面
	resp["errno"] = 0
	//this.Redirect("/login", 302)
	beego.Alert("2")
	//返回数据
}

//展示邮箱激活
func (this *UserController) ShowActivateEmail() {
	this.TplName = "registerActivate.html"
}

//处理邮箱激活业务
func (this *UserController) ActivateEmail() {
	//获取数据
	email := this.GetString("email")
	pwd := this.GetString("password")
	rpwd := this.GetString("repassword")
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//校验数据
	if email == "" || pwd == "" || rpwd == "" {
		beego.Error("输入数据不完整")
		resp["errno"] = 1
		resp["errmsg"] = "Please complete the information"
		this.TplName = "registerActivate.html"
		return
	}
	//两次密码是否一直
	if pwd != rpwd {
		beego.Error("两次密码输入不一致")
		resp["errno"] = 2
		resp["errmsg"] = "The passwords entered are different"
		this.TplName = "registerActivate.html"
		return
	}
	//校验邮箱格式
	//把字符串全部大写
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == "" {
		beego.Error("邮箱格式错误")
		resp["errno"] = 3
		resp["errmsg"] = "Mailbox format error"
		this.TplName = "registerActivate.html"
		return
	}
	config := `{"username":"dropshetest@163.com","password":"HUAZYHNAFGPIWPWQ","host":"smtp.163.com","port":25}`
	emailReg := utils.NewEMail(config)
	//内容配置
	//标题
	emailReg.Subject = "dropshe account activate"
	emailReg.From = "dropshetest@163.com"
	emailReg.To = []string{email}

	//usermail := this.Ctx.GetCookie("mail")
	emailReg.HTML = `<a href="http://192.168.71.137:8080/activate?usermail=` + email + `">Activate your dropshe account</a>`
	beego.Alert("send finish")
	//发送
	err := emailReg.Send()
	util.LogError(err)
	//返回数据
	this.Ctx.WriteString("The mail has been sent, please go to the target mailbox to activate the account")
}

////激活
func (this *UserController) Activate() {
	//获取数据
	email := this.GetString("usermail")

	if email == "" {
		beego.Error("用户名错误")
		this.Redirect("/register/activate", 302)
		return
	}
	beego.Alert("email",email)
	//处理数据   本质上是更新Activate
	o := orm.NewOrm()
	var user models.User
	user.Email = email
	err := o.Read(&user, "Email")
	if err != nil {
		util.LogError("用户名不存在")
		this.Redirect("/register/activate", 302)
		return
	}
	user.Activate = true
	o.Update(&user, "Activate")

	//返回数据
	this.Redirect("/login", 302)
}

//google创建凭据 -> OAuth 客户端 ID -> 网页应用，之后输入 JavaScript 来源、重定向 URI
//展示登录界面
func (this *UserController) ShowLogin() {
	//获取cookie数据，如果获取查到了，说明上一次记住用户名，不然的话，不记住用户名
	//email := this.Ctx.GetCookie("Email")
	////解密
	//dec, _ := base64.StdEncoding.DecodeString(email)
	//if email != "" {
	//	this.Data["email"] = string(dec)
	//	this.Data["checked"] = "checked"
	//} else {
	//	this.Data["email"] = ""
	//	this.Data["checked"] = ""
	//}

	this.TplName = "login.html"
}

//处理登录业务
func (this *UserController) HandleLogin() {
	//获取数据
	email := this.GetString("email")
	pwd1 := this.GetString("password")
	//校验数据
	if email == "" || pwd1 == "" {
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "login.html"
		return
	}
	beego.Alert(email)
	beego.Alert(pwd1)
	//处理数据
	o := orm.NewOrm()
	var user models.User
	//mdpwd := util.String2md5(pwd1)
	//pwd := util.AddSalt2string(mdpwd)
	//赋值
	//验证邮箱格式
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)

	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result != "" {
		user.Email = email
		err := o.Read(&user, "Email")
		if err != nil {
			util.LogError(err)

			resp["errno"] = 1
			resp["errmsg"] = "Mailbox not registered"
			this.TplName = "login.html"
			beego.Alert("1")
			return
		}
		if user.Password != pwd1 {
			resp["errno"] = 2
			resp["errmsg"] = "Password mistake"
			this.TplName = "login.html"
			beego.Alert("2")
			return
		}
	}
	resp["errno"] = 0
	beego.Alert("3")
	this.SetSession("id", user.Id)
	beego.Alert("4")

}

//使用google账户登录
//func (this *UserController)GoogleLogin(){
//	http.HandleFunc("/login/oauth", util.HandleGoogleLogin)
//	//获取google用户的公开信息json
//	GoogleUserJson:=util.HandleGoogleCallback
//
//	//根据获取的json存库
//
//	//写入session
//
//	//跳转？
//
//
//}
//facebook登录

//展示账户个人信息页
func (this *UserController) ShowAccount() {
	//查询用户名、电话号和默认地址
	var user models.User
	//给查询对象赋值
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	user.Id = id.(int)
	o.Read(&user, "Id")
	this.Data["user"] = user

	this.TplName = "Account.html"
}

//更新个人信息的操作（下拉框选择国家、语言、时区）
func (this *UserController) UpdateInfo() {
	//orm插入数据

	var user models.User
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	user.Id = id.(int)
	o.Read(&user, "Id")

	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//下拉框获取选中的类型
	//Beego中该页面Controller的Post()方法可通过this.input().Get("country")来获得select中country的value值，这些值就是被选中option的value值
	//o.QueryTable("Country").All(&countrys)
	//country := this.Input().Get("Country")
	country := this.GetString("country")

	timezone := this.GetString("time_zone")
	language := this.GetString("language")

	firstname := this.GetString("frist_name")
	lastname := this.GetString("last_name")
	if lastname == "" {
		beego.Alert("3")
		this.TplName = "Account.html"
		return
	}
	//this.Data["Frist name"] = firstname
	//this.Data["Last name"] = lastname
	user.FirstName = firstname
	user.LastName = lastname
	user.Country = country
	user.TimeZone = timezone
	user.Language = language
	beego.Alert("5")
	//报错1062
	//_,err:=o.Insert(&user)
	//if err!=nil{
	//	beego.Error(err)
	//	resp["errno"]=1
	//	resp["errmsg"]="Data update failed"
	//	this.TplName="Account.html"
	//	return
	//}
	_, err := o.Update(&user)
	if err != nil {
		beego.Error(err)
		resp["errno"] = 1
		resp["errmsg"] = "Data update failed"
		this.TplName = "Account.html"
		return
	}
	//image,head,err:=this.GetFile("image")
	////获取图片
	////返回值 文件二进制流  文件头    错误信息
	//if err != nil {
	//	beego.Error("图片上传失败")
	//	this.Data["errmsg"] = "图片上传失败"
	//	this.TplName = "Account.html"
	//}
	//defer image.Close()
	////校验文件大小
	//if head.Size >5000000{
	//	beego.Error("图片数据过大")
	//	this.Data["errmsg"] = "图片数据过大"
	//	this.TplName = "Account.html"
	//}
	////校验格式 获取文件后缀
	//ext := path.Ext(head.Filename)
	//if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
	//	beego.Error("上传文件格式错误")
	//	this.Data["errmsg"] = "上传文件格式错误"
	//	this.TplName = "Account.html"
	//}
	//
	//fileBuffer := make([]byte,head.Size)
	////把文件数据读入到fileBuffer中
	//image.Read(fileBuffer)
	////获取client对象
	//client,err := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	//if err!=nil{
	//	beego.Error(err)
	//}
	////上传 得到fdfsresponse.RemoteFileId就是所需字段
	//fdfsresponse,_:=client.UploadByBuffer(fileBuffer,ext[1:])
	//
	//user.Image=fdfsresponse.RemoteFileId

	beego.Alert("6")
	resp["errno"] = 0
	//this.Redirect("/sourcing",302)
	//测试看是否已经关联到user表中
	beego.Alert("7")

}

//展示更改密码页面
func (this *UserController) ShowPassword() {
	beego.Alert("email:", 2)
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")

	this.Data["user"] = user
	this.TplName = "Email.html"

}

//更改密码   建议删除current password
func (this *UserController) HandlePassword() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//获取账号名
	id := this.GetSession("id")
	beego.Alert("change password:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	beego.Alert(user.Email)

	newpwd := this.GetString("new_password")
	//确认密码
	repwd := this.GetString("confirm_password")
	//校验数据
	if newpwd == "" || repwd == "" {
		beego.Error("更改密码不能为空")
		resp["errno"] = 1
		resp["errmsg"] = "Please provide full information"
		this.TplName = "Email.html"
		return
	}
	if repwd != newpwd {
		beego.Error("两次输入密码不一致")
		resp["errno"] = 2
		resp["errmsg"] = "The passwords entered are inconsistent"
		this.TplName = "Email.html"
		return
	}
	//mdpwd := util.String2md5(newpwd)
	//pwd := util.AddSalt2string(mdpwd)
	user.Password = repwd
	beego.Alert("change8", repwd)
	_, err := o.Update(&user, "Password")
	if err != nil {
		beego.Error(err)
		beego.Alert("update password error")
	}

	beego.Alert("9")
	resp["errno"] = 5
	resp["errmsg"] = "finished"
	//this.Redirect("/account/general",302)
	beego.Alert("10")

}

//展示billing页面
func (this *UserController) ShowBilling() {
	//判定登录
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")

	var card models.CreditCard

	//给查询对象赋值

	card.Email = user.Email
	o.Read(&card, "Email")

	//username := user.LastName + " " + user.FirstName
	this.Data["user"] = user

	cardnumber := card.Number

	//数据校验
	if cardnumber == "" {

		//this.TplName = "billing_setting.html"
		beego.Alert("获取信用卡信息失败")
		//return
		//this.Data["Card Number"]=""
		this.Redirect("/setting/setting_billing", 302)
	} else {
		//before := cardnumber[:3]
		//after := cardnumber[12:]
		//cnumber := before + "********" + after
		//beego.Alert(cnumber)
		//this.Data["card_number"]=cnumber
		this.Data["card"] = card
		this.TplName = "Billing.html"
	}

}

//点击update跳转到添加卡的页面
func (this *UserController) HandleCard() {
	//重定向到hub/setting/setting_billing
	this.Redirect("/setting/setting_billing", 302)
}

func (this *UserController) ShowAddCard() {
	//判定登录
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}

	beego.Alert("show add card 1")
	this.TplName = "settingBilling.html"
}

//添加新卡
func (this *UserController) AddCard() {

	//判定登录
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	var card models.CreditCard

	card.Email = user.Email
	o.Read(&card, "Email")

	beego.Alert("add card  1")
	cardnumber := this.GetString("card_number")
	year := this.GetString("year")
	month := this.GetString("month")

	cvv := this.GetString("CVV")
	beego.Alert("get billing  2")
	beego.Alert(cardnumber, year, month, cvv)
	//创建时间
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := (time.Now().In(loc))
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)

	//cardbrand:=this.GetString("Card Brand")
	//校验数据
	if cardnumber == "" || month == "" || year == "" || cvv == "" {
		util.LogError("添加新卡失败")
		resp["errno"] = 2
		resp["errmsg"] = "Please provide full information"
		this.TplName = "setting_billing.html"
		return
	}
	if len(cardnumber) != 16 || len(year) != 4 || len(month) > 2 || len(cvv) != 3 {
		resp["errno"] = 3
		resp["errmsg"] = "Please fill in the correct credit card information"
		beego.Alert("Please fill in the correct credit card information")
		this.TplName = "setting_billing.html"
		return
	}
	//更新数据库
	card.Number = cardnumber
	card.Month = month
	card.Year = year
	card.CVC = cvv
	card.User = &user
	card.Activate = true
	card.Updatetime = t
	//card.Company=cardbrand
	_, err := o.Update(&card, "Number", "Month", "Year", "CVC", "Activate", "Updatetime", "User")
	if err != nil {
		beego.Error(err)
	}
	beego.Alert("update billing success")
	resp["errno"] = 1
	beego.Alert("redirect 1")
	//this.Redirect("/account/billing", 302)
}

//展示管理支付账户页面
func (this *UserController) ShowPayment() {
	//判定登录
	//判定登录
	id := this.GetSession("id")
	beego.Alert("Account:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	this.Data["user"] = user
	//需要判定是电子钱包还是信用卡！

	//点击EDIT直接进入billing页面

	var card models.CreditCard
	//var user models.User
	//给查询对象赋值
	card.Email = user.Email

	o.Read(&card, "Email")
	cardnumber := card.Number

	Activate := card.Activate
	//数据校验
	if cardnumber == "" {
		beego.Error("获取数据失败")
		//this.TplName = "billing_setting.html"
		beego.Alert("获取信用卡信息失败")
		//return
		this.Data["card"] = card
		this.TplName = "payment.html"
	} else {
		//before := cardnumber[:3]
		//after := cardnumber[12:]
		//cnumber := before + "********" + after
		//this.Data["Card Number"] = cnumber
		this.Data["card"] = card
		this.TplName = "payment.html"
	}
	if Activate == true {
		this.Data["Activate"] = "Activate"
	}
	if Activate == false {
		this.Data["Activate"] = "Inactivate"
	}
	//渲染页面
	this.TplName = "payment.html"
}

//管理支付账户
func (this *UserController) HandlePayment() {
	//关系到电子钱包，后续完成
}

//展示交易明细页面
func (this *UserController) ShowInvoies() {
	//关系到电子钱包，后续完成
	this.TplName = "invoices.html"
}

func (this *UserController) ShowForgetPassword() {
	this.TplName = "forgot_password.html"
}

func (this *UserController) FindPassword() {
	//从前端获取邮箱，发送该邮箱的密码到该邮箱，是否还需要手机号码用于确认信息
	email := this.GetString("email")
	phone:=this.GetString("phone")
	beego.Alert(phone)
	//number:=this.GetString("number")
	var user models.User
	user.Email = email
	o := orm.NewOrm()
	o.Read(&user, "Email")
	//确认用户的手机号码
	if phone!=user.Number{
		return
	}else {
		config := `{"username":"dropshetest@163.com","password":"HUAZYHNAFGPIWPWQ","host":"smtp.163.com","port":25}`
		emailReg := utils.NewEMail(config)
		//内容配置
		//标题
		emailReg.Subject = "dropshe forget password"
		emailReg.From = "dropshetest@163.com"
		emailReg.To = []string{email}
		emailReg.Text = user.Password

		//usermail := this.Ctx.GetCookie("email")
		emailReg.HTML = `<a href="http://192.168.71.137:8080/forgot_password?usermail=` + email + `">`+ user.Password+`</a>`
		this.Ctx.WriteString("The mail has been sent, please go to the target mailbox to get the password")
		beego.Alert("send finish",user.Password)
		//发送
		err := emailReg.Send()
		if err != nil {
			beego.Error(err)
		}

	}

}

//退出登录
func (this *UserController) Logout() {
	this.DelSession("id")
	//跳转页面
	this.Redirect("/login", 302)
}
