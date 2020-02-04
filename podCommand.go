package main

import (
	"fmt"
	"github.com/jantb/search/kube"
	"github.com/jroimartin/gocui"
	"strings"
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
	_, err := g.SetViewOnTop("podCommand")
	checkErr(err)
	v, err = g.SetCurrentView("podCommand")
	checkErr(err)
	podCommandY = 0
	pods = kube.GetPods()
	selectedPods = pods.Items
	for _, item := range selectedPods {
		fmt.Fprintln(v, item.Metadata.Name)
	}
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
	x, y := v.Cursor()
	view, err := g.View("podCommand")
	checkErr(err)
	_, maxY := g.Size()
	podCommandY++
	if y+2 > maxY {
		view.SetOrigin(0, podCommandY-maxY+1)
		return nil
	}
	v.SetCursor(x, y+1)
	return nil
}

func podCommandsUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
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
	v.SetCursor(x, y-1)
	return nil
}

func podCommandsEnter(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func editorPodCommand(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	v, e := gui.View("podCommand")
	checkErr(e)
	x, y := v.Cursor()
	v.Clear()
	switch key {
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		if len(podSearch) > 0 {
			podSearch = podSearch[:len(podSearch)-1]
			v.SetCursor(x-1, y)
		}
	case gocui.KeyDelete:

	case gocui.KeyEnter:
		return
	}

	if ch != 0 && mod == 0 {
		podSearch += string(ch)
		v.SetCursor(x+1, y)
	}

	selectedPods = nil
	for _, p := range pods.Items {
		if strings.HasPrefix(p.Metadata.Name, podSearch) {
			selectedPods = append(selectedPods, p)
		}
	}

	for _, item := range selectedPods {
		fmt.Fprintln(v, item.Metadata.Name)
	}
}
