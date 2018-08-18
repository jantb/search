package main

func getDBStatement_log() string {
	return `create table log (
  id  INTEGER PRIMARY KEY,
  time  INTEGER,
  level TEXT,
  body  TEXT,
  CONSTRAINT line_unique UNIQUE (time, body)
);
CREATE  INDEX index_timestamp ON log(time);
`
}

func initDB() {

}
