package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	example "micro/grpc/srv/proto/example"
	//"github.com/micro/go-micro/client"

	//example "github.com/micro/examples/template/srv/proto/example"
	"github.com/micro/go-grpc"
)
						//回复               传入（或理解为输出，输入）
func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	//创建一个map，用来接收前段发过来的数据
	var request map[string]interface{}
	//将r.body的数据解析到request这个map当中，解析结束之后r.body的数据就空了
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cli:=grpc.NewService()
	cli.Init()
	//客户端句柄        通过pb文件调用床架你客户端句柄的函数（      服务名          客户端默认参数            ）
	exampleClient := example.NewExampleService("go.micro.srv.srv", cli.Client())


	// call the backend service
	//客户端句柄        通过pb文件调用床架你客户端句柄的函数（      服务名                 客户端默认参数            ）
	//exampleClient := example.NewExampleService("go.micro.srv.srv", client.DefaultClient)
	//通过客户端句柄，调用函数 往里面传参（默认参数（上下文通讯），pb文件中的结构体）
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{

		Name: request["name"].(string),//断言
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	//创建一个返回的map，将值赋值到map当中
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	//将map转换成为json发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
