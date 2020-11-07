# btree
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/btree)](https://pkg.go.dev/github.com/hslam/btree)
[![Build Status](https://travis-ci.org/hslam/btree.svg?branch=master)](https://travis-ci.org/hslam/btree)
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
	tree := btree.New(2)
	str := String("Hello World")
	tree.Insert(str)
	fmt.Println(tree.Search(str))
	tree.Delete(str)
}

type String string

func (a String) Less(b btree.Item) bool {
	return a < b.(String)
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
	tree := btree.New(2)
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

func (a Int) Less(b btree.Item) bool {
	return a < b.(Int)
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


