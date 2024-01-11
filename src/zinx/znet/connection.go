package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法api
	//handleAPI ziface.HandFunc

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool

	//该连接的处理方法router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

/*处理conn读数据的Goroutine*/
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		// 创建拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg head
		headBuf := make([]byte, dp.GetHeadLen())
		//cnt, err := c.Conn.Read(buf)
		//fmt.Println(cnt)
		if _, err := io.ReadFull(c.GetTCPConnection(), headBuf); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headBuf)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}
		var dataBuf []byte
		if msg.GetDataLen() > 0 {
			dataBuf = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), dataBuf); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(dataBuf)
		//根据 dataLen 读取 data，放在msg.Data中

		// 调用当前链接业务（这里执行的是当前的conn绑定的handle方法）
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("connID", c.ConnID, " handle is error")
		//	c.ExitBuffChan <- true
		//	return
		//}
		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg, //将之前的buf 改成 msg
		}
		//从路由Routers 中找到注册绑定Conn的对应Handle
		go func(request ziface.IRequest) {
			fmt.Println(string(request.GetData()))
			//执行注册的路由方法
			fmt.Println(c.Router)
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}

}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()
	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

//从当前连接获取原始的socket TCPConn GetTCPConnection *net.TCPConn //获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取远程客户端地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection is closed when send msg")
	}
	msg := NewMessage(msgId, data)
	//将data封装包并且发送
	dp := NewDataPack()
	msgData, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	if _, err := c.Conn.Write(msgData); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}
	return nil
}
