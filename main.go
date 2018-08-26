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
	"fmt"
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

func settingsDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	v.SetCursor(x, y-1)
	return nil
}
func settingsUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	v.SetCursor(x, y+1)
	return nil
}
func settingsEnter(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("prompt")
	g.SetCurrentView("prompt")
	view, err := g.View("prompt")
	checkErr(err)
	view.Clear()
	view.SetCursor(0, 0)
	return nil
}

var settings = ""

func enterClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("prompt")
	g.SetCurrentView("settings")

	settings, err := g.View("settings")
	checkErr(err)
	_, settingsY := settings.Cursor()

	switch settingsY {
	case 0:
		storeSettings("username", v.Buffer())
	case 1:
		storeSettings("password", v.Buffer())
	case 2:
		storeSettings("openshift", v.Buffer())
	case 3:
		storeSettings("jenkins", v.Buffer())
	case 4:
		storeSettings("bitbucket", v.Buffer())
	case 5:
		storeSettings("jira", v.Buffer())
	default:
	}

	//buffer := v.Buffer()
	//fmt.Print(buffer)
	settings.Clear()
	v.SetCursor(0,0)

	fmt.Fprintln(settings, "Username       ", loadSettings("username"))
	fmt.Fprintln(settings, "Password       ", "********")
	fmt.Fprintln(settings, "Openshift url  ", loadSettings("openshift"))
	fmt.Fprintln(settings, "Jenkins url    ", loadSettings("jenkins"))
	fmt.Fprintln(settings, "Bitbucket url  ", loadSettings("bitbucket"))
	fmt.Fprintln(settings, "Jira url       ", loadSettings("jira"))
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
	if err := g.SetKeybinding("settings", gocui.KeyArrowDown, gocui.ModNone, settingsUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyArrowUp, gocui.ModNone, settingsDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyEnter, gocui.ModNone, settingsEnter); err != nil {
		return err
	}
	if err := g.SetKeybinding("prompt", gocui.KeyEnter, gocui.ModNone, enterClose); err != nil {
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
	viewPrompt(g, maxX, maxY)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
