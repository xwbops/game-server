package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	//将conn连接添加到ConnMananger中,疑问点，这里是否可以不用加判断,如果存在不做处理
	if _, ok := c.connections[conn.GetConnID()]; !ok {
		c.connections[conn.GetConnID()] = conn
	}
	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())
}

//删除链接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	//删除连接信息
	delete(c.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", c.Len())
}

//利用链接Id获取链接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (c *ConnManager) Len() int {
	return len(c.connections)
}

//删除并停止所有链接
func (c *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	//停止并删除全部的连接信息
	for connId, conn := range c.connections {
		//停止连接
		conn.Stop()
		//删除连接
		delete(c.connections, connId)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}
