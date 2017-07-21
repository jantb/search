package tail

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hpcloud/tail"
	"github.com/jantb/search/proto"
)

func tailFile(fileMonitor proto.FileMonitor, db *bolt.DB) {
	t, err := tail.TailFile(fileMonitor.Path, tail.Config{Follow: true,
		ReOpen:   true,
		Poll:     fileMonitor.Poll,
		Logger:   tail.DiscardingLogger,
		Location: &tail.SeekInfo{Offset: fileMonitor.Offset, Whence: os.SEEK_SET}})
	key := []byte{}
	var tt time.Time
	f := ""
	prevo := int64(0)
	stopo := int64(0)
	buff := bytes.Buffer{}
	var event proto.Event
	event.D = &proto.Data{}
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
			s := strings.Replace(text[:len(f)], ",", ".", -1)
			s = strings.Replace(s, "T", " ", -1)
			ti, err := time.Parse(f, s)
			if err != nil {
				ok = -1
			}
			// New event found
			if ok == 0 {
				tt = ti
				stopo = prevo
				text = prefix + text[len(f)+1:]
				event.SetData(buff.String())
				l := buff.Len()
				buff.Reset()
				found := event.Exists(db)
				if !found && l > 0 {
					fileMonitor.Offset = stopo
					fileMonitor.Store(db)
					event.Store(db)
					ok = 1
				}
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
			buff.WriteString("\n" + text)
			event.IncrementLines()

			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		event = proto.Event{
			Ts: uint64(tt.UnixNano()),
			D: &proto.Data{
				Path: fileMonitor.Path,
			},
		}
		buff.WriteString(text)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getPrefix(text string) string {
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
			s := strings.Replace(text[:len(format)], ",", ".", -1)
			s = strings.Replace(s, "T", " ", -1)
			_, err := time.Parse(format, s)
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

func TailAllFiles(db *bolt.DB) {
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
