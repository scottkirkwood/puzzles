// Exercise: Loops and Functions
package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) (z float64) {
	z = 1
	for {
		diff := (z*z - x) / (2 * z)
		z -= diff
		if math.Abs(diff) < 1e-10 {
			break
		}
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2), math.Sqrt(2))
}
