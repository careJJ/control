package routers

import (
	"pyg/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //发送短信
    beego.Router("/sendMsg",&controllers.UserController{},"post:HandleSendMsg")
	beego.Router("/register-email",&controllers.UserController{},"get:ShowEmail;post:HandleEmail")
	//激活用户
	beego.Router("/active",&controllers.UserController{},"get:Active")

	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")

	beego.Router("/index",&controllers.GoodsController{},"get:ShowIndex")

	beego.Router("/user/logout",&controllers.UserController{},"get:Logout")
	//展示用户中心页
	beego.Router("/user/userCenterInfo",&controllers.UserController{},"get:ShowUserCenterInfo")
	//收货地址页
	beego.Router("/user/site",&controllers.UserController{},"get:ShowSite;post:HandleSite")

	beego.Router("/index_sx",&controllers.GoodsController{},"get:ShowIndex_sx")
	//商品详情
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowDetail")

	beego.Router("/goodsType",&controllers.GoodsController{},"get:ShowList")

	beego.Router("/addCart",&controllers.CartController{},"post:HandleAddCart")

	beego.Router("/user/showCart",&controllers.CartController{},"get:ShowCart")

	beego.Router("/upCart",&controllers.CartController{},"post:HandleUpCart")
	beego.Router("/deleteCart",&controllers.CartController{},"post:HandDeleteCart")



    }



func guolvFunc(ctx*context.Context) {
	//过滤校验
	name := ctx.Input.Session("name")
	if name == nil {
		ctx.Redirect(302, "/login")
		return
	}
}