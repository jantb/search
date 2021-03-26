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

func (ll *LL) GetSize() int {
	ll.m.Lock()
	defer ll.m.Unlock()
	return ll.size
}
func (ll *LL) Put(line LogLine) {
	ll.m.Lock()
	defer ll.m.Unlock()
	lNode := &LNode{
		next: nil,
		prev: nil,
		val:  line,
	}

	if ll.size == 0 {
		ll.Head = lNode
		ll.Tail = ll.Head
	} else {
		curr := ll.Head

		if curr.val.Time <= line.Time {
			head := ll.Head
			head.prev = lNode
			lNode.next = head
			ll.Head = lNode
		} else {
			for curr != nil && curr.val.Time > line.Time {
				curr = curr.next
			}

			if curr != nil {
				prev := curr.prev
				prev.next = lNode
				lNode.next = curr
				lNode.prev = prev
			} else {
				tail := ll.Tail
				ll.Tail = lNode
				tail.next = lNode
				lNode.prev = tail
			}
		}
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
