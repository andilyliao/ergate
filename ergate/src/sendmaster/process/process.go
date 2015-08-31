// inputservice
package process

import (
	"fmt"
	"hash/crc32"
	"sendmaster/redis"
	"sendmaster/utils"
)

func SendTask(method string) {
	str := utils.GethealthSlavers().Values()[0]
	if str == nil || str.(string) == "" {
		//		log.Fatalf("Can't connect: %v", err)
		fmt.Println(str)
	}
	redisnum := utils.Cfg.GetInt("redisNum")
	spec := redis.GetSpec()
	num := int(crc32.ChecksumIEEE([]byte(method)))
	i := num % redisnum
	redis.SyncPublish(spec[i], str.(string), method)  //队列名是从机的名字，内容是接收到的内容
	utils.AddhealthslaveAndtask(str.(string), method) //容错使用
	utils.AddhealthtaskAndslave(method, str.(string)) //任务留存一下，以后任务做完就删除，任务未做完就保留
}
