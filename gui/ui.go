package gui

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jantb/search/proto"
	"github.com/jantb/search/searchfor"
	"github.com/jroimartin/gocui"
)

var db *bolt.DB
var g *gocui.Gui

func Run(d *bolt.DB) {
	db = d
	gg, err := gocui.NewGui(gocui.Output256)
	g = gg
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.Cursor = true
	tail.Store(false)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlG, gocui.ModNone, tailf); err != nil {
		log.Panicln(err)
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		tickerTail := time.NewTicker(50 * time.Millisecond)
		var count = atomic.Value{}
		count.Store(int64(0))
		var ts = atomic.Value{}
		ts.Store("")
		var b = []byte{}
		for {
			select {
			case <-tickerTail.C:
				if tail.Load().(bool) {
					v, _ := g.View("edit")
					vMain, _ := g.View("main")
					_, y := vMain.Size()
					buffer := v.Buffer()
					go searchfor.SearchFor([]byte(buffer), y/2, 0, resChan, db)
				}
			case <-ticker.C:
				g.Execute(func(g *gocui.Gui) error {
					v, _ := g.View("edit")

					nodecount := int64(0)
					var meta proto.Meta
					meta.Retrieve(db)
					nodecount = meta.Count

					title := ""
					if searchfor.Searching.Load() != nil && searchfor.Searching.Load().(bool) {
						title += "searching "
					}
					c := count.Load().(int64)
					if c != 0 {
						title += fmt.Sprintf("count:%d ", c)
					}
					tS := ts.Load().(string)
					title += fmt.Sprintf("%s %d", tS, nodecount)
					v.Title = title
					return nil
				})

			case b = <-resChan:
				data := make([]byte, len(b))
				copy(data, b)

				g.Execute(func(g *gocui.Gui) error {
					var res = proto.SearchRes{}
					err = res.Unmarshal(data)
					if err != nil {
						log.Panic(err)
					}
					v, err := g.View("main")
					if err != nil {
						return err
					}
					ts.Store(res.Ts)
					count.Store(res.Count)
					v.Clear()
					for i := len(res.Events) - 1; i >= 0; i-- {
						event := res.Events[i]
						fmt.Fprintf(v, "\033[38;5;87m%s\033[0m ", event.Ts)
						for i, r := range event.Data {
							s := string(r)
							if len(s) > 20000 {
								s = s[:20000]
							}
							found := false
							for h := 0; h < len(event.FoundAtIndex); h += 2 {
								if int32(i) >= event.FoundAtIndex[h] && int32(i) < event.FoundAtIndex[h+1] {
									fmt.Fprintf(v, "\033[48;5;1m%s\033[0m", s)
									found = true
								}
							}
							if !found {
								fmt.Fprintf(v, "%s", s)
							}
						}
						fmt.Fprintf(v, "\n\033[38;5;8msource:%s\033[0m\n", event.Path)
					}

					return nil
				})
			}
		}
	}()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("edit", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Editor = gocui.EditorFunc(editor)
		if _, err := g.SetCurrentView("edit"); err != nil {
			return err
		}
		v.Autoscroll = false
		v.SetCursor(0, 0)
		v.Wrap = false

	}
	if v, err := g.SetView("main", -1, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		v.Editable = false
		v.Autoscroll = true
	}
	return nil
}

var resChan = make(chan []byte)
var origin = 0
var skipItems = int64(0)
var tail = atomic.Value{}

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 50, 0, resChan, db)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 50, 0, resChan, db)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 50, 0, resChan, db)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 50, 0, resChan, db)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		//v.EditNewLine()
	case key == gocui.KeyArrowDown:
		g.Execute(func(g *gocui.Gui) error {
			vm, err := g.View("main")
			vm.Autoscroll = false
			_, y := vm.Origin()
			origin = y
			if err != nil {
				return err
			}
			origin++
			err = vm.SetOrigin(0, origin)
			if err != nil {
				origin--
			}
			if origin >= len(strings.Split(vm.Buffer(), "\n")) {
				skipItems--
				if skipItems == -1 {
					skipItems++
				}
				go searchfor.SearchFor([]byte(v.Buffer()), 50, skipItems, resChan, db)
			}
			return nil
		})

	case key == gocui.KeyArrowUp:
		g.Execute(func(g *gocui.Gui) error {
			vm, err := g.View("main")
			vm.Autoscroll = false
			_, o := vm.Origin()
			origin = o
			if err != nil {
				return err
			}
			origin--
			err = vm.SetOrigin(0, origin)
			if err != nil {
				origin++
			}
			if origin == 0 {
				skipItems++
			}
			return nil
		})

	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
func reset(g *gocui.Gui) error {
	v, _ := g.View("main")
	v.SetOrigin(0, 0)
	v.Autoscroll = true
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func tailf(g *gocui.Gui, v *gocui.View) error {
	tail.Store(!tail.Load().(bool))
	return nil
}
