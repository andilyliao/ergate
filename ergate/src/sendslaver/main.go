// sendmaster project main.go
package main

import (
	//	"fmt"
	//	"lauchpad.net/gozk"
	"sendslaver/health"
	"sendslaver/utils"
	"sendslaver/zk"
)

func Init() {
	initCfg()           //读取配置文件
	initzk()            //链接zk
	health.Inithealth() //启动心跳
}
func initzk() {
	zk.Initzk()
}
func initCfg() {
	utils.Cfg = utils.LoadConfig("conf/sendslaver.cfg")
}
func main() {
	//启动被调服务
	Init()

}
