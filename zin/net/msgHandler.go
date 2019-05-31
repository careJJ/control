package net

import (
	"zin/zinface"
	"fmt"
	"zin/config"
)
type MsgHandler struct {
	//存放路由集合的map
	Apis map[uint32]zinface.IRouter
	//负责Worker取任务的消息队列  一个worker对应一个任务队列
	TaskQueue  []chan zinface.IRequest

	//worker工作池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandler() zinface.IMsgHandler {
	//给map开辟头空间
	return &MsgHandler{
		Apis: make(map[uint32]zinface.IRouter),
	}
}

//添加路由到map集合中
func (mh *MsgHandler) AddRouter(msgID uint32, router zinface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("repeat Api msgID", msgID)
		return
	}
	mh.Apis[msgID]=router
	fmt.Println("Apd api MagID = ",msgID,"succ!!")

}


//调度路由，根据MsgID
func (mh *MsgHandler) DoMsgHandler(request zinface.IRequest) {
	router,ok:=mh.Apis[request.GetMsg().GetMsgId()]
	if !ok{
		fmt.Println("api MsgID = ",request.GetMsg().GetMsgId(),"Not FOUND!!")
		return
	}
	//根据msgID,找到对应的router进行调用
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}
//一个worker真正处理业务的 goroutine函数
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan zinface.IRequest) {
	fmt.Println(" worker ID = ", workerID , " is starting... ")

	//不断的从对应的管道 等待数据
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

//启动Worker工作池 (在整个server服务中 只启动一次)
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("WorkPool is  started..")

	//根据WorkerPoolSize 创建worker goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//开启一个workergoroutine

		//1 给当前Worker所绑定消息channel对象 开辟空间  第0个worker 就用第0个Channel
		//给channel 进行开辟空间
		mh.TaskQueue[i] = make(chan zinface.IRequest, config.Globalobject.MaxWorkerTaskLen)

		//2 启动一个Worker，阻塞等待消息从对应的管道中进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}


//将消息添加到Worker工作池中 （将消息发送给对应的消息队列）
//应该是Reader来调用的
func (mh *MsgHandler) SendMsgToTaskQueue(request zinface.IRequest) {
	//1 将消息 平均分配给worker 确定当前的request到底要给哪个worker来处理
	//1个客户端绑定一个worker来处理
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	//2 直接将 request 发送给对应的worker的taskqueue
	mh.TaskQueue[workerID] <- request
}