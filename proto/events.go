package proto

import (
	"github.com/bradfitz/slice"
	"github.com/golang/leveldb/bloom"
	"github.com/jantb/search/utils"
)

func (e *Events) Get(ts string, data string) (*Event, bool) {
	for _, ev := range e.GetEvents() {
		if data == ev.Data && ts == ev.Ts {
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
		for _, key := range utils.GetBloomKeysFromLine(ev.Data) {
			set[string(key)] = true
		}
		for _, key := range utils.GetBloomKeysFromLine(ev.Path) {
			set[string(key)] = true
		}
		ev.GenerateBloom()
	}
	keys := make([][]byte, 0, len(set))
	for k := range set {
		keys = append(keys, []byte(k))
	}

	e.Bloom = bloom.NewFilter(nil, keys, 10)
	e.BloomDirty = false
}
