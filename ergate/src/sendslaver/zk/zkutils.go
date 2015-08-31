// zk
package zk

import (
	"fmt"
	"launchpad.net/gozk/zookeeper"
	//	"log"
	"sendslaver/utils"
)

var zk *zookeeper.Conn

func GetZkCon() (res *zookeeper.Conn) {
	return res
}
func Initzk() {
	ipport := utils.Cfg.GetString("zkServer")
	zk1, session, err := zookeeper.Dial(ipport, 5e9)
	zk = zk1
	if err != nil {
		//		log.Fatalf("Can't connect: %v", err)
		fmt.Println(err)
	}
	//	defer zk.Close()

	// Wait for connection.
	event := <-session
	if event.State != zookeeper.STATE_CONNECTED {
		//		log.Fatalf("Can't connect: %v", event)
		fmt.Println(event.String())
	}
	fmt.Println(event)
	_, err = zk.Create(utils.Cfg.GetString("sendslavers")+"/"+utils.Cfg.GetString("slavename"), "0", zookeeper.EPHEMERAL, zookeeper.WorldACL(zookeeper.PERM_ALL)) //zookeeper.EPHEMERAL|zookeeper.SEQUENCE

	if err != nil {
		//		log.Fatalf("Can't create counter: %v", err)
		fmt.Println("Can't create counter: %v", err)
	} else {
		fmt.Println("Counter created!")
	}
}

//healthstr=ok/err   sendslavers/slavename里面写
func Sethealth(health string) {
	_, err := zk.Set(utils.Cfg.GetString("sendslavers")+"/"+utils.Cfg.GetString("slavename"), health, 0)
	if err != nil {
		//		log.Fatalf("Can't connect: %v", err)
		fmt.Println(err)
	}
}
