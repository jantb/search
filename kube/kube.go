package kube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func GetPods() Pods {
	output, err := exec.Command("kubectl", "get", "pods", "-o", "json").CombinedOutput()
	checkErr(err)

	var getPods = Pods{}
	err = json.Unmarshal(output, &getPods)
	checkErr(err)
	return getPods
}

func GetPodLogsStreamFastJson(podName string, insertChanJson chan []byte) {
	command := exec.Command("kubectl", "logs", "-f", "--since=200h", podName)
	pipe, err := command.StdoutPipe()
	command.Start()
	checkErr(err)
	reader := bufio.NewReader(pipe)

	var line []byte
	for {
		line, err = reader.ReadBytes(byte('\n'))
		if err != nil {
			time.Sleep(60 * time.Second)
			go func(insertChanJson chan []byte, podName string) {
				GetPodLogsStreamFastJson(podName, insertChanJson)
			}(insertChanJson, podName)
			return
		}
		insertChanJson <- line
	}
}

func checkErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
