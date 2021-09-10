package main

import (
	"github.com/jantb/search/logline"
	"github.com/valyala/fastjson"
	"go4.org/intern"
	"strings"
	"time"
)

func parseLine(line string, loglines []logline.LogLine) (logline.LogLine, bool) {
	for _, format := range formats {
		for _, regex := range format.Regex {
			match := regex.RegexCompiled.Match([]byte(line))
			if match {
				n1 := regex.RegexCompiled.SubexpNames()
				r2 := regex.RegexCompiled.FindAllStringSubmatch(line, -1)[0]
				md := map[string]string{}
				for i, n := range r2 {
					md[n1[i]] = n
				}
				if _, ok := md["timestamp"]; !ok {
					if len(loglines) == 0 {
						continue
					}
					l := loglines[len(loglines)-1]
					l.SetBody(l.GetBody() + "\n" + md["body"])
					logLine := logline.LogLine{}
					logLine.SetBody(line)
					return logLine, false
				}
				timestamp := toMillis(parseTimestamp(regex, md["timestamp"]))
				logLine := logline.LogLine{
					Time:   timestamp,
					System: intern.GetByString(md["system"]),
				}
				logLine.SetLevel(md["level"])
				logLine.SetSystem(md["system"])
				logLine.SetBody(md["body"])
				return logLine, true

			}
		}
	}
	logLine := logline.LogLine{}
	logLine.SetBody(line)
	return logLine, false
}

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

func parseTimestamp(regex Regex, timestamp string) time.Time {
	s := regex.Timestamp
	date, e := time.ParseInLocation(s, strings.Replace(timestamp, ",", ".", -1), time.Local)
	if e != nil {
		date = time.Now()
	}
	if date.Year() == 0 {
		date = date.AddDate(time.Now().Year(), 0, 0)
	}
	return date
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

func insertIntoStore(insertChan chan string) {
	for {
		length := len(insertChan)
		if length > 0 {
			var logLines []logline.LogLine
			for i := 0; i < length; i++ {
				line := <-insertChan
				logLine, found := parseLine(line, logLines)
				if !found {
					continue
				}
				logLines = append(logLines, logLine)
			}
			for _, line := range logLines {
				insertLogLinesChan <- line
			}

		} else {
			time.Sleep(time.Second)
		}
	}
}

func insertIntoStoreJsonSystem(insertChan chan []byte, system string) {
	for line := range insertChan {
		logLine := parseLineJsonFastJson(line)
		logLine.SetSystem(system)
		insertLogLinesChan <- logLine
	}
}
