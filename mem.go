package main

import (
	"container/list"
	"go4.org/intern"
	"runtime"
	"strings"
	"time"
)

type void struct{}

var ll = LL{
	l: &list.List{},
}
var realOffset = 0

func clear() {
	ll = LL{
		l: &list.List{},
	}
}

func removeLast() {
	runtime.GC()
	for mem, _, _ := memusage(); mem > 400; mem, _, _ = memusage() {
		for i := 0; i < 100000; i++ {
			ll.RemoveLast()
		}
		runtime.GC()
	}
}

func getLength() int {
	if ll.l == nil {
		return 0
	}
	return ll.GetSize()
}

func insertIntoStoreByChan(insertChan chan LogLine) {
	for line := range insertChan {
		ll.Put(line)
		bottomChan <- true
	}
}

func search(input string, limit int, offset int) (ret []LogLine, t time.Duration, count int) {
	now := time.Now()
	input = strings.TrimSpace(input)
	split := strings.Split(input, "|")
	query := input
	command := ""
	if len(split) == 2 {
		query = strings.TrimSpace(split[0])
		command = strings.TrimSpace(split[1])
	}
	tokens := strings.Split(strings.TrimSpace(query), " ")

	setOffset(offset)
	insertOffset := realOffset
	if getLength() == 0 {
		return []LogLine{}, time.Now().Sub(now), 0
	}

	skipTokens, restTokens := findTokens(tokens)

	reachedTop := false

	done := make(chan struct{})

	matchSet := make(map[*intern.Value]void)
	//	restOfQuery := strings.Join(restTokens, " ")
	for line := range ll.Iterate(done) {
		if shouldSkipLine(skipTokens, line) {
			continue
		}
		//	match := includeLine(query, line, restOfQuery)
		match, m := line.matchOrNot(restTokens, matchSet)
		if m != nil {
			matchSet[m] = member
		}

		if match {
			if insertOffset == 0 {
				ret = append(ret, line)
			} else {
				insertOffset--
			}
		}

		if len(ret) == limit+realOffset {
			reachedTop = true
			break
		}
	}
	close(done)
	for !reachedTop && len(ret) != limit+realOffset {
		ret = append(ret, LogLine{
			level:  intern.GetByString(""),
			system: intern.GetByString(""),
			Time:   0,
			body:   nil,
		})
	}

	if command == "count" {
		done := make(chan struct{})
		for line := range ll.Iterate(done) {
			if shouldSkipLine(skipTokens, line) {
				continue
			}
			match := line.matchOrNotCount(restTokens, matchSet)
			//	match := includeLine(query, line, restOfQuery)
			if match {
				count++
			}
		}
	}

	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now), count
}
func includeLine(query string, line LogLine, restOfQuery string) bool {
	return len(query) == 0 || strings.Contains(line.getLevel(), strings.ToUpper(restOfQuery)) || strings.Contains(line.getSystem(), strings.ToUpper(restOfQuery)) || strings.Contains(line.getBody(), restOfQuery)
}

func findTokens(tokens []string) ([]string, []string) {
	var skipTokens []string
	var restTokens []string
	for i, token := range tokens {
		if strings.HasPrefix(token, "level=") || strings.HasPrefix(token, "level!=") || strings.HasPrefix(token, "!") || strings.HasPrefix(token, "system=") || strings.HasPrefix(token, "system!=") {
			skipTokens = append(skipTokens, tokens[i])
			continue
		}
		restTokens = append(restTokens, strings.TrimSpace(tokens[i]))
	}
	return skipTokens, restTokens
}

func shouldSkipLine(tokens []string, line LogLine) bool {
	skip := false
	for _, token := range tokens {
		if strings.HasPrefix(token, "system=") {
			if !strings.Contains(strings.ToUpper(line.getSystem()), strings.ToUpper(strings.Split(token, "=")[1])) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "system!=") {
			if strings.Contains(strings.ToUpper(line.getSystem()), strings.ToUpper(strings.Split(token, "!=")[1])) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "level=") {
			if line.getLevel() != strings.ToUpper(strings.Split(token, "=")[1]) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "level!=") {
			if line.getLevel() == strings.ToUpper(strings.Split(token, "!=")[1]) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "!") {
			if strings.Contains(line.getBody(), token[1:]) {
				skip = true
				break
			}
			continue
		}
	}
	return skip
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
