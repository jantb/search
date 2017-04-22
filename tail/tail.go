package tail

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hpcloud/tail"
	"github.com/jantb/search/proto"
)

var regenChan = make(chan []byte, 10000)
var once sync.Once

func regenerateBloom(keys chan []byte, db *bolt.DB) {
	for {
		k := <-keys
		var e proto.Events
		e.Retrieve(ByteToint64timeTo(k), db)

		if e.BloomDirty {
			e.RegenerateBloom()
			e.Store(db)
		}
	}
}

func tailFile(fileMonitor proto.FileMonitor, db *bolt.DB) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true,
		ReOpen:                                               true,
		Poll:                                                 fileMonitor.Poll,
		Logger:                                               tail.DiscardingLogger,
		Location:                                             &tail.SeekInfo{Offset: fileMonitor.Offset, Whence: os.SEEK_SET}})
	var key []byte
	var id = int32(0)
	var tt time.Time
	f := ""
	prevo := int64(0)
	stopo := int64(0)
	for line := range t.Lines {
		text := line.Text
		prefix := getPrefix(text)
		if f == "" {
			f = findFormat(text)
		}
		if f == "" {
			continue
		}
		var ok int
		text = text[len(prefix):]
		if len(text) > len(f) {
			ti, err := time.Parse(f, strings.Replace(text[:len(f)], ",", ".", -1))
			if err != nil {
				ok = -1
			}
			// New event found
			if ok == 0 {
				tt = ti
				stopo = prevo
				text = prefix + text[len(f)+1:]

				ok = 1
			}
		}
		prevo, err = t.Tell()
		if err != nil {
			log.Fatal(err)
		}

		// Multiline entry add to last timestamp
		if ok == -1 || ok == 0 {
			if key == nil {
				continue
			}

			var events proto.Events
			events.Retrieve(tt, db)
			event, _ := events.GetById(id)
			event.SetData(event.GetData() + "\n" + text)
			event.BloomDirty = true
			event.Lines += 1
			events.BloomDirty = true
			events.Store(db)
			regenChan <- key

			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		var event = proto.Event{
			Ts:         tt.Format("2006-01-02T15:04:05.999Z07:00"),
			Path:       fileMonitor.Path,
			BloomDirty: true,
		}
		event.SetData(text)

		key = Int64timeToByte(tt.Truncate(1 * time.Minute).Unix())

		var events proto.Events
		events.Retrieve(tt, db)

		_, found := events.Get(event.Ts, text)
		if found {
			return
		}
		events.Id++
		event.Id = events.Id
		id = events.Id
		events.Events = append(events.Events, &event)
		events.SortEvents()
		events.BloomDirty = true
		events.Store(db)

		fileMonitor.Offset = stopo
		fileMonitor.Store(db)
		regenChan <- key

		var meta proto.Meta
		meta.Retrieve(db)
		meta.Count++
		meta.Store(db)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getPrefix(text string) (string) {
	prefix := ""
	if strings.HasPrefix(text, "INFO ") {
		prefix = "INFO "
	}
	if strings.HasPrefix(text, "ERROR ") {
		prefix = "ERROR "
	}
	if strings.HasPrefix(text, "WARN ") {
		prefix = "WARN "
	}
	return prefix
}

func findFormat(text string) string {
	for _, format := range formats {
		if len(text) >= len(format) {
			_, err := time.Parse(format, strings.Replace(text[:len(format)], ",", ".", -1))
			if err != nil {
				continue
			}
			return format
		}
	}
	return ""
}

var formats = []string{
	"2006-01-02 15:04:05.999999",
	"2006-01-02 15:04:05.99999",
	"2006-01-02 15:04:05.9999",
	"2006-01-02 15:04:05.999",
	"2006-01-02 15:04:05.99",
	"2006-01-02 15:04:05.9",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
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
			var buffer bytes.Buffer
			buffer.Write(f)
			fileMonitor.Unmarshal(buffer.Bytes())

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
		dir, _ := filepath.Abs(filepath.Dir(filename))
		filep := filepath.Join(dir, filepath.Base(filename))
		fileMonitor := proto.FileMonitor{
			Path:   filep,
			Offset: 0,
			Poll:   poll,
		}
		fileMonitor.Store(db)
		return
	}
}
