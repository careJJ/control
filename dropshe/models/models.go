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
	Id              int    `orm:"unique;auto"`
	Email           string `orm:"size(40)"`
	FirstName       string `orm:"size(40)"`
	LastName        string `orm:"size(40)"`
	Password        string `orm:"size(256)"`
	TimeZone        string //[]*TimeZone`orm:"reverse(many)"`
	Country         string //[]*Country`orm:"reverse(many)"`
	Language        string //[]*Language`orm:"reverse(many)"`
	Active          bool
	ShopifyStore    *ShopifyStore    `orm:"rel(one)"`      //关联店铺信息的表，以user表为主表，根据user的字段查询店铺信息，一个用户对应多个店铺
	CreditCard      *CreditCard      `orm:"reverse(many)"` //关联信用卡的表，一个user对应多张卡
	Sourcing        *Sourcing        `orm:"reverse(many)"`
	Sourcing_Demand *Sourcing_Demand `orm:"reverse(many)"`
}

//店铺
type ShopifyStore struct {
	Email       string
	Name        string
	Status      bool
	ProductSync int
	OrderSync   int
	ApiKey      string
	Secret      string
	CreateTime  string
	User        *User `orm:"reverse(one)"`
}

//信用卡
type CreditCard struct {
	Email      string //跟user中email关联
	Company    string
	Number     string `orm:"size(16)"`
	CVC        string
	Year       string
	Month      string
	State      bool //状态，是否使用
	Priority   int  //优先级，用于置顶
	Active     bool
	Updatetime time.Time
	User       *User `orm:"rel(fk)"`
}

//员工表
type ErpUser struct {
	ID       string
	Name     string
	Password string
	Power    string //权限

}

//订单列表      id,name,created_at,financial-status,total_price,email,fulfillment,sku,quantity,amount,shipping_address
type Order struct {
	Id               int64 `orm:"unique;auto"` //唯一键
	Name             string
	Created_at       time.Time
	Financial_status string
	Total_price      string

	Email string

	LineItems *[]LineItems

	Shipping_address *ShippingAddress

	User     *User
	Sourcing *Sourcing `orm:"rel(one)"`
}

//跟ProductMatch多对多?
type LineItems struct {
	Sku          string
	Id           int64 //也就是product_id
	VariantID    int64
	Title        string
	Price        string
	Quantity     int
	ProductMatch *ProductMatch
	//Variant		*[]Variant
	ErpSku  *ErpSku
}

//订单物流信息
type ShippingAddress struct {
	FirstName    string      `json:"first_name"`
	Address1     string      `json:"address1"`
	Phone        string      `json:"phone"`
	City         string      `json:"city"`
	Zip          string      `json:"zip"`
	Province     string      `json:"province"`
	Country      string      `json:"country"`
	LastName     string      `json:"last_name"`
	Address2     string      `json:"address2"`
	Company      interface{} `json:"company"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	Name         string      `json:"name"`
	CountryCode  string      `json:"country_code"`
	ProvinceCode string      `json:"province_code"`
}

//客户提出1采购需求 /Sourcing/add中产生
type Sourcing_Demand struct {
	//采购单中客户提的采购需求
	Id                 int
	Title              string
	Link               string
	Target_price       string //客户提出的目标价格
	Shipping_price     string
	Estimated_Delivery string //物流时效
	Count              int
	Status             bool		//状态
	Description        string 	//描述
	General bool  `orm:"default(true)"`//普货
	Destination string //目的地（国家）
	//EUB  bool   //是否支持易邮宝

	//订单的真实信息
	Images             string
	User               *User `orm:"rel(fk)"`
	OrderTotalPrice    string

	Order *Order
	LineItems *[]LineItems
}

//用于sourcing 页面的展示
type Sourcing struct {
	Id           int
	Link         string
	Title        string
	Images       string
	Status       bool
	Target_price string //客户目标价格，来自Sourcing_Demand
	//Product_price      string //物品价格
	Estimated_Delivery string //预计交货
	General bool
	Destination  string
	User               *User  `orm:"rel(fk)"`
	//Order Order
	ErpSku *[]ErpSku

}

type Goodsku struct {
	Name string
}

//erp采购表
type Erp_Sourcing_Demand struct {
	Sourcing_Demand *[]Sourcing_Demand `orm:"rel(one)"` //erp中的采购表跟采购需求一一对应？
}

//用于匹配
type ProductMatch struct {
	Id       int64 //即product_id
	Title    string
	Variant  *[]Variant
	ImageSrc *[]ImageSrc
}

//跟ProductMatch一对多，一个product_id对应多个variant
type Variant struct {
	Id    int64
	Sku   string
	Price string
	Title string
}

//跟ProductMatch一对多
type ImageSrc struct {
	Src     string
	FdfsSrc string
}

//erp中用于sku匹配的表
type ErpSku struct {
	Id int64  //等於product_id
	Sku string


}

//erp用于

//

func init() {
	//注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/dropshe")
	//注册表结构
	orm.RegisterModel(new(User))
	//跑起来
	orm.RunSyncdb("default", false, true)
}
