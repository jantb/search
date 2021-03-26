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
func (l LogLine) matchOrNot(query []string, matches map[*intern.Value]bool, noMatches map[*intern.Value]bool) (bool, []*intern.Value, []*intern.Value) {
	ids := map[*intern.Value]bool{}
	for k, v := range l.ids {
		ids[k] = v
	}
	var match []*intern.Value
	var noMatch []*intern.Value
	if len(query) == 0 {
		return true, match, noMatch
	}
	for k := range matches {
		if ids[k] {
			return true, match, noMatch
		}
	}

	for k := range ids {
		if noMatches[k] {
			delete(ids, k)
		}
	}

	for value := range ids {
		val := value.Get().(string)
		for _, s := range query {
			if strings.Contains(val, s) {
				match = append(match, value)
				return true, match, noMatch
			}
		}
		noMatch = append(noMatch, value)
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
