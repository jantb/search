// +build mem

package main

import (
	"sort"
	"time"
)

var store []LogLine
var prevs []LogLine
var id int

func initStore() {

}

func cleanupStore() {

}

func insertLoglinesToStore(logLines []LogLine) {
	for _, line := range logLines {
		id++
		line.Id = id
		insertSort(line)
	}
}

func insertSort(line LogLine) {
	index := sort.Search(len(store), func(i int) bool {
		if store[i].Time == line.Time {
			return store[i].Id > line.Id
		}
		return store[i].Time > line.Time
	})
	store = append(store, LogLine{})
	copy(store[index+1:], store[index:])
	line.Id = index + 1
	store[index] = line
}

func search(query string, limit int, offset int) (ret []LogLine, t time.Duration) {
	now := time.Now()
	if len(store) == 0 {
		return store, time.Now().Sub(now)
	}
	//fmt.Print(store[0])
	if len(prevs) == 0 {
		lines := store[Max(0, len(store)-limit):]
		prevs = make([]LogLine, len(lines))
		copy(prevs, lines)
		reverseLogline(prevs)
	}
	offsetLine := store[len(store)-1]
	if offset >= 0 {
		offsetLine = prevs[Min(len(prevs)-1-offset, 0)]
	}
	var matches []LogLine
	for i := len(store) - 1; i >= 0; i-- {
		line := store[i]
		if line.Time > offsetLine.Time || line.Id > offsetLine.Id {
			continue
		}

		matches = append(matches, line)

		if len(matches) == limit {
			break
		}
	}
	//fmt.Print(offsetLine)
	//for _, value := range prevs {
	//
	//}

	ret = make([]LogLine, len(matches))
	copy(ret, matches)
	reverseLogline(ret)
	prevs = ret
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
