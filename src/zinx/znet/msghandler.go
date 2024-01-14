package znet

import (
	"fmt"
	"strconv"
	"zinx/conf"
	"zinx/ziface"
)

type MsgHandle struct {
	Handlers       map[uint32]ziface.IRouter
	WorkerPoolSize uint32 //业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	mh := &MsgHandle{
		Handlers:       make(map[uint32]ziface.IRouter),
		WorkerPoolSize: conf.GameConfig.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, conf.GameConfig.WorkerPoolSize),
	}
	return mh
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Handlers[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not Found!")
		return
	}
	//执行对应处理方法
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Handlers[msgId]; ok {
		panic("repeated handle , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Handlers[msgId] = router
	fmt.Println("Add handler msgId = ", msgId)
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, conf.GameConfig.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
