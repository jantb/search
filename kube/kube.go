package kube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func GetPods() Pods {
	output, err := exec.Command("oc", "get", "pods", "-o", "json").CombinedOutput()
	checkErr(err)

	var getPods = Pods{}
	err = json.Unmarshal(output, &getPods)
	checkErr(err)
	return getPods
}

func GetPodLogsStreamFastJson(podName string, insertChanJson chan []byte, quit chan bool) {
	output, err := exec.Command("oc", "logs", podName, "--previous").CombinedOutput()
	if err == nil {
		for _, s := range strings.Split(string(output), "\n") {
			insertChanJson <- []byte(s)
		}
	}

	command := exec.Command("oc", "logs", "-f", "--since=200h", podName)
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	pipe, err := command.StdoutPipe()
	checkErr(err)
	err = command.Start()
	checkErr(err)
	reader := bufio.NewReader(pipe)

	var line []byte

	go func(quit chan bool, command *exec.Cmd) {
		select {
		case <-quit:
			_ = command.Process.Kill()
			_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
			return
		}
	}(quit, command)
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
