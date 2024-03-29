// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package btree

import (
	"testing"
)

func TestBtree(t *testing.T) {
	for d := 2; d < 9; d++ {
		for i := 0; i < 64; i++ {
			testBtree(64, i, true, d, t)
			testBtree(64, i, false, d, t)
			testBtreeM(64, i+1, true, d, t)
			testBtreeM(64, i+1, false, d, t)
		}
	}
}

func testBtree(n, j int, r bool, degree int, t *testing.T) {
	tree := New(degree)
	if r {
		for i := n - 1; i >= 0; i-- {
			tree.Insert(Int(i))
			testTraversal(tree, t)
		}
	} else {
		for i := 0; i < n; i++ {
			tree.Insert(Int(i))
			testTraversal(tree, t)
		}
	}
	if tree.Length() != n {
		t.Error("")
	}
	tree.Delete(Int(n))
	if tree.Length() != n {
		t.Error("")
	}
	testSearch(tree, j, t)
	tree.Delete(Int(j))
	testTraversal(tree, t)
	testNilNode(tree, j, t)
	if tree.Length() != n-1 {
		t.Error("")
	}
	if r {
		for i := n - 1; i >= 0; i-- {
			tree.Delete(Int(i))
			testTraversal(tree, t)
			testNilNode(tree, i, t)
		}
	} else {
		for i := 0; i < n; i++ {
			tree.Delete(Int(i))
			testTraversal(tree, t)
			testNilNode(tree, i, t)
		}
	}
	if tree.Length() != 0 {
		t.Error(tree.Length())
	}
}

func testBtreeM(n, j int, r bool, degree int, t *testing.T) {
	tree := New(degree)
	if r {
		for i := n; i > 0; i-- {
			tree.Insert(Int(i))
			testTraversal(tree, t)
			tree.Insert(Int(-i))
			testTraversal(tree, t)
		}
	} else {
		for i := 1; i < n+1; i++ {
			tree.Insert(Int(i))
			testTraversal(tree, t)
			tree.Insert(Int(-i))
			testTraversal(tree, t)
		}
	}
	if tree.Length() != n*2 {
		t.Error("")
	}
	testSearch(tree, j, t)
	tree.Delete(Int(j))
	testTraversal(tree, t)
	testNilNode(tree, j, t)
	if tree.Length() != n*2-1 {
		t.Error("")
	}
	j = -j
	testSearch(tree, j, t)
	tree.Delete(Int(j))
	testTraversal(tree, t)
	testNilNode(tree, j, t)
	if tree.Length() != n*2-2 {
		t.Error("")
	}
	if r {
		for i := n; i > 0; i-- {
			tree.Delete(Int(i))
			testTraversal(tree, t)
			testNilNode(tree, i, t)
			tree.Delete(Int(-i))
			testTraversal(tree, t)
			testNilNode(tree, -i, t)
		}
	} else {
		for i := 1; i < n+1; i++ {
			tree.Delete(Int(i))
			testTraversal(tree, t)
			testNilNode(tree, i, t)
			tree.Delete(Int(-i))
			testTraversal(tree, t)
			testNilNode(tree, -i, t)
		}
	}
	if tree.Length() != 0 {
		t.Error(tree.Length())
	}
}

func testTraversal(tree *Tree, t *testing.T) {
	count := 0
	testLength(tree.Root(), &count)
	if tree.Length() != count {
		t.Error(tree.Length(), count)
	}
	traverse(tree.Root(), t)
	testIteratorAscend(tree, t)
	testIteratorDescend(tree, t)
}

func testLength(node *Node, count *int) {
	*count += len(node.Items())
	if node != nil {
		for _, child := range node.children {
			testLength(child, count)
		}
	}
}

func traverse(node *Node, t *testing.T) {
	if node != nil {
		for _, child := range node.children {
			if child.parent != node {
				t.Error("")
			}
		}
	}
}

func testIteratorAscend(tree *Tree, t *testing.T) {
	iter := tree.Min().MinIterator()
	item := iter.Item()
	next := iter.Next()
	for iter != nil && next != nil {
		if !item.Less(next.Item()) {
			t.Error(item, next.Item())
		}
		iter = next
		item = iter.Item()
		next = iter.Next()
	}
}

func testIteratorDescend(tree *Tree, t *testing.T) {
	iter := tree.Max().MaxIterator()
	item := iter.Item()
	last := iter.Last()
	for iter != nil && last != nil {
		if !last.Item().Less(item) {
			t.Error(last.Item(), item)
		}
		iter = last
		last = iter.Last()
	}
}

func testSearch(tree *Tree, j int, t *testing.T) {
	if node := tree.SearchNode(Int(j)); node == nil {
		t.Error("")
	} else {
		node.Items()
		node.Children()
		node.Parent()
	}
	if item := tree.Search(Int(j)); item == nil {
		t.Error("")
	} else if int(item.(Int)) != j {
		t.Error("")
	}
}

func testNilNode(tree *Tree, j int, t *testing.T) {
	if item := tree.Search(Int(j)); item != nil {
		t.Error("")
	}
}

func TestInsert(t *testing.T) {
	tree := New(2)
	tree.Insert(Int(0))
	tree.Insert(Int(0))
	defer func() {
		if err := recover(); err == nil {
			t.Error("")
		}
	}()
	tree.Insert(nil)
}

func TestDegree(t *testing.T) {
	degree := 2
	tree := New(degree)
	if tree.MaxItems() != degree*2-1 {
		t.Error("")
	}
	if tree.MinItems() != degree-1 {
		t.Error("")
	}
	defer func() {
		if err := recover(); err == nil {
			t.Error("")
		}
	}()
	New(0)
}

func TestEmptyTree(t *testing.T) {
	tree := New(2)
	tree.Delete(Int(0))
	if tree.Root() != nil {
		t.Error("")
	}
	if tree.Min() != nil {
		t.Error("")
	}
	if tree.Max() != nil {
		t.Error("")
	}
	if tree.Search(Int(0)) != nil {
		t.Error("")
	}
	if tree.SearchNode(Int(0)) != nil {
		t.Error("")
	}
	if tree.SearchIterator(Int(0)) != nil {
		t.Error("")
	}
	if tree.Root().Parent() != nil {
		t.Error("")
	}
	if tree.Root().Children() != nil {
		t.Error("")
	}
	if tree.Root().Items() != nil {
		t.Error("")
	}
	if tree.Root().Iterator(0) != nil {
		t.Error("")
	}
	if tree.Root().parentIndex() != -1 {
		t.Error("")
	}
	if tree.Root().maxItems() != 0 {
		t.Error("")
	}
	if tree.Root().minItems() != 0 {
		t.Error("")
	}
	if tree.Length() != 0 {
		t.Error("")
	}
	tree.Insert(Int(1))
	tree.Insert(Int(2))
	tree.Insert(Int(3))
	if tree.SearchNode(Int(0)) != nil {
		t.Error("")
	}
	if tree.SearchIterator(Int(0)) != nil {
		t.Error("")
	}
	tree.Clear()
	if tree.Length() != 0 {
		t.Error("")
	}
}

func TestIterator(t *testing.T) {
	tree := New(2)
	iter := tree.Max().MaxIterator()
	if iter.Clone() != nil {
		t.Error("")
	} else if iter.reset(nil, 0) != nil {
		t.Error("")
	}
	tree.Insert(Int(1))
	tree.Insert(Int(2))
	tree.Insert(Int(3))
	iter = tree.Max().MaxIterator()
	clone := iter.Clone()
	if iter.Item().Less(clone.Item()) {
		t.Error("")
	} else if clone.Item().Less(iter.Item()) {
		t.Error("")
	}
}

func TestStringLess(t *testing.T) {
	a := String("a")
	b := String("b")
	if !a.Less(b) {
		t.Error("")
	}
}

func TestReplaceItem(t *testing.T) {
	tree := New(8)
	n := 1024
	for i := 0; i < n; i++ {
		tree.Insert(Int(i))
		testTraversal(tree, t)
		if tree.Length() != i+1 {
			t.Error("")
		}
	}
	for i := 0; i < n; i++ {
		tree.Insert(Int(i))
		testTraversal(tree, t)
		if tree.Length() != n {
			t.Error("")
		}
	}
	for i := 0; i < n; i++ {
		tree.Delete(Int(i))
		testTraversal(tree, t)
		if tree.Length() != n-i-1 {
			t.Error("")
		}
	}
}
