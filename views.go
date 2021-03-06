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
		v.Highlight = false
		v.Title = " Logs "
		v.Autoscroll = false
		v.Wrap = false
		v.Frame = false
	}
	return nil
}

// View: Logs
func viewPodCommand(g *gocui.Gui, maxX int, maxY int) error {
	if v, err := g.SetView("podCommand", -1, -1, maxX+1, maxY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelFgColor = gocui.ColorCyan
		v.Autoscroll = false
		v.Wrap = false
		v.Frame = false
		v.Editable = true
		v.Editor = gocui.EditorFunc(editorPodCommand)
		v.SetCursor(0, 0)
	}
	return nil
}

// View: Status
func viewStatus(g *gocui.Gui, maxX int, maxY int) error {
	if v, err := g.SetView("status", -1, maxY-3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Frame = false
		v.FgColor = gocui.ColorCyan
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
		v.Frame = false

		v.Editor = gocui.EditorFunc(editor)
		bottom.Store(true)
		if _, err := g.SetCurrentView("commands"); err != nil {
			return err
		}
	}
	return nil
}

// View: Prompt
func viewPrompt(g *gocui.Gui, maxX int, maxY int) error {

	if v, err := g.SetView("prompt", (maxX/2)-100, (maxY / 2), (maxX/2)+100, (maxY/2)+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = false
		v.Frame = true
		g.SetViewOnBottom("prompt")
	}
	return nil
}
