package main

import (
	"strings"
	"sync"
	"time"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			return make(map[string]string)
		},
	}
	poolCount = 0
)

func InternSize() int {
	return poolCount
}
func Intern(s string) string {
	m := pool.Get().(map[string]string)
	c, ok := m[s]
	if ok {
		pool.Put(m)
		return c
	}
	poolCount++
	m[s] = s
	pool.Put(m)
	return s
}

func reverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

type LogLine struct {
	level  string
	system string
	Time   int64
	body   []string
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}

func (l LogLine) getBody() string {
	return strings.Join(l.body, " ")
}

func (l *LogLine) setBody(body string) {
	s := strings.Split(body, " ")
	for _, part := range s {
		l.body = append(l.body, Intern(part))
	}
}

func (l LogLine) getLevel() string {
	return l.level
}

func (l *LogLine) setLevel(level string) {
	l.level = Intern(level)
}

func (l LogLine) getSystem() string {
	return l.system
}

func (l *LogLine) setSystem(s string) {
	l.system = Intern(s)
}
