// redis
package redis

import (
	"fmt"
	"redis"
	"sendmaster/utils"
	"strconv"
)

var spec []*redis.ConnectionSpec
var errtopic error
var sub redis.PubSubClient

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
	}
	ackip := utils.Cfg.GetString("redisIpACK")
	ackport := utils.Cfg.GetInt("redisIpACK")
	var spec1 *redis.ConnectionSpec = redis.DefaultSpec()
	spec1.Host(ackip)
	spec1.Port(ackport)
	sub, errtopic = redis.NewPubSubClientWithSpec(spec1)
	if errtopic != nil {
		//		log.Println("failed to create the sync client", e)
		fmt.Printf("Error ", errtopic)
	}
	sub.Subscribe("ACK")
	go getmessage()
}
func getmessage() {
	for {
		ch := sub.Messages("ACK")
		msg := <-ch
		fmt.Println(msg)
		slave, _ := utils.TakeSlave(string(msg)) //msg是taskid
		utils.RemhealthtaskAndslave(string(msg))
		utils.RemhealthslaveAndtask(slave.(string))
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
