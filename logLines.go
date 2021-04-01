package main

import (
	"go4.org/intern"
	"strings"
	"time"
)

var member void

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
	ids    map[*intern.Value]void
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}
func (l LogLine) matchOrNot(query []string, matches map[*intern.Value]void) (bool, *intern.Value) {

	if len(query) == 0 {
		return true, nil
	}
	for k := range matches {
		_, ok := l.ids[k]
		if ok {
			return true, nil
		}
	}

	for value := range l.ids {
		val := value.Get().(string)
		for _, s := range query {
			if strings.Contains(val, s) {
				return true, value
			}
		}
	}
	return false, nil
}

func (l LogLine) matchOrNotCount(query []string, matches map[*intern.Value]void) bool {

	if len(query) == 0 {
		return true
	}

	for k := range matches {
		_, ok := l.ids[k]
		if ok {
			return true
		}
	}

	return false
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
	l.ids = make(map[*intern.Value]void)
	l.ids[l.level] = member
	l.ids[l.system] = member
	for _, value := range l.body {
		l.ids[value] = member
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
