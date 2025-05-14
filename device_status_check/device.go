// device.go: program running on edge device
// todo: make 2 funcs(1. measure cpu, mem usage, 2. send request to server)
//
// input: No input
// output: cpu, mem usage, result in JSON format

package main

import (
	"fmt"
	"time"

	"github.com/DataDog/gopsutil/cpu"
	"github.com/DataDog/gopsutil/mem"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	measurePeriod = 3 // period to measure available CPU, Memory percent
	serverUrl = "http://[your_ip]:[your_port]"
)

type AppRunStatus struct {
	CPU float64 `json: "cpu"`
	Memory float64 `json: "memory"`
	TaskNum int `json: "taskNum"`
}

type workload struct {
	StartNum int 
	CurrentNum int
	EndNum	int
}


func measureCPUMemUsage() (cpuAvailable float64, memAvailable float64, err error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "measureCPUMemUsage",
	})

	// interval to measure cpu usage, used in cpu.Percent
	cpuMeasureInterval := 100 * time.Millisecond  
	// get available CPU core numbers
	c, cpuErr := cpu.Percent(cpuMeasureInterval, false)

	if cpuErr != nil{
		logger.Error(cpuErr)
	}

	cpuAvailable = c[0]
	
	// get memory info
	m, memErr := mem.VirtualMemory()

	if memErr != nil{
		logger.Error(memErr)
	}

	// return available memory in GB
	memAvailable = (float64(m.Available) / (1024 * 1024 * 1024)) 

	return cpuAvailable, memAvailable, nil
}
// 현재 상태
func sendDeviceStatus(cpuAvail float64, memAvail float64, taskNum int) (string, error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "sendDeviceStatus",
	})
	client := resty.New()
	data := fmt.Sprintf(`{"cpu": %f, "memory": %f, "taskNum": %d}`, cpuAvail, memAvail, taskNum)
	resp, err := client.R().SetHeader("Content-Type", "application/json").
							SetBody(data).
							Post(serverUrl)

	if err != nil {
		logger.Error(err)
		return "", err
	}
	
	// get response body as string, need to use this response to reschedule
	result := resp.String() 

	return result, nil
}

func main(){
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
	})

	ticker := time.NewTicker(measurePeriod * time.Second)
	defer ticker.Stop()

	for  i := 0; i < 10; i++ {
		<- ticker.C
		CPUAvail, MemAvail, measureErr := measureCPUMemUsage()
		if measureErr != nil {
			logger.Error(measureErr)
		}
		
		result := 10 // TODO: need to get result from application
	
		todo, sendErr := sendAppRunStatus(CPUAvail, MemAvail, result)
		if sendErr != nil {
			logger.Error(sendErr)
		}
	
		print(todo)
		
		logger.Infof("Available CPU: %.1f Cores, Available Mem: %.1fGB", CPUAvail, MemAvail)
	}
}