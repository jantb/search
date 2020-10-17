package main

import (
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
					system: md["system"],
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

func parseLineJson(line map[string]interface{}) ([]LogLine, bool) {
	var logLines []LogLine
	body := cast(line, "message")
	if len(cast(line, "stack_trace")) > 0 {
		body = cast(line, "message") + "\n" + cast(line, "stack_trace")
	}
	l := LogLine{
		Time: toMillis(parseTimestampJson(cast(line, "@timestamp"))),
	}
	l.setSystem(cast(line, "HOSTNAME"))
	l.setLevel(cast(line, "level"))
	l.setBody(body)
	logLines = append(logLines, l)
	return logLines, true
}

func cast(j map[string]interface{}, field string) string {
	if j[field] != nil {
		return j[field].(string)
	}
	return ""
}

func parseTimestamp(regex Regex, timestamp string) time.Time {
	s := regex.Timestamp
	date, e := time.ParseInLocation(s, strings.Replace(timestamp, ",", ".", -1), time.Local)
	if e != nil {
		date = time.Now().Add(-100 * time.Minute)
	}
	if date.Year() == 0 {
		date = date.AddDate(time.Now().Year(), 0, 0)
	}
	return date
}

func parseTimestampJson(timestamp string) time.Time {
	date, e := time.ParseInLocation("2006-01-02T15:04:05.999-07:00", strings.Replace(timestamp, ",", ".", -1), time.Local)
	if e != nil {
		date = time.Now().Add(-100 * time.Minute)
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
			insertLogLinesChan <- logLines
			//insertLoglinesToStore(logLines)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func insertIntoStoreJsonSystem(insertChan chan map[string]interface{}, system string) {
	for {
		length := len(insertChan)
		if length > 0 {
			for i := 0; i < length; i++ {
				line := <-insertChan
				logLines, found := parseLineJson(line)
				for j := range logLines {
					logLines[j].setSystem(system)
				}
				if !found {
					continue
				}
				insertLogLinesChan <- logLines
				//insertLoglinesToStore(logLines)
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}

func insertIntoStoreJson(insertChan chan map[string]interface{}) {
	for {
		length := len(insertChan)
		if length > 0 {
			for i := 0; i < length; i++ {
				line := <-insertChan
				logLines, found := parseLineJson(line)
				if !found {
					continue
				}
				insertLogLinesChan <- logLines
				//insertLoglinesToStore(logLines)
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}
