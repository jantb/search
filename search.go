package main

import (
	"time"
)

func reverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

type LogLine struct {
	Id     int
	Level  string
	System string
	Time   int64
	Body   string
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}
