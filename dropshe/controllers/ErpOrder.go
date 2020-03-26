package controllers

import (
	"dropshe/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//ERP 展示待处理采购单
func (this *ErpController) ShowErpSource() {
	//获取订单
	name := this.GetSession("name")
	if name == "" {
		beego.Error("用户未登录")
		this.TplName = "erplogin.html"
		return
	}

	//获取所有订单  直接查用户下单的采购需求表 Sourcing 表
	o := orm.NewOrm()
	var erpsource []models.Sourcing
	//o.Read(&erpsource)
	o.QueryTable("Sourcing").OrderBy("Id").All(&erpsource)
	//定义总容器
	var source []map[string]interface{}
	for _, v := range erpsource {
		//定义行容器
		temp := make(map[string]interface{})
		temp["title"] = v.Id
		temp["title"] = v.Title
		temp["statu"] = v.Status
		temp["totolprice"] = v.Target_price
		//temp["user"] = v.User.ShopifyStore.Name
		temp["destination"] = v.Destination
		temp["general"] = v.General
		source = append(source, temp)
	}
	this.Data["source"] = source
	this.TplName = "erporder.html"
}

//sku匹配頁面
func (this *ErpController) ShowErpAudit() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/erpsource", 302)
		return
	}
	var source models.Sourcing
	source.Id = id
	o := orm.NewOrm()
	o.QueryTable("Sourcing").Filter("id", id).One(source)
	//一个采购单对应一个order，应获取该单中所有product
	//TODO 前端修改成下拉框，每个下拉框中显示title，选中该title后显示对应的line_items
	var lines []models.LineItems
	var erpsku models.ErpSku
	o.QueryTable("LineItems").Filter("Sourcing__Id", id).RelatedSel("LineItems").All(&lines)
	//一个Line作为一个行容器，先将内容传过去,
	var sourcing map[string]interface{}

	for i := 0; i < len(lines); i++ {
		//	TODO 先插一次表还是后面再插？
		erpsku.Id = lines[i].Id
		o.Insert(&erpsku)
		pid := []*productIdSelect{
			{-1, true, 0},
			{lines[i].Id, false, lines[i].Id},
		}
		//传数据到下拉框
		this.Data["ProductId"] = pid
	}
	//把表单里的内容赋值到一个 struct
	//this.Ctx.Request.ParseForm()
	//得到该下拉框的值(返回string)
	product_ids := this.Input().Get("productid")
	var lineitem models.LineItems
	//根据title查找对应的Line
	//o.QueryTable("LineItems").Filter("LineItems__Id",product_id).One(lineitem)
	product_id, err := strconv.ParseInt(product_ids, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	lineitem.Id = product_id
	o.Read(&lineitem)
	//sourcing["Id"]=lineitem.Id
	sourcing["Title"] = lineitem.Title
	sourcing["Price"] = lineitem.Price
	sourcing["Quantity"] = lineitem.Quantity
	sourcing["VariantID"] = lineitem.VariantID
	sourcing["Sku"] = lineitem.Sku
	//查圖片
	var image models.ImageSrc
	//TODO image返回的是否集合更合适
	o.QueryTable("ImageSrc").Filter("ProductMatch__Id", product_id).RelatedSel("ImageSrc").One(&image)
	sourcing["image"] = image.FdfsSrc

	this.Data["sourcing"] = sourcing
	this.TplName = "sku.html"


	//ERPsku
	o.QueryTable("ErpSku").Filter("LineItems__Id", product_id).RelatedSel("ErpSku").One(&erpsku)
	if erpsku.Sku!=""{
		this.Data["SKU"]=erpsku.Sku
	}else {
		//获取输入的sku存库
		esku := this.GetString("SKU")
		//根據id查Erpsku
		erpsku.Sku = esku
		o.Insert(&erpsku)
	}

	//判定条件，改变sourcing的状态，可以进行交易
	var erps []models.ErpSku

	o.QueryTable("ErpSku").Filter("Sourcing__Id", id).RelatedSel("ErpSku").All(&erps)
	//判断erps的sku是否存在
	for a := 0; a < len(erps); a++ {
		qs := o.QueryTable("ErpSku").Filter("Sku", erps[a].Sku).Exist()
		if qs == true {
			source.Status = true
			//创建时间
			loc, _ := time.LoadLocation("America/Los_Angeles")
			t:=(time.Now().In(loc))

			source.Updatetime=t
			o.Update(&source)
		} else {
			source.Status = false
			return
		}

	}
	this.Redirect("/erpsource/sku", 302)

}

//下拉框类型
type productIdSelect struct {
	Name       int64
	IsSelected bool
	Value      int64
}



//代发订单（展示订单信息）
func (this *ErpController) HandlerOrder() {
	name := this.GetSession("name")
	if name == "" {
		beego.Error("用户未登录")
		this.TplName = "erplogin.html"
		return
	}
	//展示代发订单的店铺名、店铺订单号，用户名，收件人：名字，联系方式，收件地址，详细地址，邮编
	//创建总容器
	var orders []map[string]interface{}
	//创建行容器
	//var temp map[string]interface{}
	//获取所有代发订单
	var source []models.Sourcing
	//var store models.ShopifyStore
	var user models.User
	var order models.Order

	//根据sourcing查询User,由user一对一直接得到店铺信息
	//由sourcing查order，再由order得到shipingadress
	//
	o := orm.NewOrm()
	o.QueryTable("Sourcing").Filter("Ready", true).OrderBy("Updatetime").All(&source)

	for i := 0; i < len(source); i++ {
		temp := make(map[string]interface{})
		o.QueryTable("User").Filter("Sourcing__Id", source[i].Id).RelatedSel("Sourcing").One(&user)
		temp["time"] = source[i].Updatetime
		//temp["Storename"] = user.ShopifyStore.Name
		temp["user"] = user.FirstName + user.LastName
		//查order
		o.QueryTable("Order").Filter("Order__Id", source[i].Id).RelatedSel("Sourcing").One(&order)
		temp["shopify_number"] = order.Name
		temp["receiver"] = order.Shipping_address.FirstName + order.Shipping_address.LastName
		temp["zip"] = order.Shipping_address.Zip
		temp["phone"] = order.Shipping_address.Phone
		adress := order.Shipping_address.Country + order.Shipping_address.Province + order.Shipping_address.City
		temp["adress"] = adress
		temp["adress_detail"] = order.Shipping_address.Address1 + "," + order.Shipping_address.Address2

		orders = append(orders, temp)
		//
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
	this.Data["orders"] = orders

	this.TplName = "erporder.html"

	//订单详情页
}
