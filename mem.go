package main

import (
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
	done := make(chan struct{})
	for line := range tree.Iterate(done) {
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

	if command == "count" {
		done := make(chan struct{})
		for line := range tree.Iterate(done) {
			skip, restTokens := shouldSkipLine(tokens, line)

			if skip {
				continue
			}

			join := strings.Join(restTokens, " ")

			if len(query) == 0 || strings.Contains(line.getLevel(), strings.ToUpper(join)) || strings.Contains(line.getBody(), join) {
				count++
			}
		}
	}

	reverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now), count
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
