package searchfor

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/snappy"
	"github.com/jantb/search/proto"
)

var Searching atomic.Value

var stop = atomic.Value{}
var mutex = sync.Mutex{}

func SearchFor(t []byte, wantedItems int, skipItems int64, ch chan []byte, db *bolt.DB) {
	if stop.Load() == nil {
		stop.Store(false)
	}
	stop.Store(true)
	mutex.Lock()
	defer mutex.Unlock()
	stop.Store(false)
	Searching.Store(true)
	defer Searching.Store(false)
	ttt := time.Now()
	var searchRes proto.SearchRes
	count := int64(0)
	search := strings.Split(string(t), "|")
	keys := strings.Split(search[0], " ")
	search = search[1:]

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		d := tx.Bucket([]byte("Data"))
		c := b.Cursor()
		k, v := c.Last()
		for ; k != nil && count <= int64(wantedItems); k, v = c.Prev() {
			if stop.Load().(bool) {
				return nil
			}
			var buffer bytes.Buffer
			buffer.Write(v)

			var event proto.Event
			b, err := snappy.Decode(nil, buffer.Bytes())
			if err != nil {
				log.Panic(err)
			}
			err = event.Unmarshal(b)
			if err != nil {
				log.Fatal(err)
			}
			var bufferd bytes.Buffer
			bufferd.Write(d.Get(event.Data))
			data := proto.Data{}
			data.Unmarshal(bufferd.Bytes())
			event.D = &data
			if len(t) == 0 {
				if skipItems == int64(0) {
					count += int64(event.Lines) + int64(1)
					eventRes := proto.EventRes{Data: event.GetData(),
						Lines:                   event.Lines,
						Fields:                  event.Fields,
						Ts:                      event.Ts,
						Path:                    proto.GetPathFromId(proto.Itob(event.Path), db),
					}
					searchRes.Events = append(searchRes.Events, &eventRes)
					send(searchRes, ch)
					continue
				}
				skipItems--
				continue
			}
			if event.ShouldAddAndGetIndexes(keys, db) {
				if skipItems == int64(0) {
					if len(search) > 0 && strings.TrimSpace(search[0]) == "count" {
						searchRes.Count++
						if count == int64(wantedItems)-1 {
							continue
						}
					}
					count += int64(event.Lines) + int64(1)
					eventRes := proto.EventRes{Data: event.GetData(),
						Lines:                   event.Lines,
						Fields:                  event.Fields,
						FoundAtIndex:            event.GetKeyIndexes(keys),
						Ts:                      event.Ts,
						Path:                    proto.GetPathFromId(proto.Itob(event.Path), db),
					}

					searchRes.Events = append(searchRes.Events, &eventRes)
					send(searchRes, ch)
					continue
				}
				skipItems--
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	searchRes.Ts = time.Now().Sub(ttt).String()
	marshal, err := searchRes.Marshal()
	if err != nil {
		log.Panic(err)
	}
	ch <- marshal

}
func send(searchRes proto.SearchRes, ch chan []byte) {
	marshal, err := searchRes.Marshal()
	if err != nil {
		log.Panic(err)
	}
	ch <- marshal
}
