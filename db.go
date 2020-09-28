// +build !mem

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbStatement = `
CREATE TABLE IF NOT EXISTS log (
  id  INTEGER PRIMARY KEY,
  time  INTEGER,
  system TEXT, 
  level TEXT,
  body  TEXT
);

CREATE INDEX IF NOT EXISTS index_timestamp ON log(time);

CREATE TABLE IF NOT EXISTS settings(
	key TEXT PRIMARY KEY ,
	value TEXT
);
`

var prev []LogLine

func initStore() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := sql.Open("sqlite3", filepath.Join(usr.HomeDir, ".search.db"))
	db = dbs
	checkErr(err)

	_, err = db.Exec(dbStatement)
	clearDb()
	checkErr(err)
}

func cleanupStore() {
	db.Close()
}

func clearDb() {
	tx, err := db.Begin()
	checkErr(err)
	_, err = tx.Exec("delete from log")
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()
}

func insertIntoStoreByChan(insertChan chan []LogLine) {
	for {
		line := <-insertChan
		insertLoglinesToStore(line)
	}
}

func insertLoglinesToStore(logLines []LogLine) {
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("insert into log(time, system, level, body) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	defer tx.Commit()
	for _, logLine := range logLines {
		_, err = stmt.Exec(logLine.Time, logLine.System, logLine.Level, logLine.Body)
		checkErr(err)
	}
}

func search(query string, limit int, offset int) (ret []LogLine, t time.Duration) {
	tokens := strings.Split(strings.TrimSpace(query), " ")

	if strings.HasPrefix(tokens[0], "level=") {
		query = "level like '" + strings.Split(tokens[0], "=")[1] + "' "
	} else {
		query = "body like '%%" + tokens[0] + "%%' "
	}

	if len(tokens[0]) > 0 && tokens[0][0] == '!' {
		query = "body not like '%%" + tokens[0][1:] + "%%' "
	}

	for _, value := range tokens[1:] {
		if strings.HasPrefix(value, "level=") {
			query += " and level like '" + strings.Split(value, "=")[1] + "' "
		} else if len(value) > 0 && value[0] == '!' {
			query += "and body not like '%%" + value[1:] + "%%' "
		} else {
			query += "and body like '%%" + value + "%%' "
		}
	}
	now := time.Now()
	var q = ""
	if offset > 0 && len(prev) >= offset+1 {
		q = fmt.Sprintf("select id, time, system, level, body from log where (time,id) <= (%d,%d) and "+query+
			" order by time desc, id desc limit "+
			strconv.Itoa(limit), prev[len(prev)-offset-1].Time, prev[len(prev)-offset-1].Id)
		bottom.Store(false)
	} else if offset < 0 && len(prev) >= -offset+1 {
		o := -offset
		q = fmt.Sprintf("select id, time, system, level, body from log where (time,id) >= (%d,%d) and "+query+"order by time , id limit "+strconv.Itoa(limit), prev[o].Time, prev[o].Id)
		bottom.Store(false)
	} else if offset == 0 {
		prev = prev[:0]
		q = fmt.Sprintf("select id, time, system, level, body from log where " + query + "order by time desc, id desc limit " + strconv.Itoa(limit))
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
		err = rows.Scan(&line.Id, &line.Time, &line.System, &line.Level, &line.Body)
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
	prev = ret
	return ret, time.Now().Sub(now)
}

func storeSettings(key, value string) {
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("insert or replace into settings(key, value) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	defer tx.Commit()
	stmt.Exec(key, strings.TrimSpace(value))
}

func loadSettings(key string) string {
	query, err := db.Query("SELECT VALUE FROM settings WHERE key = ?", key)
	checkErr(err)
	if query.Next() {
		value := ""
		err = query.Scan(&value)
		checkErr(err)
		defer query.Close()
		return value
	}
	return ""
}
