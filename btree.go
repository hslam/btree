// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package btree implements a B-tree.
//
// According to Knuth's definition, a B-tree of order m is a tree which satisfies
// the following properties: https://en.wikipedia.org/wiki/B-tree
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

// Max returns the max node of the B-tree.
func (t *Tree) Max() *Node {
	return t.root.max()
}

// Min returns the min node of the B-tree.
func (t *Tree) Min() *Node {
	return t.root.min()
}

// Search searches the Item of the B-tree.
func (t *Tree) Search(item Item) Item {
	if t.root == nil {
		return nil
	}
	return t.root.search(item)
}

// SearchNode searches the node of the B-tree with the item.
func (t *Tree) SearchNode(item Item) *Node {
	if t.root == nil {
		return nil
	}
	n, _ := t.root.searchNode(item)
	return n
}

// SearchIterator searches the iterator of the B-tree with the item.
func (t *Tree) SearchIterator(item Item) *Iterator {
	if t.root == nil {
		return nil
	}
	n, i := t.root.searchNode(item)
	return n.Iterator(i)
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
		left.parent = t.root
		right.parent = t.root
	}
	if ok {
		t.length++
	}
	return
}

// Clear removes all items from the B-tree.
func (t *Tree) Clear() {
	t.root = nil
	t.length = 0
}

// Delete deletes the node of the B-tree with the item.
func (t *Tree) Delete(item Item) {
	var ok bool
	t.root, ok = t.root.delete(item, -1)
	if t.root != nil && t.root.parent != nil {
		t.root.parent = nil
	}
	if ok {
		t.length--
	}
}

// Node represents a node in the B-tree.
type Node struct {
	items    items
	children children
	parent   *Node
}

func newNode(maxItems int) *Node {
	return &Node{items: make([]Item, 0, maxItems), children: make([]*Node, 0, maxItems+1)}
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

// Parent returns the parent node.
func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	return n.parent
}

// Iterator returns the iterator with the item index of this node.
func (n *Node) Iterator(index int) *Iterator {
	if n == nil {
		return nil
	}
	return &Iterator{node: n, index: index, parentIndex: n.parentIndex()}
}

// MinIterator returns the iterator with the min item index of this node.
func (n *Node) MinIterator() *Iterator {
	if n == nil {
		return nil
	}
	return n.Iterator(0)
}

// MaxIterator returns the iterator with the max item index of this node.
func (n *Node) MaxIterator() *Iterator {
	if n == nil {
		return nil
	}
	return n.Iterator(len(n.items) - 1)
}

func (n *Node) parentIndex() int {
	if n == nil {
		return -1
	}
	var parentIndex = -1
	if n.parent != nil {
		i, _ := n.parent.items.search(n.items[0])
		if i < len(n.parent.children) {
			parentIndex = i
		}
	}
	return parentIndex
}

func (n *Node) maxItems() int {
	if n == nil {
		return 0
	}
	return cap(n.items)
}

func (n *Node) minItems() int {
	if n == nil {
		return 0
	}
	return cap(n.items) / 2
}

func (n *Node) search(item Item) Item {
	i, existed := n.items.search(item)
	if existed {
		return n.items[i]
	}
	if i < len(n.children) {
		return n.children[i].search(item)
	}
	return nil
}

func (n *Node) searchNode(item Item) (*Node, int) {
	i, existed := n.items.search(item)
	if existed {
		return n, i
	}
	if i < len(n.children) {
		return n.children[i].searchNode(item)
	}
	return nil, -1
}

func (n *Node) insert(item Item, nonleaf bool) (median Item, right *Node, ok bool) {
	i, existed := n.items.search(item)
	if existed {
		n.items[i] = item
		ok = false
		return
	}
	if len(n.children) == 0 || nonleaf {
		if len(n.items) < n.maxItems() {
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
			r.parent = n
			return
		}
		if right != nil {
			index, found := right.items.search(m)
			if found {
				right.children.insert(index+1, r)
				r.parent = right
			}
		}
	}
	return
}

func (n *Node) delete(item Item, parentIndex int) (root *Node, ok bool) {
	if n == nil {
		return nil, false
	}
	i, existed := n.items.search(item)
	if existed {
		if len(n.children) == 0 {
			n.items.remove(i)
			if len(n.items) > 0 {
				root = n
			}
			ok = true
			if n.parent != nil && len(n.items) < n.minItems() {
				n.rebalance(parentIndex, false)
			}
			return
		}
		leftMax := n.children[i].max()
		rightMin := n.children[i+1].min()
		if len(leftMax.items) > len(rightMin.items) {
			newSeparator := leftMax.items[len(leftMax.items)-1]
			n.items[i] = newSeparator
			item = newSeparator
		} else {
			newSeparator := rightMin.items[0]
			n.items[i] = newSeparator
			item = newSeparator
			i++
		}
	}
	root = n
	if len(n.children) > i {
		_, ok = n.children[i].delete(item, i)
		if n.parent == nil {
			if len(n.items) == 0 {
				if len(n.children) > 0 {
					root = n.children[0]
				}
			}
		} else {
			if len(n.items) < n.minItems() {
				n.rebalance(parentIndex, true)
			}
		}
	}
	return
}

func (n *Node) rebalance(parentIndex int, nonleaf bool) {
	rightSiblingItems := n.rightSiblingItems(parentIndex)
	if rightSiblingItems > n.minItems() {
		n.rotateLeft(parentIndex, nonleaf)
		return
	}
	leftSiblingItems := n.leftSiblingItems(parentIndex)
	if leftSiblingItems > n.minItems() {
		n.rotateRight(parentIndex, nonleaf)
		return
	}
	if rightSiblingItems > 0 {
		n.mergeLeft(parentIndex, nonleaf)
	} else if leftSiblingItems > 0 {
		n.mergeRight(parentIndex, nonleaf)
	}
}

func (n *Node) rightSiblingItems(parentIndex int) int {
	if parentIndex >= len(n.parent.children)-1 {
		return 0
	}
	return len(n.parent.children[parentIndex+1].items)
}

func (n *Node) leftSiblingItems(parentIndex int) int {
	if parentIndex <= 0 {
		return 0
	}
	return len(n.parent.children[parentIndex-1].items)
}

func (n *Node) rotateLeft(parentIndex int, nonleaf bool) {
	p := n.parent
	n.items.insert(len(n.items), p.items[parentIndex])
	rightSibling := p.children[parentIndex+1]
	p.items[parentIndex] = rightSibling.items[0]
	rightSibling.items.remove(0)
	if nonleaf {
		n.children.insert(len(n.children), rightSibling.children[0])
		n.children[len(n.children)-1].parent = n
		rightSibling.children.remove(0)
	}
}

func (n *Node) rotateRight(parentIndex int, nonleaf bool) {
	p := n.parent
	n.items.insert(0, p.items[parentIndex-1])
	leftSibling := p.children[parentIndex-1]
	p.items[parentIndex-1] = leftSibling.items[len(leftSibling.items)-1]
	leftSibling.items.remove(len(leftSibling.items) - 1)
	if nonleaf {
		n.children.insert(0, leftSibling.children[len(leftSibling.children)-1])
		n.children[0].parent = n
		leftSibling.children.remove(len(leftSibling.children) - 1)
	}
}

func (n *Node) mergeLeft(parentIndex int, nonleaf bool) {
	p := n.parent
	n.items.insert(len(n.items), p.items[parentIndex])
	right := p.children[parentIndex+1]
	n.items.appendRight(right.items)
	p.items.remove(parentIndex)
	p.children.remove(parentIndex + 1)
	if nonleaf {
		n.children.appendRight(right.children)
		for _, v := range right.children {
			v.parent = n
		}
	}
}

func (n *Node) mergeRight(parentIndex int, nonleaf bool) {
	p := n.parent
	leftSibling := p.children[parentIndex-1]
	leftSibling.items.insert(len(leftSibling.items), p.items[parentIndex-1])
	leftSibling.items.appendRight(n.items)
	p.items.remove(parentIndex - 1)
	p.children.remove(parentIndex)
	if nonleaf {
		leftSibling.children.appendRight(n.children)
		for _, v := range n.children {
			v.parent = leftSibling
		}
	}
}

func (n *Node) min() *Node {
	if n == nil {
		return nil
	}
	if len(n.children) > 0 {
		return n.children[0].min()
	}
	return n
}

func (n *Node) max() *Node {
	if n == nil {
		return nil
	}
	if len(n.children) > 0 {
		return n.children[len(n.children)-1].max()
	}
	return n
}

func (n *Node) split(item Item) (median Item, right *Node, ok bool) {
	ok = true
	i := n.minItems()
	median = n.items[i]
	right = newNode(n.maxItems())
	right.items = append(right.items, n.items[i+1:]...)
	n.items = n.items[:i]
	if len(n.children) > 0 {
		right.children = append(right.children, n.children[i+1:]...)
		n.children = n.children[:i+1]
	}
	for _, v := range right.children {
		v.parent = right
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

// Iterator represents an iterator in the B-tree.
type Iterator struct {
	index       int
	parentIndex int
	node        *Node
}

// Item returns the item of this iterator.
func (i *Iterator) Item() Item {
	if i == nil {
		return nil
	}
	return i.node.items[i.index]
}

// Clone returns the clone of this iterator.
func (i *Iterator) Clone() *Iterator {
	if i == nil {
		return nil
	}
	return &Iterator{node: i.node, index: i.index, parentIndex: i.parentIndex}
}

func (i *Iterator) reset(n *Node, index int) *Iterator {
	if i == nil || n == nil {
		return nil
	}
	i.index = index
	i.parentIndex = n.parentIndex()
	i.node = n
	return i
}

// Last returns the last iterator less than this iterator.
func (i *Iterator) Last() (last *Iterator) {
	if i == nil {
		return nil
	}
	n := i.node
	if len(n.children) > 0 {
		max := n.children[i.index].max()
		return i.reset(max, len(max.items)-1)
	}
	if i.index > 0 {
		i.index--
		return i
	}
	left := n
	parentIndex := i.parentIndex
	p := n.parent
	for p != nil && parentIndex == 0 {
		left = p
		parentIndex = p.parentIndex()
		p = left.parent
	}
	if parentIndex > 0 {
		return i.reset(p, parentIndex-1)
	}
	return
}

// Next returns the next iterator more than this iterator.
func (i *Iterator) Next() (next *Iterator) {
	if i == nil {
		return nil
	}
	n := i.node
	if len(n.children) > 0 && i.index < len(n.items) {
		min := n.children[i.index+1].min()
		return i.reset(min, 0)
	}
	if i.index < len(i.node.items)-1 {
		i.index++
		return i
	}
	right := n
	parentIndex := i.parentIndex
	p := n.parent
	for p != nil && parentIndex == len(p.children)-1 {
		right = p
		parentIndex = p.parentIndex()
		p = right.parent
	}
	if parentIndex > -1 && parentIndex < len(p.items) {
		return i.reset(p, parentIndex)
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

func (s *items) appendRight(i items) {
	*s = append(*s, i...)
}

func (s *items) remove(index int) {
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = nil
	*s = (*s)[:len(*s)-1]
}

func (s items) search(item Item) (index int, ok bool) {
	i, j := 0, len(s)
	for i < j {
		h := int(uint(i+j) >> 1)
		if !item.Less(s[h]) {
			i = h + 1
		} else {
			j = h
		}
	}
	if i > 0 && !s[i-1].Less(item) {
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

func (s *children) appendRight(i children) {
	*s = append(*s, i...)
}

func (s *children) remove(index int) {
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = nil
	*s = (*s)[:len(*s)-1]
}
