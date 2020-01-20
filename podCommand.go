package main

import "github.com/jroimartin/gocui"

func podCommandKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("commands", gocui.KeyCtrlP, gocui.ModNone, activatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlP, gocui.ModNone, deactivatePodCommands); err != nil {
		return err
	}
	return nil
}

func activatePodCommands(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetViewOnTop("podCommand")
	checkErr(err)
	_, err = g.SetCurrentView("podCommand")
	checkErr(err)

	return nil
}
func deactivatePodCommands(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetViewOnBottom("podCommand")
	checkErr(err)
	_, err = g.SetCurrentView("commands")
	checkErr(err)

	return nil
}
