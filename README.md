# btree
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/btree)](https://pkg.go.dev/github.com/hslam/btree)
[![Build Status](https://github.com/hslam/btree/workflows/build/badge.svg)](https://github.com/hslam/btree/actions)
[![codecov](https://codecov.io/gh/hslam/btree/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/btree)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/btree)](https://goreportcard.com/report/github.com/hslam/btree)
[![LICENSE](https://img.shields.io/github/license/hslam/btree.svg?style=flat-square)](https://github.com/hslam/btree/blob/master/LICENSE)

Package btree implements a B-tree.

#### Definition
According to Knuth's definition, a **[B-tree](https://en.wikipedia.org/wiki/B-tree "B-tree")** of order m is a tree which satisfies the following properties:
* Every node has at most m children.
* Every non-leaf node (except root) has at least ⌈m/2⌉ child nodes.
* The root has at least two children if it is not a leaf node.
* A non-leaf node with k children contains k − 1 keys.
* All leaves appear in the same level and carry no information.

## [Benchmark](http://github.com/hslam/btree-benchmark "btree-benchmark")

<img src="https://raw.githubusercontent.com/hslam/btree-benchmark/master/btree-insert.png" width = "400" height = "300" alt="insert" align=center><img src="https://raw.githubusercontent.com/hslam/btree-benchmark/master/btree-delete.png" width = "400" height = "300" alt="delete" align=center>

<img src="https://raw.githubusercontent.com/hslam/btree-benchmark/master/btree-search.png" width = "400" height = "300" alt="search" align=center><img src="https://raw.githubusercontent.com/hslam/btree-benchmark/master/btree-iterate.png" width = "400" height = "300" alt="iterate" align=center>


## Get started

### Install
```
go get github.com/hslam/btree
```
### Import
```
import "github.com/hslam/btree"
```
### Usage
#### Example
```go
package main

import (
	"fmt"
	"github.com/hslam/btree"
)

func main() {
	tree := btree.New[String](2)
	str := String("Hello World")
	tree.Insert(str)
	fmt.Println(tree.Search(str))
	tree.Delete(str)
}

type String string

func (a String) Less(b String) bool {
	return a < b
}
```

#### Output
```
Hello World
```

#### Iterator Example
```go
package main

import (
	"fmt"
	"github.com/hslam/btree"
)

func main() {
	tree := btree.New[Int](2)
	for i := 0; i < 10; i++ {
		tree.Insert(Int(i))
	}
	iter := tree.Min().MinIterator()
	for iter != nil {
		fmt.Printf("%d\t", iter.Item())
		iter = iter.Next()
	}
}

type Int int

func (a Int) Less(b Int) bool {
	return a < b
}
```
#### B-Tree
<img src="https://raw.githubusercontent.com/hslam/btree/master/btree.png" alt="btree" align=center>

#### Output
```
0	1	2	3	4	5	6	7	8	9
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)

### Author
btree was written by Meng Huang.


