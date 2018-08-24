package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"go.uber.org/atomic"
)

var formats Formats
var gui *gocui.Gui
var bottom atomic.Bool

func main() {
	bottom.Store(false)

	initStore()

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
					Time:   timestamp,
					System: md["system"],
					Level:  md["level"],
					Body:   md["body"],
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
			insertLoglinesToStore(logLines)
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
	cleanupStore()
	return gocui.ErrQuit
}
func activatePodCommands(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("podCommand")
	g.SetCurrentView("podCommand")

	return nil
}
func deactivatePodCommands(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("podCommand")
	g.SetCurrentView("commands")

	return nil
}
func activateSettings(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("settings")
	g.SetCurrentView("settings")

	return nil
}
func deactivateSettings(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("settings")
	g.SetCurrentView("commands")

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("commands", gocui.KeyCtrlP, gocui.ModNone, activatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("podCommand", gocui.KeyCtrlP, gocui.ModNone, deactivatePodCommands); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, activateSettings); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyCtrlS, gocui.ModNone, deactivateSettings); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewPodCommand(g, maxX, maxY)
	viewSettings(g, maxX, maxY)
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
