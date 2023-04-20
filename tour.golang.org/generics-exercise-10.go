// Generics
package main

import "fmt"

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func main() {
	l1 := List[int]{nil, 1}
	l2 := List[int]{&l1, 2}
	fmt.Printf("List %v, %v\n", l2, l1)
}

