package routers

import (
	"dropshe/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
//UserController  用户中心的管理
    //注册
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//邮箱激活
	beego.Router("/register-email",&controllers.UserController{},"get:ShowEmail;post:HandleEmail")
	//激活用户
	beego.Router("/active",&controllers.UserController{},"get:Active")
    //登录
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //google登录
	beego.Router("/login/oauth",&controllers.UserController{},"get,post:HandleGoogleLogin")


	//个人信息页面
    beego.Router("/accout/genereal",&controllers.UserController{},"get:ShowAccount;post:UpdateInfo")
	//退出登录
	beego.Router("/user/logout",&controllers.UserController{},"get:Logout")
    //信用卡操作
    beego.Router("/accout/billing",&controllers.UserController{},"get:ShowBilling;post:HandleCard")
	//添加新卡
	beego.Router("hub/setting/setting_billing",&controllers.UserController{},"post:AddCard")
	//管理支付方式
	beego.Router("accout/payment",&controllers.UserController{},"get:ShowPayment;post:HandlePayment")
	//展示交易明细
	beego.Router("accout/invoices",&controllers.UserController{},"get:ShowInvoies")

//HubController   商品页面的管理



//CoreController 核心功能抓取与需求的反馈
	//展示订单页面
	beego.Router("/order/drop",&controllers.CoreController{},"get:ShowOrder;post:HandleOrder")
	//交易短信验证

	//展示shopify店铺页面
	beego.Router("/stores",&controllers.CoreController{},"get:ShowStore;post:HandleStore")

	//添加店铺
	beego.Router("/stores/add",&controllers.CoreController{},"post:AddStore")
	//一键同步订单

    //订单反馈
    beego.Router("/sourcing/add",&controllers.CoreController{},"post:ShowSource")
	//客户采购需求列表页面（或者理解为待确定采购列表）
	beego.Router("/Souring",&controllers.CoreController{},"get:ShowSourcing;post:HandleSouring")
    //删除sourcing
	beego.Router("/Souring/delete",&controllers.CoreController{},"get:DeleteSouring")


	//CartController    购物车功能
	//展现购物车页面



//LogisticsController      物流功能



//PayController   支付功能
	//展示支付页面

	//支付

	//交易短信验证

//ChatController 即时通讯



//PurseController 电子钱包





//ErpController ERP系统

  //ERP员工端的处理   问题：权限管理，业务员和仓库收发员

	//ERP系统的注册
	//beego.Router("/stores",&controllers.ErpController{},"get:ShowErpRegister,post:HandleErpRegister")
	////ERP系统的登录
	//beego.Router("/stores",&controllers.ErpController{},"get:ShowErpLogin,post:HandleErpLogin")
	////ERP首页的展示
	//beego.Router("/stores",&controllers.ErpController{},"get:ShowErpIndex")


  //ERP订单模块

  	//ERP待处理页展示,采购单
	beego.Router("/erpsource",&controllers.ErpController{},"get:ShowErpSource")
	//ERP信息审核页面展示（sku匹配頁面）
	beego.Router("/erpsource/sku",&controllers.ErpController{},"get:ShowErpAudit")
	//ERP sku匹配（）
	beego.Router("/auditsku",&controllers.ErpController{},"post:AuditSku")

	//SKU采购（考虑到客户订单的支付）

  //ERP仓库模块
	//SKU的库存管理页面展示

	//SKU的库存管理操作









}
