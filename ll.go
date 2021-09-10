package main

import (
	"container/list"
	"github.com/jantb/search/logline"
	"sync"
)

type LL struct {
	l       *list.List
	m       sync.Mutex
	systems []*list.Element
}

func (ll *LL) Init() {
	ll.l.Init()
}

func (ll *LL) GetSize() int {
	ll.m.Lock()
	defer ll.m.Unlock()
	return ll.l.Len()
}

func (ll *LL) Put(line logline.LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()

	curr := ll.l.Front()
	for i := range ll.systems {
		if ll.systems[i].Value.(logline.LogLine).System == line.System {
			curr = ll.systems[i]
			break
		}
	}
	element := curr
	if curr == nil {
		ll.l.PushFront(line)
	} else if curr.Value.(logline.LogLine).Time <= line.Time {
		for curr != nil && curr.Value.(logline.LogLine).Time <= line.Time {
			curr = curr.Prev()
		}

		if curr == nil {
			element = ll.l.PushFront(line)
		} else {
			element = ll.l.InsertAfter(line, curr)
		}
	} else {
		for curr != nil && curr.Value.(logline.LogLine).Time > line.Time {
			curr = curr.Next()
		}
		if curr != nil {
			element = ll.l.InsertBefore(line, curr)
		} else {
			element = ll.l.PushBack(line)
		}
	}
	if element != nil {
		for i := range ll.systems {
			if ll.systems[i].Value.(logline.LogLine).System == line.System {
				ll.systems[i] = element
				return
			}
		}
		ll.systems = append(ll.systems, element)
	}
}

func (ll *LL) Iterate(done <-chan struct{}) <-chan logline.LogLine {
	out := make(chan logline.LogLine)
	go func(out chan logline.LogLine) {
		ll.iterate(done, out)
		close(out)
	}(out)
	return out
}

func (ll *LL) iterate(done <-chan struct{}, ch chan<- logline.LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()
	for i := ll.l.Front(); i != nil; i = i.Next() {
		select {
		case ch <- i.Value.(logline.LogLine):
		case <-done:
			return
		}
	}
}

func (ll *LL) RemoveLast() {
	ll.m.Lock()
	defer ll.m.Unlock()
	ll.l.Remove(ll.l.Back())
}
