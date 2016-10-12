package proto

import (
	"github.com/golang/leveldb/bloom"
	"github.com/jantb/search/utils"
	"strconv"
	"strings"
	"github.com/golang/snappy"
	"log"
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
				if !(strings.Contains(event.GetData(), key) || strings.Contains(event.Path, key)) {
					add = false
					continue
				}
			}
		} else {
			if strings.Contains(key, "<") {
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
				if (bloom.Filter(event.Bloom).MayContain([]byte(key[1:])) && (strings.Contains(event.GetData(), key[1:]) || strings.Contains(event.Path, key[1:]))) {
					add = false
					break
				}
			} else if !bloom.Filter(event.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.GetData(), key) || strings.Contains(event.Path, key)) {
				add = false
				continue
			}
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

func (e *Event) GenerateBloom() {
	if e.BloomDirty {
		e.Fields = e.Fields[:0]
		set := make(map[string]bool)
		for _, key := range utils.GetBloomKeysFromLine(e.GetData()) {
			set[string(key)] = true
			if strings.ContainsRune(string(key), '=') {
				split := strings.Split(string(key), "=")
				set[string(split[0])] = true
				set[string(split[1])] = true
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
		e.Bloom = bloom.NewFilter(nil, keys, 10)
		e.BloomDirty = false
	}
}

var data string

func (e *Event) SetData(text string) {
	e.Data = snappy.Encode(nil, []byte(text))
	data = text
}

func (e *Event) GetData() string {
	if data == "" {
		b, err := snappy.Decode(nil, e.Data)
		if err != nil {
			log.Panic(err)
		}
		data = string(b)
	}
	return data
}
