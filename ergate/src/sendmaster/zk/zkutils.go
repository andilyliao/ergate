// zk
package zk

import (
	"fmt"
	"launchpad.net/gozk/zookeeper"
	//	"log"
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	"sendmaster/process"
	"sendmaster/utils"
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

	_, err = zk.Create(utils.Cfg.GetString("sendslavers"), "0", 0, zookeeper.WorldACL(zookeeper.PERM_ALL)) //zookeeper.EPHEMERAL|zookeeper.SEQUENCE

	if err != nil {
		//		log.Fatalf("Can't create counter: %v", err)
		fmt.Println("Can't create counter: %v", err)
	} else {
		fmt.Println("Counter created!")
	}
	watchslavers()
}
func watchslavers() {
	_, _, session1, err := zk.ChildrenW(utils.Cfg.GetString("sendslavers"))
	if err != nil {
		//		log.Fatalf("Can't connect: %v", err)
		fmt.Println(err)
	}
	event := <-session1
	alivedset := hashset.New()
	deadset := hashset.New()
	alivedset.Clear()
	deadset.Clear()
	if event.Type == zookeeper.EVENT_CHILD {
		slavers, _, err2 := zk.Children(utils.Cfg.GetString("sendslavers"))
		if err2 != nil {
			//		log.Fatalf("Can't connect: %v", err)
			fmt.Println(err2)
		}
		for _, slaver := range slavers {
			//			utils.AddlivedSlavers(slaver)
			alivedset.Add(slaver)
		}
	}
	for _, oldslave := range utils.GetlivedSlavers().Values() {
		if oldslave == nil {
			break
		}
		alivedset.Contains(oldslave) //判断获取的活着的slave里面是否包换之前的活着的列表里面的slave，如果没有，说明过去的节点死了
		deadset.Add(oldslave)
	}
	utils.SetlivedSlavers(alivedset)
	//容错，重新下发死了的节点的任务
	for _, slave := range deadset.Values() {
		utils.RemlivedSlaver(slave.(string))
		utils.RemhealthSlaver(slave.(string))
		taskid, _ := utils.TakeTask(slave.(string))
		utils.RemhealthslaveAndtask(slave.(string))
		utils.RemhealthtaskAndslave(taskid.(string))
		//重新发送任务
		process.SendTask(taskid.(string))
	}
	watchslavers()
}

func GetHealthSlaverfromzk() {
	strs, _, err := zk.Children(utils.Cfg.GetString("sendslavers"))
	if err != nil {
		//		log.Fatalf("Can't connect: %v", err)
		fmt.Println(err)
	}
	flag := false
	res := hashset.New()
	var err3 error
	var healthstr string
	for _, str := range strs {
		healthstr, _, err3 = zk.Get(utils.Cfg.GetString("sendslavers") + "/" + str) //healthstr=ok/err
		if err3 != nil {
			//		log.Fatalf("Can't connect: %v", err)
			fmt.Println(err3)
		}
		if healthstr == "ok" {
			flag = true
			res.Add(str)
			break
		}
	}
	if flag == false {
		err3 = errors.New("no have healthy slaver now!!!!!!")
	}
	utils.SethealthSlavers(res)
}
