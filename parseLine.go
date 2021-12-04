package main

import (
	"github.com/jantb/search/logline"
	"github.com/valyala/fastjson"
	"strings"
	"time"
)

func parseLineJsonFastJson(line []byte) logline.LogLine {
	body := fastjson.GetString(line, "message")
	stack := fastjson.GetString(line, "stack_trace")
	if len(stack) > 0 {
		body = body + "\n" + stack
	}
	timestamp := fastjson.GetString(line, "@timestamp")
	if len(timestamp) == 0 {
		timestamp = fastjson.GetString(line, "timestamp")
	}
	l := logline.LogLine{
		Time: toMillis(parseTimestampJson(timestamp)),
	}
	l.SetSystem(fastjson.GetString(line, "HOSTNAME"))
	l.SetLevel(fastjson.GetString(line, "level"))
	l.SetBody(body)
	return l
}

func parseTimestampJson(timestamp string) time.Time {
	date, e := time.ParseInLocation("2006-01-02T15:04:05.999-07:00", strings.Replace(timestamp, ",", ".", -1), time.Local)
	if e != nil {
		date, e = time.ParseInLocation("2006-01-02T15:04:05.999Z", strings.Replace(timestamp, ",", ".", -1), time.Local)
		if e != nil {
			date = time.Now().Add(-200 * time.Hour)
		}
	}
	if date.Year() == 0 {
		date = date.AddDate(time.Now().Year(), 0, 0)
	}
	return date
}

func insertIntoStoreJsonSystem(insertChan chan []byte, system string) {
	for line := range insertChan {
		logLine := parseLineJsonFastJson(line)
		logLine.SetSystem(system)
		insertLogLinesChan <- logLine
	}
}
