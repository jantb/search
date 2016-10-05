package  tail

import (
	"time"
	"os"
	"github.com/hpcloud/tail"
	"log"
	"encoding/binary"
	"github.com/boltdb/bolt"
	"github.com/jantb/search/proto"
)

func TailFile(fileMonitor proto.FileMonitor, db *bolt.DB) {
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

				var events proto.Events
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
			continue
		}

		var event = proto.Event{
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
				var events proto.Events
				eventsb, _ = events.Marshal()
			}
			var events proto.Events
			events.Unmarshal(eventsb)

			_, found := events.Get(event.Ts, event.Data)
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

func int64timeToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

