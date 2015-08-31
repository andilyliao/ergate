// redis
package redis

import (
	"fmt"
	"redis"
	"sendslaver/exec"
	"sendslaver/utils"
	"strconv"
	"strings"
)

var spec []*redis.ConnectionSpec
var errtopic error
var sub []redis.PubSubClient

func GetSpec() []*redis.ConnectionSpec {
	return spec
}
func InitRedis() {
	redisnum := utils.Cfg.GetInt("redisNum")
	for i := 0; i < redisnum; i++ {
		ip := utils.Cfg.GetString("redisIp" + strconv.Itoa(i))
		port := utils.Cfg.GetInt("redisPort" + strconv.Itoa(i))
		spec[i] = redis.DefaultSpec()
		spec[i].Host(ip)
		spec[i].Port(port)
		sub[i], errtopic = redis.NewPubSubClientWithSpec(spec[i])
		if errtopic != nil {
			//			log.Println("failed to create the sync client", e)
			fmt.Printf("Error ", errtopic)
		}
		sub[i].Subscribe(utils.Cfg.GetString("slavename")) //队列名为从的名字
		go getmessage(i)
	}
}
func getmessage(i int) {
	for {
		ch := sub[i].Messages(utils.Cfg.GetString("slavename"))
		msg := <-ch
		fmt.Println(msg)
		cmds := strings.Split(string(msg), ",")
		//第1个是随机值，这个需要注意，这个值发送给任务shell后需要任务shell保留，但是操作任务的时候不需要这个值，目的是任务完成后为了把完整的method返回ACK
		exec.ExecCMD(cmds[0], cmds[1:])
	}

}
func SyncPublish(spec *redis.ConnectionSpec, channel string, message string) {

	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		//		log.Println("failed to create the sync client", e)
		fmt.Printf("Error ", err)
		return
	}

	for i := 0; i < 100; i++ {
		msg := []byte(message)
		rcvCnt, err := client.Publish(channel, msg)
		if err != nil {
			fmt.Printf("Error on Publish - %s", err)
		} else {
			fmt.Printf("Message sent to %d subscribers\n", rcvCnt)
		}
	}

	client.Quit()
}
