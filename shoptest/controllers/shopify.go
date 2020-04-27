package controllers

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"regexp"

	//"github.com/astaxie/beego/orm"
	"shoptest/util"
	"shoptest/models"
)

type CoreController struct{
	beego.Controller
}

func RespFunc(this *beego.Controller, resp map[string]interface{}) {
	//1.把容器传递给前段
	this.Data["json"] = resp
	//2.指定传递方式 也就是值前端的ajax
	this.ServeJSON()
}

//展示登录界面
func (this *CoreController) ShowLogin() {
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
func (this *CoreController) HandleLogin() {
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
	defer RespFunc(&this.Controller,resp)

	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result != "" {
		user.Email = email
		err := o.Read(&user, "Email")
		if err != nil {
			beego.Error(err)

			resp["errno"]=1
			resp["errmsg"] = "Mailbox not registered"
			this.TplName = "login.html"
			beego.Alert("1")
			return
		}
		if user.Password != pwd1 {
			resp["errno"]=2
			resp["errmsg"] = "Password mistake"
			this.TplName = "login.html"
			beego.Alert("2")
			return
		}
	}
	resp["errno"]=0
	beego.Alert("3")
	this.SetSession("id", user.Id)
	beego.Alert("4")

}



//展示一键获取订单的页面
func (this *CoreController) ShowOrder() {
	beego.Alert("showorder")
	id := this.GetSession("id")
	beego.Alert("Order:",id)
	if id == nil{
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Id=id.(int)
	o.Read(&user,"Id")
	this.Data["user"]=user
	this.TplName = "orders.html"

}



//一键同步订单
func (this *CoreController) HandleOrder() {
	beego.Alert("HandleOrder")
	id := this.GetSession("id")
	beego.Alert("Order:",id)
	if id == nil{
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	beego.Alert("HandleOrder",id)
	o:=orm.NewOrm()
	var user models.User
	user.Id=id.(int)
	o.Read(&user,"Id")
	beego.Alert("user email",user.Email)
	//跟店铺建立一对一关系
	this.Data["user"]=user
	var store models.ShopifyStore
	store.Email=user.Email
	o.Read(&store,"Email")
	//o.Read(&user.ShopifyStore)
	//创建json的容器
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)

	//得到所有订单id
	orderid := util.GetOrderid(store.Name, store.ApiKey, store.Secret)
	//直接获取所需订单所有信息  返回json字符串
	orderjson := util.GetOrderJson1(store.Name, store.ApiKey, store.Secret)
	//如果不行就定义行容器，拼接order信息再发送
	this.Data["order"] = orderjson

	//Order存库操作
	var order models.Order
	var lineitem []models.LineItems

	//var ship models.ShippingAddress
	//o1:=util.Orders{}
	//l1:=[]util.LineItems{}
	//s1:=util.ShippingAddress{}
	for f := 0; f < len(orderid); f++ {
		o1, s1, l1 := util.GetOrderStruct1(store.Name, store.ApiKey, store.Secret,f)
		//order赋值
		order.Id = o1.ID
		order.Name = o1.Name
		order.Email = o1.Email
		order.Financial_status = o1.FinancialStatus
		order.Total_price = o1.TotalPrice
		order.Created_at = o1.CreatedAt

		order.Shipping_address.Name = s1.Name
		order.Shipping_address.Zip = s1.Zip
		order.Shipping_address.LastName = s1.LastName
		order.Shipping_address.FirstName = s1.FirstName
		order.Shipping_address.Country = s1.Country
		//Company原本是interface
		order.Shipping_address.Company = s1.Company.(string)
		order.Shipping_address.Address1 = s1.Address1
		order.Shipping_address.Address2 = s1.Address2
		order.Shipping_address.City = s1.City
		order.Shipping_address.CountryCode = s1.CountryCode
		order.Shipping_address.Latitude = s1.Latitude
		order.Shipping_address.Longitude = s1.Longitude
		order.Shipping_address.Phone = s1.Phone
		order.Shipping_address.Province = s1.Province
		order.Shipping_address.ProvinceCode = s1.ProvinceCode
		o.Insert(&order)

		o.QueryTable("LineItems").Filter("Order__Id", order.Id).RelatedSel().All(&lineitem)
		for g := 0; g < len(l1); g++ {
			lineitem[g].Id = l1[g].ID
			lineitem[g].Quantity = l1[g].Quantity
			lineitem[g].VariantID = l1[g].VariantID
			lineitem[g].Price = l1[g].Price
			lineitem[g].Title = l1[g].Title
			lineitem[g].Sku = l1[g].Sku
			o.Insert(&lineitem)
		}

	}

	//获取产品匹配表ProductMatch
	var promatch models.ProductMatch
	var provanriant []models.Variant
	var imgsrc []models.ImageSrc
	a := util.Products{}
	d1 := []util.Variant{}
	imgs1 := []util.Image{}
	pid := util.GetProductId(store.Name, store.ApiKey, store.Secret)
	//p:=util.ProductMatch{}
	//	var imgstr []string

	for i := 0; i < len(pid); i++ {
		a.ID, a.Title, d1, imgs1 = util.GetProductMatch(store.Name, store.ApiKey, store.Secret, i)
		promatch.Id = a.ID
		promatch.Title = a.Title
		//先插入id和title，再一对多查表插入
		o.Insert(&promatch)
		//查表 variant  ProductMatch是主表 variant子表，variant建表时，自动生成ProductMatch_id字段，等于ProductMatch的ProductId
		o.QueryTable("Variant").Filter("ProductMatch__Id", a.ID).RelatedSel().All(&provanriant)
		for c := 0; c < len(d1); c++ {
			provanriant[c].Id = d1[c].Id
			provanriant[c].Sku = d1[c].Sku
			provanriant[c].Price = d1[c].Price
			provanriant[c].Title = d1[c].Title
			o.Insert(&provanriant)
		}
		//查表ImageSrc,
		o.QueryTable("ImageSrc").Filter("ProductMatch__Id", a.ID).RelatedSel().All(&imgsrc)
		for e := 0; e < len(imgs1); e++ {
			//img链接
			imgsrc[e].Src = imgs1[e].Src
			//img存FDFS
			fdfssrc := util.FdsUploadImage(&this.Controller, imgs1[e].Src)
			imgsrc[e].FdfsSrc = fdfssrc
			o.Insert(&imgsrc)
		}
	}

	//定义行容器
	var orders []map[string]interface{}
	//定义行容器

	var ordershow []models.Order
	//var productshow []models.ProductMatch

	//
	o.QueryTable("Order").All(&ordershow)
	for i1 := 0; i1 < len(ordershow); i1++ {
		t := make(map[string]interface{})
		t["Id"] = ordershow[i1].Id
		t["Name"] = ordershow[i1].Name
		t["Create_at"] = ordershow[i1].Created_at
		t["Total_price"] = ordershow[i1].Created_at
		t["Status"] = ordershow[i1].Financial_status
		var line []models.LineItems
		o.QueryTable("LineItems").Filter("Order__Id", ordershow[i1].Id).RelatedSel().All(&line)
		t["line_items"] = line
		//var ship models.ShippingAddress
		//
		o.QueryTable("Order").Filter("Order__Id", ordershow[i1].Id).RelatedSel().All(&ordershow)
		//City //	Zip //	Province //	Country //	LastName
		t["City"] = ordershow[i1].Shipping_address.City
		t["Zip"] = ordershow[i1].Shipping_address.Zip
		t["Province"] = ordershow[i1].Shipping_address.Province
		t["Country"] = ordershow[i1].Shipping_address.Country
		t["LastName"] = ordershow[i1].Shipping_address.LastName
		orders = append(orders, t)

	}
	for _, v1 := range orders {
		var second []map[string]interface{}
		for _, v2 := range v1["line_items"].([]models.LineItems) {
			t := make(map[string]interface{})

			t["Sku"] = v2.Sku
			t["Quantity"] = v2.Quantity
			t["Title"] = v2.Title
			t["Price"] = v2.Price
			//获取相应图片  一个product_id对应多个img,只需要一个img
			var imgsrc []models.ImageSrc
			o.QueryTable("ImageSrc").Filter("ProductMatch__Id", v2.Id).RelatedSel("imgsrc").All(&imgsrc)
			//TODO 区分图片
			t["Image"] = imgsrc[0].FdfsSrc
			second = append(second, t)

		}
		//将二级容器放到总容器
		//TODO  用line_items可能错
		v1["line_items"] = second

	}

	this.Data["order"] = orders

	this.TplName = "order.html"
	//defer RespFunc(&this.Controller,resp)

	//获取对应页的数据   获取几条数据     起始位置
	//ORM多表查询的时候默认是惰性查询 关联查询之后，如果关联的字段为空，数据查询不到
	//分页方案

}
