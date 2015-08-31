// inputservice
package inputservice

import (
	"net/http"
	"sendmaster/process"
)

func Do(w http.ResponseWriter, req *http.Request) {
	method := req.URL.Query().Get("method") //本期任务method就当做taskid使用了，所以 method每次发送过来的时候必须有一个特别的值
	process.SendTask(method)
}
