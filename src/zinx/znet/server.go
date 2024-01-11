package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/conf"
	"zinx/ziface"
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
	msgHandler ziface.IMsgHandle
}

//============== 实现 ziface.IServer 里的全部接口方法 ========
//开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		conf.GameConfig.Version,
		conf.GameConfig.MaxConn,
		conf.GameConfig.MaxPacketSize)
	//开启一个go去做服务器listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		//获取一个tcp的 addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		// 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		//已经监听成功
		fmt.Println("start Zinx server ", s.Name, " success, now listenning...")
		var cid uint32 = 0
		// 启动server网络链接业务
		for {
			//阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, cid, s.msgHandler)
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
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//TODO Server.Stop()将其他需要清理的连接信息或者其他信息也要一并停止或者清理
}

//开启业务服务器方法
func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

//创建一个服务器句柄
func NewServer() ziface.IServer {
	//先初始化全局配置文件
	conf.GameConfig.Reload()
	s := &Server{
		Name:      conf.GameConfig.Name,
		IPVersion: "tcp4",
		IP:        conf.GameConfig.Host,
		Port:      conf.GameConfig.TcpPort,
		//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
		msgHandler: NewMsgHandle(), //msgHandler 初始化
	}
	return s
}
