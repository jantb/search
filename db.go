// +build !mem

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB
var dbStatement = `
create table log (
  id  INTEGER PRIMARY KEY,
  time  INTEGER,
  level TEXT,
  body  TEXT
);
CREATE  INDEX index_timestamp ON log(time);
`

func initStore() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(filepath.Join(usr.HomeDir, ".search.db"))

	dbs, err := sql.Open("sqlite3", filepath.Join(usr.HomeDir, ".search.db"))
	db = dbs
	checkErr(err)

	_, err = db.Exec(dbStatement)
	checkErr(err)
}

func cleanupStore() {
	db.Close()
}

func insertLoglinesToStore(logLines []LogLine) {
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("insert into log(time, level, body) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	defer tx.Commit()
	for _, logLine := range logLines {
		_, err = stmt.Exec(logLine.Time, logLine.Level, logLine.Body)
		checkErr(err)
	}
}

func search(query string, limit int, offset int, prev []LogLine) (ret []LogLine, t time.Duration) {
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

//CREATE VIRTUAL TABLE log_idx USING fts5(level, body, content='log', content_rowid='id');
//
//-- Triggers to keep the FTS index up to date.
//CREATE TRIGGER log_ai AFTER INSERT ON log BEGIN
//INSERT INTO log_idx(rowid, level, body) VALUES (new.id, new.level, new.body);
//END;
//CREATE TRIGGER log_ad AFTER DELETE ON log BEGIN
//INSERT INTO log_idx(fts_idx, rowid, level, body) VALUES('delete', old.id, old.level, old.body);
//END;
//CREATE TRIGGER log_au AFTER UPDATE ON log BEGIN
//INSERT INTO log_idx(fts_idx, rowid, level, body) VALUES('delete', old.id, old.level, old.body);
//INSERT INTO log_idx(rowid, level, body) VALUES (new.id, new.level, new.body);
//END;
