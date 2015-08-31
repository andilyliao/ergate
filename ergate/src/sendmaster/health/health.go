// health project main.go
package health

import (
	"sendmaster/zk"
	"time"
)

func Inithealth() {
	zk.GetHealthSlaverfromzk()
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			go zk.GetHealthSlaverfromzk()
		}
	}
}
