package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//开启业务服务器方法
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	//AddRouter(router IRouter)
	AddRouter(msgId uint32, router IRouter)
	//得到链接管理器
	GetConnMgr() IConnManager
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(connection IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(connection IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
