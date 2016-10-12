package tail

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"github.com/hpcloud/tail"
	"github.com/jantb/search/proto"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var regenChan = make(chan []byte, 10000)
var once sync.Once

func regenerateBloom(keys chan []byte, db *bolt.DB) {
	for {
		k := <-keys
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Events"))
			b, _ = b.CreateBucketIfNotExists(Int64timeToByte(ByteToint64timeTo(k).Truncate(24 * time.Hour).Unix()))

			by := b.Get(k)
			var e proto.Events
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
}

func tailFile(fileMonitor proto.FileMonitor, db *bolt.DB) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true,
		ReOpen:   true,
		Poll:     fileMonitor.Poll,
		Logger:   tail.DiscardingLogger,
		Location: &tail.SeekInfo{fileMonitor.Offset, os.SEEK_SET}})
	var key []byte
	var dayKey []byte
	var prevData string
	var prevTs string

	f := ""
	prevo := int64(0)
	stopo := int64(0)
	for line := range t.Lines {
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
				tt = ti
				stopo = prevo
				if tt.Before(time.Now().Truncate(time.Hour * 24).AddDate(0, -1, 0)) {
					key = nil
					continue
				}
				text = text[len(f) + 1:]
				ok = 1
			}
		}
		o, err := t.Tell()
		if err != nil {
			log.Fatal(err)
		}
		prevo = o
		// Multiline entry add to last timestamp
		if ok == -1 || ok == 0 {
			if key == nil {
				continue
			}
			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Events"))
				b, _ = b.CreateBucketIfNotExists(dayKey)
				by := b.Get(key)

				var events proto.Events
				err := events.Unmarshal(by)
				if err != nil {
					log.Fatal(err)
				}

				event, _ := events.Get(prevTs, prevData)
				event.SetData(event.GetData() + "\n" + text)
				prevData = event.GetData()
				event.BloomDirty = true
				event.Lines += 1
				events.BloomDirty = true

				by, err = events.Marshal()
				if err != nil {
					log.Fatal(err)
				}
				b.Put(key, by)
				regenChan <- key
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		var event = proto.Event{
			Ts:         tt.Format(time.RFC3339),
			Path:       fileMonitor.Path,
			BloomDirty: true,
		}
		event.SetData(text)

		key = Int64timeToByte(tt.Truncate(1 * time.Minute).Unix())
		dayKey = Int64timeToByte(tt.Truncate(24 * time.Hour).Unix())
		prevData = event.GetData()
		prevTs = event.Ts
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Events"))
			b, _ = b.CreateBucketIfNotExists(dayKey)
			eventsb := b.Get(key)
			if eventsb == nil {
				var events proto.Events
				eventsb, _ = events.Marshal()
			}
			var events proto.Events
			events.Unmarshal(eventsb)

			_, found := events.Get(event.Ts, event.GetData())
			if found {
				return nil
			}

			events.Events = append(events.Events, &event)
			events.SortEvents()

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
				var meta proto.Meta
				b, _ := meta.Marshal()
				by = b
			}
			var meta proto.Meta
			meta.Unmarshal(by)
			meta.Count++
			by, _ = meta.Marshal()
			b.Put([]byte("Meta"), by)
			regenChan <- key
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		panic(err)
	}
}
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

var formats = []string{"2006/01/02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.0",
	"2006-01-02 15:04:05.00",
	"2006-01-02 15:04:05.0000",
	"2006-01-02 15:04:05.00000",
	"2006-01-02 15:04:05.000000",
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
}

func Int64timeToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func ByteToint64timeTo(bytes []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(bytes)), int64(0))
}

func TailAllFiles(db *bolt.DB) {
	once.Do(func() {
		go regenerateBloom(regenChan, db)
	})
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		c := b.Cursor()
		for k, f := c.First(); k != nil; k, f = c.Next() {
			fileMonitor := proto.FileMonitor{}
			fileMonitor.Unmarshal(f)

			go tailFile(fileMonitor, db)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
func AddFileToTail(filename string, poll bool, db *bolt.DB) {
	file, err := os.Open(filename)
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
			dir, _ := filepath.Abs(filepath.Dir(filename))
			filep := filepath.Join(dir, filepath.Base(filename))
			fileMonitor := proto.FileMonitor{
				Path:   filep,
				Offset: 0,
				Poll:   poll,
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
