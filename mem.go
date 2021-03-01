package main

import (
	"go4.org/intern"
	"strings"
	"time"
)

var tree = Tree{}
var realOffset = 0

func clear() {
	tree = Tree{}
}

func getLength() int {
	return tree.size
}

func insertIntoStoreByChan(insertChan chan LogLine) {
	for line := range insertChan {
		tree.Put(line)
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
	for line := range tree.Iterate(done) {
		if shouldSkipLine(skipTokens, line) {
			continue
		}

		if includeLine(query, line, restOfQuery) {
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
		for line := range tree.Iterate(done) {
			if shouldSkipLine(tokens, line) {
				continue
			}

			if includeLine(query, line, restOfQuery) {
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
