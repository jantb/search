package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os/user"
	"flag"
	"os"
	"path/filepath"
	"github.com/nsf/termbox-go"
	"time"
	"syscall"
	"search/proto"
	"search/tail"
	"sync/atomic"
)

var filename = flag.String("add", "", "Filename to monitor")
var poll = flag.Bool("poll", false, "use poll")
var db *bolt.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, _ := os.OpenFile("x.txt", os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
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
				fileMonitor := proto.FileMonitor{
					Path:filep,
					Offset:0,
					Poll: *poll,
				}
				by, err := fileMonitor.Marshal()
				if err != nil {
					log.Fatal(err)
				}
				b.Put([]byte(filep), by)
				return nil
			})
			return
		}
	}

	edit_box := New()
	edit_box.eventChan = make(chan proto.SearchRes)
	edit_box.quitSearch = make(chan bool)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		c := b.Cursor()
		for k, f := c.First(); k != nil; k, f = c.Next() {
			fileMonitor := proto.FileMonitor{}
			fileMonitor.Unmarshal(f)
			if err != nil {
				log.Fatal(err)
			}
			go tail.TailFile(fileMonitor, db)
		}
		return nil
	})

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	go func(e *EditBox) {
		for {
			time.Sleep(time.Millisecond * 1000)
			if atomic.LoadInt64(&edit_box.seek) == int64(0) {
				e.Search()
			}
		}
	}(edit_box)
	go func(edit_box *EditBox) {
		var searchRes proto.SearchRes
		for {
			searchRes = <-edit_box.eventChan
			edit_box.mtx.Lock()
			edit_box.count = searchRes.Count
			edit_box.stats = searchRes.Ts
			edit_box.events = searchRes.Events
			edit_box.mtx.Unlock()
			redraw_all(*edit_box)
		}
	}(edit_box)

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
			case termbox.KeyPgup:
				edit_box.ScrollUp();
			case termbox.KeyPgdn:
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
	}
}

