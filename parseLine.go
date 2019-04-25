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
				r2 := regex.RegexCompiled.FindAllStringSubmatch(string(line), -1)[0]
				md := map[string]string{}
				for i, n := range r2 {
					md[n1[i]] = n
				}
				if _, ok := md["timestamp"]; !ok {
					if len(loglines) == 0 {
						continue
					}
					loglines[len(loglines)-1].Body += "\n" + md["body"]
					return LogLine{Body: line}, false
				}
				timestamp := toMillis(parseTimestamp(regex, md["timestamp"]))
				return LogLine{
					Time:   timestamp,
					System: md["system"],
					Level:  md["level"],
					Body:   md["body"],
				}, true

			}
		}
	}
	return LogLine{Body: line}, false
}
func parseLineJson(line map[string]interface{}) ([]LogLine, bool) {
	var logLines []LogLine

	l := LogLine{
		Time:   toMillis(parseTimestampJson(cast(line, "@timestamp"))),
		System: cast(line, "consumerId"),
		Level:  cast(line, "level"),
		Body:   cast(line, "message"),
	}
	s := cast(line, "stack_trace")
	logLines = append(logLines, l)
	if len(s) == 0 {
		return logLines, true
	}

	i := strings.Split(s, "\n")
	for _, v := range i {
		l := LogLine{
			Time: toMillis(parseTimestampJson(cast(line, "@timestamp"))),
			Body: v,
		}
		logLines = append(logLines, l)
	}

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
	checkErr(e)
	if date.Year() == 0 {
		date = date.AddDate(time.Now().Year(), 0, 0)
	}
	return date
}
func parseTimestampJson(timestamp string) time.Time {
	date, e := time.ParseInLocation("2006-01-02T15:04:05.999-07:00", strings.Replace(timestamp, ",", ".", -1), time.Local)
	checkErr(e)
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
			insertLoglinesToStore(logLines)
		} else {
			time.Sleep(time.Second)
		}

		if bottom.Load() {
			v, e := gui.View("commands")
			checkErr(e)
			renderSearch(v, 0)
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
				insertLoglinesToStore(logLines)
			}
		} else {
			time.Sleep(time.Second)
		}

		if bottom.Load() {
			v, e := gui.View("commands")
			checkErr(e)
			renderSearch(v, 0)
		}
	}
}
