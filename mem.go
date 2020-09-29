// +build mem

package main

import (
	"sort"
	"strings"
	"time"
)

var store []LogLine
var realOffset = 0
var id int

func initStore() {

}

func cleanupStore() {

}
func clearDb() {

}

func storeSettings(key, value string) {

}

func loadSettings(key string) string {
	return ""
}

func getLength() int {
	return len(store)
}

func insertIntoStoreByChan(insertChan chan []LogLine) {
	for {
		line := <-insertChan
		insertLoglinesToStore(line)
		bottomChan <- true
	}
}

func insertLoglinesToStore(logLines []LogLine) {
	for _, line := range logLines {
		id++
		line.Id = id
		insertSort(line)
	}
}

func insertSort(line LogLine) {
	index := sort.Search(len(store), func(i int) bool {
		if store[i].Time == line.Time {
			return store[i].Id > line.Id
		}
		return store[i].Time > line.Time
	})
	store = append(store, LogLine{})
	copy(store[index+1:], store[index:])
	line.Id = index + 1
	store[index] = line
}

func search(query string, limit int, offset int) (ret []LogLine, t time.Duration) {
	now := time.Now()
	query = strings.TrimSpace(query)
	tokens := strings.Split(strings.TrimSpace(query), " ")

	setOffset(offset)

	if len(store) == 0 {
		return store, time.Now().Sub(now)
	}

	for i := len(store) - 1; i >= 0; i-- {
		line := store[i]

		skip := false
		var restTokens []string
		for _, token := range tokens {
			if strings.HasPrefix(token, "level=") {
				if line.Level != strings.ToUpper(strings.Split(token, "=")[1]) {
					skip = true
				}
				continue
			}

			if strings.HasPrefix(token, "level!=") {
				if line.Level == strings.ToUpper(strings.Split(token, "!=")[1]) {
					skip = true
				}
				continue
			}

			if strings.HasPrefix(token, "!") {
				if strings.Contains(line.Body, token[1:]) {
					skip = true
				}
				continue
			}
			restTokens = append(restTokens, token)
		}

		if skip {
			continue
		}

		join := strings.Join(restTokens, " ")

		if len(query) == 0 || strings.Contains(line.Level, strings.ToUpper(join)) || strings.Contains(line.Body, join) {
			ret = append(ret, line)
		}

		if len(ret) == limit+realOffset {
			break
		}
	}
	if len(ret) > realOffset {
		ret = ret[realOffset:]
	}
	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now)
}

func setOffset(offset int) {
	if realOffset == 0 {
		if offset > 0 {
			realOffset += offset
		}
	}
	if realOffset+offset < 0 {
		realOffset = 0
	} else {
		realOffset += offset
	}
}
