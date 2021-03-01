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
	output, err := exec.Command("oc", "get", "pods", "-o", "json").CombinedOutput()
	checkErr(err)

	var getPods = Pods{}
	err = json.Unmarshal(output, &getPods)
	checkErr(err)
	return getPods
}

func GetPodLogsStreamFastJson(podName string, insertChanJson chan []byte) {
	output, err := exec.Command("oc", "logs", podName, "--previous").CombinedOutput()
	if err == nil {
		for _, s := range strings.Split(string(output), "\n") {
			insertChanJson <- []byte(s)
		}
	}

	command := exec.Command("oc", "logs", "-f", "--since=200h", podName)
	pipe, err := command.StdoutPipe()
	command.Start()
	checkErr(err)
	reader := bufio.NewReader(pipe)

	var line []byte
	for {
		line, err = reader.ReadBytes(byte('\n'))
		if err != nil {
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
