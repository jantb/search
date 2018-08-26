package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func activateSettings(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("settings")
	g.SetCurrentView("settings")

	return nil
}
func deactivateSettings(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("settings")
	g.SetCurrentView("commands")

	return nil
}

func settingsDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	v.SetCursor(x, y-1)
	return nil
}

func settingsUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	v.SetCursor(x, y+1)
	return nil
}

func settingsEnter(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("prompt")
	g.SetCurrentView("prompt")
	view, err := g.View("prompt")
	checkErr(err)
	view.Clear()
	view.SetCursor(0, 0)
	return nil
}

func enterClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("prompt")
	g.SetCurrentView("settings")

	settings, err := g.View("settings")
	checkErr(err)
	_, settingsY := settings.Cursor()

	switch settingsY {
	case 0:
		storeSettings("username", v.Buffer())
	case 1:
		storeSettings("password", v.Buffer())
	case 2:
		storeSettings("openshift", v.Buffer())
	case 3:
		storeSettings("jenkins", v.Buffer())
	case 4:
		storeSettings("bitbucket", v.Buffer())
	case 5:
		storeSettings("jira", v.Buffer())
	default:
	}

	//buffer := v.Buffer()
	//fmt.Print(buffer)
	settings.Clear()
	v.SetCursor(0, 0)

	fmt.Fprintln(settings, "Username       ", loadSettings("username"))
	fmt.Fprintln(settings, "Password       ", "********")
	fmt.Fprintln(settings, "Openshift url  ", loadSettings("openshift"))
	fmt.Fprintln(settings, "Jenkins url    ", loadSettings("jenkins"))
	fmt.Fprintln(settings, "Bitbucket url  ", loadSettings("bitbucket"))
	fmt.Fprintln(settings, "Jira url       ", loadSettings("jira"))
	return nil
}

func settingsKeybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, activateSettings); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyCtrlS, gocui.ModNone, deactivateSettings); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyArrowDown, gocui.ModNone, settingsUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyArrowUp, gocui.ModNone, settingsDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyEnter, gocui.ModNone, settingsEnter); err != nil {
		return err
	}
	if err := g.SetKeybinding("prompt", gocui.KeyEnter, gocui.ModNone, enterClose); err != nil {
		return err
	}
	return nil
}
