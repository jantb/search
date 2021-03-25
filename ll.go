package main

import "sync"

type LL struct {
	Head *LNode
	Tail *LNode
	size int
	m    sync.Mutex
}

type LNode struct {
	next *LNode
	prev *LNode
	val  LogLine
}

func (ll *LL) Put(line LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()
	if ll.size == 0 {
		ll.Head = &LNode{
			next: nil,
			prev: nil,
			val:  line,
		}
		ll.Tail = ll.Head
	} else {
		curr := ll.Head

		if curr.val.Time <= line.Time {
			ll.prepend(&LNode{
				next: nil,
				prev: nil,
				val:  line,
			})
			return
		} else {
			for curr != nil && curr.val.Time > line.Time {
				curr = curr.next
			}
			l := &LNode{
				next: nil,
				prev: nil,
				val:  line,
			}
			if curr == nil {
				tail := ll.Tail
				ll.Tail = l
				tail.next = l
				l.prev = tail
			} else {
				prev := curr.prev
				prev.next = l
				l.next = curr
				l.prev = prev
			}
		}
	}

	ll.size++
}

func (ll *LL) prepend(lNode *LNode) {
	if ll.size == 0 {
		ll.Head = lNode
		ll.Tail = lNode
	} else {
		head := ll.Head
		head.prev = lNode
		lNode.next = head
		ll.Head = lNode
	}
	ll.size++
}

func (ll *LL) Iterate(done <-chan struct{}) <-chan LogLine {
	out := make(chan LogLine)
	go func() {
		ll.iterate(done, out)
		close(out)
	}()
	return out
}

func (ll *LL) iterate(done <-chan struct{}, ch chan<- LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()
	for i := ll.Head; i != nil; i = i.next {
		select {
		case ch <- i.val:
		case <-done:
			return
		}
	}
}

func (ll *LL) RemoveLast() {
	ll.m.Lock()
	defer ll.m.Unlock()
	ll.Tail.prev.next = nil
	ll.size--
}
