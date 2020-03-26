package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//外键始终在子表上
/*#一个用户对应一个简介；一个简介对应一个用户；
one2one:User(子表) -> Profile（主表）;one2one:Profile -> User

#一个邮件对应一个用户；一个用户有多个邮件；
one2many:Post(子表) -> User（主表）;many2one:User -> Post

#一个邮件对应多个标签；一个标签对应多个邮件；
many2many:Post(子表) -> Tag（主表）;many2many:Tag -> Post*/

//用户信息表
type User struct {
	Id           int    `orm:"pk;unique;auto"`
	Email        string `orm:"size(40)"`
	FirstName    string `orm:"size(40)"`
	LastName     string `orm:"size(40)"`
	Password     string `orm:"size(256)"`
	TimeZone     string `orm:"size(40)"`
	Country      string `orm:"size(40)"`
	Language     string `orm:"size(40)"`
	Activate       bool
	Number int
	StoreName string`orm:"size(40)"`
	//头像
	Image string
	//ShopifyStore *ShopifyStore `orm:"rel(one)",null`//关联店铺信息的表，以user表为主表，根据user的字段查询店铺信息，一个用户对应多个店铺
	CreditCard     []*CreditCard     `orm:"reverse(many)"` //关联信用卡的表，一个user对应多张卡
	Sourcing       []*Sourcing       `orm:"reverse(many)"`
	SourcingDemand []*SourcingDemand `orm:"reverse(many)"`
}

//店铺
type ShopifyStore struct {
	Id          int64  `orm:"unique;auto"`
	Email       string `orm:"size(40)"`
	Name        string `orm:"size(40)"`
	Status      bool
	ProductSync int
	OrderSync   int
	ApiKey      string `orm:"size(256)"`
	Secret      string `orm:"size(256)"`
	CreateTime  time.Time
	//User        *User `orm:"reverse(one)"`
	Order       []*Order `orm:"reverse(many)"`
}

//信用卡
type CreditCard struct {
	Id         int64  `orm:"unique;auto"`
	Email      string `orm:"size(40)"` //跟user中email关联
	Company    string `orm:"size(40)"`
	Number     string `orm:"size(16)"`
	CVC        string `orm:"size(40)"`
	Year       string `orm:"size(4)"`
	Month      string `orm:"size(2)"`
	State      bool   //状态，是否使用
	Priority   int    //优先级，用于置顶
	Activate     bool
	Updatetime time.Time
	User       *User `orm:"rel(fk)"`
}

//员工表
type ErpUser struct {
	ID       int64  `orm:"unique;auto"`
	Name     string `orm:"size(40)"`
	Password string `orm:"size(40)"`
	Power    string `orm:"size(40)"` //权限

}

//订单列表      id,name,created_at,financial-status,total_price,email,fulfillment,sku,quantity,amount,shipping_address
type Order struct {
	Id               int64  `orm:"unique;auto"` //唯一键
	Name             string `orm:"size(40)"`
	Created_at       time.Time
	Financial_status string `orm:"size(40)"`
	Total_price      string `orm:"size(40)"`

	Email string `orm:"size(40)"`

	LineItems []*LineItems `orm:"reverse(many)"`

	Shipping_address *ShippingAddress `orm:"rel(one)"`
	Sourcing *Sourcing `orm:"rel(one)"`
	SourcingDemand *SourcingDemand `orm:"rel(one)"`
	ShopifyStore *ShopifyStore `orm:"rel(fk)"`
}

//跟ProductMatch多对多?
type LineItems struct {
	Sku       string `orm:"size(100)"`
	Id        int64  //也就是product_id
	VariantID int64
	Title     string `orm:"size(100)"`
	Price     string `orm:"size(40)"`
	Quantity  int
	//ProductMatch *ProductMatch
	Variant		[]*Variant `orm:"reverse(many)"`
	ErpSku  *ErpSku  `orm:"reverse(one)"`
	Order *Order `orm:"rel(fk)"`
	SourcingDemand *SourcingDemand `orm:"rel(fk)"`
}


//订单物流信息
type ShippingAddress struct {
	Id int64   `orm:"unique;auto"`
	FirstName    string  `json:"first_name"`
	Address1     string  `json:"address1"`
	Phone        string  `json:"phone"`
	City         string  `json:"city"`
	Zip          string  `json:"zip"`
	Province     string  `json:"province"`
	Country      string  `json:"country"`
	LastName     string  `json:"last_name"`
	Address2     string  `json:"address2"`
	Company      string  `json:"company"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Name         string  `json:"name"`
	CountryCode  string  `json:"country_code"`
	ProvinceCode string  `json:"province_code"`
	Order        *Order  `orm:"reverse(one)"`//一对一的反向关系
}

//客户提出1采购需求 /Sourcing/add中产生
type SourcingDemand struct {
	//采购单中客户提的采购需求
	Id                 int
	Title              string `orm:"size(100)"`
	Link               string `orm:"size(100)"`
	Target_price       string `orm:"size(40)"` //客户提出的目标价格
	Shipping_price     string `orm:"size(100)"`
	Estimated_Delivery string `orm:"size(100)"` //物流时效
	Quantity           int
	Status             bool   //状态
	Description        string `orm:"size(256)"`     //描述
	General            bool   `orm:"default(true)"` //普货
	Destination        string `orm:"size(40)"`      //目的地（国家）
	//EUB  bool   //是否支持易邮宝
	//订单的真实信息
	Images          string `orm:"size(256)"`
	User            *User  `orm:"rel(fk)"`
	OrderTotalPrice string `orm:"size(40)"`

	Order *Order  `orm:"reverse(one)"`
	LineItems  []*LineItems  `orm:"reverse(many)"`
}

//用于sourcing 页面的展示
type Sourcing struct {
	Id           int
	Link         string `orm:"size(100)"`
	Title        string `orm:"size(100)"`
	Images       string `orm:"size(256)"`
	Status       bool
	Target_price string `orm:"size(40)"` //客户目标价格，来自Sourcing_Demand
	//Product_price      string //物品价格
	Estimated_Delivery string `orm:"size(100)"` //预计交货
	General            bool
	Destination        string `orm:"size(40)"`
	Ready              bool
	Updatetime         time.Time
	User               *User `orm:"rel(fk)"`
	Order              *Order `orm:"reverse(one)"`
	ErpSku []*ErpSku `orm:"reverse(many)"`
}

//erp采购表
//type Erp_Sourcing_Demand struct {
//	SourcingDemand []*SourcingDemand `orm:"rel(one)"` //erp中的采购表跟采购需求一一对应？
//}

//用于匹配
type ProductMatch struct {
	Id    int64  //即product_id
	Title string `orm:"size(100)"`
	Variant  []*Variant   `orm:"reverse(many)"`
	ImageSrc []*ImageSrc   `orm:"reverse(many)"`

}

//跟ProductMatch一对多，一个product_id对应多个variant
type Variant struct {
	Id    int64
	Sku   string `orm:"size(100)"`
	Price string `orm:"size(40)"`
	Title string `orm:"size(100)"`
	ProductMatch *ProductMatch  `orm:"rel(fk)"`
	LineItems *LineItems   `orm:"rel(fk)"`
}

//跟ProductMatch一对多
type ImageSrc struct {
	Id int64   `orm:"unique;auto"`
	Src     string `orm:"size(100)"`
	FdfsSrc string `orm:"size(100)"`
	ProductMatch *ProductMatch  `orm:"rel(fk)"`
}

//erp中用于sku匹配的表
type ErpSku struct {
	Id  int64  //等於product_id
	Sku string `orm:"size(100)"`
	Sourcing *Sourcing `orm:"rel(fk)"`
	LineItems *LineItems `orm:"rel(one)"`
}



func init() {
	//注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/shoptest")
	//注册表结构   	orm.RegisterModel(new(User),new(ShopifyStore),new(CreditCard),new(ErpUser),new(Order),new(LineItems),new(ShippingAddress),new(Sourcing_Demand),new(Sourcing),new(Goodsku),new(Erp_Sourcing_Demand),new(ProductMatch),new(Variant),new(ImageSrc),new(ErpSku))
	//	orm.RegisterModel(new(User),new(ShopifyStore),new(CreditCard),new(ErpUser),new(Order),new(LineItems),new(ShippingAddress),new(SourcingDemand),new(Sourcing),new(ProductMatch),new(Variant),new(ImageSrc),new(ErpSku))
	orm.RegisterModel(new(User),new(ShopifyStore),new(CreditCard),new(ErpUser),new(Order),new(LineItems),new(ShippingAddress),new(SourcingDemand),new(Sourcing),new(ProductMatch),new(Variant),new(ImageSrc),new(ErpSku))

	//跑起来
	orm.RunSyncdb("default", false, true)


}

