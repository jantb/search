package main

import (
	"fmt"
	"github.com/jantb/search/kube"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

var podCommandY = 0
var podSearch = ""
var pods = kube.Pods{}
var selectedPods []kube.Items

func podCommandKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("commands", gocui.KeyCtrlP, gocui.ModNone, activatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlP, gocui.ModNone, deactivatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlS, gocui.ModNone, deactivateSettings); err != nil {
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

	return nil
}

func activatePodCommands(g *gocui.Gui, v *gocui.View) error {
	podSearch = ""
	_, err := g.SetViewOnTop("podCommand")
	checkErr(err)
	v, err = g.SetCurrentView("podCommand")
	checkErr(err)
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
	} else {
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
	} else {
		v.SetCursor(strings.Index(selectedPods[len(selectedPods)].Metadata.Name, podSearch)+len(podSearch), len(selectedPods)-1)
		podCommandY = len(selectedPods) - 1
	}
	return nil
}

func podCommandsEnter(g *gocui.Gui, v *gocui.View) error {
	insertChanJson := make(chan map[string]interface{}, 10000)
	_, y := v.Cursor()
	if len(selectedPods) > y {

		go func(insertChanJson chan map[string]interface{}, podName string) {
			kube.GetPodLogsStream(podName, insertChanJson)
		}(insertChanJson, selectedPods[podCommandY].Metadata.Name)
		go insertIntoStoreJsonSystem(insertChanJson, selectedPods[podCommandY].Metadata.Name)
	}

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
	} else {
		v.SetCursor(strings.Index(selectedPods[len(selectedPods)-1].Metadata.Name, podSearch)+len(podSearch), len(selectedPods)-1)
		podCommandY = len(selectedPods) - 1
	}

	printPods(v)
}

func printPods(v *gocui.View) {
	for _, item := range selectedPods {
		duration := time.Now().Sub(item.Metadata.CreationTimestamp)
		fmt.Fprintf(v, "%s %s\n", item.Metadata.Name, fmtDuration(duration))
	}
}
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
