package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type world struct {
	cur     cube
	next    cube
	cycles  int
	extents minMax
}

type minMax struct {
	min vec3
	max vec3
}

type cube struct {
	w map[vec3]bool
}

type vec3 struct {
	x, y, z int
}

// String is the stringer interface for world
func (w world) String() string {
	return w.cur.String()
}

func (w *world) cycle() {
	w.next = newCube()
	coord := vec3{}
	for z := w.extents.min.z - 1; z <= w.extents.max.z+1; z++ {
		coord.z = z
		for y := w.extents.min.y - 1; y <= w.extents.max.y+1; y++ {
			coord.y = y
			for x := w.extents.min.x - 1; x <= w.extents.max.x+1; x++ {
				coord.x = x
				if w.cur.alive(coord) {
					w.next.set(coord, true)
				}
			}
		}
	}
	w.cur = w.next
	w.extents = w.cur.calcExtents()
	w.cycles++
}

func newCube() cube {
	return cube{
		w: make(map[vec3]bool),
	}
}

func (c *cube) alive(coord vec3) bool {
	count := 0
	for z := -1; z <= 1; z++ {
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				if c.get(vec3{coord.x + x, coord.y + y, coord.z + z}) {
					count++
					if count > 3 {
						break
					}
				}
			}
		}
	}
	active := c.get(coord)
	if active {
		return count == 2 || count == 3
	}
	return count == 3
}

func (c *cube) countActive() int {
	count := 0
	extents := c.calcExtents()
	coord := vec3{}
	for z := extents.min.z; z <= extents.max.z; z++ {
		coord.z = z
		for y := extents.min.y; y <= extents.max.y; y++ {
			coord.y = y
			for x := extents.min.x; x <= extents.max.x; x++ {
				coord.x = x
				if c.get(coord) {
					count++
				}
			}
		}
	}
	return count
}

func (c *cube) set(coord vec3, val bool) {
	c.w[coord] = val
}

func (c *cube) get(coord vec3) bool {
	return c.w[coord]
}

func minMaxMaxout() (m minMax) {
	m.min.x = 999
	m.min.y = 999
	m.min.z = 999
	m.max.x = -999
	m.max.y = -999
	m.max.z = -999
	return m
}

func (c cube) calcExtents() minMax {
	m := minMaxMaxout()
	for k, v := range c.w {
		if v {
			if k.x < m.min.x {
				m.min.x = k.x
			}
			if k.y < m.min.y {
				m.min.y = k.y
			}
			if k.z < m.min.z {
				m.min.z = k.z
			}
			if k.x > m.max.x {
				m.max.x = k.x
			}
			if k.y > m.max.y {
				m.max.y = k.y
			}
			if k.z > m.max.z {
				m.max.z = k.z
			}
		}
	}
	return m
}

// String is the stringer interface
func (c cube) String() string {
	extents := c.calcExtents()
	lines := []string{}
	for z := extents.min.z; z <= extents.max.z; z++ {
		lines = append(lines, fmt.Sprintf("\nz=%d", z))
		for y := extents.min.y; y <= extents.max.y; y++ {
			line := ""
			for x := extents.min.x; x <= extents.max.x; x++ {
				if c.get(vec3{x, y, z}) {
					line += "#"
				} else {
					line += "."
				}
			}
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func read(fname string) (cube, error) {
	ret := newCube()
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	coord := vec3{}
	for scanner.Scan() {
		line := scanner.Text()
		coord.x = 0
		for _, ch := range line {
			if ch == '#' {
				ret.set(coord, true)
			}
			coord.x++
		}

		coord.y++
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func main() {
	flag.Parse()

	fmt.Printf("Day 9\n")

	infile := "day9.input"
	if *testFileFlag {
		infile = "day9-sample.input"
	}
	start, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("%s\n", start)
	w := world{cur: start, extents: start.calcExtents()}
	for i := 0; i < 6; i++ {
		w.cycle()
	}
	fmt.Printf("%s\n", w)

	fmt.Printf("Count active %d\n", w.cur.countActive())
}
