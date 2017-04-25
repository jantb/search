package proto

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/boltdb/bolt"
	"github.com/golang/leveldb/bloom"
	"github.com/golang/snappy"
	"github.com/jantb/search/utils"
)

func (event *Event) ShouldAddAndGetIndexes(keys []string) bool {
	add := true
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if event.BloomDirty {
			if strings.Contains(key, "<") {
				split := strings.Split(key, "<")
				if !strings.Contains(event.GetData(), split[0]) {
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
				if !strings.Contains(event.GetData(), split[0]) {
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
				if strings.Contains(event.GetData(), key[1:]) {
					add = false
					break
				}
			} else {
				fmt.Println(key)
				if !(strings.Contains(event.GetData(), key) || strings.Contains(event.Path, key)) {
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
			if bloom.Filter(event.Bloom).MayContain([]byte(key[1:])) && (strings.Contains(event.GetData(), key[1:]) || strings.Contains(event.Path, key[1:])) {
				add = false
				break
			}
		} else if !bloom.Filter(event.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.GetData(), key) || strings.Contains(event.Path, key)) {
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

func (e *Event) GetKeys() [][]byte {
	if len(e.Keys) != 0 {
		return e.Keys
	}
	e.Fields = e.Fields[:0]
	set := make(map[string]bool)
	for _, key := range utils.GetBloomKeysFromLine(e.GetData()) {
		set[string(key)] = true
		if strings.ContainsRune(string(key), '=') {
			split := strings.Split(string(key), "=")
			set[split[0]] = true
			set[split[1]] = true
			e.Fields = append(e.Fields, &Field{Key: split[0], Value: split[1]})
		}
	}
	for _, key := range utils.GetBloomKeysFromLine(e.Path) {
		set[string(key)] = true
	}
	keys := make([][]byte, 0, len(set))
	for k := range set {
		keys = append(keys, []byte(k))
	}
	e.Keys = keys
	return keys
}
func (e *Event) BloomUpdate() {
	e.Bloom = bloom.NewFilter(nil, e.GetKeys(), 10)
	e.BloomDirty = false
}

func (e *Event) GetKey() []byte {
	e.Id = xxhash.New64().Sum(e.Data)
	var buffer bytes.Buffer
	buffer.Write([]byte(e.Ts))
	buffer.Write(e.Id)
	return buffer.Bytes()
}

func (e *Event) SetData(text string) {
	e.Data = []byte(text)
	if len(e.Data) > 20000 {
		e.Data = e.Data[:20000]
	}
}

func (e *Event) GetData() string {
	if len(e.Data) == 0 {
		return ""
	}

	return string(e.Data)
}

func (e *Event) Store(db *bolt.DB) {

	marshal, err := e.Marshal()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		b.Put(e.GetKey(), snappy.Encode(nil, marshal))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Event) Exists(key []byte, db *bolt.DB) bool {
	found := false
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		found = b.Get(key) != nil
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
		b, err := snappy.Decode(nil, eventsb)
		if err != nil {
			log.Panic(err)
		}

		e.Unmarshal(b)
	}
}
