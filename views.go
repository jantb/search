package main

import (
	"github.com/jroimartin/gocui"
)

// View: Logs
func viewLogs(g *gocui.Gui, maxX int, maxY int) error {
	if v, err := g.SetView("logs", -1, -1, maxX+1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.Title = " Logs "
		v.Autoscroll = true
		v.Frame = false
	}
	return nil
}

// View: Commands
func viewCommands(g *gocui.Gui, maxX int, maxY int) error {

	if v, err := g.SetView("commands", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = false
		v.Frame = true

		v.Editor = gocui.EditorFunc(editor)
		if _, err := g.SetCurrentView("commands"); err != nil {
			return err
		}
	}
	return nil
}
