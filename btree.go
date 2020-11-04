// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package btree implements a B-tree.
//
// B-tree properties: https://en.wikipedia.org/wiki/B-tree
//
// 1. Every node has at most m children.
//
// 2. Every non-leaf node (except root) has at least ⌈m/2⌉ child nodes.
//
// 3. The root has at least two children if it is not a leaf node.
//
// 4. A non-leaf node with k children contains k − 1 keys.
//
// 5. All leaves appear in the same level and carry no information.
//
package btree

import (
	"fmt"
	"sort"
)

// Item represents a value in the tree.
type Item interface {
	// Less compares whether the current item is less than the given Item.
	Less(than Item) bool
}

// Int implements the Item interface for int.
type Int int

// Less returns true if int(a) < int(b).
func (a Int) Less(b Item) bool {
	return a < b.(Int)
}

// String implements the Item interface for string.
type String string

// Less returns true if string(a) < string(b).
func (a String) Less(b Item) bool {
	return a < b.(String)
}

// Tree represents a B-tree.
type Tree struct {
	degree int
	length int
	root   *Node
}

// New returns a new B-tree with the given degree.
func New(degree int) *Tree {
	if degree <= 1 {
		panic("bad degree")
	}
	return &Tree{degree: degree}
}

// Length returns the number of items currently in the B-tree.
func (t *Tree) Length() int {
	return t.length
}

// Root returns the root node of the B-tree.
func (t *Tree) Root() *Node {
	return t.root
}

// MaxItems returns the max number of items to allow per Node.
func (t *Tree) MaxItems() int {
	return t.degree*2 - 1
}

// MinItems returns the min number of items to allow per node (ignored for the root node).
func (t *Tree) MinItems() int {
	return t.degree - 1
}

// Insert inserts the item into the B-tree.
func (t *Tree) Insert(item Item) {
	if item == nil {
		panic("nil item being inserted to tree")
	}
	if t.root == nil {
		t.root = newNode(t.MaxItems())
		t.root.items = append(t.root.items, item)
		t.length++
		return
	}
	median, right, ok := t.root.insert(item, false)
	if median != nil {
		left := t.root
		t.root = newNode(t.MaxItems())
		t.root.items = append(t.root.items, median)
		t.root.children = append(t.root.children, left, right)
	}
	if ok {
		t.length++
	}
	return
}

// Node represents a node in the B-tree.
type Node struct {
	items    items
	children children
}

func newNode(maxItems int) *Node {
	return &Node{items: make([]Item, 0, maxItems), children: make([]*Node, 0, maxItems+1)}
}

// MaxItems returns the max number of items to allow per Node.
func (n *Node) MaxItems() int {
	if n == nil {
		return 0
	}
	return cap(n.items)
}

// MinItems returns the min number of items to allow per node (ignored for the root node).
func (n *Node) MinItems() int {
	if n == nil {
		return 0
	}
	return cap(n.items) / 2
}

// Items returns the items of this node.
func (n *Node) Items() []Item {
	if n == nil {
		return nil
	}
	return n.items
}

// Children returns the children of this node.
func (n *Node) Children() []*Node {
	if n == nil {
		return nil
	}
	return n.children
}

func (n *Node) insert(item Item, nonleaf bool) (median Item, right *Node, ok bool) {
	i, existed := n.items.search(item)
	if existed {
		n.items[i] = item
		ok = false
		return
	}
	if len(n.children) == 0 || nonleaf {
		if len(n.items) < n.MaxItems() {
			n.items.insert(i, item)
			ok = true
			return
		}
		return n.split(item)
	}
	median, right, ok = n.children[i].insert(item, false)
	if median != nil {
		m := median
		r := right
		median, right, ok = n.insert(median, true)
		index, found := n.items.search(m)
		if found {
			n.children.insert(index+1, r)
			return
		}
		if right != nil {
			index, found := right.items.search(m)
			if found {
				right.children.insert(index+1, r)
			}
		}
	}
	return
}

func (n *Node) split(item Item) (median Item, right *Node, ok bool) {
	ok = true
	i := n.MinItems()
	median = n.items[i]
	right = newNode(n.MaxItems())
	right.items = append(right.items, n.items[i+1:]...)
	n.items = n.items[:i]
	if len(n.children) > 0 {
		right.children = append(right.children, n.children[i+1:]...)
		n.children = n.children[:i+1]
	}
	if item.Less(median) {
		index, _ := n.items.search(item)
		n.items.insert(index, item)
	} else {
		index, _ := right.items.search(item)
		right.items.insert(index, item)
	}
	return
}

type items []Item

func (s *items) insert(index int, item Item) {
	*s = append(*s, nil)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = item
}

func (s *items) search(item Item) (index int, ok bool) {
	i := sort.Search(len(*s), func(i int) bool {
		return item.Less((*s)[i])
	})
	if i > 0 && !(*s)[i-1].Less(item) {
		return i - 1, true
	}
	return i, false
}

type children []*Node

func (s *children) insert(index int, node *Node) {
	*s = append(*s, nil)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = node
}

func (s *children) string() (str string) {
	for i := 0; i < len(*s); i++ {
		if (*s)[i] != nil {
			str += fmt.Sprintf("%v", (*s)[i].items)
		} else {
			str += fmt.Sprintf("  ")
		}

	}
	return
}
