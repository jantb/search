package main

import (
	"fmt"
	"github.com/jantb/search/kafka"
	"github.com/jantb/search/kube"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

var podCommandY = 0
var podSearch = ""
var pods = kube.Pods{}
var selectedPods []kube.Items
var quitChans []chan bool

func podCommandKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("commands", gocui.KeyCtrlP, gocui.ModNone, activatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlP, gocui.ModNone, deactivatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyArrowDown, gocui.ModNone, podCommandsDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyArrowUp, gocui.ModNone, podCommandsUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyEnter, gocui.ModNone, podCommandsEnter); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlA, gocui.ModNone, podCommandsCTRLA); err != nil {
		return err
	}

	return nil
}

func activatePodCommands(g *gocui.Gui, v *gocui.View) error {
	podSearch = ""
	_, err := g.SetViewOnTop("podCommand")
	checkErr(err)
	v, err = g.SetCurrentView("podCommand")
	checkErr(err)
	v.Clear()
	podCommandY = 0
	v.SetCursor(0, 0)
	pods = kube.GetPods()
	selectedPods = pods.Items
	printPods(v)
	return nil
}

func deactivatePodCommands(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetViewOnBottom("podCommand")
	checkErr(err)
	_, err = g.SetCurrentView("commands")
	checkErr(err)
	podCommandY = 0
	return nil
}

func podCommandsDown(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	view, err := g.View("podCommand")
	checkErr(err)
	_, maxY := g.Size()
	if len(selectedPods)-2 < podCommandY {
		return nil
	}
	podCommandY++
	if y+2 > maxY {
		view.SetOrigin(0, podCommandY-maxY+1)
		return nil
	}

	if len(selectedPods) > y {
		v.SetCursor(strings.Index(selectedPods[y+1].Metadata.Name, podSearch)+len(podSearch), y+1)
	} else if len(selectedPods) > 0 {
		v.SetCursor(strings.Index(selectedPods[len(selectedPods)].Metadata.Name, podSearch)+len(podSearch), len(selectedPods)-1)
		podCommandY = len(selectedPods) - 1
	}
	return nil
}

func podCommandsUp(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	view, err := g.View("podCommand")
	checkErr(err)
	podCommandY--
	if y-1 < 0 {
		view.SetOrigin(0, podCommandY)
		if podCommandY < 0 {
			podCommandY = 0
		}
		return nil
	}
	if len(selectedPods) > y {
		v.SetCursor(strings.Index(selectedPods[y-1].Metadata.Name, podSearch)+len(podSearch), y-1)
	} else if len(selectedPods) > 0 {
		v.SetCursor(strings.Index(selectedPods[len(selectedPods)].Metadata.Name, podSearch)+len(podSearch), len(selectedPods)-1)
		podCommandY = len(selectedPods) - 1
	}
	return nil
}

func podCommandsEnter(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	if len(selectedPods) > y {
		name := selectedPods[podCommandY].Metadata.Name
		insertChanJson := make(chan []byte)
		quit := make(chan bool)
		quitChans = append(quitChans, quit)
		go func(insertChanJson chan []byte, podName string, quit chan bool) {
			kube.GetPodLogsStreamFastJson(podName, insertChanJson, quit)
		}(insertChanJson, name, quit)
		go insertIntoStoreJsonSystem(insertChanJson, name)
	}

	return nil
}

func podCommandsCTRLA(g *gocui.Gui, v *gocui.View) error {

	if len(selectedPods) > 0 {
		for i, _ := range selectedPods {
			insertChanJson := make(chan []byte)
			quit := make(chan bool)
			quitChans = append(quitChans, quit)
			podName := selectedPods[i].Metadata.Name
			go func(insertChanJson chan []byte, podName string, quit chan bool) {
				kube.GetPodLogsStreamFastJson(podName, insertChanJson, quit)
			}(insertChanJson, podName, quit)
			go insertIntoStoreJsonSystem(insertChanJson, podName)
		}
	}

	return nil
}
func demo(g *gocui.Gui, v *gocui.View) error {
	go kafka.KafkaRead(insertLogLinesChan)
	return nil
}

func editorPodCommand(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	v, e := gui.View("podCommand")
	checkErr(e)
	_, y := v.Cursor()
	v.Clear()
	switch key {
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		if len(podSearch) > 0 {
			podSearch = podSearch[:len(podSearch)-1]
		}
	case gocui.KeyDelete:

	case gocui.KeyEnter:
		return
	}

	if ch != 0 && mod == 0 {
		podSearch += string(ch)
	}

	selectedPods = nil
	for _, p := range pods.Items {
		if len(podSearch) == 0 || strings.Contains(p.Metadata.Name, podSearch) {
			selectedPods = append(selectedPods, p)
		}
	}
	if len(selectedPods) > y {
		v.SetCursor(strings.Index(selectedPods[y].Metadata.Name, podSearch)+len(podSearch), y)
	} else if len(selectedPods) > 0 {
		v.SetCursor(strings.Index(selectedPods[len(selectedPods)-1].Metadata.Name, podSearch)+len(podSearch), len(selectedPods)-1)
		podCommandY = len(selectedPods) - 1
	}

	printPods(v)
}

func printPods(v *gocui.View) {
	for _, item := range selectedPods {
		duration := time.Now().Sub(item.Metadata.CreationTimestamp)
		restartcount := 0
		reason := ""
		finishedAt := ""
		if len(item.Status.ContainerStatuses) > 0 {
			containerStatuses := item.Status.ContainerStatuses[len(item.Status.ContainerStatuses)-1]
			restartcount = containerStatuses.RestartCount
			reason = containerStatuses.LastState.Terminated.Reason
			finishedAt = containerStatuses.LastState.Terminated.FinishedAt.String()
			if containerStatuses.LastState.Terminated.FinishedAt.IsZero() {
				finishedAt = ""
			}
		}

		fmt.Fprintf(v, "%-70s %-10s %5d %40s %s\n", item.Metadata.Name, fmtDuration(duration), restartcount, reason, finishedAt)
	}
}
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
