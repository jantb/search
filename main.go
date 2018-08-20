package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/atomic"
	"io"
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
var bottom atomic.Bool

func main() {
	bottom.Store(false)
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

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	defer g.Close()
	insertChan := make(chan string, 10000)
	go insertIntoDb(insertChan)
	go readFromPipe(insertChan)

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
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
					Time:  timestamp,
					Level: md["level"],
					Body:  md["body"],
				}, true

			}
		}
	}
	return LogLine{Body: line}, false
}

func insertIntoDb(insertChan chan string) {
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
			insertLoglinesToDb(logLines)
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

func insertLoglinesToDb(logLines []LogLine) {
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

func readFromPipe(insertChan chan string) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}
	readFormats()
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			fmt.Print("No more to read, terminating")
			break
		}
		insertChan <- string(line)
	}
}

func toMillis(time time.Time) int64 {
	return time.UnixNano() / 1000000
}

type Formats []struct {
	Title     string  `json:"title"`
	Multiline bool    `json:"multiline"`
	Regex     []Regex `json:"regex"`
}

type Regex struct {
	Name          string `json:"name"`
	Regex         string `json:"regex"`
	RegexCompiled *regexp.Regexp
	Timestamp     string `json:"timestamp"`
}

func readFormats() {
	e := json.Unmarshal([]byte(format), &formats)
	checkErr(e)
	for i, format := range formats {
		for ii, regex := range format.Regex {
			r, _ := regexp.Compile(regex.Regex)
			formats[i].Regex[ii].RegexCompiled = r
		}
	}
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
	viewStatus(g, maxX, maxY)
	viewCommands(g, maxX, maxY)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
