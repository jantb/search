package proto

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bradfitz/slice"
	"github.com/golang/leveldb/bloom"
	"github.com/golang/snappy"
)

func (e *Events) Get(ts string, data string) (*Event, bool) {
	for _, ev := range e.GetEvents() {
		if data == ev.GetData() && ts == ev.Ts {
			return ev, true
		}
	}
	return &Event{}, false
}

func (e *Events) Store(db *bolt.DB) {
	if len(e.Events) == 0 {
		return
	}
	ts, _ := time.Parse("2006-01-02T15:04:05.999Z07:00", e.Events[0].Ts)
	key := Int64timeToByte(ts.Truncate(1 * time.Minute).Unix())
	dayKey := Int64timeToByte(ts.Truncate(24 * time.Hour).Unix())
	marshal, err := e.Marshal()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		b, _ = b.CreateBucketIfNotExists(dayKey)
		if err != nil {
			log.Fatal(err)
		}
		b.Put(key, snappy.Encode(nil, marshal))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
func (e *Events) Retrieve(ts time.Time, db *bolt.DB) bool {
	var eventsb []byte

	key := Int64timeToByte(ts.Truncate(1 * time.Minute).Unix())
	dayKey := Int64timeToByte(ts.Truncate(24 * time.Hour).Unix())
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		b, _ = b.CreateBucketIfNotExists(dayKey)
		var buffer bytes.Buffer
		buffer.Write(b.Get(key))
		eventsb = buffer.Bytes()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(eventsb) != 0 {
		b, err := snappy.Decode(nil, eventsb)
		if err != nil {
			log.Panic(err)
		}

		e.Unmarshal(b)
		return true
	}
	return false
}

func (e *Events) GetById(id int32) (*Event, bool) {
	for _, ev := range e.GetEvents() {
		if e.Id == id {
			return ev, true
		}
	}
	return &Event{}, false
}
func (e *Events) SortEvents() {
	slice.Sort(e.Events, func(i, j int) bool {
		return e.Events[i].Ts < e.Events[j].Ts
	})
}
func (e *Events) RegenerateBloom() {
	set := make(map[string]bool)

	for _, ev := range e.Events {
		keys := ev.GenerateBloom()
		for _, k := range keys {
			set[string(k)] = true
		}
	}
	keys := make([][]byte, 0, len(set))
	for k := range set {
		keys = append(keys, []byte(k))
	}

	e.Bloom = bloom.NewFilter(nil, keys, 10)
	e.BloomDirty = false
}

func Int64timeToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func ByteToint64timeTo(bytes []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(bytes)), int64(0))
}
