package main

import (
	"github.com/jroimartin/gocui"
	"go.uber.org/atomic"
	"log"
)

var formats Formats
var gui *gocui.Gui
var bottom atomic.Bool

func main() {
	bottom.Store(false)

	initStore()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	defer g.Close()
	insertChan := make(chan string, 10000)
	go insertIntoStore(insertChan)
	go readFromPipe(insertChan)

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	cleanupStore()
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	podCommandKeybindings(g)
	settingsKeybindings(g)
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewPodCommand(g, maxX, maxY)
	viewSettings(g, maxX, maxY)
	viewLogs(g, maxX, maxY)
	viewStatus(g, maxX, maxY)
	viewCommands(g, maxX, maxY)
	viewPrompt(g, maxX, maxY)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
