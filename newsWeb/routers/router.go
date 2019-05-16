package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/index",&controllers.ArticleController{},"get,post:ShowIndex")
	beego.Router("/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/content",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/delete",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/deleteType",&controllers.ArticleController{},"get:HandleDeleteType")
	beego.Router("/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")




}
