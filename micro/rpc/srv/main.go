package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"micro/rpc/srv/handler"
	"micro/rpc/srv/subscriber"

	example "micro/rpc/srv/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.srv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	//实例化一个机构体
	ex:=new(handler.Example)

	// Register Handler
	example.RegisterExampleHandler(service.Server(), ex)

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.srv", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.srv", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
