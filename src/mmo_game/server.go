package main

import (
	"fmt"
	"mmo_game/api"
	"mmo_game/core"
	"zinx/ziface"
	"zinx/zlog"
	"zinx/znet"
)

//当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	zlog.Info("colsexxxxxxxxxxxxxxxxxxxxx")
	//获取当前连接的Pid属性
	pid, _ := conn.GetProperty("pid")

	//根据pid获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	if pid != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", pid, " left =====")

}

//当客户端建立连接的时候的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayId给客户端，走Msg:1消息
	player.SyncPid()
	//同步当有的玩家的初始化坐标信息给客户端，走Msg: 200消息
	player.BroadCastStartPosition()
	//========将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)
	//========================================
	//=================将该连接绑定属性Pid===============
	conn.SetProperty("pid", player.Pid)
	//===============================================
	//==============同步周边玩家上线信息，与现实周边玩家信息========
	player.SyncSurrounding()
	//=======================================================
	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}
func main() {
	//创建服务器句柄
	s := znet.NewServer()
	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	// ========== 注册 hook 函数 =====
	s.SetOnConnStop(OnConnectionLost)
	// ==============================
	s.AddRouter(2, &api.WorldChatRouter{})
	s.AddRouter(3, &api.MoveApi{}) //移动
	//启动服务
	s.Serve()
}
