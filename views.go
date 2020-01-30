package main

import (
	"encoding/json"
	"fmt"
	"github.com/jantb/search/kube"
	"github.com/jroimartin/gocui"
	"os/exec"
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
		v.Editable = false

		fmt.Fprintln(v, "Pods")
		output, err := exec.Command("kubectl", "get", "pods", "-o", "json").CombinedOutput()
		checkErr(err)

		var getPods = kube.GetPods{}
		err = json.Unmarshal(output, &getPods)
		checkErr(err)
		for _, item := range getPods.Items {
			fmt.Fprintln(v, item.Metadata.Name)
		}

		v.SetCursor(0, 0)
	}
	return nil
}

// View: Settings
func viewSettings(g *gocui.Gui, maxX int, maxY int) error {
	if v, err := g.SetView("settings", 0, 0, 120, 7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Settings"
		v.Highlight = true
		v.SelFgColor = gocui.ColorYellow
		v.Autoscroll = false
		v.Wrap = false
		v.Frame = true
		v.Editable = false

		fmt.Fprintln(v, "Username       ", loadSettings("username"))
		fmt.Fprintln(v, "Password       ", "********")
		fmt.Fprintln(v, "Openshift url  ", loadSettings("openshift"))
		fmt.Fprintln(v, "Jenkins url    ", loadSettings("jenkins"))
		fmt.Fprintln(v, "Bitbucket url  ", loadSettings("bitbucket"))
		fmt.Fprintln(v, "Jira url       ", loadSettings("jira"))
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
