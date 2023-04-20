// Exercise: Slices

package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	// Setup
	p := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		p[y] = make([]uint8, dx)
	}

	// Draw
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			p[y][x] = uint8((x + y) / 2)
		}
	}
	return p
}

func main() {
	pic.Show(Pic)
}
