package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type floor struct {
	dayNum int
	tiles  map[hex]bool // false = white, true = black
	next   map[hex]bool
}

type hex struct {
	q, r int
}

const (
	E = iota
	NE
	NW
	W
	SW
	SE
	NUM_DIRECTIONS
)

var axialDirections = []hex{
	hex{+1, 0}, hex{+1, -1}, hex{0, -1},
	hex{-1, 0}, hex{-1, +1}, hex{0, +1},
}

func newFloor() floor {
	return floor{
		dayNum: 1,
		tiles:  map[hex]bool{},
		next:   map[hex]bool{},
	}
}

func (f *floor) generation() {
	f.next = make(map[hex]bool, len(f.tiles))
	q1, r1, q2, r2 := f.extents()
	for r := r1 - 1; r <= r2+1; r++ {
		for q := q1 - 1; q <= q2+1; q++ {
			hx := hex{q, r}
			black := f.tiles[hx]
			n := f.countAdjacentBlack(hx)
			if black {
				if !(n == 0 || n > 2) {
					f.next[hx] = true
				}
			} else if n == 2 {
				f.next[hx] = true
			}
		}
	}
	f.tiles = f.next
}

func (f floor) extents() (q1, r1, q2, r2 int) {
	q1, r1 = 9999999, 9999999
	q2, r2 = -9999999, -9999999
	for hx, black := range f.tiles {
		if !black {
			continue
		}
		if hx.q < q1 {
			q1 = hx.q
		}
		if hx.r < r1 {
			r1 = hx.r
		}
		if hx.q > q2 {
			q2 = hx.q
		}
		if hx.r > r2 {
			r2 = hx.r
		}
	}
	return q1, r1, q2, r2
}

func (f floor) countAdjacentBlack(hx hex) int {
	count := 0
	for _, dir := range axialDirections {
		hx2 := hx
		hx2.move(dir)
		if f.tiles[hx2] {
			count++
		}
	}
	return count
}

func (f *floor) flip(hx hex) {
	f.tiles[hx] = !f.tiles[hx]
}

func (f floor) countBlack() int {
	count := 0
	for _, black := range f.tiles {
		if black {
			count++
		}
	}
	return count
}

func (h *hex) move(hx hex) {
	h.q += hx.q
	h.r += hx.r
}

func letterToDirection(dir string) int {
	switch dir {
	case "e":
		return E
	case "ne":
		return NE
	case "nw":
		return NW
	case "w":
		return W
	case "sw":
		return SW
	case "se":
		return SE
	}
	fmt.Printf("Invalid direction %q\n", dir)
	return NUM_DIRECTIONS
}

var directionRx = regexp.MustCompile(`(e|se|sw|w|nw|ne)`)

func read(fname string) ([][]hex, error) {
	ret := [][]hex{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ret = append(ret, []hex{})
		for _, ch := range directionRx.FindAllString(line, -1) {
			ret[len(ret)-1] = append(ret[len(ret)-1], axialDirections[letterToDirection(ch)])
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("Unable to parse decimal %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func main() {
	flag.Parse()

	fmt.Printf("Day 24\n")

	infile := "day24.input"
	if *testFileFlag {
		infile = "day24-sample.input"
	}
	directions, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("%d\n", len(directions))
	f := newFloor()
	for _, moves := range directions {
		start := hex{0, 0}
		for _, move := range moves {
			start.move(move)
		}
		f.flip(start)
	}
	fmt.Printf("Black tiles %d\n", f.countBlack())
	for day := 1; day <= 100; day++ {
		f.generation()
		fmt.Printf("Day %d: %d\n", day, f.countBlack())
	}
}
