package gui

import (
	"fmt"
	"log"

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

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	go func() {
		for {
			select {
			case res := <-resChan:

				g.Execute(func(g *gocui.Gui) error {
					v, err := g.View("main")
					if err != nil {
						return err
					}
					v.Clear()
					for i := len(res.Events) - 1; i >= 0; i-- {
						event := res.Events[i]
						fmt.Fprintf(v, "\033[38;5;87m%s\033[0m ", event.Ts)
						for i, r := range event.Data {
							found := false
							for h := 0; h < len(event.FoundAtIndex); h += 2 {
								if int32(i) >= event.FoundAtIndex[h] && int32(i) < event.FoundAtIndex[h+1] {
									fmt.Fprintf(v, "\033[48;5;1m%s\033[0m", string(r))
									found = true
								}
							}
							if !found {
								fmt.Fprintf(v, "%s", string(r))
							}
						}
						fmt.Fprintf(v, "\n\033[38;5;8msource:%s\033[0m\n", event.Path)
					}
					v.Size()

					return nil
				})

				g.Execute(func(g *gocui.Gui) error {
					v, err := g.View("edit")
					if err != nil {
						return err
					}

					nodecount := int64(0)

					db.View(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("Meta"))
						by := b.Get([]byte("Meta"))
						var meta proto.Meta
						meta.Unmarshal(by)
						nodecount = meta.Count
						return nil
					})
					title := ""
					if res.Count != 0 {
						title = fmt.Sprintf("count:%d ", res.Count)
					}
					title += fmt.Sprintf("%s %d/%d", res.Ts, len(res.Events), nodecount)
					v.Title = title
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

var resChan = make(chan proto.SearchRes)
var quitChan = make(chan bool)
var origin = 0

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 1000, 0, resChan, quitChan, db)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 1000, 0, resChan, quitChan, db)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 1000, 0, resChan, quitChan, db)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		g.Execute(reset)
		go searchfor.SearchFor([]byte(v.Buffer()), 1000, 0, resChan, quitChan, db)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		//v.EditNewLine()
	case key == gocui.KeyArrowDown:
		g.Execute(func(g *gocui.Gui) error {
			v, err := g.View("main")
			v.Autoscroll = false
			_, y := v.Origin()
			origin = y
			if err != nil {
				return err
			}
			origin++
			err = v.SetOrigin(0, origin)
			if err != nil {
				origin--
			}

			return nil
		})

	case key == gocui.KeyArrowUp:
		g.Execute(func(g *gocui.Gui) error {
			v, err := g.View("main")
			v.Autoscroll = false
			_, y := v.Origin()
			origin = y
			if err != nil {
				return err
			}
			origin--
			err = v.SetOrigin(0, origin)
			if err != nil {
				origin++
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
