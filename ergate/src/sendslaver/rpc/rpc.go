// rpc
package rpc

import (
	"net/http"
	"runtime"
	"sendslaver/inputservice"
)

var serIp string

func SetServerIp(serverIp string) {
	serIp = serverIp
}
func Init() {
	runtime.GOMAXPROCS(1)
	http.HandleFunc("/finishtask", inputservice.Do) //http://127.0.0.1:8080/xxx?method=cccc&bbb=aaa
	http.ListenAndServe(":"+serIp, nil)
}
