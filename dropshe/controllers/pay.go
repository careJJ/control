package controllers

import (
	//"encoding/json"
	//"fmt"
	//"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	//"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"dropshe/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/stripe/stripe-go/paymentintent"
	"net/http"

	"html/template"
	

	stripe "github.com/stripe/stripe-go"
	//"math/rand"
	//"os"
	//"regexp"
	//"time"
	////"github.com/astaxie/beego/utils"
	//"github.com/plutov/paypal"
	//"dropshe/models"
)

type PayController struct {
	beego.Controller
}

type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

/*
//生成图片验证码
func CheckImages(this *beego.Controller) {
	cap := captcha.New()
	//通过句柄调用 字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		beego.Info("没有字体文件")
		panic(err.Error())
	}
	//设置图片的大小
	cap.SetSize(91, 41)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	//生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)

	beego.Info(str)

	a := *img      //解引用

	email := this.GetSession("email")

	//连接redis
	bm, err := util.RedisOpen(util.Server_name, util.Redis_addr,
		util.Redis_port, util.Redis_dbnum)
	if err != nil {
		beego.Error("连接数据库失败")
	}
	//存入redis，有效期时间5min
	bm.Put(email.(string), str, time.Second*300)
	//将图片转成base64
	imgbuff:=bytes.NewBuffer(nil)
	//将图片写入到buff
	jpeg.Encode(imgbuff,a,nil)
	//开辟存储空间
	dist:=make([]byte,50000)
	//将buff转成base64
	base64.StdEncoding.Encode(dist,imgbuff.Bytes())
}

*/

//交易短信验证  缺前端的ajax
//func (this *PayController) HandleSenMsg() {
//	//接受数据
//
//	countrycode := this.GetString("Country Code")
//	phone := this.GetString("Phone")
//	mixphone := "+" + countrycode + phone
//	resp := make(map[string]interface{})
//
//	defer RespFunc(&this.Controller, resp)
//	//返回json格式数据
//	//校验数据
//	if countrycode == "" || phone == "" {
//		beego.Error("获取电话号码失败")
//		//2.给容器赋值
//		resp["errno"] = 1
//		resp["errmsg"] = "获取电话号码错误"
//		return
//	}
//	//检查电话号码格式是否正确
//	reg, _ := regexp.Compile(`\+(9[976]\d|8[987530]\d|6[987]\d|5[90]\d|42\d|3[875]\d|
//2[98654321]\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|
//4[987654310]|3[9643210]|2[70]|7|1)\d{1,14}$`)
//	result := reg.FindString(mixphone)
//	if result == "" {
//		beego.Error("电话号码格式错误")
//		//2.给容器赋值
//		resp["errno"] = 2
//		resp["errmsg"] = "电话号码格式错误"
//		return
//	}
//	//发送短信   SDK调用
//	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIu4sh9mfgqjjr", "sTPSi0Ybj0oFyqDTjQyQNqdq9I9akE")
//	if err != nil {
//		beego.Error("电话号码格式错误")
//		//2.给容器赋值
//		resp["errno"] = 3
//		resp["errmsg"] = "初始化短信错误"
//		return
//	}
//	//生成6位数随机数
//	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
//	//验证码
//	vcode := fmt.Sprintf("%06d", rnd.Int31n(1000000))
//	//阿里云短息服务
//	request := requests.NewCommonRequest()
//	request.Method = "POST"
//	request.Scheme = "https" // https | http
//	request.Domain = "dysmsapi.aliyuncs.com"
//	request.Version = "2017-05-25"
//	request.ApiName = "SendSms"
//	request.QueryParams["RegionId"] = "cn-hangzhou"
//	request.QueryParams["PhoneNumbers"] = phone
//	request.QueryParams["SignName"] = "dropshe"
//	request.QueryParams["TemplateCode"] = "SMS_164275022"
//	request.QueryParams["TemplateParam"] = "{\"code\":" + vcode + "}"
//
//	response, err := client.ProcessCommonRequest(request)
//	if err != nil {
//		beego.Error("电话号码格式错误")
//		//2.给容器赋值
//		resp["errno"] = 4
//		resp["errmsg"] = "短信发送失败"
//		return
//	}
//	//json数据解析
//	var message Message
//	json.Unmarshal(response.GetHttpContentBytes(), &message)
//	if message.Message != "OK" {
//		beego.Error("电话号码格式错误")
//		//2.给容器赋值
//		resp["errno"] = 6
//		resp["errmsg"] = message.Message
//		return
//	}
//
//	resp["errno"] = 5
//	resp["errmsg"] = "发送成功"
//	resp["code"] = vcode
//}

//NewClient-> NewRequest-> SendWithAuth
//func (this *PayController)HandlePayPal(){
//	c,err:=paypal.NewClient("","",paypal.APIBaseLive)
//	if err!=nil{
//		beego.Error(err)
//		this.Redirect("/pay",302)
//		return
//	}
//
//	c.SetLog(os.Stdout) // Set log to terminal stdout
//
//	accessToken, err := c.GetAccessToken()
//
//}
//type CheckoutData struct {
//	ClientSecret string
//}

func (this *PayController) HandlePayStripe() {
	beego.Alert("HandlePay")
	sourceidlist := this.GetStrings("sourceidList[]")
	beego.Alert(sourceidlist)
	o := orm.NewOrm()
	var order models.Order
	var prices float64
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	for i := 0; i < len(sourceidlist); i++ {
		o.QueryTable("Order").Filter("SourceId", sourceidlist[i]).One(&order)
		if order.TotalPrice==0{
			resp["errno"]=1
			resp["errmsg"]="请选择等待支付的订单"
			return
		}
		price := order.TotalPrice
		prices += price
	}
	//得到选择订单的价格
	beego.Alert(prices)

	this.Data["Prices"]=prices
	finalprice:=int64(prices*100)


	stripe.Key = "sk_test_TyAe7P109fmUiQBHTf0xGn9G00EaocYpqU"

	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(finalprice),
		Currency: stripe.String(string(stripe.CurrencyHKD)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
	}
	// Verify your integration in this guide by including this parameter
	params.AddMetadata("integration_check", "accept_a_payment")

	pi, _ := paymentintent.New(params)


	checkoutTmpl := template.Must(template.ParseFiles("views/checkout.html"))

	http.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
		intent := pi
		data := CheckoutData{
			ClientSecret: intent.ClientSecret,
		}
		checkoutTmpl.Execute(w, data)
	})

	http.ListenAndServe(":3000", nil)


}


