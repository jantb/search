package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
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

func renderSearch(v *gocui.View, offset int) {
	gui.Update(func(g *gocui.Gui) error {
		view, e := gui.View("logs")
		checkErr(e)
		_, y := view.Size()
		logLinesPrev = s(v.Buffer(), y, offset, logLinesPrev)
		view.Clear()
		for _, value := range logLinesPrev {
			fmt.Fprintf(view, "%s %s\n", value.getTime(), value.Body)
		}
		return nil
	})
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
