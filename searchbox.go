package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
	"time"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
	"strings"
)

var mutex = &sync.Mutex{}
var mutex2 = &sync.Mutex{}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

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
	events         []Event
	text           []byte
	seek           int64
	line_voffset   int
	eventChan      chan []Event
	quitSearch     chan bool
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
	stats          time.Duration
	storeLine      time.Duration
	count          int
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

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset + size)
	eb.Search()
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset + size)
	eb.Search()
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursor_boffset]
	eb.Search()
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
	go eb.Search()
	eb.MoveCursorOneRuneForward()
}

func (eb *EditBox) ScrollUp() {
	mutex.Lock()
	eb.seek++
	mutex.Unlock()
	eb.Search()
}

func (eb *EditBox) ScrollDown() {
	mutex.Lock()
	if eb.seek > 0 {
		eb.seek--
	}
	mutex.Unlock()

	eb.Search()
}
func (eb *EditBox) Follow() {
	mutex.Lock()
	eb.seek = int64(0)
	mutex.Unlock()

	eb.Search()
}

func (eb *EditBox) Search() {
	_, h := termbox.Size()
	mutex.Lock()
	close(eb.quitSearch)
	eb.quitSearch = make(chan bool)
	mutex.Unlock()
	go SearchFor(eb.text, h - 2, eb.seek, eb.eventChan, eb.quitSearch)
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

var edit_box EditBox

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

func redraw_all() {
	mutex.Lock()
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
		var meta Meta
		meta.Unmarshal(by)
		nodecount = meta.Count
		return nil
	})
	count := ""
	if edit_box.count > 0 {
		count = fmt.Sprintf("Count: %d", edit_box.count)
	}
	ns := fmt.Sprintf("%s Parse line: %s Search: %s Events: %d", count, edit_box.storeLine, edit_box.stats, nodecount)
	for i, r := range ns {
		termbox.SetCell(w - len(ns) + i, h - 1, r, coldef, coldef)
	}
	termbox.Flush()

	mutex.Unlock()
}

