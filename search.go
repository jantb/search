package main

import (
	"github.com/boltdb/bolt"
	"log"
	"github.com/hpcloud/tail"
	"time"
	"github.com/golang/leveldb/bloom"
	"strings"
	"os"
	"encoding/json"
	"strconv"
)


func (e *Events) Get(ts string, data string) (*Event, bool) {
	for _, ev := range e.GetEvents() {
		if data == ev.Data && ts == ev.Ts {
			return ev, true
		}
	}
	return &Event{}, false
}
func (e *Events) RegenerateBloom() {
	keys := [][]byte{}
	for _, ev := range e.Events {
		keys = append(keys, getBloomKeysFromLine(ev.Data)...)
		keys = append(keys, getBloomKeysFromLine(ev.Path)...)
	}
	e.Bloom = bloom.NewFilter(nil, keys, 10)
}
func tailFile(fileMonitor FileMonitor) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true,
		ReOpen:true,
		Poll: fileMonitor.Poll,
		Logger:tail.DiscardingLogger,
		Location:&tail.SeekInfo{fileMonitor.Offset, os.SEEK_SET}})
	var key []byte
	var prevData string
	var prevTs string
	formats := []string{"2006/01/02 15:04:05",
		"2006-01-02 15:04:05.000",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano}
	f := ""
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

				by := b.Get(key)

				var events Events
				err := events.Unmarshal(by)
				if err != nil {
					log.Fatal(err)
				}

				event, _ := events.Get(prevTs, prevData)
				event.Data += "\n" + text
				prevData = event.Data
				keys := getBloomKeysFromLine(event.Data)
				keys = append(keys, getBloomKeysFromLine(fileMonitor.Path)...)
				fields := [][]byte{}
				for _, key := range keys {
					if strings.ContainsRune(string(key), '=') {
						split := strings.Split(string(key), "=")
						event.Fields = append(event.Fields, &Field{Key:split[0], Value:split[1]})
						fields = append(fields, []byte(split[0]))
						fields = append(fields, []byte(split[1]))
					}
				}

				keys = append(keys, fields...)
				event.Bloom = bloom.NewFilter(nil, keys, 10)
				event.Lines += 1
				by, err = events.Marshal()
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
			Ts: tt.Format(time.RFC3339),
			Data: text,
			Path: fileMonitor.Path,
		}
		fields := [][]byte{}
		for _, key := range keys {
			if strings.ContainsRune(string(key), '=') {
				split := strings.Split(string(key), "=")
				event.Fields = append(event.Fields, &Field{Key:split[0], Value:split[1]})
				fields = append(fields, []byte(split[0]))
				fields = append(fields, []byte(split[1]))
			}
		}
		keys = append(keys, fields...)
		event.Bloom = bloom.NewFilter(nil, keys, 10)

		key = []byte(tt.Truncate(1 * time.Minute).Format(time.RFC3339))
		prevData = event.Data
		prevTs = event.Ts
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Events"))
			eventsb := b.Get(key)
			if eventsb == nil {
				var events Events
				eventsb, _ = events.Marshal()
			}
			var events Events
			events.Unmarshal(eventsb)
			_, found := events.Get(event.Ts, event.Data)
			if found {
				return nil
			}
			events.Events = append(events.Events, &event)
			events.RegenerateBloom()
			by, err := events.Marshal()
			if err != nil {
				log.Panic(err)
			}
			b.Put(key, by)

			fileMonitor.Offset = stopo
			b = tx.Bucket([]byte("Files"))
			by, err = fileMonitor.Marshal()
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
	mutex2.Lock()
	defer mutex2.Unlock()
	ttt := time.Now()
	var eventsRet []Event
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

				var events Events
				err := events.Unmarshal(v)
				if err != nil {
					log.Fatal(err)
				}
				keys := strings.Split(string(t), " ")
				noInSet := false
				if len(t) != 0 {
					for _, key := range keys {
						if strings.Contains(key, "<") {
							split := strings.Split(key, "<")
							if !bloom.Filter(events.Bloom).MayContain([]byte(split[0])) {
								continue
							}
						}
						if strings.Contains(key, ">") {
							split := strings.Split(key, ">")
							if !bloom.Filter(events.Bloom).MayContain([]byte(split[0])) {
								continue
							}
						}
						if !bloom.Filter(events.Bloom).MayContain([]byte(key)) {
							noInSet = true
							break
						}
					}
					if noInSet {
						continue
					}
				}
				for i := len(events.Events) - 1; i >= 0; i-- {
					event := events.Events[i]
					if len(t) == 0 {
						if seek == int64(0) {
							count++
							eventsRet = append(eventsRet, *event)
							continue
						}
						seek--
						continue
					}

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
							for _, f := range event.Fields {
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

						} else if strings.Contains(key, ">") {
							split := strings.Split(key, ">")
							if !bloom.Filter(event.Bloom).MayContain([]byte(split[0])) {
								add = false
								continue
							}
							val := ""
							for _, f := range event.Fields {
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

						} else if key[:1] == "!" {
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
							eventsRet = append(eventsRet, *event)
							continue
						}
						seek--
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	edit_box.stats = time.Now().Sub(ttt)
	ch <- eventsRet
}