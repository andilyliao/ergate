// sendmaster project main.go
package main

import (
	//	"fmt"
	//	"lauchpad.net/gozk"
	"sendmaster/redis"
	"sendmaster/rpc"
	"sendmaster/utils"
	"sendmaster/zk"
)

func Init() {
	initCfg()         //读取配置文件
	initzk()          //链接zk
	redis.InitRedis() //链接redis的pub服务器
	initHttpServer()  //启动http服务，接收数据
}
func initHttpServer() {
	rpc.Init()
}
func initzk() {
	zk.Initzk()
}
func initCfg() {
	utils.Cfg = utils.LoadConfig("conf/sendmaster.cfg")
	rpc.SetServerIp(utils.Cfg.GetString("serverPort"))
}
func main() {
	//启动被调服务
	Init()
}
