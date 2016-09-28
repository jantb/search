package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os/user"
	"flag"
	"os"
	"path/filepath"
	"github.com/nsf/termbox-go"
	"encoding/json"
	"time"
)

var filename = flag.String("add", "", "Filename to monitor")
var poll = flag.Bool("poll", false, "use poll")
var db *bolt.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := bolt.Open(filepath.Join(usr.HomeDir, ".search.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = dbs
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Events"))
		tx.CreateBucketIfNotExists([]byte("Files"))
		tx.CreateBucketIfNotExists([]byte("Meta"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}

		if !fi.IsDir() {
			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Files"))
				dir, _ := filepath.Abs(filepath.Dir(*filename))
				filep := filepath.Join(dir, filepath.Base(*filename))
				fileMonitor := FileMonitor{
					Path:filep,
					Offset:0,
					Poll: *poll,
				}
				by, err := json.Marshal(fileMonitor)
				if err != nil {
					log.Fatal(err)
				}
				b.Put([]byte(filep), by)
				return nil
			})
			return
		}
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		c := b.Cursor()
		for k, f := c.First(); k != nil; k, f = c.Next() {
			fileMonitor := FileMonitor{}
			json.Unmarshal(f, &fileMonitor)
			if err != nil {
				log.Fatal(err)
			}
			go tailFile(fileMonitor)
		}
		return nil
	})

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	edit_box.eventChan = make(chan []Event)
	edit_box.quitSearch = make(chan bool)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if edit_box.seek == int64(0) {
				edit_box.Search( )
			}
			redraw_all()
		}
	}()
	go func() {
		for {
			edit_box.events = <-edit_box.eventChan
			redraw_all()
		}
	}()

	edit_box.Search()
	mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				edit_box.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				edit_box.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				edit_box.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				edit_box.DeleteRuneForward()
			case termbox.KeyTab:
				edit_box.InsertRune('\t')
			case termbox.KeyArrowUp:
				edit_box.ScrollUp();
			case termbox.KeyArrowDown:
				edit_box.ScrollDown();
			case termbox.KeySpace:
				edit_box.InsertRune(' ')
			case termbox.KeyCtrlG:
				edit_box.Follow()
			case termbox.KeyCtrlK:
				edit_box.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				edit_box.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				edit_box.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					edit_box.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}

}

