package main

import (
	"github.com/boltdb/bolt"
	"log"
	"github.com/hpcloud/tail"
	"time"
	"github.com/golang/leveldb/bloom"
	"crypto/md5"
	"strings"
	"os"
	"encoding/json"
	"strconv"
)

type Meta struct {
	Count int64
}

func tailFile(fileMonitor FileMonitor) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true, ReOpen:true, Poll: fileMonitor.Poll, Logger:tail.DiscardingLogger, Location:&tail.SeekInfo{fileMonitor.Offset, os.SEEK_SET}})
	var key []byte
	formats := []string{"2006/01/02 15:04:05", "2006-01-02 15:04:05.000", time.ANSIC,time.UnixDate,time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z , time.RFC3339, time.RFC3339Nano}
	f := ""
	h := md5.New()
	prevo := int64(0)
	stopo := int64(0)
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
			if f == "" {
				continue
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
				_, err := event.Unmarshal(by)
				if err != nil {
					log.Fatal(err)
				}
				event.Data += "\n" + text
				keys := getBloomKeysFromLine(event.Data)
				keys = append(keys, getBloomKeysFromLine(fileMonitor.Path)...)
				event.Fields=[]Field{}
				fields :=[][]byte{}
				for _, key := range keys {
					if strings.ContainsRune(string(key), '=') {
						split := strings.Split(string(key), "=")
						event.Fields = append(event.Fields, Field{Key:split[0],Value:split[1]})
						fields = append(fields, []byte(split[0]))
						fields = append(fields, []byte(split[1]))
					}
				}

				keys = append(keys, fields...)
				event.Bloom = bloom.NewFilter(nil, keys, 10)
				event.Lines = event.Lines + 1
				by, err = event.Marshal(nil)
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
		keys := getBloomKeysFromLine(text)
		keys = append(keys, getBloomKeysFromLine(fileMonitor.Path)...)


		var event = Event{
			Ts: tt,
			Data: text,
			Path: fileMonitor.Path,
		}
		fields :=[][]byte{}
		for _, key := range keys {
			if strings.ContainsRune(string(key), '=') {
				split := strings.Split(string(key), "=")
				event.Fields = append(event.Fields, Field{Key:split[0],Value:split[1]})
				fields = append(fields, []byte(split[0]))
				fields = append(fields, []byte(split[1]))
			}
		}
		keys = append(keys, fields...)
		event.Bloom = bloom.NewFilter(nil, keys, 10)

		by, err := event.Marshal(nil)
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
				_, err := event.Unmarshal(bt)
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
			by, err := fileMonitor.Marshal(nil)
			if err != nil {
				log.Fatal(err)
			}
			b.Put([]byte(fileMonitor.Path), by)
			b = tx.Bucket([]byte("Meta"))
			by = b.Get([]byte("Meta"))
			if by == nil {
				b, _ := json.Marshal(Meta{})
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
				_, err := event.Unmarshal(v)
				if err != nil {
					log.Fatal(err)
				}
				if len(t) == 0 {
					if seek == int64(0) {
						count++
						events = append(events, event)
						continue
					}
					seek--
					continue
				}

				keys := strings.Split(string(t), " ")
				add := true
				for _, key := range keys {
					if strings.TrimSpace(key) == "" {
						continue
					}

					if strings.Contains(key, "<") {
						split := strings.Split(key, "<")
						if !bloom.Filter(event.Bloom).MayContain([]byte(split[0])) {
							add = false
							continue
						}
						val := ""
						for _, f :=range event.Fields  {
							if split[0] == f.Key {
								val = f.Value
							}
						}
						i, err := strconv.Atoi(split[1])
						if err != nil {
							add = false
							continue
						}
						i2, err := strconv.Atoi(val)
						if err != nil {
							add = false
							continue
						}
						if i2 >= i {
							add = false
							continue
						}

					}else if strings.Contains(key, ">") {
						split := strings.Split(key, ">")
						if !bloom.Filter(event.Bloom).MayContain([]byte(split[0])) {
							add = false
							continue
						}
						val := ""
						for _, f :=range event.Fields  {
							if split[0] == f.Key {
								val = f.Value
							}
						}
						i, err := strconv.Atoi(split[1])
						if err != nil {
							add = false
							continue
						}
						i2, err := strconv.Atoi(val)
						if err != nil {
							add = false
							continue
						}
						if i2 <= i {
							add = false
							continue
						}

					}else if key[:1] == "!" {
						if bloom.Filter(event.Bloom).MayContain([]byte(key[1:])) {
							add = false
							break
						}
					} else {
						if !bloom.Filter(event.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.Data, key) || strings.Contains(event.Path, key)) {
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

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	ch <- events
}