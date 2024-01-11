package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	Handlers map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	mh := &MsgHandle{
		Handlers: make(map[uint32]ziface.IRouter),
	}
	return mh
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Handlers[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not Found!")
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
