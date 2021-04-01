package main

import (
	"container/list"
	"sync"
)

type LL struct {
	l *list.List
	m sync.Mutex
}

func (ll *LL) Init() {
	ll.l.Init()
}

func (ll *LL) GetSize() int {
	ll.m.Lock()
	defer ll.m.Unlock()
	return ll.l.Len()
}
func (ll *LL) Put(line LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()

	curr := ll.l.Front()

	if curr == nil || curr.Value.(LogLine).Time <= line.Time {
		ll.l.PushFront(line)
	} else {
		for curr != nil && curr.Value.(LogLine).Time > line.Time {
			curr = curr.Next()
		}
		if curr != nil {
			ll.l.InsertBefore(line, curr)
		} else {
			ll.l.PushBack(line)
		}
	}
}

func (ll *LL) Iterate(done <-chan struct{}) <-chan LogLine {
	out := make(chan LogLine)
	go func(out chan LogLine) {
		ll.m.Lock()
		ll.iterate(done, out)
		ll.m.Unlock()
		close(out)
	}(out)
	return out
}

func (ll *LL) iterate(done <-chan struct{}, ch chan<- LogLine) {
	for i := ll.l.Front(); i != nil; i = i.Next() {
		select {
		case ch <- i.Value.(LogLine):
		case <-done:
			return
		}
	}
}

func (ll *LL) RemoveLast() {
	ll.m.Lock()
	defer ll.m.Unlock()
	ll.RemoveLast()
}
