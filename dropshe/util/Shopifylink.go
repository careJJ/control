package util

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-shopify-master/shopify"
	"io/ioutil"
	"net/http"
	"strconv"
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
	ShippingLines   []ShippingLines `json:"shipping_lines"`
	Note            interface{}     `json:"note"`
	//Customer               Customer               `json:"customer"`
}
type ShippingLines struct {
	ID                            int64              `json:"id"`
	Title                         string             `json:"title"`
	Price                         string             `json:"price"`
	Code                          string             `json:"code"`
	Source                        string             `json:"source"`
	Phone                         interface{}        `json:"phone"`
	RequestedFulfillmentServiceID interface{}        `json:"requested_fulfillment_service_id"`
	DeliveryCategory              interface{}        `json:"delivery_category"`
	CarrierIdentifier             interface{}        `json:"carrier_identifier"`
	DiscountedPrice               string             `json:"discounted_price"`
	PriceSet                      PriceSet           `json:"price_set"`
	DiscountedPriceSet            DiscountedPriceSet `json:"discounted_price_set"`
	DiscountAllocations           []interface{}      `json:"discount_allocations"`
	TaxLines                      []interface{}      `json:"tax_lines"`
}
//用于获取时效和收货地址

//关联shippingLines，暂不知何物，支付金额？
type DiscountedPriceSet struct {
	ShopMoney        ShopMoney        `json:"shop_money"`
	PresentmentMoney PresentmentMoney `json:"presentment_money"`
}


type Counts struct {
	Count int `json:"count"`
}

type Va struct {
	Variant Variant1 `json:"variant"`
}

//用于解析variantid
type Variant1 struct {
	ID                   int64       `json:"id"`
	ProductID            int64       `json:"product_id"`
	Title                string      `json:"title"`
	Price                string      `json:"price"`
	Sku                  string      `json:"sku"`
	Position             int         `json:"position"`
	InventoryPolicy      string      `json:"inventory_policy"`
	CompareAtPrice       interface{} `json:"compare_at_price"`
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
	ImageID              int64       `json:"image_id"`
	Weight               float64     `json:"weight"`
	WeightUnit           string      `json:"weight_unit"`
	InventoryItemID      int64       `json:"inventory_item_id"`
	InventoryQuantity    int         `json:"inventory_quantity"`
	OldInventoryQuantity int         `json:"old_inventory_quantity"`
	RequiresShipping     bool        `json:"requires_shipping"`
	AdminGraphqlAPIID    string      `json:"admin_graphql_api_id"`
}

type GetProductImage struct {
	Image struct {
		ID                int64       `json:"id"`
		ProductID         int64       `json:"product_id"`
		Position          int         `json:"position"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		Alt               interface{} `json:"alt"`
		Width             int         `json:"width"`
		Height            int         `json:"height"`
		Src               string      `json:"src"`
		VariantIds        []int64     `json:"variant_ids"`
		AdminGraphqlAPIID string      `json:"admin_graphql_api_id"`
	} `json:"image"`
}

type Ships struct {
	Orders []Ordership `json:"orders"`
}


type Ordership struct {
	Email string  `json:"email"`
	ShippingLines   []ShippingLines `json:"shipping_lines"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}


const (
	GetOrderId    = "admin/api/2019-04/orders.json?fields=id"
	OrderInfoAPI2019  = ".myshopify.com/admin/api/2019-04/orders.json?fields=id,name,created_at,financial-status,total_price,email,line_items,shipping_address,ShippingLines,note&page="
	ProductInfoAPI2019a = ".myshopify.com/admin/api/2019-04/products.json?ids="
	ProductInfoAPI2019b ="&fields=id,title,variants,images"

	produtid      = "admin/api/2019-04/products.json?ids="
	GetProuctImg  = ""
	OrderCount    = "admin/api/2020-04/orders/count.json"
	ProductCount  = ".myshopify.com/admin/api/2020-04/products/count.json"
	ordercountAPI2019 = ".myshopify.com/admin/api/2019-04/orders/count.json"
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



//获取product_id的切片
//func GetProductId(name, apikey, secret string) []string {
//	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
//	productIDurl := url + GetProduct_id
//
//	resp, err := http.Get(productIDurl)
//	if err != nil {
//		beego.Error(err)
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		beego.Error(err)
//	}
//	var J map[string][]map[string]int
//	//
//	if err := json.Unmarshal(body, &J);
//		err != nil {
//		fmt.Println(err)
//	}
//
//	orders := J["products"]
//	var ids []string
//	for _, order := range orders {
//		id := order["id"]
//		ids = append(ids, fmt.Sprintf("%d", id))
//	}
//	//返回图片id的切片
//	return ids
//	//return(strings.Join(ids, ","))
//}

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
//func GetProductMatch(name, apikey, secret string, i int) (int64, string, []Variant, []Image) {
//	///admin/api/2020-01/products/4468589985930/images.json?fields=src
//	pid := GetProductId(name, apikey, secret)
//	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
//	productinfo := ProductInfo{}
//	a := Products{}
//	//p:=ProductMatch{}
//	d := Variant{}
//	imgs := Image{}
//	d1 := []Variant{}
//	imgs1 := []Image{}
//	purl := url + produtid + pid[i] + "&fields=id,title,variants,images"
//
//	resp, err := http.Get(purl)
//	if err != nil {
//		beego.Error(err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		beego.Error(err)
//	}
//	if err := json.Unmarshal(body, &productinfo);
//		err != nil {
//		fmt.Println(err)
//	}
//	a = productinfo.Products[0]
//	for _, b := range a.Variants {
//		d.Sku = b.Sku
//		d.Id = b.ID
//		d.Price = b.Price
//		d1 = append(d1, d)
//		//fmt.Println(reflect.TypeOf(b))
//	}
//
//	for _, c := range a.Images {
//		imgs.Src = c.Src
//		imgs1 = append(imgs1, imgs)
//	}
//	return a.ID, a.Title, d1, imgs1
//
//}

type OrdersCount struct {
	Count int `json:"count"`
}

func GetOrderCount(name, apikey, secret string) int {
	url := "https://" + apikey + ":" + secret + "@" + name + ".myshopify.com/"
	curl := url + OrderCount
	count := OrdersCount{}
	resp, err := http.Get(curl)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
	}
	if err := json.Unmarshal(body, &count);
		err != nil {
		beego.Error(err)
	}
	c := count.Count
	return c
}

func GetOrderCountTest(name,key,password string) int {
	url := "https://"+key+":"+password+"@"+name+ordercountAPI2019
	//idurl := url + GetOrderId
	resp, err := http.Get(url)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	J := Counts{}

	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}

	return J.Count
}

func GetOrderStructTest(name,key,password string,i int) (OrderInfo) {
	//ids := GetOrderid11()
	O := OrderInfo{}
	//o1 := Orders{}
	s := strconv.Itoa(i)
	url := "https://"+key+":"+password+"@"+name+OrderInfoAPI2019+ s
	//orderurl := url + GetOrderInfo + "&ids=" + ids[i]
	resp, err := http.Get(url)
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
	return O
}

func GetProductStruct(name,key,password string,i int) ProductInfo {
	pid := strconv.Itoa(i)
	url := "https://"+key+":"+password+"@"+name+ProductInfoAPI2019a+ pid + ProductInfoAPI2019b
	productinfo := ProductInfo{}
	resp, err := http.Get(url)
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
	return productinfo
}

func GetSrc(i, p int64) (string) {
	vid := strconv.FormatInt(i, 10)
	url := "https://1318a69cafc1b553ca362b0ad295ffad:a0712bbc46b4efe98b5a0218fef68eff@beachmolly.myshopify.com/admin/api/2020-01/variants/" + vid + ".json"
	variant := Va{}
	resp, err := http.Get(url)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
	}
	if err := json.Unmarshal(body, &variant);
		err != nil {
		fmt.Println(err)
	}
	beego.Alert(variant)
	imgid := strconv.FormatInt(variant.Variant.ImageID, 10)
	pid := strconv.FormatInt(p, 10)

	imgurl := "https://1318a69cafc1b553ca362b0ad295ffad:a0712bbc46b4efe98b5a0218fef68eff@beachmolly.myshopify.com/admin/api/2020-01/products/" + pid + "/images/" + imgid + ".json"
	images := GetProductImage{}
	imgresp, err := http.Get(imgurl)
	if err != nil {
		beego.Error(err)
	}
	defer imgresp.Body.Close()
	imgbody, err := ioutil.ReadAll(imgresp.Body)
	if err != nil {
		beego.Error(err)
	}
	if err := json.Unmarshal(imgbody, &images);
		err != nil {
		fmt.Println(err)
	}
	src := images.Image.Src

	return src

}

func GetAllProductImages() {

}

func GetProductCount(name,key,password string)int{
	url := "https://"+key+":"+password+"@"+name+ProductCount
	resp, err := http.Get(url)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	J := Counts{}

	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}

	return J.Count
}

func GetShippingAdress(i string)(Ships){
	url := "https://1318a69cafc1b553ca362b0ad295ffad:a0712bbc46b4efe98b5a0218fef68eff@beachmolly.myshopify.com/admin/api/2019-04/orders.json?fields=email,shipping_address,shipping_lines&ids="+i
	resp, err := http.Get(url)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	J := Ships{}

	if err := json.Unmarshal(body, &J);
		err != nil {
		fmt.Println(err)
	}
	return J
}