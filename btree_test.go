package btree

import (
	"testing"
)

func TestBtree(t *testing.T) {
	for d := 2; d < 16; d++ {
		for i := 0; i < 128; i++ {
			testBtree(128, i, true, d, t)
			testBtree(128, i, false, d, t)
		}
	}
}

func testBtree(n, j int, r bool, degree int, t *testing.T) {
	tree := New(degree)
	if r {
		for i := n - 1; i >= 0; i-- {
			tree.Insert(Int(i))
		}
	} else {
		for i := 0; i < n; i++ {
			tree.Insert(Int(i))
		}
	}
	if tree.Length() != n {
		t.Error("")
	}
	testSearch(tree, j, t)
	tree.Delete(Int(j))
	testNilNode(tree, j, t)
	if tree.Length() != n-1 {
		t.Error("")
	}
	if r {
		for i := n - 1; i >= 0; i-- {
			tree.Delete(Int(i))
		}
	} else {
		for i := 0; i < n; i++ {
			tree.Delete(Int(i))
		}
	}
	if tree.Length() != 0 {
		t.Error(tree.Length())
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
	if tree.Root().Parent() != nil {
		t.Error("")
	}
	if tree.Root().Children() != nil {
		t.Error("")
	}
	if tree.Root().Items() != nil {
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
	tree.Clear()
	if tree.Length() != 0 {
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
