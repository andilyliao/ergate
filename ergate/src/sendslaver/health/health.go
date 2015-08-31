// health project main.go
package health

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	//	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"sendslaver/zk"
	"time"
)

func Inithealth() {
	gethealth()
	t := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t.C:
			go gethealth()
		}
	}
}
func gethealth() {
	var res string
	r1 := getMem()
	r2 := getCPU()
	if r1 && r2 {
		res = "ok"
	} else {
		res = "err"
	}
	zk.Sethealth(res)
}
func getMem() bool {
	var flag bool
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total/1024, v.Free/1024, v.UsedPercent/1024)

	// convert to JSON. String() is also implemented
	if v.Free/(1024*1024*1024) >= 1 {
		flag = true
	} else {
		flag = false
	}
	return flag
}

func getCPU() bool {
	var flag bool
	i := float64(counts() * 100) // 个数
	j := cpuPercents(false)      //负载
	//	cpuPercents(true)
	//	cpuload()
	if (i - j) >= 100 {
		flag = true
	} else {
		flag = false
	}
	return flag
}
func counts() int {
	v, err := cpu.CPUCounts(true)
	if err != nil {
		fmt.Println("error %v", err)
	}
	if v == 0 {
		fmt.Println("could not get CPU counts: %v", v)
	}
	fmt.Println(v)
	return v
}

func cpuPercents(percpu bool) float64 {
	numcpu := runtime.NumCPU()
	if runtime.GOOS != "windows" {
		v, err := cpu.CPUPercent(time.Millisecond, percpu)
		if err != nil {
			fmt.Println("error %v", err)
		}
		if (percpu && len(v) != numcpu) || (!percpu && len(v) != 1) {
			fmt.Println("wrong number of entries from CPUPercent: %v", v)
		}
	}
	duration := time.Duration(10) * time.Microsecond
	v, err := cpu.CPUPercent(duration, percpu)
	if err != nil {
		fmt.Println("error %v", err)
	}
	var res float64
	res = 0.0
	for _, percent := range v {
		if percent < 0.0 || percent > 100.0*float64(numcpu) {
			fmt.Println("CPUPercent value is invalid: %f", percent)
			res += percent
		}
	}
	fmt.Println(res)
	return res / float64(numcpu)
}

//func cpuload() {
//	v, err := load.LoadAvg()
//	if err != nil {
//		fmt.Println("error %v", err)
//	}

//	empty := &load.LoadAvgStat{}
//	if v == empty {
//		fmt.Println("error load: %v", v)
//	}
//	fmt.Println(v)
//}
