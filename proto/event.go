package proto

import (
	"strconv"
	"strings"
	"github.com/golang/leveldb/bloom"
)

func (event *Event) ShouldAddAndGetIndexes(keys []string) (bool) {
	add := true
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if event.BloomDirty {
			if key[:1] == "!" {
				if strings.Contains(event.Data, key[1:]) {
					add = false
					break
				}
			} else {
				if !(strings.Contains(event.Data, key) || strings.Contains(event.Path, key)) {
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
			if bloom.Filter(event.Bloom).MayContain([]byte(key[1:])) {
				add = false
				break
			}
		} else if !bloom.Filter(event.Bloom).MayContain([]byte(key)) || !(strings.Contains(event.Data, key) || strings.Contains(event.Path, key)) {
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
		index := strings.Index(event.Data, key)
		text := event.Data
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
