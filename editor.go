package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"golang.org/x/sync/semaphore"
)

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	logs, e := gui.View("logs")
	checkErr(e)
	ox, oy := logs.Origin()
	_, sy := logs.Size()
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
		v.MoveCursor(1, 0, false)
		return
	case gocui.KeyArrowUp:
		logs.SetOrigin(ox, oy-1)
		bottom.Store(false)
		if oy-1 <= 0 {
			renderSearch(v, 1)
		}
		return
	case gocui.KeyArrowDown:
		logs.SetOrigin(ox, Min(oy+1, len(logs.BufferLines())-sy))
		if oy == len(logs.BufferLines())-sy {
			renderSearch(v, -1)
		}
		return
	case gocui.KeyPgup:
		logs.SetOrigin(ox, oy-10)
		if oy-1 == 0 {
			renderSearch(v, 10)
		}
		renderSearch(v, 10)
		return
	case gocui.KeyPgdn:
		logs.SetOrigin(ox, Min(oy+10, len(logs.BufferLines())-sy))
		if oy == len(logs.BufferLines())-sy {
			renderSearch(v, -10)
		}
		renderSearch(v, -10)
		return
	case gocui.KeyEnter:
		renderSearch(v, math.MinInt32)
		return
	}
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
	}
}
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

var printBlue = color.New(color.FgBlue).Sprint
var printCyan = color.New(color.FgCyan).Sprint
var printRed = color.New(color.FgRed).Sprint
var printYellow = color.New(color.FgYellow).Sprint
var printGreen = color.New(color.FgGreen).Sprint
var printWhite = color.New(color.FgWhite).Sprint
var printYellowBack = color.New(color.BgYellow, color.Faint).Add(color.FgBlue).Sprint
var runeTL, runeTR, runeBL, runeBR = '┌', '┐', '└', '┘'
var runeH, runeV = '─', '│'
var renderSearchSemaphore = semaphore.NewWeighted(int64(1))

func renderSearch(v *gocui.View, offset int) {
	if renderSearchSemaphore.TryAcquire(1) {
		gui.Update(func(g *gocui.Gui) error {
			defer renderSearchSemaphore.Release(1)
			logs, e := gui.View("logs")
			checkErr(e)
			x, y := logs.Size()
			l, t := search(v.Buffer(), y, offset)
			logs.Clear()
			for _, value := range l {
				buffer := strings.TrimSpace(v.Buffer())
				levelFunc := printWhite
				switch value.getLevel() {
				case "ERROR":
					levelFunc = printRed
				case "WARN":
					levelFunc = printYellow
				case "INFO":
					levelFunc = printGreen
				case "DEBUG":
					levelFunc = printWhite
				}
				line := fmt.Sprintf("%s %s %s %s", printCyan(value.getTime().Format("2006-01-02T15:04:05")), printYellow(value.getSystem()), levelFunc(value.getLevel()), highlight(buffer, strings.TrimSpace(value.getBody())))
				lines := strings.Split(line, "\n")
				for _, value := range split([]rune(strings.TrimSpace(lines[0])), len(lines[0])-len(Strip(lines[0]))+x-1) {
					fmt.Fprintln(logs, string(value))
				}
				for _, line := range lines[1:] {
					for _, value := range split([]rune(strings.TrimSpace(line)), x-4) {
						for i := 0; i < 4; i++ {
							fmt.Fprint(logs, " ") // continuation
						}
						fmt.Fprintln(logs, string(value))
					}
				}
			}
			status, e := gui.View("status")
			checkErr(e)
			status.Clear()

			for i := 0; i < x-100; i++ {
				fmt.Fprint(status, " ")
			}
			_, sy := logs.Size()
			if bottom.Load() && len(l) > 0 {
				lastMessageDuration := time.Now().Sub(l[len(l)-1].getTime())
				logs.SetOrigin(0, len(logs.BufferLines())-sy)
				fmt.Fprintf(status, "┌─%10s──Follow mode, last message: %s ago──total lines: %d", t, fmt.Sprint(lastMessageDuration.Round(time.Second)), getLength())
			} else {
				fmt.Fprintf(status, "┌─%10s", t)
			}
			cx, _ := v.Cursor()
			for i := cx; i < x; i++ {
				fmt.Fprint(status, "─")
			}
			return nil
		})
	}
}

func split(buf []rune, lim int) [][]rune {
	var chunk []rune
	chunks := make([][]rune, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
func highlight(buffer string, line string) string {
	if len(buffer) > 0 {
		tokens := strings.Split(strings.TrimSpace(buffer), " ")
		for _, value := range tokens {
			line = strings.Replace(line, value, printYellowBack(value), -1)
		}
	}
	return line
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
