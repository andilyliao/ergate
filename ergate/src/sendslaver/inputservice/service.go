// inputservice
package inputservice

import (
	//	"fmt"
	"net/http"
	"redis"
	redis1 "sendslaver/redis"
	"sendslaver/utils"
)

func Do(w http.ResponseWriter, req *http.Request) {
	taskid := req.URL.Query().Get("taskid") // 通知主删除这个任务
	var spec *redis.ConnectionSpec = redis.DefaultSpec()
	ackip := utils.Cfg.GetString("redisIpACK")
	ackport := utils.Cfg.GetInt("redisIpACK")
	spec.Host(ackip)
	spec.Port(ackport)
	redis1.SyncPublish(spec, "ACK", taskid)
}
