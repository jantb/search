package main

import (
	"go4.org/intern"
	"strings"
	"time"
)

var ll = LL{}
var realOffset = 0

func clear() {
	ll = LL{}
}

func removeLast() {
	for mem, _, _ := memusage(); mem > 500; mem, _, _ = memusage() {
		ll.RemoveLast()
	}
}

func getLength() int {
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
	restOfQuery := strings.Join(restTokens, " ")

	reachedTop := false

	done := make(chan struct{})

	matchSet := make(map[*intern.Value]bool)
	noMatchSet := make(map[*intern.Value]bool)

	for line := range ll.Iterate(done) {
		if shouldSkipLine(skipTokens, line) {
			continue
		}
		match, m, n := line.matchOrNot(restOfQuery, matchSet, noMatchSet)
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
			close(done)
			reachedTop = true
			break
		}
	}
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
			if shouldSkipLine(tokens, line) {
				continue
			}

			if line.matches(query, restOfQuery) {
				count++
			}
		}
	}

	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now), count
}

func findTokens(tokens []string) ([]string, []string) {
	var skipTokens []string
	var restTokens []string
	for i, token := range tokens {
		if strings.HasPrefix(token, "level=") || strings.HasPrefix(token, "level!=") || strings.HasPrefix(token, "!") || strings.HasPrefix(token, "system=") || strings.HasPrefix(token, "system!=") {
			skipTokens = append(skipTokens, tokens[i])
			continue
		}
		restTokens = append(restTokens, tokens[i])
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
