package main

import (
	"container/list"
	"go4.org/intern"
	"runtime"
	"strings"
	"time"
)

var ll []LL
var realOffset = 0

func clear() {
	ll = []LL{}
}

func removeLast() {
	for mem, _, _ := memusage(); mem > 500; mem, _, _ = memusage() {
		for i := 0; i < 1000; i++ {
			for i := range ll {
				l := &ll[i]
				l.RemoveLast()
			}
		}
		runtime.GC()
	}
}

func getLength() int {
	c := 0
	for x := range ll {
		c += ll[x].GetSize()
	}
	return c
}

func insertIntoStoreByChan(insertChan chan LogLine) {
	for line := range insertChan {
		found := false
		for i := range ll {
			l := &ll[i]
			if l.System == line.getSystem() {
				l.Put(line)
				found = true
				break
			}
		}
		if !found {
			ll = append(ll, LL{
				System: line.getSystem(),
				l:      &list.List{},
			})

			for i := range ll {
				l := &ll[i]
				if l.System == line.getSystem() {
					l.Put(line)
					break
				}
			}
		}

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

	matchSet := make(map[*intern.Value]bool)
	noMatchSet := make(map[*intern.Value]bool)

	for line := range iterate(done) {
		if shouldSkipLine(skipTokens, line) {
			continue
		}
		match, m, n := line.matchOrNot(restTokens, matchSet, noMatchSet)
		for _, value := range m {
			matchSet[value] = true
		}
		for _, value := range n {
			noMatchSet[value] = true
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
		for line := range iterate(done) {
			if shouldSkipLine(skipTokens, line) {
				continue
			}
			match, _, _ := line.matchOrNot(restTokens, matchSet, noMatchSet)
			if match {
				count++
			}
		}
	}

	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now), count
}

func iterate(done <-chan struct{}) <-chan LogLine {
	var channels []<-chan LogLine

	for i := range ll {
		channels = append(channels, (&ll[i]).Iterate(done))
	}

	out := make(<-chan LogLine)
	if len(channels) > 0 {
		out = channels[0]
		for _, channel := range channels[1:] {
			out = Merge(out, channel)
		}
	}

	return out
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
