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
	maxItems int
	length   int
	root     *Node
}

// New returns a new B-tree with the max number of items.
func New(maxItems int) *Tree {
	if maxItems < 2 {
		panic("bad maxItems")
	}
	return &Tree{maxItems: maxItems}
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
	return t.maxItems
}

// MinItems returns the min number of items to allow per node (ignored for the root node).
func (t *Tree) MinItems() int {
	return t.maxItems / 2
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
	median, left, right, ok := t.root.insert(item, false)
	if median != nil {
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

func (n *Node) insert(item Item, nonleaf bool) (median Item, left *Node, right *Node, ok bool) {
	if len(n.children) == 0 || nonleaf {
		if len(n.items) >= n.MaxItems() {
			return n.split(item)
		}
		_, ok = n.items.insert(item)
		return
	}
	for i := 0; i < len(n.items); i++ {
		if item.Less(n.items[i]) {
			median, left, right, ok = n.children[i].insert(item, false)
			if median != nil {
				n.children.insert(i, left, right)
				return n.insert(median, true)
			}
			return
		}
		if !n.items[i].Less(item) {
			n.items[i] = item
			return nil, nil, nil, false
		} else if i == len(n.items)-1 {
			median, left, right, ok = n.children[i+1].insert(item, false)
			if median != nil {
				n.children.insert(i+1, left, right)
				return n.insert(median, true)
			}
			return
		}
	}
	return
}

func (n *Node) split(item Item) (median Item, left *Node, right *Node, ok bool) {
	index := n.MinItems()
	median = n.items[index]
	compare := 0
	if item.Less(median) {
		compare = -1
	} else if median.Less(item) {
		compare = 1
	}
	if compare == 0 {
		n.items[index] = item
		return nil, nil, nil, false
	}
	right = newNode(n.MaxItems())
	right.items = append(right.items, n.items[index+1:]...)
	n.items = n.items[:index]
	if len(n.children) > 0 {
		right.children = append(right.children, n.children[index+1:]...)
		n.children = n.children[:index+1]
	}
	left = n
	if compare < 0 {
		left.items.insert(item)
		ok = true
		return
	}
	right.items.insert(item)
	ok = true
	return
}

type items []Item

func (s *items) insert(item Item) (index int, ok bool) {
	for i := 0; i < len(*s); i++ {
		if item.Less((*s)[i]) {
			*s = append(*s, nil)
			copy((*s)[i+1:], (*s)[i:])
			(*s)[i] = item
			return i, true
		}
		if !(*s)[i].Less(item) {
			(*s)[i] = item
			return i, false
		}
	}
	*s = append(*s, item)
	return len(*s) - 1, true
}

type children []*Node

func (s *children) insert(index int, left, right *Node) {
	(*s)[index] = left
	*s = append(*s, nil)
	if index+1 < len(*s) {
		copy((*s)[index+2:], (*s)[index+1:])
	}
	(*s)[index+1] = right
}
