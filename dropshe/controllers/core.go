package controllers

import (
	"dropshe/models"
	"dropshe/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"context"
	"strings"
	"time"
)

type CoreController struct {
	beego.Controller
}

//展示商店页面
func (this *CoreController) ShowStore() {
	//给查询对象赋值
	email := this.GetSession("email")
	if email == "" {
		this.TplName = "login.html"
		return
	}
	var user models.User
	user.Email = email.(string)

	//一对一查店铺
	o := orm.NewOrm()
	var store models.ShopifyStore
	//一对一查询，user
	//在User表中给store添加rel(one)，就表示跟store建立一对一关系（ 一个用户对应一个店铺；一个店铺对应一个用户），user是子表，user表中自动生成了store_id字段，store表中的反向关系可以省略。
	//一：可以在查询了子表user,o.Read(user)后，直接查询o.Read(user.Store)。
	//二：级联查询（主流方法）：user=&User{},O.QueryTable("user").Filter("Id", 1).RelatedSel().One(user),这样查询出user表中id为1的所有store数据，可以直接得到user.store。
	//三：reverse查询，通过子表条件查询主表，此时并没有获取另一个表的数据。store:=store{},O.QueryTable("profile").Filter("User__Id", 1).One(&store)。
	// 或者store := []*Store{}，O.QueryTable("store").Filter("User__Name", "ming").One(&stores)，for _, a := range stores，此处的a就是名字是ming的store。

	o.QueryTable("User").Filter("Email", user.Email).RelatedSel().One(user)

	name := user.ShopifyStore.Name
	store.Name = name
	if name == "" {
		beego.Error("请绑定商店")
		return
		this.Redirect("/stores/add", 302)
	}
	if name != "" {

		t := store.CreateTime
		psync := store.ProductSync
		osync := store.OrderSync

		this.Data["CreateTime"] = t
		this.Data["Product Sync"] = psync
		this.Data["Order Sync"] = osync
		this.Data["Status"] = "Active"
	}
	this.TplName = "stores.html"

}

//跳转到ADD页面
func (this *CoreController) HandleStore() {
	this.Redirect("/stores/add", 302)
}

//添加商店的操作
func (this *CoreController) AddStore() {

	//给查询对象赋值
	email := this.GetSession("email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Email=email.(string)
	o.Read(&user)
	//跟店铺建立一对一关系

	var store models.ShopifyStore
	store.User.Id=user.Id
	store.Email = email.(string)
	name := this.GetString("name")
	//暂时定为password
	password := this.GetString("password")
	apikey := this.GetString("Api key")
	//判定店铺是否已经存在
	qs := o.QueryTable("ShopifyStore").Filter("Name", name).Exist()
	if qs == true {
		beego.Error("该店铺已存在")
		this.Redirect("/stores/add", 302)
	}
	//创建时间  店铺的添加时间
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := (time.Now().In(loc))
	store.CreateTime = t.String()

	//传入店铺名字和password，之前应该先出一份安全声明。

	client := util.LinkStore(name, password)
	s, _, err := client.Shop.Get(context.Background())
	if err != nil {
		beego.Error("连接shopify店铺失败")
	}
	store.Name = s.Name
	store.Secret = password
	store.ApiKey = apikey
	o.Insert(&store)

	this.Redirect("/stores", 302)
}

//展示一键获取订单的页面
func (this *CoreController) ShowOrder() {
	this.TplName = "order.html"
}

//一键同步订单
func (this *CoreController) HandleOrder() {
	//获取店铺，查店铺信息
	email := this.GetSession("email")
	var user models.User
	user.Email = email.(string)
	//判断是否已经绑定店铺，若未绑定，则跳转到添加店铺页面
	o := orm.NewOrm()
	//o.Read(&user.ShopifyStore)
	o.QueryTable("User").Filter("Email", email).RelatedSel().One(user)
	//创建json的容器
	resp := make(map[string]interface{})
	if user.ShopifyStore.Name == "" {
		this.Redirect("/stores/add", 302)
	}
	//if user.ShopifyStore.Name=="beachmolly mygymgirl"	||user.ShopifyStore.Name=="mygymgirl beachmolly "{
	//	//得到beachmolly的所有订单相关json
	//	beachOrderJson:=util.GetBeachmollyOrderInfo()
	//	//得到Mygymgirl的所有订单相关json
	//	mygymOrderJson:=util.GetMygymgirlOrderInfo()
	//	resp["order"]=beachOrderJson+mygymOrderJson
	name := user.ShopifyStore.Name
	apikey := user.ShopifyStore.ApiKey
	secret := user.ShopifyStore.Secret
	//得到所有订单id
	orderid := util.GetOrderid(name, apikey, secret)
	//直接获取所需订单所有信息  返回json字符串
	orderjson := util.GetOrderJson(name, apikey, secret)
	//如果不行就定义行容器，拼接order信息再发送
	resp["order"] = orderjson

	//Order存库操作
	var order models.Order
	var lineitem models.LineItems
	//var ship models.ShippingAddress
	//o1:=util.Orders{}
	//l1:=[]util.LineItems{}
	//s1:=util.ShippingAddress{}
	for f := 0; f < len(orderid); f++ {
		o1, s1, l1 := util.GetOrderStruct(name, apikey, secret, f)
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
		order.Shipping_address.Company = s1.Company
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
			lineitem.Id = l1[g].ID
			lineitem.Quantity = l1[g].Quantity
			lineitem.VariantID = l1[g].VariantID
			lineitem.Price = l1[g].Price
			lineitem.Title = l1[g].Title
			lineitem.Sku = l1[g].Sku
			o.Insert(&lineitem)
		}

	}

	//获取产品匹配表ProductMatch
	var promatch models.ProductMatch
	var provanriant models.Variant
	var imgsrc models.ImageSrc
	a := util.Products{}
	d1 := []util.Variant{}
	imgs1 := []util.Image{}
	pid := util.GetProductId(name, apikey, secret)
	//p:=util.ProductMatch{}
	//	var imgstr []string

	for i := 0; i < len(pid); i++ {
		a.ID, a.Title, d1, imgs1 = util.GetProductMatch(name, apikey, secret, i)
		promatch.Id = a.ID
		promatch.Title = a.Title
		//先插入id和title，再一对多查表插入
		o.Insert(&promatch)
		//查表 variant  ProductMatch是主表 variant子表，variant建表时，自动生成ProductMatch_id字段，等于ProductMatch的ProductId
		o.QueryTable("Variant").Filter("ProductMatch__Id", a.ID).RelatedSel().All(&provanriant)
		for c := 0; c < len(d1); c++ {
			provanriant.Id = d1[c].Id
			provanriant.Sku = d1[c].Sku
			provanriant.Price = d1[c].Price
			provanriant.Title = d1[c].Title
			o.Insert(&provanriant)
		}
		//查表ImageSrc,
		o.QueryTable("ImageSrc").Filter("ProductMatch__Id", a.ID).RelatedSel().All(&imgsrc)
		for e := 0; e < len(imgs1); e++ {
			//img链接
			imgsrc.Src = imgs1[e].Src
			//img存FDFS
			fdfssrc := util.FdsUploadImage(&this.Controller, imgs1[e].Src)
			imgsrc.FdfsSrc = fdfssrc
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
		o.QueryTable("Order").Filter("Order__Id", ordershow[i1].Id).RelatedSel().One(ordershow)
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

//从订单页跳转到sourcing agent订单的采购代理
func (this *CoreController) ShowSource() {
	//判定登录
	email := this.GetSession("email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	id, err := this.GetInt64("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/order", 302)
		return
	}
	var order models.Order
	//获取订单详情
	o := orm.NewOrm()
	//建立跟Order的对应关系
	order.Id = id
	o.Read(&order)
	var line []models.LineItems
	var image models.ImageSrc
	var source models.Sourcing_Demand
	//建立Sourcing_Demand跟Order的一对一关系
	source.Order.Id=id
	//根据订单id查订单详细商品
	o.QueryTable("LineItems").Filter("Order__Id", id).RelatedSel("LineItems").All(&line)

	//每个line对应一个product_id，每个product_id对应多个image，只需获取对应的一张图
	//根据订单商品查详细图片
	var images []string
	for i := 0; i < len(line); i++ {
		o.QueryTable("ImageSrc").Filter("ProductMatch__Id", line[i].Id).RelatedSel("ImageSrc").All(&image)

		images = append(images, image.FdfsSrc)
	}
	//var user models.User
	//user.Email=email.(string)


	//目的地，在ERP中使用
	source.Destination = order.Shipping_address.Country

	imagepath := util.FdsUploadImage(&this.Controller, "images")

	images = append(images, imagepath)
	imgsrc := strings.Join(images, ",")
	source.Images = imgsrc
	//返回文件ID 存入数据库
	//source.Images=fdfsresponse.RemoteFileId

	//product_price 是客户的期待价格
	product_price := this.GetString("Target product price(USD)")
	ship_price := this.GetString("Target shipping price(USD)")
	delivery := this.Input().Get("Estimated Delivery")
	//count, err := this.GetInt("SOURCING QTY(Option)")
	//if err != nil {
	//	beego.Error(err)
	//	return
	//}
	//eub,err:=this.GetBool("Epacket")
	title := this.GetString("Product Title")
	link := this.GetString("Souring Porduct Link")
	description := this.GetString("Description")
	general, err := this.GetBool("General Cargo")
	if err != nil {
		beego.Error(err)
		return
	}
	//添加两个勾选框，内容是普货还是特货，下拉框发货国家
	//标题
	source.Title = title
	//采购链接
	source.Link = link
	//采购数量 TODO 没有参考价值？
	//source.Count = count
	source.Target_price = product_price
	source.Estimated_Delivery = delivery
	source.Shipping_price = ship_price
	source.General = general
	//是否易邮宝国家？
	//source.EUB=eub

	source.Description = description
	o.Insert(&source)
	this.Redirect("/sourcing", 302)
}

//获取对应页的数据   获取几条数据     起始位置
//ORM多表查询的时候默认是惰性查询 关联查询之后，如果关联的字段为空，数据查询不到
/*	//处理数据
	//查询文章数据
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//给查询条件赋值
	article.Id = id
	//查询
	o.Read(&article)

	//多对多查询一
	//o.LoadRelated(&article,"Users")

	//高级查询   首先要指定表  多对多查询二   获取用户名   为了使用高级查询
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"] = users

	//给更新条件赋值
	article.ReadCount += 1
	o.Update(&article)

	//返回数据
	this.Data["article"] = article

	//插入多对多关系  根据用户名获取用户对象
	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")

	//多对多的插入操作
	//获取ORM对象

	//获取被插入数据的对象  文章

	//获取多对多操作对象
	m2m := o.QueryM2M(&article,"Users")

	//用多对多操作对象插入
	m2m.Add(user)*/

//展示所有采购需求的页面”/sourcing“
func (this *CoreController) ShowSourcing() {
	email := this.GetSession("email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	var user models.User
	var source_demand []models.Sourcing_Demand
	user.Email = email.(string)
	userid := user.Id

	//一个客户对应自己的多个采购需求，一对多，用户是主表，采购需求是子表，自动产生user_id字段，根据user_id字段查找数据
	/*一：级联查询：var posts []*Post ⇥ O.QueryTable("post").Filter("User__Id", 1).RelatedSel().All(&posts) ⇥ for _, v := range posts {}，v就是post所有数据。可以直接得到：v.User.Name
	二：reverse查询：var user User ⇥ err := O.QueryTable("user").Filter("Post__Title", "paper1").Limit(1).One(&user)*/
	o := orm.NewOrm()
	//根据user查source_Demand表
	o.QueryTable("Sourcing_Demand").Filter("User__Id", userid).RelatedSel().All(&source_demand)

	//遍历
	//创建总容器
	var source []map[string]interface{}
	var sourcing models.Sourcing

	for _, v := range source_demand {
		//创建行容器
		temp:=make(map[string]interface{})
		//先判断订单的状态



		if v.Status == false {
			temp["Status"] = "Pending review"
			//关闭支付按钮
			temp["pay"]=false
		}
		if v.Status == true {
			temp["Status"] = "pending payment"
			//关闭支付按钮
			temp["pay"] = true
		}
		temp["Id"]=v.Id
		temp["title"]=v.Title
		temp["targetprice"]=v.Target_price
		temp["delivery"]=v.Estimated_Delivery
		//获得图片切片
		image := util.String2slice(v.Images)
		//只展示一张图片
		temp["image"]=image[0]
	source=append(source,temp)
		//插入数据到Sourcing表
		sourcing.Id = v.Id
		sourcing.Title = v.Title
		sourcing.Images = v.Images
		sourcing.Target_price = v.Target_price
		sourcing.Estimated_Delivery = v.Estimated_Delivery
		sourcing.Status = v.Status
		sourcing.Link = v.Link
		sourcing.Destination=v.Destination
		sourcing.General=v.General
		o.Insert(&sourcing)
	}

	this.Data["source"]=source
	this.TplName = "Sourcing.html"
}

//直接提采购需求的页面


//Sourcing页面的操作，Action暂定为删除和支付（加入购物车），在状态更新为Pending payment才可以进行支付

//删除Sourcing  前端添加删除按钮和重定向
func (this *CoreController) DeleteSouring() {

	id, err := this.GetInt("Id")
	if err != nil {
		beego.Error("获取sourcing id失败")
		this.Redirect("/Sourcing", 302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var source models.Sourcing
	source.Id = id
	//删除数据
	o.Delete(&source, "Id")

	this.TplName = "Sourcing.html"

}

//支付  sku的匹配后
func (this *CoreController) PaySoucing() {
	//获取订单的状态，状态待支付传true给前端用来开放按钮
	email := this.GetSession("email")
	if email == "" {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	var user models.User
	user.Email = email.(string)
	//uid:=user.Id
	id, err := this.GetInt("Id")
	if err != nil {
		beego.Error("获取sourcing id失败")
		this.Redirect("/Sourcing", 302)
		return
	}
	o := orm.NewOrm()
	var source models.Sourcing
	//o.QueryTable("source").Filter("User__id",uid).All(&source)
	source.Id = id
	//检验status是否可以支付
	if source.Status == true {
		//添加到购物车

	}

}

//搜索  按订单名字搜索

//分页
