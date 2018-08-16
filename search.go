package main

import (
	"database/sql"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strconv"
	"strings"
	"time"
)

var timestamps []time.Time
var ids []int

func searchUp(query string, offset int) {

	gui.Update(func(g *gocui.Gui) error {
		view, e := gui.View("logs")
		checkErr(e)
		_, y := view.Size()
		if len(timestamps) < offset {
			return nil
		}
		timestamp := timestamps[offset-1].Format("2006-01-02 15:04:05.999999999")
		id := strconv.Itoa(ids[offset-1])

		var q = "select id, time, body from log where (time,id) < ('" + timestamp + "'," + id + ") and body like '%" +
			strings.TrimSpace(query) +
			"%' order by time desc, id desc limit " + strconv.Itoa(y-1)
		rows, err := db.Query(q)
		fmt.Print(q)
		checkErr(err)
		view.Rewind()
		view.Clear()
		var lines []string
		timestamps = timestamps[:0]
		ids = ids[:0]
		var count = 0
		for rows.Next() {
			var id int
			var timestamp time.Time
			var body string
			err = rows.Scan(&id, &timestamp, &body)
			if err != nil {
				log.Fatal(err)
			}
			lines = append(lines, fmt.Sprintf("%s %d %s", timestamp.Format("2006-01-02T15:04:05.999999999"), id, body))
			timestamps = append(timestamps, timestamp)
			ids = append(ids, id)
			count++
			if count == 1 {
				fmt.Print(" ", id)
			}
		}
		defer rows.Close()

		for i := len(lines) - 1; i >= 0; i-- {
			fmt.Fprintln(view, lines[i])
		}
		return nil
	},
	)

}
func searchDown(query string, offset int) {

	gui.Update(func(g *gocui.Gui) error {
		view, e := gui.View("logs")
		checkErr(e)
		_, y := view.Size()
		if len(timestamps) < offset {
			return nil
		}
		timestamp := timestamps[y-2].Format("2006-01-02 15:04:05.999999999")
		id := strconv.Itoa(ids[y-2])

		var q = "select id, time, body from log where (time,id) > ('" + timestamp + "'," + id + ") and body like '%" +
			strings.TrimSpace(query) +
			"%' order by time, id limit " + strconv.Itoa(y-1)
		rows, err := db.Query(q)
		fmt.Print(q)
		checkErr(err)
		view.Rewind()
		view.Clear()
		timestamps = timestamps[:0]
		ids = ids[:0]
		var count = 0
		for rows.Next() {
			var id int
			var timestamp time.Time
			var body string
			err = rows.Scan(&id, &timestamp, &body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(view, "%s %d %s\n", timestamp.Format("2006-01-02T15:04:05.999999999"), id, body)
			timestamps = append(timestamps, timestamp)
			ids = append(ids, id)
			count++
			if count == 1 {
				fmt.Print(" ", id)
			}
		}
		defer rows.Close()
		reverseNumbers(ids)
		reverseTime(timestamps)
		return nil
	},
	)

}
func reverseNumbers(numbers []int) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}
func reverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func reverseTime(time []time.Time) {
	for i, j := 0, len(time)-1; i < j; i, j = i+1, j-1 {
		time[i], time[j] = time[j], time[i]
	}
}
func search(query string) {

	gui.Update(func(g *gocui.Gui) error {
		view, e := gui.View("logs")
		checkErr(e)
		_, y := view.Size()
		var q = "select id, time, body from log where body like '%" +
			strings.TrimSpace(query) +
			"%' order by time desc, id desc limit " + strconv.Itoa(y-1)

		rows, err := db.Query(q)

		checkErr(err)
		view.Rewind()
		view.Clear()
		var lines []string
		timestamps = timestamps[:0]
		ids = ids[:0]
		for rows.Next() {
			var id int
			var timestamp time.Time
			var body string
			err = rows.Scan(&id, &timestamp, &body)
			if err != nil {
				log.Fatal(err)
			}
			lines = append(lines, fmt.Sprintf("%s %s", timestamp.Format("2006-01-02T15:04:05.999999999"), body))
			timestamps = append(timestamps, timestamp)
			ids = append(ids, id)
		}
		defer rows.Close()

		for i := len(lines) - 1; i >= 0; i-- {
			fmt.Fprintln(view, lines[i])
		}
		return nil
	},
	)

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
	db, err := sql.Open("sqlite3", "./.search.db")

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
	timestamps = timestamps[:0]
	ids = ids[:0]
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
	defer db.Close()
	if len(ret) != len(prev) && len(prev) > 0 {
		return prev
	}
	return ret
}
