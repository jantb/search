package main

import "github.com/jroimartin/gocui"

var podCommandY = 0

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
	_, err = g.SetCurrentView("podCommand")
	checkErr(err)
	podCommandY = 0
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
