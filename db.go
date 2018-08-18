package main

func getDBStatement_log() string {
	return `create table log (
  id  INTEGER PRIMARY KEY,
  time  INTEGER,
  level TEXT,
  body  TEXT
);
CREATE  INDEX index_timestamp ON log(time);

`
}

/*
CREATE VIRTUAL TABLE log_idx USING fts5(level, body, content='log', content_rowid='id');

-- Triggers to keep the FTS index up to date.
CREATE TRIGGER log_ai AFTER INSERT ON log BEGIN
  INSERT INTO fts_idx(rowid, level, body) VALUES (new.id, new.level, new.body);
END;
CREATE TRIGGER log_ad AFTER DELETE ON log BEGIN
  INSERT INTO fts_idx(fts_idx, rowid, level, body) VALUES('delete', old.id, old.level, old.body);
END;
CREATE TRIGGER log_au AFTER UPDATE ON log BEGIN
  INSERT INTO fts_idx(fts_idx, rowid, level, body) VALUES('delete', old.id, old.level, old.body);
  INSERT INTO fts_idx(rowid, level, body) VALUES (new.id, new.level, new.body);
END;
*/
func initDB() {

}
