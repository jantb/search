package proto

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/boltdb/bolt"
	"github.com/golang/leveldb/bloom"
	"github.com/jantb/search/utils"
)

func (event *Event) ShouldAddAndGetIndexes(keys []string) bool {
	add := true
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if strings.Contains(key, "<") {
			split := strings.Split(key, "<")
			if !bloom.Filter(event.D.Bloom).MayContain([]byte(split[0])) {
				add = false
				continue
			}
			val := ""
			for _, f := range event.D.Fields {
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
			if !bloom.Filter(event.D.Bloom).MayContain([]byte(split[0])) {
				add = false
				continue
			}
			val := ""
			for _, f := range event.D.Fields {
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
			if bloom.Filter(event.D.Bloom).MayContain([]byte(key[1:])) && (strings.Contains(event.GetData(), key[1:]) || strings.Contains(event.D.Path, key[1:])) {
				add = false
				break
			}
		} else if !bloom.Filter(event.D.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.GetData(), key) || strings.Contains(event.D.Path, key)) {
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
		index := strings.Index(event.GetData(), key)
		text := event.GetData()
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

func (e *Event) GetKeys(db *bolt.DB) [][]byte {
	e.D.Fields = e.D.Fields[:0]
	set := make(map[string]bool)
	for _, key := range utils.GetBloomKeysFromLine(e.GetData()) {
		set[string(key)] = true
		if strings.ContainsRune(string(key), '=') {
			split := strings.Split(string(key), "=")
			set[split[0]] = true
			set[split[1]] = true
			e.D.Fields = append(e.D.Fields, &Field{Key: split[0], Value: split[1]})
		}
	}
	for _, key := range utils.GetBloomKeysFromLine(e.D.Path) {
		set[string(key)] = true
	}
	keys := make([][]byte, 0, len(set))
	for k := range set {
		keys = append(keys, []byte(k))
	}
	return keys
}
func (e *Event) BloomUpdate(db *bolt.DB) {
	e.D.Bloom = bloom.NewFilter(nil, e.GetKeys(db), 10)
}

func (e *Event) GetLines() int32 {
	return e.D.Lines
}

func (e *Event) IncrementLines() {
	e.D.Lines++
}

func (e *Event) GenerateKey() []byte {
	if e.D == nil {
		return []byte{}
	}

	var buffer bytes.Buffer
	buffer.Write(xxhash.New64().Sum([]byte(e.D.Data + e.D.Path)))
	key := buffer.Bytes()
	e.Data = key
	return key
}

func (e *Event) SetData(text string) {
	e.D.Data = text
}

func (e *Event) GetData() string {
	return e.D.Data
}

func (e *Event) Store(db *bolt.DB) {
	found := false
	e.GenerateKey()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Data"))
		found = b.Get(e.Data) != nil
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if !found {
		e.BloomUpdate(db)
		da, _ := e.D.Marshal()
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Data"))
			b.Put(e.Data, da)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		var meta Meta
		meta.Retrieve(db)
		meta.Unique++
		meta.Store(db)
	}
	e.D = nil
	marshal, err := e.Marshal()
	var meta Meta
	meta.Retrieve(db)
	meta.Count++

	meta.Store(db)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		b.Put([]byte(e.Ts+string(e.Data)), marshal)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Event) Exists(db *bolt.DB) bool {
	found := false
	e.GenerateKey()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		found = b.Get([]byte(e.Ts+string(e.Data))) != nil
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return found
}

func (e *Event) Retrieve(key []byte, db *bolt.DB) {
	var eventsb []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		var buffer bytes.Buffer
		buffer.Write(b.Get(key))
		eventsb = buffer.Bytes()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(eventsb) != 0 {
		e.Unmarshal(eventsb)
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Data"))
			var buffer bytes.Buffer
			buffer.Write(b.Get(key))
			data := Data{}
			data.Unmarshal(buffer.Bytes())
			e.D = &data
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
}
