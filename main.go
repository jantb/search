package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"go.uber.org/atomic"
	"log"
	"os"
)

var formats Formats
var gui *gocui.Gui
var bottom atomic.Bool

func main() {
	bottom.Store(false)

	initStore()

	g, err := gocui.NewGui(gocui.Output256)
	checkErr(err)
	gui = g
	defer g.Close()
	insertChan := make(chan string, 10000)
	insertChanJson := make(chan map[string]interface{}, 10000)
	go insertIntoStore(insertChan)
	go insertIntoStoreJson(insertChanJson)
	go readFromPipe(insertChan, insertChanJson)

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

func clean(g *gocui.Gui, v *gocui.View) error {
	clearDb()
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	err := podCommandKeybindings(g)
	checkErr(err)
	err = settingsKeybindings(g)
	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, clean); err != nil {
		return err
	}
	checkErr(err)
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	err := viewPodCommand(g, maxX, maxY)
	checkErr(err)
	err = viewSettings(g, maxX, maxY)
	checkErr(err)
	err = viewLogs(g, maxX, maxY)
	checkErr(err)
	err = viewStatus(g, maxX, maxY)
	checkErr(err)
	err = viewCommands(g, maxX, maxY)
	checkErr(err)
	err = viewPrompt(g, maxX, maxY)
	checkErr(err)
	return nil
}

func checkErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
