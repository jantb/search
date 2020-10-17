package main

import (
	"strings"
	"time"
)

var bst = Tree{}
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
	return bst.size
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
		bst.Put(line)
	}
}

func search(query string, limit int, offset int) (ret []LogLine, t time.Duration) {
	now := time.Now()
	query = strings.TrimSpace(query)
	tokens := strings.Split(strings.TrimSpace(query), " ")

	setOffset(offset)
	insertOffset := realOffset
	if getLength() == 0 {
		return []LogLine{}, time.Now().Sub(now)
	}
	done := make(chan struct{})
	ch := bst.Iterate(done)

	for line := range ch {
		skip, restTokens := shouldSkipLine(tokens, line)

		if skip {
			continue
		}

		join := strings.Join(restTokens, " ")

		if len(query) == 0 || strings.Contains(line.getLevel(), strings.ToUpper(join)) || strings.Contains(line.getBody(), join) {
			if insertOffset == 0 {
				ret = append(ret, line)
			} else {
				insertOffset--
			}
		}

		if len(ret) == limit+realOffset {
			close(done)
			break
		}
	}
	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now)
}

func shouldSkipLine(tokens []string, line LogLine) (bool, []string) {
	skip := false
	var restTokens []string
	for _, token := range tokens {
		if strings.HasPrefix(token, "level=") {
			if line.getLevel() != strings.ToUpper(strings.Split(token, "=")[1]) {
				skip = true
			}
			continue
		}

		if strings.HasPrefix(token, "level!=") {
			if line.getLevel() == strings.ToUpper(strings.Split(token, "!=")[1]) {
				skip = true
			}
			continue
		}

		if strings.HasPrefix(token, "!") {
			if strings.Contains(line.getBody(), token[1:]) {
				skip = true
			}
			continue
		}
		restTokens = append(restTokens, token)
	}
	return skip, restTokens
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
