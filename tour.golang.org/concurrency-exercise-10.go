// Exercise: Equivalent Binary Trees
package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	walkHelper(t, ch)
	close(ch)
}

func walkHelper(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walkHelper(t.Left, ch)
	ch <- t.Value
	walkHelper(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for c1 := range ch1 {
		c2 := <-ch2
		if c2 != c1 {
			return false
		}
	}
	return true
}

func Same2(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		c1, ok := <-ch1
		if !ok {
			if _, ok := <-ch2; !ok {
				return true
			}
			return false
		}
		c2, ok := <-ch2
		if !ok {
			return false
		}
		if c2 != c1 {
			return false
		}
	}
	return true
}

func main() {
	if Same2(tree.New(1), tree.New(1)) {
		fmt.Printf("Trees are the same as expected\n")
	} else {
		fmt.Printf("Oops, trees should not be different\n")
	}
	if Same2(tree.New(1), tree.New(2)) {
		fmt.Printf("Oops, trees should be different\n")
	} else {
		fmt.Printf("Trees are different as expected\n")
	}
}
