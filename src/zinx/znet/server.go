package znet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"zinx/conf"
	"zinx/ziface"
	"zinx/zlog"
)

//iServer 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//协议版本号
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	//Router ziface.IRouter
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager

	// =======================
	//新增两个hook函数原型

	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)

	// =======================
}

//创建一个服务器句柄
func NewServer() ziface.IServer {
	// 初始化日志选项
	zlog.SetOptions(zlog.WithDisableCaller(true))
	//先初始化全局配置文件
	//conf.GameConfig.Reload()
	s := &Server{
		Name:      conf.GameConfig.Name,
		IPVersion: "tcp4",
		IP:        conf.GameConfig.Host,
		Port:      conf.GameConfig.TcpPort,
		//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
		MsgHandler: NewMsgHandle(),   //msgHandler 初始化
		ConnMgr:    NewConnManager(), //创建ConnManager
	}
	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========
//开启网络服务
func (s *Server) Start() {
	zlog.Infof("[START] Server listenner at IP: %s, Port %d, is starting", s.IP, s.Port)
	zlog.Infof("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d",
		conf.GameConfig.Version,
		conf.GameConfig.MaxConn,
		conf.GameConfig.MaxPacketSize)
	//开启一个go去做服务器listener业务
	go func() {
		//0 启动worker工作池机制
		s.MsgHandler.StartWorkerPool()
		//获取一个tcp的 addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			zlog.Error("resolve tcp addr err: ", err)
			return
		}
		// 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			zlog.Error("listen ", s.IPVersion, " err ", err)
			return
		}
		//已经监听成功
		zlog.Info("start Zinx server ", s.Name, " success, now listenning...")
		var cid uint32 = 0
		// 启动server网络链接业务
		for {
			//阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				zlog.Error("Accept err: ", err)
				continue
			}
			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			//=============
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= conf.GameConfig.MaxConn {
				zlog.Info("conn closed, connections has been exceeded: ", conf.GameConfig.MaxConn)
				conn.Close()
				continue
			}
			//=============
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
			//TODO Server.Start()设置服务器最大链接控制，如果超过最大连接，那么则关闭新连接
			//TODO Server.Start()处理该新连接请求的业务方法，此时应该有handler和conn是绑定的
			//我们这里暂时做一个最大512字节的回显服务
			//go func() {
			//	//不断的循环从客户端获取数据
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err ", err)
			//			continue
			//		}
			//		//回显
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err ", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

//关闭网络服务
func (s *Server) Stop() {
	zlog.Info("[STOP] Zinx server , name ", s.Name)
	//TODO Server.Stop()将其他需要清理的连接信息或者其他信息也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//开启业务服务器方法
func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	c := make(chan os.Signal, 1)
	//(监听指定信号 ctrl+c kill信号)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	zlog.Infof("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.Name, sig)
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

//获取链接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		zlog.Info("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		zlog.Info("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
