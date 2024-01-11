package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"zinx/ziface"
)

type Conf struct {
	TcpServer        ziface.IServer //当前Zinx的全局Server对象
	Host             string         //当前服务器主机IP
	TcpPort          int            //当前服务器主机监听端口号
	Name             string         //当前服务器名称
	Version          string         //当前Zinx版本号
	MaxPacketSize    uint32         //都需数据包的最大值
	MaxConn          int            //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32         //业务工作Worker池的数量
	MaxWorkerTaskLen uint32         //业务工作Worker对应负责的任务队列最大任务存储数量
}

func (c *Conf) Reload() {
	fmt.Println(os.Getwd())
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GameConfig)
	if err != nil {
		panic(err)
	}
}

/*
	定义一个全局的对象
*/
var GameConfig *Conf

/*
	提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GameConfig = &Conf{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024, //业务工作Worker对应负责的任务队列最大任务存储数量
	}

	//从配置文件中加载一些用户配置的参数
	GameConfig.Reload()
}
