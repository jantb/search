/*
Copyright (c) 2015, Emir Pasic
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

-------------------------------------------------------------------------------

AVL Tree:

Copyright (c) 2017 Benjamin Scher Purcell <benjapurcell@gmail.com>

Permission to use, copy, modify, and distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/
package main

import (
	"fmt"
	"strings"
)

type Tree struct {
	Root *Node
	size int
}

type Node struct {
	Key      string
	val      LogLine
	Parent   *Node
	Children [2]*Node
	b        int8
}

func (t *Tree) Put(val LogLine) {
	key := fmt.Sprintf("%d%d", val.Time, t.size)
	t.put(key, val, nil, &t.Root)
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Iterate(done <-chan struct{}) <-chan LogLine {
	out := make(chan LogLine)
	go func() {
		t.iterate(done, out)
		close(out)
	}()
	return out
}

func (t *Tree) iterate(done <-chan struct{}, ch chan<- LogLine) {
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		select {
		case ch <- it.Value():
		case <-done:
			return
		}
	}
}

func (n *Node) Prev() *Node {
	return n.walk1(0)
}

func (n *Node) Next() *Node {
	return n.walk1(1)
}
func (n *Node) walk1(a int) *Node {
	if n == nil {
		return nil
	}

	if n.Children[a] != nil {
		n = n.Children[a]
		for n.Children[a^1] != nil {
			n = n.Children[a^1]
		}
		return n
	}

	p := n.Parent
	for p != nil && p.Children[a] == n {
		n = p
		p = p.Parent
	}
	return p
}

func (t *Tree) put(key string, value LogLine, p *Node, qp **Node) bool {
	q := *qp
	if q == nil {
		t.size++
		*qp = &Node{Key: key, val: value, Parent: p}
		return true
	}

	c := strings.Compare(key, q.Key)
	if c == 0 {
		q.Key = key
		q.val = value
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	var fix bool
	fix = t.put(key, value, q, &q.Children[a])
	if fix {
		return putBalance(int8(c), qp)
	}
	return false
}

func putBalance(c int8, t **Node) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.Children[(c+1)/2].b == c {
		s = singlerotation(c, s)
	} else {
		s = doublerotRotation(c, s)
	}
	*t = s
	return false
}

func singlerotation(c int8, s *Node) *Node {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doublerotRotation(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.Children[a]
	s.Children[a] = rotate(-c, s.Children[a])
	p := rotate(c, s)

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.Children[a]
	s.Children[a] = r.Children[a^1]
	if s.Children[a] != nil {
		s.Children[a].Parent = s
	}
	r.Children[a^1] = s
	r.Parent = s.Parent
	s.Parent = r
	return r
}

type Iterator struct {
	tree     *Tree
	node     *Node
	position position
}

func (iterator *Iterator) Next() bool {
	switch iterator.position {
	case begin:
		iterator.position = between
		iterator.node = iterator.tree.Right()
	case between:
		iterator.node = iterator.node.Prev()
	}

	if iterator.node == nil {
		iterator.position = end
		return false
	}
	return true
}

func (iterator *Iterator) Value() LogLine {
	if iterator.node == nil {
		return LogLine{}
	}
	return iterator.node.val
}

func (iterator *Iterator) Key() string {
	if iterator.node == nil {
		return ""
	}
	return iterator.node.Key
}

func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = end
}

func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}
func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}
func (iterator *Iterator) Prev() bool {
	switch iterator.position {
	case end:
		iterator.position = between
		iterator.node = iterator.tree.Right()
	case between:
		iterator.node = iterator.node.Prev()
	}

	if iterator.node == nil {
		iterator.position = begin
		return false
	}
	return true
}

func (tree *Tree) Iterator() ReverseIteratorWithKey {
	return &Iterator{tree: tree, node: nil, position: begin}
}

func (t *Tree) Right() *Node {
	return t.bottom(1)
}
func (t *Tree) Left() *Node {
	return t.bottom(0)
}
func (t *Tree) bottom(d int) *Node {
	n := t.Root
	if n == nil {
		return nil
	}

	for c := n.Children[d]; c != nil; c = n.Children[d] {
		n = c
	}
	return n
}

type ReverseIteratorWithKey interface {
	// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
	// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	Prev() bool

	// End moves the iterator past the last element (one-past-the-end).
	// Call Prev() to fetch the last element if any.
	End()

	// Last moves the iterator to the last element and returns true if there was a last element in the container.
	// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	Last() bool

	IteratorWithKey
}

type IteratorWithKey interface {
	// Next moves the iterator to the next element and returns true if there was a next element in the container.
	// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
	// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Value returns the current element's value.
	// Does not modify the state of the iterator.
	Value() LogLine

	// Key returns the current element's key.
	// Does not modify the state of the iterator.
	Key() string

	// Begin resets the iterator to its initial state (one-before-first)
	// Call Next() to fetch the first element if any.
	Begin()

	// First moves the iterator to the first element and returns true if there was a first element in the container.
	// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	First() bool
}
