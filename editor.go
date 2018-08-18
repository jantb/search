package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"golang.org/x/sync/semaphore"
	"time"
)

var logLinesPrev []LogLine

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {

	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
		moveAhead(v)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
		moveAhead(v)
	case gocui.KeyArrowRight:
		x, _ := v.Cursor()
		x2, _ := v.Origin()
		x += x2
		buf := v.Buffer()
		if buf != "" && len(buf) > (x+2) {
			v.MoveCursor(1, 0, false)
		}
		return
	case gocui.KeyArrowUp:
		renderSearch(v, 1)
		return
	case gocui.KeyArrowDown:
		renderSearch(v, -1)
		return
	case gocui.KeyPgup:
		renderSearch(v, 10)
		return
	case gocui.KeyPgdn:
		renderSearch(v, -10)
		return
	case gocui.KeyEnter:
		renderSearch(v, 0)
		return
	}
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
	}
}

var printBlue = color.New(color.FgBlue).Sprint
var printRed = color.New(color.FgRed).Sprint
var runeTL, runeTR, runeBL, runeBR = '┌', '┐', '└', '┘'
var runeH, runeV = '─', '│'
var renderSearchSemaphore = semaphore.NewWeighted(int64(1))

func renderSearch(v *gocui.View, offset int) {
	if renderSearchSemaphore.TryAcquire(1) {
		gui.Update(func(g *gocui.Gui) error {
			defer renderSearchSemaphore.Release(1)
			view, e := gui.View("logs")
			checkErr(e)
			x, y := view.Size()
			l, t := s(v.Buffer(), y, offset, logLinesPrev)
			logLinesPrev = l
			view.Clear()
			for _, value := range logLinesPrev {
				fmt.Fprintf(view, "%s %s %s\n", printBlue(value.getTime().Format("2006-01-02T15:04:05.999")), printRed(value.Level), value.Body)
			}
			view, e = gui.View("status")
			checkErr(e)
			view.Clear()

			for i := 0; i < x-100; i++ {
				fmt.Fprint(view, " ")
			}
			if bottom.Load() && len(logLinesPrev) > 0 {
				lastMessageDuration := time.Now().Sub(logLinesPrev[len(logLinesPrev)-1].getTime())
				fmt.Fprintf(view, "┌─%s──Follow mode, last message: %s ago──total lines: %d", t, fmt.Sprint(lastMessageDuration.Round(time.Second)), logLinesPrev[len(logLinesPrev)-1].Id)
			} else {
				fmt.Fprintf(view, "┌─%s──", t)
			}
			cx, _ := v.Cursor()
			for i := cx; i < x; i++ {
				fmt.Fprint(view, "─")
			}
			return nil
		})
	}
}

func moveAhead(v *gocui.View) {
	cX, _ := v.Cursor()
	oX, _ := v.Origin()
	if cX < 10 && oX > 0 {
		newOX := oX - 10
		forward := 10
		if newOX < 0 {
			forward += newOX
			newOX = 0
		}
		v.SetOrigin(newOX, 0)
		v.MoveCursor(forward, 0, false)
	}
}
