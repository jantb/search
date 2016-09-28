package main

import (
	"github.com/boltdb/bolt"
	"log"
	"github.com/hpcloud/tail"
	"time"
	"encoding/json"
	"github.com/golang/leveldb/bloom"
	"crypto/md5"
	"strings"
	"os"
	"fmt"
)

type Event struct {
	Ts    time.Time
	Data  string
	Lines int
	Path  string
	Bloom bloom.Filter
}
type FileMonitor struct {
	Path   string
	Offset int64
	Poll   bool
}
type Meta struct {
	Count int64
}

func tailFile(fileMonitor FileMonitor) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true, ReOpen:true, Poll: fileMonitor.Poll, Logger:tail.DiscardingLogger, Location:&tail.SeekInfo{fileMonitor.Offset, os.SEEK_SET}})
	var key []byte
	formats := []string{"2006/01/02 15:04:05", "2006-01-02 15:04:05.000"}
	f := ""
	h := md5.New()
	prevo :=int64(0)
	stopo :=int64(0)
	for line := range t.Lines {
		var ok int
		var tt time.Time
		text := line.Text

		if f == "" {
			for _, format := range formats {
				if len(text) >= len(format) {
					_, err := time.Parse(format, text[:len(format)])
					if err != nil {
						continue
					}
					f = format
				}
			}
			if f== "" {
				fmt.Println("Unknown date format")
				fmt.Println(text)
				os.Exit(-1)
			}
		}
		if len(text) > len(f) {
			ti, err := time.Parse(f, text[:len(f)])
			if err != nil {
				ok = -1
			}
			if ok == 0 {
				ok = 1
				tt = ti
				text = text[len(f) + 1:]
				stopo = prevo
			}
		}
		o, err := t.Tell()
		if err != nil {
			log.Fatal(err)
		}
		prevo = o
		if ok == -1 || ok == 0 {
			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Events"))

				var event Event
				by := b.Get(key)
				err := json.Unmarshal(by, &event)
				if err != nil {
					log.Fatal(err)
				}
				event.Data += "\n" + text
				event.Bloom = bloom.NewFilter(nil, getBloomKeysFromLine(event.Data), 10)
				event.Lines = event.Lines + 1
				by, err = json.Marshal(event)
				if err != nil {
					log.Fatal(err)
				}

				b.Put(key, by)
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		filter := bloom.NewFilter(nil, getBloomKeysFromLine(text), 10)

		var event = Event{
			Ts: tt,
			Data: text,
			Bloom: filter,
			Path: fileMonitor.Path,
		}
		by, err := json.Marshal(event)
		if err != nil {
			log.Fatal(err)
		}

		h.Reset()
		key = []byte(tt.Format(time.RFC3339) + string(h.Sum([]byte(event.Data))))
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Events"))
			bt := b.Get(key)
			if bt != nil {
				var event Event
				err := json.Unmarshal(bt, &event)
				if err != nil {
					log.Fatal(err)
				}
				if event.Lines == 0 {
					return nil
				}
			}

			b.Put(key, by)

			fileMonitor.Offset = stopo
			b = tx.Bucket([]byte("Files"))
			by, err := json.Marshal(fileMonitor)
			if err != nil {
				log.Fatal(err)
			}
			b.Put([]byte(fileMonitor.Path), by)
			b = tx.Bucket([]byte("Meta"))
			by = b.Get([]byte("Meta"))
			if by == nil {
				b, _ :=json.Marshal(Meta{})
				by = b
			}
			meta := Meta{}
			json.Unmarshal(by, &meta)
			meta.Count++
			by, _ = json.Marshal(meta)
			b.Put([]byte("Meta"), by)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}
func SearchFor(t []byte, s int, seek int64, ch chan []Event, quit chan bool) {
	var events []Event
	count := 0
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		c := b.Cursor()
		k, v := c.Last()

		for ; k != nil && count < s; k, v = c.Prev() {
			select {
			case <-quit:
				return nil
			default:

				var event Event
				err := json.Unmarshal(v, &event)
				if err != nil {
					log.Fatal(err)
				}

				if len(t) == 0 {
					events = append(events, event)
					continue
				}

				keys := strings.Split(string(t), " ")
				add := true
				for _, key := range keys {
					if strings.TrimSpace(key) == "" {
						continue
					}
					if key[:1] == "!" {
						if event.Bloom.MayContain([]byte(key[1:])) {
							add = false
							break
						}
					} else {
						if !event.Bloom.MayContain([]byte(key)) || !strings.Contains(event.Data, key) {
							add = false
							continue
						}
					}
				}
				if add {
					if seek == int64(0) {
						count++
						events = append(events, event)
						continue
					}
					seek--
				}
			}
		}
		ch <- events
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}