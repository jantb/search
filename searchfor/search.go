package searchfor

import (
	"github.com/boltdb/bolt"
	"log"
	"time"
	"github.com/golang/leveldb/bloom"
	"strings"
	"github.com/jantb/search/proto"
	"github.com/jantb/search/tail"
)

var Searching int32

func shouldNotContinueBasedOnBucketFilter(keys []string, bloomArray []byte) bool {
	noInSet := false
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if strings.Contains(key, "<") {
			split := strings.Split(key, "<")
			if !bloom.Filter(bloomArray).MayContain([]byte(split[0])) {
				noInSet = true
				continue
			}
		} else if strings.Contains(key, ">") {
			split := strings.Split(key, ">")
			if !bloom.Filter(bloomArray).MayContain([]byte(split[0])) {
				noInSet = true
				continue
			}
		} else if key[:1] == "!" {
			if bloom.Filter(bloomArray).MayContain([]byte(key[1:])) {
				noInSet = true
				break
			}
		} else if !bloom.Filter(bloomArray).MayContain([]byte(key)) {
			noInSet = true
			break
		}
	}
	return noInSet
}

func SearchFor(t []byte, wantedItems int, skipItems int64, ch chan proto.SearchRes, quit chan bool, db *bolt.DB) {
	Searching = 1
	defer func() {
		Searching = 0
	}()
	ttt := time.Now()
	var searchRes proto.SearchRes
	count := int64(0)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Events"))
		for tim := time.Now().Truncate(time.Hour * 24); tim.After(time.Now().Truncate(time.Hour * 24).AddDate(0,-1,0)) ; tim = tim.Add(time.Hour * -24) {

			b := b.Bucket(tail.Int64timeToByte(tim.Unix()))
			if b == nil {
				continue
			}
			c := b.Cursor()
			k, v := c.Last()

			for ; k != nil && count < int64(wantedItems); k, v = c.Prev() {
				select {
				case <-quit:
					return nil
				default:
					var events proto.Events
					err := events.Unmarshal(v)
					if err != nil {
						log.Fatal(err)
					}

					search := strings.Split(string(t), "|")
					keys := strings.Split(search[0], " ")
					search = search[1:]

					if len(t) != 0 {
						if shouldNotContinueBasedOnBucketFilter(keys, events.Bloom) {
							continue
						}
					}
					for i := len(events.Events) - 1; i >= 0; i-- {
						event := events.Events[i]
						if len(t) == 0 {
							if skipItems == int64(0) {
								count += int64(event.Lines) + int64(1)
								eventRes := proto.EventRes{Data:event.Data,
									Lines: event.Lines,
									Fields:event.Fields,
									Ts: event.Ts,
									Path: event.Path,
								}
								searchRes.Events = append(searchRes.Events, &eventRes)
								continue
							}
							skipItems--
							continue
						}
						add := event.ShouldAddAndGetIndexes(keys)
						if add {
							if skipItems == int64(0) {
								if len(search) > 0 && strings.TrimSpace(search[0]) == "count" {
									searchRes.Count++
									if count == int64(wantedItems) - 1 {
										continue
									}
								}
								count += int64(event.Lines) + int64(1)
								eventRes := proto.EventRes{Data:event.Data,
									Lines: event.Lines,
									Fields:event.Fields,
									FoundAtIndex: event.GetKeyIndexes(keys),
									Ts: event.Ts,
									Path: event.Path,
								}

								searchRes.Events = append(searchRes.Events, &eventRes)
								continue
							}
							skipItems--
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	searchRes.Ts = time.Now().Sub(ttt).String()
	ch <- searchRes
}