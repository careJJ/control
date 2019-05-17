package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/models"
)

type GoodsController struct {
	beego.Controller
}

func(this*GoodsController)ShowIndex(){
	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}else {
		this.Data["name"] = ""
	}

	//获取类型信息并传递给前段
	//获取一级菜单
	o := orm.NewOrm()
	//接受对象
	var oneClass []models.TpshopCategory
	//查询
	o.QueryTable("TpshopCategory").Filter("Pid",0).All(&oneClass)


	//获取第二级
	var types []map[string]interface{}//定义总容器
	for _,v := range oneClass{
		//行容器
		t := make(map[string]interface{})

		var secondClass []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid",v.Id).All(&secondClass)
		t["t1"] = v  //一级菜单对象
		t["t2"] = secondClass  //二级菜单集合
		//把行容器加载到总容器中
		types = append(types,t)
	}

	//获取第三季菜单
	for _,v1 := range types{
		//循环获取二级菜单
		var erji []map[string]interface{} //定义二级容器
		for _,v2 := range v1["t2"].([]models.TpshopCategory){
			t := make(map[string]interface{})
			var thirdClass []models.TpshopCategory
			//获取三级菜单
			o.QueryTable("TpshopCategory").Filter("Pid",v2.Id).All(&thirdClass)
			t["t22"] = v2  //二级菜单
			t["t23"] = thirdClass   //三级菜单
			erji = append(erji,t)
		}
		//把二级容器放到总容器中
		v1["t3"] = erji
	}

	this.Data["types"] = types
	this.TplName = "index.html"
}

func (this*GoodsController)ShowIndex_sx(){
	o:=orm.NewOrm()
	//获取所有类型
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"]=goodsTypes
	//获取轮播图
	var goodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&goodsBanners)
	this.Data["goodsBanners"]=goodsBanners
	//获取所有促销商品
	var promotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionBanners)
	this.Data["promotions"]=promotionBanners
	//获取首页展示(总容器）
	var goods []map[string]interface{}

	for _,v:=range goodsTypes {
		var textGoods []models.IndexTypeGoodsBanner
		var imageGoods []models.IndexTypeGoodsBanner
		beego.Info(v.Id)
		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType__Id", v.Id).OrderBy("Index")

		//获取文字和图片商品
		qs.Filter("DisplayType", 0).All(&textGoods)
		qs.Filter("DisplayType", 1).All(&imageGoods)

		beego.Info(textGoods)
		beego.Info(imageGoods)
		//定义行容器
		temp := make(map[string]interface{})
		temp["goodsType"] = v
		temp["textGoods"] = textGoods
		temp["imageGoods"] = imageGoods

		goods = append(goods, temp)

	}
	this.Data["goods"]=goods
	this.TplName="index_sx.html"
}

func(this*GoodsController)ShowDetail(){
	//获取数据
	id,err := this.GetInt("Id")
	//校验数据
	if err != nil {
		beego.Error("商品链接错误")
		this.Redirect("/index_sx",302)
		return
	}
	//处理数据
	//根据id获取商品有关数据
	o := orm.NewOrm()
	var goodsSku models.GoodsSKU

	//获取商品详情
	o.QueryTable("GoodsSKU").RelatedSel("Goods","GoodsType").Filter("Id",id).One(&goodsSku)

	//获取同一类型的新品推荐
	var newGoods []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",goodsSku.GoodsType.Name)
	qs.OrderBy("-Time").Limit(2,0).All(&newGoods)
	this.Data["newGoods"] = newGoods
	//传递数据
	this.Data["goodsSku"] = goodsSku
	this.TplName = "detail.html"
}

//商品列表
func (this*GoodsController)ShowList(){
	//获取数据
	id,err:=this.GetInt("id")
	//教研数据
	if err!=nil{
		beego.Error("类型不存在")
		this.Redirect("/index_sx",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var goods []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType_Id",id).All(&goods)

	//返回数据
	this.Data["goods"]=goods
	this.TplName="list.html"
//1.不足五页  有几页显示几页
// 2.大于五页 前三页   1 2 3 4 5     12345
//	3.后三页    10页       67 8  9  10
//	4.中间页码    6  10页     6-2   6-1  6  6+1  6+2

}