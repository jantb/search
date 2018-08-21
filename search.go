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
	Id    int
	Level string
	Time  int64
	Body  string
}

func (l LogLine) getTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}

func s(query string, limit int, offset int, prev []LogLine) (ret []LogLine, t time.Duration) {
	tokens := strings.Split(strings.TrimSpace(query), " ")
	query = "body like '%%" + tokens[0] + "%%' "
	if len(tokens[0]) > 0 && tokens[0][0] == '!' {
		query = "body not like '%%" + tokens[0][1:] + "%%' "
	}

	for _, value := range tokens[1:] {
		if len(value) > 0 && value[0] == '!' {
			query += "and body not like '%%" + value[1:] + "%%' "
		} else {
			query += "and body like '%%" + value + "%%' "
		}
	}
	now := time.Now()
	var q = ""
	if offset > 0 && len(prev) >= offset+1 {
		q = fmt.Sprintf("select id, time, level, body from log where (time,id) <= (%d,%d) and "+query+
			" order by time desc, id desc limit "+
			strconv.Itoa(limit), prev[len(prev)-offset-1].Time, prev[len(prev)-offset-1].Id)
		bottom.Store(false)
	} else if offset < 0 && len(prev) >= -offset+1 {
		o := -offset
		q = fmt.Sprintf("select id, time, level, body from log where (time,id) >= (%d,%d) and "+query+"order by time , id limit "+strconv.Itoa(limit), prev[o].Time, prev[o].Id)
		bottom.Store(false)
	} else if offset == 0 {
		prev = prev[:0]
		q = fmt.Sprintf("select id, time, level, body from log where " + query + "order by time desc, id desc limit " + strconv.Itoa(limit))
		bottom.Store(true)
	} else {
		return prev, time.Now().Sub(now)
	}
	rows, err := db.Query(q)

	if err != nil {
		return prev, time.Now().Sub(now)
	}
	defer rows.Close()

	for rows.Next() {
		line := LogLine{}
		err = rows.Scan(&line.Id, &line.Time, &line.Level, &line.Body)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, line)
	}
	if offset >= 0 {
		reverseLogline(ret)
	}

	if len(ret) != len(prev) && len(prev) > 0 {
		return prev, time.Now().Sub(now)
	}
	return ret, time.Now().Sub(now)
}
