package main

import (
	"github.com/micro/go-log"
	"net/http"

	"github.com/micro/go-web"
	"sss/IhomeWeb/handler"
	"github.com/julienschmidt/httprouter"
	_ "sss/IhomeWeb/model"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":22333"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/example/call", handler.ExampleCall)
	rou := httprouter.New()

	rou.NotFound = http.FileServer(http.Dir("html"))
	//模板
	rou.GET("/api/v1.0/areas", handler.GetArea)
	//欺骗浏览器
	rou.GET("/api/v1.0/house/index", handler.GetIndex)

	rou.GET("/api/v1.0/session", handler.GetSession)

	rou.GET("/api/v1.0/imagecode/:uuid",handler.GetImageCd)

	rou.GET("/api/v1.0/smscode/:mobile",handler.GetSmscd)

	rou.POST("/api/v1.0/users",handler.PostRet)

	//登陆服务
	rou.POST("/api/v1.0/sessions", handler.PostLogin)

	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
