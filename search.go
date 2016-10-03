package main

import (
	"github.com/boltdb/bolt"
	"log"
	"github.com/hpcloud/tail"
	"time"
	"github.com/golang/leveldb/bloom"
	"strings"
	"os"
	"strconv"
	"github.com/bradfitz/slice"
	"encoding/binary"
)

func (e *Events) Get(ts string, data string) (*Event, bool) {
	for _, ev := range e.GetEvents() {
		if data == ev.Data && ts == ev.Ts {
			return ev, true
		}
	}
	return &Event{}, false
}
func (e *Events) sortEvents() {
	slice.Sort(e.Events, func(i, j int) bool {
		return e.Events[i].Ts < e.Events[j].Ts
	})
}
func (e *Events) RegenerateBloom() {
	set := make(map[string]bool)

	for _, ev := range e.Events {
		for _, key := range getBloomKeysFromLine(ev.Data) {
			set[string(key)] = true
		}
		for _, key := range getBloomKeysFromLine(ev.Path) {
			set[string(key)] = true
		}
	}
	keys := make([][]byte, 0, len(set))
	for k := range set {
		keys = append(keys, []byte(k))
	}

	e.Bloom = bloom.NewFilter(nil, keys, 10)
	e.BloomDirty = false
}
func (e *Event) GenerateBloom() {
	if e.BloomDirty {
		set := make(map[string]bool)
		for _, key := range getBloomKeysFromLine(e.Data) {
			set[string(key)] = true
			if strings.ContainsRune(string(key), '=') {
				split := strings.Split(string(key), "=")
				set[string(split[0])] = true
				set[string(split[1])] = true
			}
		}
		for _, key := range getBloomKeysFromLine(e.Path) {
			set[string(key)] = true
		}
		keys := make([][]byte, 0, len(set))
		for k := range set {
			keys = append(keys, []byte(k))
		}
		e.Bloom = bloom.NewFilter(nil, keys, 10)
		e.BloomDirty = false
	}
}

func (event *Event) shouldAddAndGetIndexes(keys []string) (bool) {
	add := true
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if event.BloomDirty {
			if key[:1] == "!" {
				if strings.Contains(event.Data, key[1:]) {
					add = false
					break
				}
			} else {
				if !(strings.Contains(event.Data, key) || strings.Contains(event.Path, key)) {
					add = false
					continue
				}
			}
		} else if strings.Contains(key, "<") {
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
		} else if !bloom.Filter(event.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.Data, key) || strings.Contains(event.Path, key)) {
			add = false
			continue
		}
	}
	return add
}

func (event Event) GetKeyIndexes(keys []string) []int32 {
	var keyIndexes []int32
	for _, key := range keys {
		if key == "" {
			continue
		}
		index := strings.Index(event.Data, key)
		text := event.Data
		indexPrev := 0
		for ; index != -1; index = strings.Index(text[indexPrev:], key) {
			index += indexPrev
			keyIndexes = append(keyIndexes, int32(index))
			index += len(key)
			keyIndexes = append(keyIndexes, int32(index))
			indexPrev = index
		}
	}
	return keyIndexes
}

var formats = []string{"2006/01/02 15:04:05",
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

func findFormat(text string) string {
	for _, format := range formats {
		if len(text) >= len(format) {
			_, err := time.Parse(format, text[:len(format)])
			if err != nil {
				continue
			}
			return format
		}
	}
	return ""
}

func int64timeToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
func tailFile(fileMonitor FileMonitor, edit_box *EditBox ) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true,
		ReOpen:true,
		Poll: fileMonitor.Poll,
		Logger:tail.DiscardingLogger,
		Location:&tail.SeekInfo{fileMonitor.Offset, os.SEEK_SET}})
	var key []byte
	var prevData string
	var prevTs string

	f := ""
	prevo := int64(0)
	stopo := int64(0)
	for line := range t.Lines {

		td := time.Now()
		var tt time.Time
		text := line.Text
		if f == "" {
			f = findFormat(text)
		}
		if f == "" {
			continue
		}
		var ok int
		if len(text) > len(f) {
			ti, err := time.Parse(f, text[:len(f)])
			if err != nil {
				ok = -1
			}
			// New event found
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



		// Multiline entry add to last timestamp
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
				event.BloomDirty = true
				event.Lines += 1
				events.BloomDirty = true
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
			edit_box.Lock()
			edit_box.storeLine = time.Now().Sub(td)
			edit_box.Unlock()
			continue
		}

		var event = Event{
			Ts: tt.Format(time.RFC3339),
			Data: text,
			Path: fileMonitor.Path,
			BloomDirty:true,
		}

		key = int64timeToByte(tt.Truncate(1 * time.Minute).Unix())
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
			events.sortEvents()

			events.BloomDirty = true
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
				var meta Meta
				b, _ := meta.Marshal()
				by = b
			}
			var meta Meta
			meta.Unmarshal(by)
			meta.Count++
			by, _ = meta.Marshal()
			b.Put([]byte("Meta"), by)

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
		edit_box.Lock()
		edit_box.storeLine = time.Now().Sub(td)
		edit_box.Unlock()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func regenerateBloom(k []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		by := b.Get(k)
		var e Events
		err := e.Unmarshal(by)
		if err != nil {
			log.Fatal(err)
		}
		if e.BloomDirty {
			e.RegenerateBloom()
			by, err = e.Marshal()
			if err != nil {
				log.Fatal(err)
			}
			b.Put(k, by)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
func eventRegenBloom(k []byte, i int) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		by := b.Get(k)
		var e Events
		err := e.Unmarshal(by)
		if err != nil {
			log.Fatal(err)
		}
		if i > len(e.Events) - 1 {
			return nil
		}
		e.Events[i].GenerateBloom()
		by, err = e.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		b.Put(k, by)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
func shouldNotContinueBasedOnBucketFilter(keys []string, bloomArray []byte) bool {
	noInSet := false
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if strings.Contains(key, "<") {
			split := strings.Split(key, "<")
			if !bloom.Filter(bloomArray).MayContain([]byte(split[0])) {
				noInSet = true
				continue
			}
		} else if strings.Contains(key, ">") {
			split := strings.Split(key, ">")
			if !bloom.Filter(bloomArray).MayContain([]byte(split[0])) {
				noInSet = true
				continue
			}
		} else if key[:1] == "!" {
			if bloom.Filter(bloomArray).MayContain([]byte(key[1:])) {
				noInSet = true
				break
			}
		} else if !bloom.Filter(bloomArray).MayContain([]byte(key)) {
			noInSet = true
			break
		}
	}
	return noInSet
}

func SearchFor(t []byte, wantedItems int, skipItems int64, ch chan SearchRes, quit chan bool) {
	mutex2.Lock()
	defer mutex2.Unlock()
	ttt := time.Now()
	var searchRes SearchRes
	count := int64(0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		c := b.Cursor()
		k, v := c.Last()

		for ; k != nil && count < int64(wantedItems); k, v = c.Prev() {
			select {
			case <-quit:
				return nil
			default:
				var events Events
				err := events.Unmarshal(v)
				if err != nil {
					log.Fatal(err)
				}

				search := strings.Split(string(t), "|")
				keys := strings.Split(search[0], " ")
				search = search[1:]
				if events.BloomDirty {
					go regenerateBloom(k);
				} else {
					if len(t) != 0 {
						if shouldNotContinueBasedOnBucketFilter(keys, events.Bloom) {
							continue
						}
					}
				}
				for i := len(events.Events) - 1; i >= 0; i-- {
					event := events.Events[i]
					if len(t) == 0 {
						if skipItems == int64(0) {
							count += int64(event.Lines) + int64(1)
							eventRes := EventRes{Data:event.Data,
								Lines: event.Lines,
								Fields:event.Fields,
								Ts: event.Ts,
								Path: event.Path,
							}
							searchRes.Events = append(searchRes.Events, &eventRes)
							continue
						}
						skipItems--
						continue
					}
					if event.BloomDirty {
						go eventRegenBloom(k, i)
					}

					add := event.shouldAddAndGetIndexes(keys)
					if add {
						if skipItems == int64(0) {
							if len(search) > 0 && strings.TrimSpace(search[0]) == "count" {
								searchRes.Count++
								if count == int64(wantedItems) - 1 {
									continue
								}
							}
							count += int64(event.Lines) + int64(1)
							eventRes := EventRes{Data:event.Data,
								Lines: event.Lines,
								Fields:event.Fields,
								FoundAtIndex: event.GetKeyIndexes(keys),
								Ts: event.Ts,
								Path: event.Path,
							}

							searchRes.Events = append(searchRes.Events, &eventRes)
							continue
						}
						skipItems--
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	searchRes.Ts = time.Now().Sub(ttt).String()
	ch <- searchRes
}