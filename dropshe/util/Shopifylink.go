package util

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-shopify-master/shopify"
	"io/ioutil"
	"net/http"
	"time"
)

type ProductInfo struct {
	Products []Products `json:"products"`
}

type Variants struct {
	ID                   int64       `json:"id"`
	ProductID            int64       `json:"product_id"`
	Title                string      `json:"title"`
	Price                string      `json:"price"`
	Sku                  string      `json:"sku"`
	Position             int         `json:"position"`
	InventoryPolicy      string      `json:"inventory_policy"`
	CompareAtPrice       string      `json:"compare_at_price"`
	FulfillmentService   string      `json:"fulfillment_service"`
	InventoryManagement  string      `json:"inventory_management"`
	Option1              string      `json:"option1"`
	Option2              string      `json:"option2"`
	Option3              interface{} `json:"option3"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	Taxable              bool        `json:"taxable"`
	Barcode              interface{} `json:"barcode"`
	Grams                int         `json:"grams"`
	ImageID              interface{} `json:"image_id"`
	Weight               float64     `json:"weight"`
	WeightUnit           string      `json:"weight_unit"`
	InventoryItemID      int64       `json:"inventory_item_id"`
	InventoryQuantity    int         `json:"inventory_quantity"`
	OldInventoryQuantity int         `json:"old_inventory_quantity"`
	RequiresShipping     bool        `json:"requires_shipping"`
	AdminGraphqlAPIID    string      `json:"admin_graphql_api_id"`
}
type Images struct {
	ID                int64         `json:"id"`
	ProductID         int64         `json:"product_id"`
	Position          int           `json:"position"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	Alt               interface{}   `json:"alt"`
	Width             int           `json:"width"`
	Height            int           `json:"height"`
	Src               string        `json:"src"`
	VariantIds        []interface{} `json:"variant_ids"`
	AdminGraphqlAPIID string        `json:"admin_graphql_api_id"`
}

type Products struct {
	ID       int64      `json:"id"`
	Title    string     `json:"title"`
	Variants []Variants `json:"variants"`
	Images   []Images   `json:"images"`
}

//用于接受product信息的表
type ProductMatch struct {
	ProductId int64
	Title     string
	Variant   []Variant
	Image     []Image
}

//跟ProductMatch一对多，一个product_id对应多个variant
type Variant struct {
	Id    int64
	Sku   string
	Price string
	Title string
}

//跟ProductMatch一对多
type Image struct {
	Src string
}

type OrderInfo struct {
	Orders []Orders `json:"orders"`
}
type ShopMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}
type PresentmentMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}
type PriceSet struct {
	ShopMoney        ShopMoney        `json:"shop_money"`
	PresentmentMoney PresentmentMoney `json:"presentment_money"`
}
type TotalDiscountSet struct {
	ShopMoney        ShopMoney        `json:"shop_money"`
	PresentmentMoney PresentmentMoney `json:"presentment_money"`
}
type AmountSet struct {
	ShopMoney        ShopMoney        `json:"shop_money"`
	PresentmentMoney PresentmentMoney `json:"presentment_money"`
}
type DiscountAllocations struct {
	Amount                   string    `json:"amount"`
	DiscountApplicationIndex int       `json:"discount_application_index"`
	AmountSet                AmountSet `json:"amount_set"`
}
type OriginLocation struct {
	ID           int64  `json:"id"`
	CountryCode  string `json:"country_code"`
	ProvinceCode string `json:"province_code"`
	Name         string `json:"name"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	City         string `json:"city"`
	Zip          string `json:"zip"`
}
type LineItems struct {
	ID                         int64                 `json:"id"`
	VariantID                  int64                 `json:"variant_id"`
	Title                      string                `json:"title"`
	Quantity                   int                   `json:"quantity"`
	Sku                        string                `json:"sku"`
	VariantTitle               string                `json:"variant_title"`
	Vendor                     string                `json:"vendor"`
	FulfillmentService         string                `json:"fulfillment_service"`
	ProductID                  int64                 `json:"product_id"`
	RequiresShipping           bool                  `json:"requires_shipping"`
	Taxable                    bool                  `json:"taxable"`
	GiftCard                   bool                  `json:"gift_card"`
	Name                       string                `json:"name"`
	VariantInventoryManagement string                `json:"variant_inventory_management"`
	Properties                 []interface{}         `json:"properties"`
	ProductExists              bool                  `json:"product_exists"`
	FulfillableQuantity        int                   `json:"fulfillable_quantity"`
	Grams                      int                   `json:"grams"`
	Price                      string                `json:"price"`
	TotalDiscount              string                `json:"total_discount"`
	FulfillmentStatus          interface{}           `json:"fulfillment_status"`
	PriceSet                   PriceSet              `json:"price_set"`
	TotalDiscountSet           TotalDiscountSet      `json:"total_discount_set"`
	DiscountAllocations        []DiscountAllocations `json:"discount_allocations"`
	AdminGraphqlAPIID          string                `json:"admin_graphql_api_id"`
	TaxLines                   []interface{}         `json:"tax_lines"`
	OriginLocation             OriginLocation        `json:"origin_location"`
}
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
type Orders struct {
	ID              int64           `json:"id"`
	Email           string          `json:"email"`
	CreatedAt       time.Time       `json:"created_at"`
	TotalPrice      string          `json:"total_price"`
	FinancialStatus string          `json:"financial_status"`
	Name            string          `json:"name"`
	LineItems       []LineItems     `json:"line_items"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

const (
	GetOrderId    = "admin/api/2019-04/orders.json?fields=id"
	GetOrderInfo  = "admin/api/2019-04/orders.json?fields=id,name,created_at,financial-status,total_price,email,line_items,shipping_address"
	GetProduct_id = "admin/api/2019-04/products.json?fields=id"
	produtid      = "admin/api/2019-04/products.json?ids="
	GetProuctImg  = ""
)

//对接shopify
//返回链接的套接字
func LinkStore(storename, password string) *shopify.Client {

	client, err := shopify.NewClient(nil, shopify.ShopURL("https://"+storename+".myshopify.com/admin"), shopify.Token(password))

	if err != nil {
		beego.Error("连接shopify失败")
	}
	//s, _, err := client.Shop.Get(context.Background())
	//return s
	return client
}

//连接Beachmolly,返回链接的套接字

//func LinkBeachmolly() *shopify.Client {
//
//	client, err := shopify.NewClient(nil, shopify.ShopURL("https://beachmolly.myshopify.com/admin"), shopify.Token("8b2405755ddb36e3bd784ebc72e2fc99"))
//
//	if err != nil {
//		beego.Error("连接shopify失败")
//
//	}
//	//s, _, err := client.Shop.Get(context.Background())
//	//return s
//	return client
//}

//返回order id切片的方法
func GetOrderid(name, apikey, secret string) []string {
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	idurl := url + GetOrderId
	resp, err := http.Get(idurl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	var J map[string][]map[string]int

	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}
	orders := J["orders"]
	var ids []string
	for _, order := range orders {
		id := order["id"]
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	return ids
}

//获取Beachmolly所需订单信息json   name,created_at,financial-status,total_price,email,fulfillment,sku,quantity,amount,shipping_address
//获得完整的json直接穿前端，再根据product_id查表productMatch得到图片
func GetOrderJson(name, apikey, secret string) (str string) {
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	infoUrl := url + GetOrderInfo
	resp, err := http.Get(infoUrl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("获取MGID失败")
	}
	var J map[string][]map[string]interface{}
	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}
	for _, b := range J {
		if c, err := json.Marshal(b); err == nil {
			str := string(c[:])
			return str
			//fmt.Println(STR)
		}
	}
	return
}

//也许要检查一下reflect.Valueof 在使用之前返回   这种情况，将得到一个零值，因为反射不能看到未导出的函数
func GetOrderStruct(name, apikey, secret string, i int) (Orders, ShippingAddress, []LineItems) {
	ids := GetOrderid(name, apikey, secret)
	O := OrderInfo{}
	o1 := Orders{}
	l1 := []LineItems{}
	l := LineItems{}
	//ship := ShippingAddress{}
	s1 := ShippingAddress{}
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	orderurl := url + GetOrderInfo + "&ids=" + ids[i]
	resp, err := http.Get(orderurl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	if err := json.Unmarshal(body, &O);
		err != nil {
		fmt.Println(err)
	}
	o1 = O.Orders[0]
	for _, b := range o1.LineItems {

		l.Sku = b.Sku
		l.Title = b.Title
		l.ID = b.ProductID
		l.Price = b.Price
		l.VariantID = b.VariantID
		l.Quantity = b.Quantity
		l1 = append(l1, l)
	}
	s1 = o1.ShippingAddress
	return o1, s1, l1
}

//获取product_id的切片
func GetProductId(name, apikey, secret string) []string {
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	productIDurl := url + GetProduct_id

	resp, err := http.Get(productIDurl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		beego.Error(err)
	}
	var J map[string][]map[string]int
	//
	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}

	orders := J["products"]
	var ids []string
	for _, order := range orders {
		id := order["id"]
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	//返回图片id的切片
	return ids
	//return(strings.Join(ids, ","))
}

//返回BM产品的第一张图片

//func DealBMsrc() {
//	//imabuf:=make (chan []string)
//	pid := GetBMProductJson()
//	len := len(pid)
//	//buf := make(chan []string, len)
//	for i := 0; i < len; i++ {
//		ima := GetBmSrc(i)
//		//fmt.Println(ima)
//		var image []string
//		image = append(image, fmt.Sprintf("%s", ima[0]))
//		fmt.Println(strings.Join(image, ","))
//	}
//}

//获得ProductMatch结构体
func GetProductMatch(name, apikey, secret string, i int) (int64, string, []Variant, []Image) {
	///admin/api/2020-01/products/4468589985930/images.json?fields=src
	pid := GetProductId(name, apikey, secret)
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	productinfo := ProductInfo{}
	a := Products{}
	//p:=ProductMatch{}
	d := Variant{}
	imgs := Image{}
	d1 := []Variant{}
	imgs1 := []Image{}
	purl := url + produtid + pid[i] + "&fields=id,title,variants,images"

	resp, err := http.Get(purl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
	}
	if err := json.Unmarshal(body, &productinfo);
		err != nil {
		fmt.Println(err)
	}
	a = productinfo.Products[0]
	for _, b := range a.Variants {
		d.Sku = b.Sku
		d.Id = b.ID
		d.Price = b.Price
		d1 = append(d1, d)
		//fmt.Println(reflect.TypeOf(b))
	}

	for _, c := range a.Images {
		imgs.Src = c.Src
		imgs1 = append(imgs1, imgs)
	}
	return a.ID, a.Title, d1, imgs1

}
