package searchbox

import (
	"github.com/mattn/go-runewidth"
	"unicode/utf8"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	"github.com/gdamore/tcell/termbox"
	"github.com/jantb/search/searchfor"
	"github.com/jantb/search/proto"
	"time"
	"sync/atomic"
	"os"
	"os/signal"
)

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x + lx, y + ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func rune_advance_len(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos % tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += rune_advance_len(r, voffset)
	}
	return
}

func byte_slice_grow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text) - size]
	return text
}

func byte_slice_insert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byte_slice_grow(text, n)
	text = text[:n]
	copy(text[offset + len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 8

type EditBox struct {
	events         []*proto.EventRes
	text           []byte
	seek           int64
	line_voffset   int
	eventChan      chan proto.SearchRes
	quitSearch     chan bool
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
	stats          string
	count          int64
	follow         int32
}

func (eb EditBox) Seek() int64 {
	return eb.seek
}
func (eb *EditBox) SeekSet(seek int64) {
	eb.seek = seek
}

func (eb *EditBox) SeekInc() {
	eb.seek++
}
func (eb *EditBox) SeekDec() {
	if eb.seek > 0 {
		eb.seek--
	}
}
// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x + w - 1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.line_voffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x + rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x + rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
		next:
		t = t[size:]
	}

	if eb.line_voffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.line_voffset != 0 {
		threshold = width - ht
	}
	if eb.cursor_voffset - eb.line_voffset >= threshold {
		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
	}

	if eb.line_voffset != 0 && eb.cursor_voffset - eb.line_voffset < ht {
		eb.line_voffset = eb.cursor_voffset - ht
		if eb.line_voffset < 0 {
			eb.line_voffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursor_boffset = boffset
	eb.cursor_voffset, eb.cursor_coffset = voffset_coffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursor_boffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursor_boffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward(db *bolt.DB) {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset + size)
	eb.Search(db)
}

func (eb *EditBox) DeleteRuneForward(db *bolt.DB) {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset + size)
	eb.Search(db)
}

func (eb *EditBox) DeleteTheRestOfTheLine(db *bolt.DB) {
	eb.text = eb.text[:eb.cursor_boffset]
	eb.Search(db)
}

func (eb *EditBox) InsertRune(r rune, db *bolt.DB) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
	eb.MoveCursorOneRuneForward()
	eb.Search(db)
}

func (eb *EditBox) ScrollUp(db *bolt.DB) {
	eb.SeekInc()
	eb.Search(db)
}

func (eb *EditBox) ScrollDown(db *bolt.DB) {
	eb.SeekDec()
	eb.Search(db)
}
func (eb *EditBox) Follow(db *bolt.DB) {
	eb.follow ^= 1
	eb.seek = int64(0)
	eb.Search(db)
}
func New() *EditBox {
	return &EditBox{

	}
}
func (eb *EditBox) Search(db *bolt.DB) {
	_, h := termbox.Size()
	close(eb.quitSearch)
	eb.quitSearch = make(chan bool)
	eb.count = 0
	go searchfor.SearchFor(eb.text, h - 2, eb.seek, eb.eventChan, eb.quitSearch, db)
}

func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

func insertNewlineAtIInString(in string, i int) (string, int) {
	if i < 21 {
		return in, 0
	}
	split := strings.Split(in, "\n")
	c := 0;
	for j := 0; j < len(split); j++ {
		inn := split[j]
		if len(inn) < i {
			continue
		}
		x := strings.LastIndex(inn[:i], " ")
		if x == -1 {
			continue
		}
		s := []rune(inn)
		s[x] = '\n'
		split[j] = string(s)
		c++
	}

	return strings.Join(split, "\n"), c
}

func redraw_all(edit_box *EditBox, db *bolt.DB) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	var edit_box_width = w
	midy := h - 1
	midx := 0
	fill(midx, midy - 1, edit_box_width, 1, termbox.Cell{Ch: '─', Fg:termbox.ColorBlue})

	edit_box.Draw(midx, midy, edit_box_width, 1)
	termbox.SetCursor(midx + edit_box.CursorX(), midy)
	previ := h - 2

	for i, event := range edit_box.events {
		offset := 21;
		text := event.Data
		i = previ - 2

		for n := 1; n > 0; {
			text, n = insertNewlineAtIInString(text, w - offset)
			i -= n
		}

		i -= int(event.Lines)
		previ = i
		for index, r := range event.Ts {
			if i < h - 2 && i >= 0 {
				termbox.SetCell(index, i, r, termbox.ColorGreen, coldef)
			}
		}

		pastOffset := 0

		for ir, r := range text {
			if r == '\n' {
				i += 1
				pastOffset = ir + 1
				continue
			}
			if i < h - 2 && i >= 0 {
				x := ir + offset - pastOffset
				termbox.SetCell(x, i, r, coldef, coldef)

				for h := 0; h < len(event.FoundAtIndex); h += 2 {
					if int32(ir) >= event.FoundAtIndex[h]&& int32(ir) < event.FoundAtIndex[h + 1] {
						termbox.SetCell(x, i, r, coldef, termbox.ColorRed)
					}
				}

			}
		}
		ns := fmt.Sprintf("Source: %s", event.Path)
		for ix, r := range ns {
			termbox.SetCell(offset + ix, i + 1, r, termbox.ColorCyan, coldef)
		}

	}
	nodecount := int64(0)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Meta"))
		by := b.Get([]byte("Meta"))
		var meta proto.Meta
		meta.Unmarshal(by)
		nodecount = meta.Count
		return nil
	})
	count := ""
	searching := ""
	if atomic.LoadInt32(&searchfor.Searching) != int32(0) {
		searching = "Searching..."
	}
	if edit_box.count > 0 {
		count = fmt.Sprintf("Count: %d", edit_box.count)
	}
	ns := fmt.Sprintf("%s %s Search: %s Events: %d", searching,count, edit_box.stats, nodecount)
	for i, r := range ns {
		termbox.SetCell(w - len(ns) + i, h - 1, r, coldef, coldef)
	}
	termbox.Flush()

}
func Run(db *bolt.DB) {
	edit_box := New()
	edit_box.eventChan = make(chan proto.SearchRes)
	edit_box.quitSearch = make(chan bool)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	searchChan := make(chan bool)
	redrawChan := make(chan bool)
	go func(e *EditBox) {
		for {
			time.Sleep(time.Millisecond * 100)
			if atomic.LoadInt32(&edit_box.follow) == int32(1)&& atomic.LoadInt32(&searchfor.Searching) == int32(0){
				searchChan <- true
			} else {
				redrawChan <- true
			}
		}
	}(edit_box)

	eventChan := make(chan termbox.Event)
	go func() {
		for {
			event := termbox.PollEvent()
			eventChan <- event
		}
	}()
	// register signals to channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	mainloop:
	for {
		select {
		case ev := <-eventChan:
			switch  ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlC:
					break mainloop
				case termbox.KeyArrowLeft, termbox.KeyCtrlB:
					edit_box.MoveCursorOneRuneBackward()
				case termbox.KeyArrowRight, termbox.KeyCtrlF:
					edit_box.MoveCursorOneRuneForward()
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					edit_box.DeleteRuneBackward(db)
				case termbox.KeyDelete, termbox.KeyCtrlD:
					edit_box.DeleteRuneForward(db)
				case termbox.KeyTab:
					edit_box.InsertRune('\t', db)
				case termbox.KeyArrowUp:
					edit_box.ScrollUp(db);
				case termbox.KeyArrowDown:
					edit_box.ScrollDown(db);
				case termbox.KeyPgup:
					edit_box.ScrollUp(db);
				case termbox.KeyPgdn:
					edit_box.ScrollDown(db);
				case termbox.KeySpace:
					edit_box.InsertRune(' ', db)
				case termbox.KeyCtrlG:
					edit_box.Follow(db)
				case termbox.KeyCtrlK:
					edit_box.DeleteTheRestOfTheLine(db)
				case termbox.KeyHome, termbox.KeyCtrlA:
					edit_box.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					edit_box.MoveCursorToEndOfTheLine()
				default:
					if ev.Ch != 0 {
						edit_box.InsertRune(ev.Ch, db)
					}
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		case searchRes := <-edit_box.eventChan:
			edit_box.count = searchRes.Count
			edit_box.stats = searchRes.Ts
			edit_box.events = searchRes.Events
			redraw_all(edit_box, db)
		case <-redrawChan:
			redraw_all(edit_box, db)
		case <-searchChan:
			edit_box.Search(db)
		case <-sigChan:
			break mainloop
		}
	}
}

