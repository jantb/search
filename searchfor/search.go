package searchfor

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jantb/search/proto"
)
import (
	"github.com/hashicorp/golang-lru/simplelru"
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

	lru, err := simplelru.NewLRU(10000, func(key interface{}, value interface{}) {})
	lru2, err := simplelru.NewLRU(10000, func(key interface{}, value interface{}) {})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
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

			err := event.Unmarshal(buffer.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			getData(d, &event, lru)
			if len(t) == 0 {
				if skipItems == int64(0) {
					count++
					eventRes := proto.EventRes{Data: event.GetData(),
						Lines:                   event.GetLines(),
						Fields:                  event.D.Fields,
						Ts:                      event.Ts,
						Path:                    event.D.Path,
					}
					searchRes.Events = append(searchRes.Events, &eventRes)
					//send(searchRes, ch)
					continue
				}
				skipItems--
				continue
			}
			add := addevent(lru2, event, keys)

			if add {
				if skipItems == int64(0) {
					if len(search) > 0 && strings.TrimSpace(search[0]) == "count" {
						searchRes.Count++
						if count == int64(wantedItems)-1 {
							continue
						}
					}
					count ++
					eventRes := proto.EventRes{Data: event.GetData(),
						Lines:                   event.GetLines(),
						Fields:                  event.D.Fields,
						FoundAtIndex:            event.GetKeyIndexes(keys),
						Ts:                      event.Ts,
						Path:                    event.D.Path,
					}

					searchRes.Events = append(searchRes.Events, &eventRes)
					//send(searchRes, ch)
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
func addevent(lru2 *simplelru.LRU, event proto.Event, keys []string) bool {
	add := false
	r, f := lru2.Get(string(event.Data))
	if !f {
		add = event.ShouldAddAndGetIndexes(keys)
		lru2.Add(string(event.Data), add)
	} else {
		add = r.(bool)
	}
	return add
}
func getData(d *bolt.Bucket, event *proto.Event, lru *simplelru.LRU) {

	dl, found := lru.Get(string(event.Data))
	data := proto.Data{}
	if !found {
		var bufferd bytes.Buffer
		bufferd.Write(d.Get(event.Data))
		data.Unmarshal(bufferd.Bytes())
		lru.Add(string(event.Data), data)
	} else {
		data = dl.(proto.Data)
	}

	event.D = &data
}
func send(searchRes proto.SearchRes, ch chan []byte) {
	marshal, err := searchRes.Marshal()
	if err != nil {
		log.Panic(err)
	}
	ch <- marshal
}
