package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var db *sql.DB

var formats Formats
var gui *gocui.Gui

func main() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(filepath.Join(usr.HomeDir, ".search.db"))

	dbs, err := sql.Open("sqlite3", filepath.Join(usr.HomeDir, ".search.db"))
	db = dbs
	checkErr(err)
	defer db.Close()

	_, err = db.Exec(getDBStatement_log())
	checkErr(err)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	defer g.Close()
	go readFromPipe()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	//tx, err := db.Begin()
	//checkErr(err)
	//stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer stmt.Close()
	//for i := 0; i < 100; i++ {
	//	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//tx.Commit()
	//
	//rows, err := db.Query("select id, name from foo")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id int
	//	var name string
	//	err = rows.Scan(&id, &name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(id, name)
	//}
	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//stmt, err = db.Prepare("select name from foo where id = ?")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer stmt.Close()
	//var name string
	//err = stmt.QueryRow("3").Scan(&name)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(name)
	//
	//_, err = db.Exec("delete from foo")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//rows, err = db.Query("select id, name from foo")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id int
	//	var name string
	//	err = rows.Scan(&id, &name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(id, name)
	//}
	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}
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

func readFromPipe() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}
	readFormats()
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			fmt.Print("No more to read, terminating")
			break
		}
		if input == '\n' {
			line := string(output)
			for _, format := range formats {
				for _, regex := range format.Regex {
					r, _ := regexp.Compile(regex.Regex)
					match := r.Match([]byte(line))
					if match {
						n1 := r.SubexpNames()
						r2 := r.FindAllStringSubmatch(line, -1)[0]
						md := map[string]string{}
						for i, n := range r2 {
							md[n1[i]] = n
						}
						timestamp := toMillis(parseTimestamp(regex, md["timestamp"]))
						insertLineToDb("insert into log(time, body) values(?, ?)", timestamp, md["body"])
					}
				}
			}
			output = output[:0]
		}

		output = append(output, input)
	}
}

func toMillis(time time.Time) int64 {
	return time.UnixNano() / 1000000
}

func insertLineToDb(statement string, args ...interface{}) {
	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	checkErr(err)
	//fmt.Println(s.RowsAffected())
	tx.Commit()
}

type Formats []struct {
	Title     string  `json:"title"`
	Multiline bool    `json:"multiline"`
	Regex     []Regex `json:"regex"`
}
type Regex struct {
	Name      string `json:"name"`
	Regex     string `json:"regex"`
	Timestamp string `json:"timestamp"`
}

func readFormats() {
	bytes, err := ioutil.ReadFile("formats.json")
	checkErr(err)
	e := json.Unmarshal(bytes, &formats)
	checkErr(e)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewLogs(g, maxX, maxY)
	viewCommands(g, maxX, maxY)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
