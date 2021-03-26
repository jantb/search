package main

import (
	"go4.org/intern"
	"strings"
	"time"
)

func reverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

type LogLine struct {
	level  *intern.Value
	system *intern.Value
	Time   int64
	body   []*intern.Value
	ids    map[*intern.Value]bool
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}
func (l LogLine) matchOrNot(query string,
	matches map[*intern.Value]bool,
	noMatches map[*intern.Value]bool) (bool, []*intern.Value, []*intern.Value) {

	for k := range l.ids {
		if noMatches[k] {
			delete(l.ids, k)
		}
	}

	var match []*intern.Value
	var noMatch []*intern.Value
	for k := range matches {
		if l.ids[k] {
			return true, match, noMatch
		}
	}

	for value := range l.ids {
		val := value.Get().(string)
		found := false
		for _, s := range strings.Split(query, " ") {
			if strings.Contains(val, s) {
				found = true
				break
			}
		}
		if found {
			match = append(match, value)
			return true, match, noMatch
		} else {
			noMatch = append(noMatch, value)
		}
	}
	return false, match, noMatch
}

func (l LogLine) getBody() string {
	if l.body == nil {
		return ""
	}

	var sb strings.Builder
	for i, value := range l.body {
		if i == 0 {
			sb.WriteString(value.Get().(string))
		} else {
			sb.WriteString(" ")
			sb.WriteString(value.Get().(string))
		}
	}
	return sb.String()
}

func (l *LogLine) setBody(body string) {
	s := strings.Split(body, " ")
	for _, part := range s {
		l.body = append(l.body, intern.GetByString(part))
	}
	l.ids = make(map[*intern.Value]bool)
	l.ids[l.level] = true
	l.ids[l.system] = true
	for _, value := range l.body {
		l.ids[value] = true
	}
}

func (l LogLine) getLevel() string {
	return l.level.Get().(string)
}

func (l *LogLine) setLevel(level string) {
	l.level = intern.GetByString(level)
}

func (l LogLine) getSystem() string {
	return l.system.Get().(string)
}

func (l *LogLine) setSystem(s string) {
	l.system = intern.GetByString(s)
}
