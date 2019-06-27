package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"micro/grpc/srv/handler"
	//"micro/grpc/srv/subscriber"
	"github.com/micro/go-grpc"
	example "micro/grpc/srv/proto/example"
)

func main() {
	// New Service

		service := grpc.NewService(
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
	//micro.RegisterSubscriber("go.micro.srv.srv", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.srv", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
