package controllers

import (
	"dropshe/models"
	"dropshe/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"math"
	"reflect"
	"strconv"
	"time"
	//"github.com/stripe/stripe-go/product"
)

type CoreController struct {
	beego.Controller
}

//展示商店页面
func (this *CoreController) ShowStore() {
	//给查询对象赋值
	id := this.GetSession("id")
	beego.Alert("Show Store:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	var store models.ShopifyStore
	store.Email = user.Email
	//一对一查询，user
	//在User表中给store添加rel(one)，就表示跟store建立一对一关系（ 一个用户对应一个店铺；一个店铺对应一个用户），user是子表，user表中自动生成了store_id字段，store表中的反向关系可以省略。
	//一：可以在查询了子表user,o.Read(user)后，直接查询o.Read(user.Store)。
	//二：级联查询（主流方法）：user=&User{},O.QueryTable("user").Filter("Id", 1).RelatedSel().One(user),这样查询出user表中id为1的所有store数据，可以直接得到user.store。
	//三：reverse查询，通过子表条件查询主表，此时并没有获取另一个表的数据。store:=store{},O.QueryTable("profile").Filter("User__Id", 1).One(&store)。
	// 或者store := []*Store{}，O.QueryTable("store").Filter("User__Name", "ming").One(&stores)，for _, a := range stores，此处的a就是名字是ming的store。
	o.Read(&user, "Email")
	//o.QueryTable("User").Filter("Email", user.Email).RelatedSel().One(user)
	beego.Alert(user.Email)
	o.Read(&store, "Email")
	name := store.Name
	this.Data["user"] = user
	this.Data["store"] = store
	beego.Alert("store name:", name)
	if name == "" {
		beego.Error("请绑定商店")
		//this.Data["errno"]=0
		this.Redirect("/stores/add", 302)
		beego.Alert("store  2")
		//return
	}
	if name != "" {
		beego.Alert("store 3")
		this.Data["store"] = store

	}
	this.TplName = "Stores.html"

}

//跳转到ADD页面
//func (this *CoreController) HandleStore() {
//	this.Redirect("/stores/add", 302)
//}
//展示添加店铺页面
func(this *CoreController)ShowAddStore(){
	id := this.GetSession("id")
	beego.Alert("ADD Store:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
		this.TplName="StoresAdd.html"
}

//添加商店的操作
func (this *CoreController) AddStore() {
	beego.Alert("ADD store 1")
	//this.TplName = "StoresAdd.html"
	//给查询对象赋值
	id := this.GetSession("id")
	beego.Alert("ADD Store:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	//跟店铺建立一对一关系
	this.Data["user"] = user
	var store models.ShopifyStore
	//store.User.Id=user.Id

	store.Email = user.Email
	o.Read(&store, "Email")
	name := this.GetString("shop_url")
	//暂时定为password
	password := this.GetString("Password")
	apikey := this.GetString("APIKEY")
	beego.Alert("get 2")
	beego.Alert("name:", name)

	beego.Alert("apikey:", apikey)
	beego.Alert("password:", password)
	beego.Alert("exist 3")
	if name==""||apikey==""||password==""{
		resp["errno"]=1
		resp["errmsg"]="Please complete the required store information"
		return
	}
	qs:=o.QueryTable("ShopifyStore").Filter("Name",name).Exist()
	//如果已经存在，就更新
	if qs==true{
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := (time.Now().In(loc))
		store.CreateTime = t
		store.Name = name
		store.Secret = password
		store.ApiKey = apikey
		store.Id=int64(user.Id)
		store.User=&user
		_, err := o.Update(&store)
		if err != nil {
			beego.Error(err)
			return
		}
		_,err1:=o.Update(&user, "StoreName")
		if err1!=nil{
			beego.Error(err1)
			return
		}
	}else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := (time.Now().In(loc))
		store.CreateTime = t
		store.Name = name
		store.Secret = password
		store.ApiKey = apikey
		store.Id=int64(user.Id)
		store.User=&user
		_, err := o.Insert(&store)
		if err != nil {
			beego.Error(err)
			return
		}
		_,err1:=o.Update(&user, "StoreName")
		if err1!=nil{

			beego.Error(err1)
			return
		}
	}
	beego.Alert("finish")
	this.Redirect("/stores", 302)
	resp["errno"]=0
}

//展示一键获取订单的页面
func (this *CoreController) ShowOrder() {
	id := this.GetSession("id")
	beego.Alert("Order:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	//跟店铺建立一对一关系
	this.Data["user"] = user
	var store models.ShopifyStore

	o.QueryTable("ShopifyStore").Filter("User__Id",id).One(&store)
	beego.Alert(store.Name)
	//resp := make(map[string]interface{})
	//defer RespFunc(&this.Controller,resp)
	//获取总页码
	name:=store.Name
	key:=store.ApiKey
	pwd:=store.Secret

	ordercount := util.GetOrderCountTest(name,key,pwd)
	//productcount:=util.GetProductCount()
	this.Data["ordercount"] = ordercount
	//this.Data["productcount"]=productcount
	pageSize := 50
	pageCount := int(math.Ceil(float64(ordercount) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {

	}
	beego.Alert("aaaaaaaaaaaaaaaaaaaaaa", pageIndex)
	pages := PageEdit(pageCount, pageIndex)
	this.Data["pages"] = pages
	this.Data["pageCount"] = pageCount
	//获取上一页，下一页的值
	var prePage, nextPage int
	//设置范围
	if pageIndex-1 <= 0 {
		prePage = 1
	} else {
		prePage = pageIndex - 1
	}
	if pageIndex+1 >= pageCount {
		nextPage = pageCount
	} else {
		nextPage = pageIndex + 1
	}
	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage

	var orders1 []util.Orders
	beego.Alert("get order begin")
	O := util.GetOrderStructTest(name,key,pwd,pageIndex)
	//O是五十个orderinfo的集合
	for i := 0; i < len(O.Orders); i++ {
		orders1 = append(orders1, O.Orders[i])
	}
	beego.Alert("get order finish")

	var lineitem []util.LineItems

	//var source models.SourcingDemand
	var orders []map[string]interface{}
	for a := 0; a < len(orders1); a++ {
		t := make(map[string]interface{})
		createtime := orders1[a].CreatedAt
		totalprice := orders1[a].TotalPrice
		shipping := orders1[a].ShippingAddress
		name := orders1[a].Name
		note := orders1[a].Note
		order_id := orders1[a].ID
		t["Name"] = name
		t["CreatedAt"] = createtime
		t["ShippingAddress"] = shipping
		t["TotalPrice"] = totalprice
		t["order_id"] = order_id
		t["Ordernote"] = note

		qs := o.QueryTable("SourcingDemand").Filter("Orderid", order_id).Filter("Ordernumber", name).Exist()
		if qs {
			t["Nameexit"] = true
		} else {
			t["Nameexit"] = false
		}

		var line []util.LineItems
		var shipline util.ShippingLines

		for b := 0; b < len(orders1[a].LineItems); b++ {
			line = append(line, orders1[a].LineItems[b])
			shipline.Title = orders1[a].ShippingLines[0].Title

		}
		t["ShippingLines"] = shipline.Title
		t["line_items"] = line
		orders = append(orders, t)
	}

	for _, v1 := range orders {
		var second []map[string]interface{}
		for _, v2 := range v1["line_items"].([]util.LineItems) {
			t := make(map[string]interface{})
			t["Sku"] = v2.Sku
			t["Quantity"] = v2.Quantity
			t["Title"] = v2.Title
			t["Price"] = v2.Price
			//product_id
			t["Pid"] = v2.ProductID
			qs := o.QueryTable("SourcingDemand").Filter("Sku", v2.Sku).Exist()
			//beego.Alert(qs)
			if qs {
				t["Skuexit"] = true
			} else {
				t["Skuexit"] = false
			}

			//店名
			t["store"] = v2.OriginLocation.Name
			// o.QueryTable("SourcingDemand").Filter("User__Id", user.Id).Filter("Sku",v2.Sku).One(&source)
			//this.Data["source"]=source

			//t["Vid"] = v2.VariantID
			//variant_id := v2.VariantID
			//image := util.GetSrc(variant_id,v2.ProductID)
			//beego.Alert("image:",v2.ProductID,variant_id,image)
			//product := util.GetProductStruct(int(v2.ProductID))
			//image:=product.Products[0].Images[0].Src
			//t["image"] = image
			second = append(second, t)
		}
		//将二级容器放到总容器
		//TODO  用line_items可能错
		v1["line_items"] = second
	}
	this.Data["orders"] = orders
	beego.Alert("line begin ")
	this.Data["item"] = lineitem
	beego.Alert("line finish")
	//如果不行就定义行容器，拼接order信息再发送

	this.TplName = "orders.html"
	beego.Alert("show order test")
}

//点击sync 刷新订单页
//func(this *CoreController)HandleOrder(){
//	this.Redirect("/order/drop",302)
//}

//从订单页跳转到sourcing agent订单的采购代理
func (this *CoreController) ShowSourceAgent() {
	//判定登录
	beego.Alert("ShowSource")
	id := this.GetSession("id")
	beego.Alert("Sourcing add:", id)
	if id == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = id.(int)
	o.Read(&user, "Id")
	//跟店铺建立一对一关系
	this.Data["user"] = user

	this.Data["user"] = user
	var store models.ShopifyStore

	o.QueryTable("ShopifyStore").Filter("User__Id",id).One(&store)
	beego.Alert(store.Name)
	//resp := make(map[string]interface{})
	//defer RespFunc(&this.Controller,resp)
	//获取总页码
	name:=store.Name
	key:=store.ApiKey
	pwd:=store.Secret
	productid, err := this.GetInt("product_Id")
	beego.Alert("sourcing agent get product_id", productid)

	if err != nil {
		beego.Error(err)
		this.Redirect("/order/drop", 302)
		return
	}
	sku := this.GetString("sku")
	store2 := this.GetString("store")
	order_id := this.GetString("order_id")
	beego.Alert(store2, order_id)
	//根据productid获取图片
	product := util.GetProductStruct(name,key,pwd,productid)
	beego.Alert("get product struct finish")

	title := product.Products[0].Title
	image := product.Products[0].Images[0].Src

	this.Data["Title"] = title
	this.Data["Sku"] = sku
	this.Data["image"] = image
	//this.Data["pid"]=productid
	beego.Alert("image:", image)

	beego.Alert(title)
	this.TplName = "SourcingAgent.html"
	beego.Alert("show SourcingAgent 1")

}

//提交采购需求
func (this *CoreController) HandleSourceAgent() {

	userid := this.GetSession("id")

	if userid == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = userid.(int)
	o.Read(&user, "Id")
	beego.Alert("userid", userid)
	var store models.ShopifyStore

	o.QueryTable("ShopifyStore").Filter("User__Id",userid).One(&store)
	beego.Alert(store.Name)
	//resp := make(map[string]interface{})
	//defer RespFunc(&this.Controller,resp)
	//获取总页码
	name:=store.Name
	key:=store.ApiKey
	pwd:=store.Secret
	url := this.GetString("url")
	ordernumber := this.GetString("number")
	beego.Alert("ordernumber", ordernumber)
	//上传图片
	sellingprice, err := this.GetFloat("j")
	order_id, err := this.GetInt64("order_id")
	beego.Alert(order_id, sellingprice)
	productPrice1 := this.GetString("price")
	productPrice, err := strconv.ParseFloat(productPrice1, 64)
	shippingPrice1 := this.GetString("shipping_price")
	shippingPrice, err := strconv.ParseFloat(shippingPrice1, 64)
	delivery := this.GetString("user_delivery")
	quantity := this.GetString("num")
	description := this.GetString("description")
	general := this.GetString("sex")
	pid := this.GetString("product_Id")
	title := this.GetString("product_name")
	sku := this.GetString("sku")
	stores := this.GetString("store")
	productid, err := strconv.Atoi(pid)
	if err != nil {
		beego.Error(err)
	}
	oid := this.GetString("order_id")
	ship := util.GetShippingAdress(oid)
	email := ship.Orders[0].Email
	adress := ship.Orders[0].ShippingAddress
	line := ship.Orders[0].ShippingLines[0].Title
	beego.Alert("shipping_line", line)
	beego.Alert("shipping_line type", reflect.TypeOf(line))
	//获取该订单的地址和时效
	product := util.GetProductStruct(name,key,pwd,productid)
	sourcImage := product.Products[0].Images[0].Src
	q, err := strconv.Atoi(quantity)
	if err != nil {
		beego.Error(err)
	}
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//普货和特货
	//cargo:=this.GetString("cargo")
	var sourcingAgent models.SourcingDemand
	var skulist models.ErpSku
	//var erpsku models.ErpSku
	if general == "Generalcargo" {
		sourcingAgent.General = true
	} else {
		sourcingAgent.General = false
	}
	skulist.Sku = sku
	o.Read(&skulist, "Sku")

	//该id等同product_id
	//sourcingAgent.Id = id
	o.Read(&sourcingAgent, "Id")
	//o.QueryTable("Skulist").Filter()
	skulist.Title = title
	skulist.Image = sourcImage
	//skulist.ErpSku.Id=skulist.Id

	_, err2 := o.Insert(&skulist)
	if err2 != nil {
		beego.Error(err2)
	}
	sourcingAgent.SellingPrice = sellingprice
	sourcingAgent.Ordernumber = ordernumber
	sourcingAgent.Orderid = order_id
	sourcingAgent.Title = title
	sourcingAgent.Link = url
	sourcingAgent.Quantity = q
	sourcingAgent.Target_price = productPrice
	sourcingAgent.Shipping_price = shippingPrice
	sourcingAgent.Estimated_Delivery = delivery
	sourcingAgent.Description = description
	sourcingAgent.Sku = sku
	sourcingAgent.User = &user
	sourcingAgent.Productid = pid
	sourcingAgent.Store = stores
	sourcingAgent.SourcImage = sourcImage
	sourcingAgent.Status = 1
	sourcingAgent.Email = email
	//地址
	sourcingAgent.FirstName = adress.FirstName
	sourcingAgent.LastName = adress.LastName
	sourcingAgent.Address1 = adress.Address1
	sourcingAgent.Phone = adress.Phone
	sourcingAgent.Country = adress.Country
	sourcingAgent.Province = adress.Province
	sourcingAgent.City = adress.City
	sourcingAgent.Zip = adress.Zip
	//sourcingAgent.Company = adress.Company.(string)
	sourcingAgent.Address2 = adress.Address2
	sourcingAgent.Name = adress.Name
	sourcingAgent.ShipLine = line

	//添加创建时间
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := (time.Now().In(loc))
	beego.Alert(t)
	sourcingAgent.Created_at = t
	//普货特货
	//sourcingAgent.General=cargo
	_, err1 := o.Insert(&sourcingAgent)
	if err1 != nil {
		beego.Error(err1)
	}

	//this.Redirect("/order/drop", 302)
	resp["errno"] = 1
	beego.Alert("handlesourcing finish")
}

//采购单上传图片

//func (this *CoreController) SourcUploadImage() []string {
//	//pid := this.GetString("product_Id")
//
//	resp := make(map[string]interface{})
//	defer RespFunc(&this.Controller, resp)
//	image, head, err := this.GetFile("file")
//	//获取图片
//	//返回值 文件二进制流  文件头    错误信息
//	if err != nil {
//		beego.Error("图片上传失败")
//		resp["errno"] = 2
//		resp["errmsg"] = "Picture upload failed"
//		//return
//		//this.TplName = "SourcingAgent.html"
//	}
//
//	defer image.Close()
//	//校验文件大小
//	if head.Size > 5000000 {
//		beego.Error("图片数据过大")
//		resp["errno"] = 3
//		resp["errmsg"] = "Do not exceed 5MB in image size"
//		//return
//		//this.TplName = "Account.html"阿萨德AAA
//	}
//	//校验格式 获取文件后缀
//	ext := path.Ext(head.Filename)
//	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
//		beego.Error("上传文件格式错误")
//		resp["errno"] = 4
//		resp["errmsg"] = `Please upload the image file in the format of "JPG,jpeg,PNG"`
//		//return
//		//this.TplName = "Account.html"
//	}
//
//	fileBuffer := make([]byte, head.Size)
//	//把文件数据读入到fileBuffer中
//	image.Read(fileBuffer)
//	//beego.Alert(fileBuffer)
//	//this.Data["file"]=fileBuffer
//	//获取client对象
//	beego.Alert("open fdfs")
//	client, err := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
//	if err != nil {
//		beego.Error(err)
//	}
//	////上传 得到fdfsresponse.RemoteFileId就是所需字段
//	fdfsresponse, _ := client.UploadByBuffer(fileBuffer, ext[1:])
//	var imageurl []string
//	imageurl = append(imageurl, fdfsresponse.RemoteFileId)
//	//imageurl := fdfsresponse.RemoteFileId
//	beego.Alert(imageurl)
//	this.Data["data"] = imageurl
//	beego.Alert("6")
//	resp["errno"] = 5
//	resp["data"] = imageurl
//	//this.Redirect("/sourcing",302)
//	//测试看是否已经关联到user表中
//	beego.Alert("7")
//	return imageurl
//}

//添加图片闭包

//func (this *CoreController) Handle(imageurl string) func(imageurl string) []string {
//	beego.Alert(imageurl)
//	image := func(imageurl string) []string {
//		//imageurl := fdfsresponse.RemoteFileId
//		var images []string
//		beego.Alert(imageurl)
//		//this.Data["data"] = imageurl
//		beego.Alert("6")
//		images = append(images, imageurl)
//		return images
//	}
//	return image
//}

//func (this *CoreController) Upload() {
//	resp := make(map[string]interface{})
//	defer RespFunc(&this.Controller, resp)
//	image, head, err := this.GetFile("file")
//	beego.Alert("aaaaaaa get image")
//	//返回值 文件二进制流  文件头    错误信息
//	if err != nil {
//		beego.Error("图片上传失败")
//		resp["errno"] = 2
//		resp["errmsg"] = "Picture upload failed"
//		//return
//		//this.TplName = "SourcingAgent.html"
//	}
//	defer image.Close()
//	//校验文件大小
//	if head.Size > 5000000 {
//		beego.Error("图片数据过大")
//		resp["errno"] = 3
//		resp["errmsg"] = "Do not exceed 5MB in image size"
//		//return
//		//this.TplName = "Account.html"
//	}
//	//校验格式 获取文件后缀
//	ext := path.Ext(head.Filename)
//	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
//		beego.Error("上传文件格式错误")
//		resp["errno"] = 4
//		resp["errmsg"] = `Please upload the image file in the format of "JPG,jpeg,PNG"`
//		//return
//		//this.TplName = "Account.html"
//	}
//	fileBuffer := make([]byte, head.Size)
//	//把文件数据读入到fileBuffer中
//	image.Read(fileBuffer)
//	//beego.Alert(fileBuffer)
//	//this.Data["file"]=fileBuffer
//	//获取client对象
//	beego.Alert("open fdfs")
//	client, err := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
//	if err != nil {
//		beego.Error(err)
//	}
//	////上传 得到fdfsresponse.RemoteFileId就是所需字段
//	fdfsresponse, _ := client.UploadByBuffer(fileBuffer, ext[1:])
//	//var imageurl []string
//	imageurl := fdfsresponse.RemoteFileId
//
//	images := this.Handle(imageurl)
//	resp["data"] = images
//	resp["errno"] = 5
//}

//采购单删除图片
//func (this *CoreController) SourcDeleteImage() {
//	//传code=1删除该图片
//	beego.Alert("delete image")
//	this.Data["code"] = 1
//
//}

//展示所有采购需求的页面

func (this *CoreController) ShowSourcingtest() {
	userid := this.GetSession("id")

	if userid == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = userid.(int)
	o.Read(&user, "Id")
	this.Data["user"] = user
	beego.Alert("ShowSourcing", userid)
	//直接根据user的与一对多查找所有与该用户相关的采购单，按照create_time去排序
	var sourcedemand []models.SourcingDemand
	o.QueryTable("SourcingDemand").Filter("User__Id", userid).OrderBy("-Created_at").All(&sourcedemand)
	this.Data["sourcedemand"] = sourcedemand
	sourcingcount, err := o.QueryTable("SourcingDemand").Filter("User__Id", userid).Count()
	if err != nil {
		beego.Error(err)
	}

	beego.Alert(sourcingcount)
	this.Data["sourcingcount"] = sourcingcount
	pageSize := 50
	pageCount := int(math.Ceil(float64(sourcingcount) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
	}
	pages := PageEdit(pageCount, pageIndex)
	this.Data["pages"] = pages
	this.Data["pageCount"] = pageCount
	//获取上一页，下一页的值
	var prePage, nextPage int
	//设置范围
	if pageIndex-1 <= 0 {
		prePage = 1
	} else {
		prePage = pageIndex - 1
	}
	if pageIndex+1 >= pageCount {
		nextPage = pageCount
	} else {
		nextPage = pageIndex + 1
	}
	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage

	this.TplName = "Sourcing.html"
	beego.Alert("ShowSourcing finish")

}

//展示所有采购需求的页面”/sourcing“
//func (this *CoreController) ShowSourcing() {
//	email := this.GetSession("email")
//	if email == "" {
//		beego.Error("用户未登录")
//		this.TplName = "login.html"
//		return
//	}
//	var user models.User
//	var source_demand []models.SourcingDemand
//	user.Email = email.(string)
//	userid := user.Id
//
//	//一个客户对应自己的多个采购需求，一对多，用户是主表，采购需求是子表，自动产生user_id字段，根据user_id字段查找数据
//	/*一：级联查询：var posts []*Post ⇥ O.QueryTable("post").Filter("User__Id", 1).RelatedSel().All(&posts) ⇥ for _, v := range posts {}，v就是post所有数据。可以直接得到：v.User.Name
//	二：reverse查询：var user User ⇥ err := O.QueryTable("user").Filter("Post__Title", "paper1").Limit(1).One(&user)*/
//	o := orm.NewOrm()
//	//根据user查source_Demand表
//	o.QueryTable("Sourcing_Demand").Filter("User__Id", userid).RelatedSel().All(&source_demand)
//
//	//遍历
//	//创建总容器
//	var source []map[string]interface{}
//	var sourcing models.SourcingDemand
//
//	for _, v := range source_demand {
//		//创建行容器
//		temp := make(map[string]interface{})
//		//先判断订单的状态
//
//		if v.Status == false {
//			temp["Status"] = "Pending review"
//			//关闭支付按钮
//			temp["pay"] = false
//		}
//		if v.Status == true {
//			temp["Status"] = "pending payment"
//			//开启支付按钮
//			temp["pay"] = true
//		}
//		temp["Id"] = v.Id
//		temp["title"] = v.Title
//		temp["targetprice"] = v.Target_price
//		temp["delivery"] = v.Estimated_Delivery
//		//获得图片切片
//		image := util.String2slice(v.Images)
//		//只展示一张图片
//		temp["image"] = image[0]
//		source = append(source, temp)
//		//插入数据到Sourcing表
//		sourcing.Id = v.Id
//		sourcing.Title = v.Title
//		sourcing.Images = v.Images
//		sourcing.Target_price = v.Target_price
//		sourcing.Estimated_Delivery = v.Estimated_Delivery
//		sourcing.Status = v.Status
//		sourcing.Link = v.Link
//		sourcing.Destination = v.Destination
//		sourcing.General = v.General
//		o.Insert(&sourcing)
//	}
//
//	this.Data["source"] = source
//	this.TplName = "Sourcing.html"
//}

//直接提采购需求的页面

//Sourcing页面的操作，Action暂定为删除和支付（加入购物车），在状态更新为Pending payment才可以进行支付

//删除Sourcing  前端添加删除按钮和重定向 暂定url中添加Sourcing_id跟sku？
func (this *CoreController) DeleteSourcing() {
	beego.Alert("DeleteSourcing")
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	id, err := this.GetInt64("sourcing_id")
	if err != nil {
		beego.Error("获取sourcing id失败")
		this.Redirect("/sourcing", 302)
		return
	}
	beego.Alert(id)
	//处理数据
	o := orm.NewOrm()
	var source models.SourcingDemand
	source.Id = id
	//删除数据
	if source.Status!=3{
		o.Delete(&source, "Id")
	}else{
		resp["errno"]=1
		resp["errmsg"]="This order cannot be deleted"
	}

	//this.Redirect("/sourcing",302)
}

//搜索采购单，暂定由SKU或者ordernumber来搜索  暂不分页
func (this *CoreController) SourcingSearch() {
	beego.Alert("SourcingSearch")
	userid := this.GetSession("id")

	if userid == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = userid.(int)
	o.Read(&user, "Id")
	this.Data["user"] = user

	var source []models.SourcingDemand

	source_status := this.GetString("source_status")
	searchkey := this.GetString("search_key")
	searchvalue := this.GetString("search_value")
	if searchvalue == "" {
		this.Redirect("/sourcing", 302)
	}
	beego.Alert("sourceSearchby", source_status, searchkey, searchvalue)
	if source_status == "all" {
		if searchvalue == "" {
			o.QueryTable("SourcingDemand").All(&source)
			this.Data["sourcedemand"] = source
			beego.Alert("searchbyvalue1")
			this.TplName = "Sourcing.html"
		} else {
			switch searchkey {
			case "order_number":
				o.QueryTable("SourcingDemand").Filter("Ordernumber", searchvalue).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue2")
				this.TplName = "Sourcing.html"
			case "sku":
				o.QueryTable("SourcingDemand").Filter("Sku", searchvalue).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue3")
				this.TplName = "Sourcing.html"
			}
		}
	} else if source_status == "1" {
		if searchvalue == "" {
			o.QueryTable("SourcingDemand").Filter("Status", 2).All(&source)
			this.Data["sourcedemand"] = source
			beego.Alert("searchbyvalue4")
			this.TplName = "Sourcing.html"
		} else {
			switch searchkey {
			case "order_number":
				o.QueryTable("SourcingDemand").Filter("Ordernumber", searchvalue).Filter("Status", 2).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue5")
				this.TplName = "Sourcing.html"
			case "sku":
				o.QueryTable("SourcingDemand").Filter("Sku", searchvalue).Filter("Status", 2).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue6")
				this.TplName = "Sourcing.html"
			}
		}
	} else if source_status == "2" {
		if searchvalue == "" {
			o.QueryTable("SourcingDemand").Filter("Status", 2).All(&source)
			this.Data["sourcedemand"] = source
			beego.Alert("searchbyvalue7")
			this.TplName = "Sourcing.html"
		} else {
			switch searchkey {
			case "order_number":
				o.QueryTable("SourcingDemand").Filter("Ordernumber", searchvalue).Filter("Status", 2).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue8")
				this.TplName = "Sourcing.html"
			case "sku":
				o.QueryTable("SourcingDemand").Filter("Sku", searchvalue).Filter("Status", 2).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue9")
				this.TplName = "Sourcing.html"
			}
		}
	}else if source_status == "3" {
		if searchvalue == "" {
			o.QueryTable("SourcingDemand").Filter("Status", 3).All(&source)
			this.Data["sourcedemand"] = source
			beego.Alert("searchbyvalue10")
			this.TplName = "Sourcing.html"
		} else {
			switch searchkey {
			case "order_number":
				o.QueryTable("SourcingDemand").Filter("Ordernumber", searchvalue).Filter("Status", 3).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue11")
				this.TplName = "Sourcing.html"
			case "sku":
				o.QueryTable("SourcingDemand").Filter("Sku", searchvalue).Filter("Status", 3).All(&source)
				this.Data["sourcedemand"] = source
				beego.Alert("searchbyvalue12")
				this.TplName = "Sourcing.html"
			}
		}
	}
}

//独立于beego框架的
//分页
func PageEdit(pageCount int, pageIndex int) []int {
	//不足五页
	var pages []int
	if pageCount < 5 {
		for i := 1; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else if pageIndex <= 3 {
		for i := 1; i <= 5; i++ {
			pages = append(pages, i)
		}
	} else if pageIndex >= pageCount-2 {
		for i := pageCount - 4; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else {
		for i := pageIndex - 2; i <= pageIndex+2; i++ {
			pages = append(pages, i)
		}
	}
	return pages
}

//排序
//按时间排序，从新到旧
//编辑，更改sourcingAgent 页面
func (this *CoreController) ShowUpdateSourcingAgent() {
	//从url拿id，根据id查表
	userid := this.GetSession("id")

	if userid == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Id = userid.(int)
	o.Read(&user, "Id")
	beego.Alert("userid", userid)
	this.Data["user"] = user
	var source models.SourcingDemand
	id := this.GetString("sourcing_id")
	sid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	source.Id = sid
	o.Read(&source, "Id")
	//o.QueryTable("SourcingDemand").Filter("Id",sid).One(&source)
	beego.Alert(source)
	//title:=source.Title
	//sku:=source.Sku
	//image:=source.SourcImage
	this.Data["source"] = source
	this.TplName = "SourcingAgentUpdate.html"
	beego.Alert("show update finish")
}

//更新sourcingAgent的操作
func (this *CoreController) UpdateSourcingAgent() {
	userid := this.GetSession("id")
	if userid == nil {
		beego.Error("用户未登录")
		this.TplName = "login.html"
		return
	}
	id := this.GetString("sourcing_id")
	beego.Alert(id)
	o := orm.NewOrm()
	var user models.User
	user.Id = userid.(int)
	var sourcedemand models.SourcingDemand
	sid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	sourcedemand.Id = sid
	o.Read(&sourcedemand, "Id")
	url := this.GetString("url")

	productPrice1 := this.GetString("price")
	productPrice, err := strconv.ParseFloat(productPrice1, 64)
	shippingPrice1 := this.GetString("shipping_price")
	shippingPrice, err := strconv.ParseFloat(shippingPrice1, 64)
	delivery := this.GetString("user_delivery")
	quantity := this.GetString("num")
	description := this.GetString("Description")
	general := this.GetString("sex")
	beego.Alert("new url", url, productPrice, shippingPrice, delivery, quantity, general)
	//title := this.GetString("product_name")
	//sku := this.GetString("sku")

	q, err := strconv.Atoi(quantity)
	if err != nil {
		beego.Error(err)
	}
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//普货和特货
	//cargo:=this.GetString("cargo")

	if general == "Generalcargo" {
		sourcedemand.General = true
	} else {
		sourcedemand.General = false
	}

	//sourcedemand.Title = title
	sourcedemand.Link = url
	sourcedemand.Quantity = q
	sourcedemand.Target_price = productPrice
	sourcedemand.Shipping_price = shippingPrice
	sourcedemand.Estimated_Delivery = delivery
	sourcedemand.Description = description
	sourcedemand.User = &user

	a, err := o.Update(&sourcedemand)
	if err != nil {
		beego.Alert(a)
		beego.Error(err)
	}
	resp["errno"] = 1
	beego.Alert("update finish")
}

func (this *CoreController) PayHandleTest() {
	beego.Alert("1111")
	stripe.Key = "sk_test_TyAe7P109fmUiQBHTf0xGn9G00EaocYpqU"

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(2005),
		Currency: stripe.String(string(stripe.CurrencyHKD)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}
	pi, _ := paymentintent.New(params)
	this.Data["result"] = pi
	beego.Alert("pay")

}
func (this *CoreController) ShowSourcingDetail() {
	beego.Alert("ShowSourcingDetail")
	id, err := this.GetInt64("sourcing_id")
	if err != nil {
		beego.Error("获取sourcing id失败")
		this.Redirect("/sourcing", 302)
		return
	}
	beego.Alert(id)
	//处理数据
	o := orm.NewOrm()
	var order models.Order
	order.SourceId = id
	o.Read(&order, "SourceId")
	beego.Alert(order)
	this.Data["order"] = order
	this.TplName = "SourcingDetail.html"

}

//type List struct {
//	sourceidList []string `json:"sourceidList"`
//
//}
type CheckoutData struct {
	ClientSecret string
}

func (this *CoreController) HandlePay() {
	beego.Alert("HandlePay")
	sourceidlist := this.GetStrings("sourceidList[]")
	beego.Alert(sourceidlist)
	o := orm.NewOrm()
	var order models.Order
	var prices float64
	var finalprice int64
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	for i := 0; i < len(sourceidlist); i++ {
		o.QueryTable("Order").Filter("SourceId", sourceidlist[i]).One(&order)
		if order.TotalPrice == 0 {
			resp["errno"] = 1
			resp["errmsg"] = "Please select the order to be paid"
			return
		}
		if order.Status==3{
			resp["errno"]=2
			resp["errmsg"]="The order has been paid successfully"
			return
		}
		price := order.TotalPrice
		prices += price
		//得到选择订单的价格
		beego.Alert(prices)
		finalprice = int64(prices * 100)

	}
	resp["errno"]=10
	pi := util.StripeTest(finalprice)
	beego.Alert(pi.ClientSecret)
	resp["clientSecret"] = pi.ClientSecret

}

func (this *CoreController) Pay() {
	sourceidlist := this.GetStrings("sourceidList[]")
	beego.Alert(sourceidlist)
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	o := orm.NewOrm()
	var order models.Order
	var source models.SourcingDemand
	for i := 0; i < len(sourceidlist); i++ {
		o.QueryTable("Order").Filter("SourceId", sourceidlist[i]).One(&order)
		order.Status = 3
		o.Update(&order)
	}
	for i := 0; i < len(sourceidlist); i++ {
		o.QueryTable("SourcingDemand").Filter("Id", sourceidlist[i]).One(&source)
		source.Status = 3
		o.Update(&source)
	}
	this.Redirect("/sourcing",302)
}

