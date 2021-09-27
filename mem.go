package main

import (
	"container/list"
	"github.com/jantb/search/logline"
	"go4.org/intern"
	"os"
	"runtime"
	"strings"
	"time"
)

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

func insertIntoStoreByChan(insertChan chan logline.LogLine) {
	file, e := os.ReadFile(".ignoreSearch")
	var tokens []string
	if e == nil {
		query := string(file)
		tokens, _ = findTokens(strings.Split(strings.TrimSpace(query), " "))
	}

	for line := range insertChan {
		if line.GetLevel() != "error" && shouldSkipLine(tokens, line) {
			continue
		}
		ll.Put(line)
		bottomChan <- true
	}
}

func search(input string, limit int, offset int) (ret []logline.LogLine, t time.Duration, count int) {
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
		return []logline.LogLine{}, time.Now().Sub(now), 0
	}

	skipTokens, restTokens := findTokens(tokens)

	reachedTop := false

	done := make(chan struct{})

	matchSet := make(map[*intern.Value]logline.Void)
	//	restOfQuery := strings.Join(restTokens, " ")
	for line := range ll.Iterate(done) {
		if shouldSkipLine(skipTokens, line) {
			continue
		}
		//	match := includeLine(query, line, restOfQuery)
		match, m := line.MatchOrNot(restTokens, matchSet)
		if m != nil {
			matchSet[m] = logline.Member
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
		ret = append(ret, logline.LogLine{
			Level:  intern.GetByString(""),
			System: intern.GetByString(""),
			Time:   0,
			Body:   nil,
		})
	}

	if command == "count" {
		done := make(chan struct{})
		for line := range ll.Iterate(done) {
			if shouldSkipLine(skipTokens, line) {
				continue
			}
			match := line.MatchOrNotCount(restTokens, matchSet)
			//	match := includeLine(query, line, restOfQuery)
			if match {
				count++
			}
		}
	}

	logline.ReverseLogline(ret)
	bottom.Store(realOffset == 0)
	return ret, time.Now().Sub(now), count
}
func includeLine(query string, line logline.LogLine, restOfQuery string) bool {
	return len(query) == 0 || strings.Contains(line.GetLevel(), strings.ToUpper(restOfQuery)) || strings.Contains(line.GetSystem(), strings.ToUpper(restOfQuery)) || strings.Contains(line.GetBody(), restOfQuery)
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

func shouldSkipLine(tokens []string, line logline.LogLine) bool {
	skip := false
	for _, token := range tokens {
		if strings.HasPrefix(token, "system=") {
			if !strings.Contains(strings.ToUpper(line.GetSystem()), strings.ToUpper(strings.Split(token, "=")[1])) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "system!=") {
			if strings.Contains(strings.ToUpper(line.GetSystem()), strings.ToUpper(strings.Split(token, "!=")[1])) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "level=") {
			if line.GetLevel() != strings.ToUpper(strings.Split(token, "=")[1]) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "level!=") {
			if line.GetLevel() == strings.ToUpper(strings.Split(token, "!=")[1]) {
				skip = true
				break
			}
			continue
		}

		if strings.HasPrefix(token, "!") {
			if strings.Contains(line.GetBody(), token[1:]) {
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
