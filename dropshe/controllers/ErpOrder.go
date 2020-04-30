package controllers

import (
	"dropshe/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

//ERP 展示待处理采购单
func (this *ErpController) ShowErpSource() {
	//获取订单
	//name := this.GetSession("name")
	//if name == nil {
	//	beego.Error("用户未登录")
	//	this.TplName = "erplogin.html"
	//	return
	//}
	//获取所有订单  直接查用户下单的采购需求表 Sourcing 表
	o := orm.NewOrm()
	var erpsource []models.SourcingDemand
	//o.Read(&erpsource)
	//按创建顺序排列？
	o.QueryTable("SourcingDemand").Filter("Status", 1).OrderBy("-Created_at").All(&erpsource)
	count, err := o.QueryTable("SourcingDemand").Count()

	//beego.Alert(erpsource)
	this.Data["source"] = erpsource
	pageSize := 200
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		beego.Error(err)
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

	//直接传递整个Sourcing
	this.Data["source"] = erpsource
	this.TplName = "Erpsourcing.html"

}

//Erpsourcing的搜索功能
func (this *ErpController) SearchBySku() {
	beego.Alert("SearchBySku")
	sku := this.GetString("seachbysku")
	beego.Alert(sku)
	o := orm.NewOrm()
	var source []models.SourcingDemand
	o.QueryTable("SourcingDemand").Filter("Sku", sku).OrderBy("-Created_at").All(&source)
	beego.Alert(source)
	this.Data["source"] = source
	this.TplName = "Erpsourcing.html"

}

//按店铺搜索
func (this *ErpController) SearchByStore() {

	store := this.GetString("seachbystore")
	o := orm.NewOrm()
	var source []models.SourcingDemand
	o.QueryTable("SourcingDemand").Filter("Store", store).OrderBy("-Created_at").All(&source)
	this.Data["source"] = source
	this.TplName = "Erpsourcing.html"

}

//sku匹配頁面
func (this *ErpController) ShowErpAudit() {
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	//o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")
	//this.Data["user"]=user
	id := this.GetString("sourcing_id")
	beego.Alert(id)
	var source models.SourcingDemand
	sourcing_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	source.Id = sourcing_id
	o := orm.NewOrm()
	o.Read(&source, "Id")
	beego.Alert(source)
	this.Data["source"] = source
	this.TplName = "ErpsourcingMate.html"

}

//SKU匹配
func (this *ErpController) HandleErpEdit() {
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	//o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")

	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	beego.Alert("HandleErpEdit")
	id := this.GetString("sourcing_id")
	beego.Alert(id)
	o := orm.NewOrm()

	//获取输入的匹配sku
	esku := this.GetString("ErpSku")
	beego.Alert("aaaaaaaaaaaaaaaaaaaaaaaaaaaaa", esku)
	if esku == "" {

	}
	sex := this.GetString("sex")
	shippricestring := this.GetString("shipprice")
	shipprice,err:=strconv.ParseFloat(shippricestring,64)
	pricestring := this.GetString("price")
	price,err:=strconv.ParseFloat(pricestring,64)
	link := this.GetString("url")
	shipmethod := this.GetString("shippingmethod")
	beego.Alert(shipmethod)
	var sourcedemand models.SourcingDemand

	var erpsku models.ErpSku
	//var sku models.Skuli
	sourcing_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	sourcedemand.Id = sourcing_id

	o.Read(&sourcedemand, "Id")
	if sourcedemand.Status == 2 {
		resp["errno"] = 3
		resp["errmsg"] = "该订单已经存在"
		return
	}
	qs := o.QueryTable("ErpSku").Filter("Sku", sourcedemand.Sku).Exist()
	if qs == true {
		o.QueryTable("ErpSku").Filter("Sku", sourcedemand.Sku).One(&erpsku)
		erpsku.Erpsku = esku
		o.Insert(&erpsku)
	}

	beego.Alert(sex, shipprice, price, link)
	var order models.Order
	//添加创建时间
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := (time.Now().In(loc))
	order.SellingPrice=sourcedemand.SellingPrice
	order.SourceId = sourcing_id
	order.Title = sourcedemand.Title
	order.Ordernumber = sourcedemand.Ordernumber
	order.Orderid = sourcedemand.Orderid
	order.Sku = sourcedemand.Sku
	order.Link = sourcedemand.Link
	order.Target_price = sourcedemand.Target_price
	order.Shipping_price = sourcedemand.Shipping_price
	order.Estimated_Delivery = sourcedemand.Estimated_Delivery
	order.Quantity = sourcedemand.Quantity
	order.Status = 2
	order.General = sourcedemand.General
	order.Created_at = t
	order.Store = sourcedemand.Store
	order.SourcImage = sourcedemand.SourcImage
	//该单的操作员
	//order.Staff=user.Name

	order.FirstName = sourcedemand.FirstName
	order.LastName = sourcedemand.LastName
	order.Address1 = sourcedemand.Address1
	order.Email = sourcedemand.Email
	order.Phone = sourcedemand.Phone
	order.City = sourcedemand.City
	order.Zip = sourcedemand.Zip
	order.Province = sourcedemand.Province
	order.Country = sourcedemand.Country
	order.ShipLine = sourcedemand.ShipLine
	order.ErpPrice = price
	order.ErpShipPrice = shipprice
	order.ShipMethod = shipmethod
	order.SourceLink = link
	order.Erpsku = esku

	//
	order.TotalPrice=shipprice+float64(sourcedemand.Quantity)*price
	_, err1 := o.Insert(&order)
	if err1 != nil {
		beego.Error(err1)
		resp["errno"] = 1
		resp["errmsg"] = "插入order失败"
		return
	}
	sourcedemand.Status = 2
	//sourcedemand.ErpPrice=price
	//sourcedemand.ErpShipPrice = shipprice
	//sourcedemand.ShipMethod = shipmethod
	//sourcedemand.SourceLink = link
	//sourcedemand.Erpsku = esku
	_, err2 := o.Update(&sourcedemand)
	if err2 != nil {
		beego.Error(err2)
		resp["errno"] = 2
		resp["errmsg"] = "更新采购单状态失败"
		return
	}
	this.Redirect("/erpsourcing", 302)

	beego.Alert("order insert finish")

}
//展示已经匹配的采购单
func (this *ErpController)ShowErpSourcestatus2(){
	//获取订单
	//name := this.GetSession("name")
	//if name == nil {
	//	beego.Error("用户未登录")
	//	this.TplName = "erplogin.html"
	//	return
	//}
	//获取所有订单  直接查用户下单的采购需求表 Sourcing 表
	o := orm.NewOrm()
	var erpsource []models.SourcingDemand
	//o.Read(&erpsource)
	//按创建顺序排列？
	o.QueryTable("SourcingDemand").Filter("Status", 2).OrderBy("-Created_at").All(&erpsource)
	count, err := o.QueryTable("SourcingDemand").Count()

	//beego.Alert(erpsource)
	this.Data["source"] = erpsource
	pageSize := 200
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		beego.Error(err)
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

	//直接传递整个Sourcing
	this.Data["source"] = erpsource
	this.TplName = "Erpsourcing.html"
}


//管理已审核订单
func (this *ErpController) ShowErpOrder_status2() {
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")
	//this.Data["user"]=user
	var order []models.Order
	o.QueryTable("Order").Filter("Status",2).OrderBy("-Created_at").All(&order)
	count, err := o.QueryTable("Order").Count()

	//beego.Alert(erpsource)

	pageSize := 200
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		beego.Error(err)
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

	this.Data["order"] = order
	beego.Alert(order)
	this.TplName = "erporder.html"

}


func (this *ErpController) ShowErpOrder_status3() {
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")
	//this.Data["user"]=user
	var order []models.Order
	o.QueryTable("Order").Filter("Status",3).OrderBy("-Created_at").All(&order)
	count, err := o.QueryTable("Order").Count()

	//beego.Alert(erpsource)

	pageSize := 200
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		beego.Error(err)
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

	this.Data["order"] = order
	beego.Alert(order)
	this.TplName = "erporder.html"

}



//更改订单数据
func (this *ErpController) HandleOrder() {
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	//o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")
	//this.Data["user"]=user
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)

	id,err:=this.GetInt64("order_id")
	if err!=nil{
		beego.Error(err)
	}
	var order models.Order
	order.Id=id
	o:=orm.NewOrm()
	o.Read(&order,"Id")
	price,err:=this.GetFloat("price")
	if price!=0{
		order.ErpPrice=price
	}
	shipprice,err:=this.GetFloat("shipprice")
	if shipprice!=0{
		order.ErpShipPrice=shipprice
	}

	erpsku:=this.GetString("ErpSku")
	if erpsku!=""{
		order.Erpsku=erpsku
	}
	url:=this.GetString("url")
	if url!=""{
		order.SourceLink=url
	}
	shippingmethod:=this.GetString("shippingmethod")
	if shippingmethod!=""{
		order.ShipMethod=shippingmethod
	}

	waybillnumber:=this.GetString("waybillnumber")
	order.WaybillNumber=waybillnumber

	o.Update(&order)

	resp["errno"]=0

}

//erporder 的搜索功能
//根据店铺搜索
func (this *ErpController) OrderSearch() {
	var order []models.Order
	o := orm.NewOrm()
	beego.Alert("OrderSearch")
	order_status := this.GetString("order_status")
	searchkey := this.GetString("search_key")

	searchvalue := this.GetString("search_value")
	beego.Alert("OrderSearchby", order_status,searchkey, searchvalue)
	if order_status == "all" {
		if searchvalue==""{
			o.QueryTable("Order").All(&order)
			this.Data["order"] = order
			this.TplName = "erporder.html"
		}else {
			switch searchkey {
			case "id":
				o.QueryTable("Order").Filter("Id", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "order_number":
				o.QueryTable("Order").Filter("Ordernumber", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "email":
				o.QueryTable("Order").Filter("Email", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "sku":
				o.QueryTable("Order").Filter("Sku", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "erpsku":
				o.QueryTable("Order").Filter("Erpsku", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			}
		}

	} else if order_status == "2" {
		if searchvalue==""{
			o.QueryTable("Order").Filter("Status",2).All(&order)
			this.Data["order"] = order
			this.TplName = "erporder.html"
		}else {
			switch searchkey {
			case "id":
				o.QueryTable("Order").Filter("Id", searchvalue).Filter("Status",2).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "order_number":
				o.QueryTable("Order").Filter("Ordernumber", searchvalue).Filter("Status",2).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "email":
				o.QueryTable("Order").Filter("Email", searchvalue).Filter("Status",2).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "sku":
				o.QueryTable("Order").Filter("Sku", searchvalue).Filter("Status",2).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "erpsku":
				o.QueryTable("Order").Filter("Erpsku", searchvalue).Filter("Status",2).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			}
		}

	}else if order_status=="3"{
		if searchvalue==""{
			o.QueryTable("Order").Filter("Status",3).All(&order)
			this.Data["order"] = order
			this.TplName = "erporder.html"
		}else {
			switch searchkey {
			case "id":
				o.QueryTable("Order").Filter("Id", searchvalue).Filter("Status",3).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "order_number":
				o.QueryTable("Order").Filter("Ordernumber", searchvalue).Filter("Status",3).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "email":
				o.QueryTable("Order").Filter("Email", searchvalue).Filter("Status",3).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "sku":
				o.QueryTable("Order").Filter("Sku", searchvalue).Filter("Status",3).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			case "erpsku":
				o.QueryTable("Order").Filter("Erpsku", searchvalue).All(&order)
				this.Data["order"] = order
				this.TplName = "erporder.html"
			}
		}
	}
}
//根据时间排序
func (this *ErpController) OrderSortByTime() {

}

//编辑订单
func (this *ErpController) ShowErpOrderEdit() {
	beego.Alert("ErpOrderEdit")
	//sessionid:=this.GetSession("Id")
	//if sessionid==nil{
	//	beego.Error("请登录")
	//	return
	//}
	o := orm.NewOrm()
	//var  user models.ErpUser
	//user.ID=sessionid.(int64)
	//o.Read(&user,"Id")
	//获取订单id
	orderid := this.GetString("order_id")
	if orderid == "" {
		return
	}
	beego.Alert(orderid)
	var order models.Order

	id, err := strconv.ParseInt(orderid, 10, 64)
	order.Id = id
	if err != nil {
		beego.Error(err)
	}
	o.Read(&order, "Id")

	this.Data["order"] = order
	this.TplName = "ErpOrderEdit.html"
}


//func (this *ErpController) HandleErpOrderEdit() {
//	//sessionid:=this.GetSession("Id")
//	//if sessionid==nil{
//	//	beego.Error("请登录")
//	//	return
//	//}
//	//o := orm.NewOrm()
//	//var  user models.ErpUser
//	//user.ID=sessionid.(int64)
//	//o.Read(&user,"Id")
//
//	resp := make(map[string]interface{})
//	defer RespFunc(&this.Controller, resp)
//	beego.Alert("HandleErpEdit")
//	id := this.GetString("sourcing_id")
//	beego.Alert(id)
//	o := orm.NewOrm()
//
//	//获取输入的匹配sku
//	esku := this.GetString("ErpSku")
//	beego.Alert("aaaaaaaaaaaaaaaaaaaaaaaaaaaaa", esku)
//	if esku == "" {
//
//	}
//	sex := this.GetString("sex")
//	shipprice := this.GetString("shipprice")
//	price := this.GetString("price")
//	link := this.GetString("url")
//	shipmethod := this.GetString("shippingmethod")
//	beego.Alert(shipmethod)
//	var sourcedemand models.SourcingDemand
//
//	var erpsku models.ErpSku
//	//var sku models.Skuli
//	sourcing_id, err := strconv.ParseInt(id, 10, 64)
//	if err != nil {
//		beego.Error(err)
//	}
//	sourcedemand.Id = sourcing_id
//
//	o.Read(&sourcedemand, "Id")
//	if sourcedemand.Status == 2 {
//		resp["errno"] = 3
//		resp["errmsg"] = "该订单已经存在"
//		return
//	}
//	qs := o.QueryTable("ErpSku").Filter("Sku", sourcedemand.Sku).Exist()
//	if qs == true {
//		o.QueryTable("ErpSku").Filter("Sku", sourcedemand.Sku).One(&erpsku)
//		erpsku.Erpsku = esku
//		o.Insert(&erpsku)
//	}
//}
