package controllers

import (
	"dropshe/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
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
		temp["user"] = v.User.ShopifyStore.Name
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
	o.QueryTable("ImageSrc").Filter("ProductMatch__Id", product_id).RelatedSel("ImageSrc").All(&image)
	sourcing["image"] = image.FdfsSrc

	this.Data["sourcing"] = sourcing
	this.TplName = "sku.html"

}

//下拉框类型
type productIdSelect struct {
	Name       int64
	IsSelected bool
	Value      int64
}

//匹配
func (this *ErpController) AuditSku() {
	//ERPsku
	o := orm.NewOrm()
	esku := this.GetString("SKU")
	//根據id查Erpsku

	o.QueryTable("ErpSku").Filter("LineItems__Id", product_id).RelatedSel("ErpSku").One(&esku)
	erpsku.Sku = esku
	o.Insert(&esku)
	//判定条件，改变sourcing的状态，可以进行交易
	var erps []models.ErpSku

	o.QueryTable("ErpSku").Filter("Sourcing__Id", id).RelatedSel("ErpSku").All(&erps)
	for a := 0; a < len(erps); a++ {
		qs := o.QueryTable("ErpSku").Filter("Sku", erps[a].Sku).Exist()
		if qs == true {
			source.Status = true
			o.Update(&source)
		} else {
			source.Status = false
			return
		}

	}
	this.Redirect("/erpsource/sku", 302)

}

//待处理订单的处理（审核、详情？）
func (this *ErpController) HandlerOrder() {
	//根据订单跳转到信息审核的页面

}
