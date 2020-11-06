# btree
Package btree implements a B-tree.

**[Properties](https://en.wikipedia.org/wiki/B-tree "properties")**
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

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)

### Author
btree was written by Meng Huang.


