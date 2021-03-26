package main

import (
	"github.com/valyala/fastjson"
	"go4.org/intern"
	"strings"
	"time"
)

func parseLine(line string, loglines []LogLine) (LogLine, bool) {
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
					l.setBody(l.getBody() + "\n" + md["body"])
					logLine := LogLine{}
					logLine.setBody(line)
					return logLine, false
				}
				timestamp := toMillis(parseTimestamp(regex, md["timestamp"]))
				logLine := LogLine{
					Time:   timestamp,
					system: intern.GetByString(md["system"]),
				}
				logLine.setLevel(md["level"])
				logLine.setSystem(md["system"])
				logLine.setBody(md["body"])
				return logLine, true

			}
		}
	}
	logLine := LogLine{}
	logLine.setBody(line)
	return logLine, false
}

func parseLineJsonFastJson(line []byte) LogLine {
	body := fastjson.GetString(line, "message")
	stack := fastjson.GetString(line, "stack_trace")
	if len(stack) > 0 {
		body = body + "\n" + stack
	}
	timestamp := fastjson.GetString(line, "@timestamp")
	if len(timestamp) == 0 {
		timestamp = fastjson.GetString(line, "timestamp")
	}
	l := LogLine{
		Time: toMillis(parseTimestampJson(timestamp)),
	}
	l.setSystem(fastjson.GetString(line, "HOSTNAME"))
	l.setLevel(fastjson.GetString(line, "level"))
	l.setBody(body)
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
			var logLines []LogLine
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
		logLine.setSystem(system)
		insertLogLinesChan <- logLine
	}
}
