package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os/user"
	"flag"
	"os"
	"path/filepath"
	"time"
	"syscall"
	"sync/atomic"
	"os/signal"
	"github.com/gdamore/tcell/termbox"
	"github.com/jantb/search/proto"
	"github.com/jantb/search/tail"
	"golang.org/x/net/context"
	"github.com/jantb/search/searchfor"
	"net"
	"google.golang.org/grpc"
)

var filename = flag.String("add", "", "Filename to monitor")
var poll = flag.Bool("poll", false, "use poll")
var db *bolt.DB

const (
	port = ":50051"
)
type server struct{}

func (s *server) Process(ctx context.Context, in *proto.SearchConf) (*proto.SearchRes, error) {
	channel := make(chan proto.SearchRes)
	quitChan := make(chan bool)
	searchfor.SearchFor(in.Text, in.Size_, in.Skipped, quitChan, channel,db)
	return <-channel, nil
}

func main() {

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		proto.RegisterSearchServer(s, &server{})
		s.Serve(lis)
	}()
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
	searchChan := make(chan bool)
	go func(e *EditBox) {
		for {
			time.Sleep(time.Millisecond * 100)
			if atomic.LoadInt64(&edit_box.seek) == int64(0) {
				searchChan <- true
			}
		}
	}(edit_box)

	eventChan := make(chan termbox.Event)
	go func() {
		for {
			event := termbox.PollEvent()
			eventChan <- event
		}
	}()
	// register signals to channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	mainloop:
	for {
		select {
		case ev := <-eventChan:
			switch  ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlC:
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
		case searchRes := <-edit_box.eventChan:
			edit_box.count = searchRes.Count
			edit_box.stats = searchRes.Ts
			edit_box.events = searchRes.Events
			redraw_all(edit_box)
		case <-searchChan:
			edit_box.Search()
		case <-sigChan:
			break mainloop
		}
	}

}

