package routers

import (
	"shoptest/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//登录
	beego.Router("/login",&controllers.CoreController{},"get:ShowLogin;post:HandleLogin")



    beego.Router("/order", &controllers.CoreController{},"get:ShowOrder;post:HandleOrder")
    //beego.Router()

}
