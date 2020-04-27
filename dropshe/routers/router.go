package routers

import (
	"dropshe/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//UserController  用户中心的管理
	//注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	//邮箱激活
	beego.Router("/register/activate", &controllers.UserController{}, "get:ShowActivateEmail;post:ActivateEmail")
	//激活用户
	beego.Router("/activate", &controllers.UserController{}, "get:Activate")
	//登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")

	//google登录
	//beego.Router("/login/oauth",&controllers.UserController{},"get,post:HandleGoogleLogin")

	//个人信息页面
	beego.Router("/account/general", &controllers.UserController{}, "get:ShowAccount;post:UpdateInfo")
	//更改密码
	beego.Router("/account/password", &controllers.UserController{}, "get:ShowPassword;post:HandlePassword")
	//beego.Router("/account/changepassword",&controllers.UserController{},"post:HandlePassword")
	//更改头像
	//beego.Router("account/image",&controllers.UserController{},"post:SentImage")
	//退出登录
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")
	//信用卡操作
	beego.Router("/account/billing", &controllers.UserController{}, "get:ShowBilling;post:HandleCard")
	//添加新卡
	beego.Router("/setting/setting_billing", &controllers.UserController{}, "get:ShowAddCard;post:AddCard")
	//beego.Router("/setting/settingbilling",&controllers.UserController{},"post:AddCard")

	//管理支付方式

	beego.Router("/account/payment", &controllers.UserController{}, "get:ShowPayment;post:HandlePayment")
	//展示交易明细
	beego.Router("/account/invoices", &controllers.UserController{}, "get:ShowInvoies")
	//找回密码
	beego.Router("/forgot_password", &controllers.UserController{}, "get:ShowForgetPassword;post:FindPassword")

	//HubController   商品页面的管理

	//CoreController 核心功能抓取与需求的反馈
	//展示订单页面
	//beego.Router("/order/drop",&controllers.CoreController{},"get:ShowOrder;post:HandleOrder")
	//交易短信验证

	//展示shopify店铺页面
	beego.Router("/stores", &controllers.CoreController{}, "get:ShowStore")

	//展示添加店铺
	beego.Router("/stores/add", &controllers.CoreController{}, "get:ShowAddStore")
	beego.Router("/stores/addstore", &controllers.CoreController{}, "post:AddStore")

	//一键同步订单

	//测试订单同步
	beego.Router("/order/drop", &controllers.CoreController{}, "get:ShowOrder")

	//订单反馈
	beego.Router("/sourcing/add", &controllers.CoreController{}, "get:ShowSourceAgent;post:HandleSourceAgent")

	//添加和删除图片
	//beego.Router("/sourcing/api/upload", &controllers.CoreController{}, "post:Upload")
	//beego.Router("/sourcing/api/delImage", &controllers.CoreController{}, "post:SourcDeleteImage")

	//客户采购需求列表页面（或者理解为待确定采购列表）
	beego.Router("/sourcing", &controllers.CoreController{}, "get:ShowSourcingtest")
	//删除sourcing
	beego.Router("/sourcing/delete", &controllers.CoreController{}, "post:DeleteSourcing")
	//sourcing 搜索  SourcingSearch
	beego.Router("/sourcing/search", &controllers.CoreController{}, "get:SourcingSearch")
	//编辑sourcingAgent
	beego.Router("/sourcing/update", &controllers.CoreController{}, "get:ShowUpdateSourcingAgent;post:UpdateSourcingAgent")

	//LogisticsController      物流功能

	//PayController   支付功能
	//展示支付页面
	beego.Router("/sourcing/detail", &controllers.CoreController{}, "get:ShowSourcingDetail")
	//支付
	//beego.Router("/pay/stripe", &controllers.PayController{}, "post:RequestStripe")
	//交易短信验证
	beego.Router("/sourcing/pay", &controllers.CoreController{}, "post:HandlePay")

	//beego.Router("/checkout",&controllers.CoreController{})

	//ChatController 即时通讯

	//PurseController 电子钱包

	//ErpController ERP系统

	//ERP员工端的处理   问题：权限管理，业务员和仓库收发员

	//ERP系统的注册
	//beego.Router("/erpregister",&controllers.ErpController{},"get:ShowErpRegister,post:HandleErpRegister")
	////ERP系统的登录
	//beego.Router("/erplogin",&controllers.ErpController{},"get:ShowErpLogin,post:HandleErpLogin")
	////ERP首页的展示
	//beego.Router("/erpindex",&controllers.ErpController{},"get:ShowErpIndex")

	beego.Router("/pay/stripe",&controllers.PayController{},"post:HandlePayStripe")
	beego.Router("/pay/after",&controllers.CoreController{},"post:Pay")


	//SKU采购（考虑到客户订单的支付）

	//ERP仓库模块
	//SKU的库存管理页面展示

	//SKU的库存管理操作

	//

	//ShowErpSourceList
	beego.Router("/erpsourcing", &controllers.ErpController{}, "get:ShowErpSource")

	beego.Router("/erpsourcing/searchbysku", &controllers.ErpController{}, "get:SearchBySku")

	beego.Router("/erpsourcing/searchbystore", &controllers.ErpController{}, "get:SearchByStore")

	beego.Router("/erpsourcing/edit", &controllers.ErpController{}, "get:ShowErpAudit")

	beego.Router("/erpsourcing/mate", &controllers.ErpController{}, "post:HandleErpEdit")
//展示已经匹配好的订单 默认
	beego.Router("/erporder", &controllers.ErpController{}, "get:ShowErpOrder_status2")
	//展示已经支付完待处理的订单
	beego.Router("/erporder/status3", &controllers.ErpController{}, "get:ShowErpOrder_status3")


	//修改已经匹配好的订单
	beego.Router("/erporder/edit", &controllers.ErpController{}, "post:HandleOrder")

	beego.Router("/erporder/search", &controllers.ErpController{}, "post:OrderSearch")
	//展示订单的细节
	beego.Router("/erporder/detail", &controllers.ErpController{}, "get:ShowErpOrderEdit")




}
