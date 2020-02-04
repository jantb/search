package kube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetPods() Pods {
	output, err := exec.Command("kubectl", "get", "pods", "-o", "json").CombinedOutput()
	checkErr(err)

	var getPods = Pods{}
	err = json.Unmarshal(output, &getPods)
	checkErr(err)
	return getPods
}

func GetPodLogs(podName string, insertChanJson chan map[string]interface{}) {
	output, err := exec.Command("kubectl", "logs", "--since=200h", podName).CombinedOutput()
	checkErr(err)
	for _, line := range strings.Split(string(output), "\n") {
		var j map[string]interface{}
		if err := json.Unmarshal([]byte(line), &j); err != nil {
		} else {
			insertChanJson <- j
		}
	}
}

func GetPodLogsStream(podName string, insertChanJson chan map[string]interface{}) {
	command := exec.Command("kubectl", "logs", "-f", "--since=200h", podName)
	pipe, err := command.StdoutPipe()
	command.Start()
	checkErr(err)
	reader := bufio.NewReader(pipe)

	var line string
	for {
		line, err = reader.ReadString('\n')
		var j map[string]interface{}
		if err := json.Unmarshal([]byte(line), &j); err != nil {
		} else {
			insertChanJson <- j
		}
		if err != nil {
			break
		}
	}
}

func checkErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
