package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func reverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

type LogLine struct {
	Id   int
	Time int64
	Body string
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}

func s(query string, limit int, offset int, prev []LogLine) (ret []LogLine) {

	var q = ""
	if offset > 0 && len(prev) >= offset+1 {
		q = fmt.Sprintf("select id, time, body from log where (time,id) <= (%d,%d) and body like '%%"+
			strings.TrimSpace(query)+
			"%%' order by time desc, id desc limit "+strconv.Itoa(limit), prev[len(prev)-offset-1].Time, prev[len(prev)-offset-1].Id)
	} else if offset < 0 && len(prev) >= -offset+1 {
		o := -offset
		q = fmt.Sprintf("select id, time, body from log where (time,id) >= (%d,%d) and body like '%%"+
			strings.TrimSpace(query)+
			"%%' order by time , id limit "+strconv.Itoa(limit), prev[o].Time, prev[o].Id)
	} else if offset == 0 {
		prev = prev[:0]
		q = fmt.Sprintf("select id, time, body from log where  body like '%%" +
			strings.TrimSpace(query) +
			"%%' order by time desc, id desc limit " + strconv.Itoa(limit))
	} else {
		return prev
	}

	rows, err := db.Query(q)

	checkErr(err)
	for rows.Next() {
		line := LogLine{}
		err = rows.Scan(&line.Id, &line.Time, &line.Body)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, line)
	}
	if offset >= 0 {
		reverseLogline(ret)
	}
	defer rows.Close()
	if len(ret) != len(prev) && len(prev) > 0 {
		return prev
	}
	return ret
}
