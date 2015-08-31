// exec
package exec

import (
	"fmt"
	"os/exec"
)

func ExecCMD(name string, args []string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	d, err1 := cmd.Output()
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(string(d))
}
