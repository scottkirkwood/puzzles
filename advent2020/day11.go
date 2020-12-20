package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	testFileFlag = flag.Bool("t", false, "Use the test file")
)

type world struct {
	cur  floor
	next floor
}

func newWorld(f floor) world {
	return world{cur: f}
}

func (w world) String() string {
	return w.cur.String()
}

func (w *world) numOccupied() int {
	return w.cur.numOccupied()
}

func (w *world) cycle() bool {
	w.next = newFloor(w.cur.width, w.cur.height)
	for y := 0; y < w.cur.height; y++ {
		for x := 0; x < w.cur.width; x++ {
			w.next.set(x, y, w.cur.want(x, y))
		}
	}
	same := w.cur.equals(w.next)
	w.cur = w.next
	return same
}

type floor struct {
	width  int
	height int
	rows   [][]byte
}

func newFloor(w, h int) floor {
	f := floor{
		width:  w,
		height: h,
		rows:   make([][]byte, h),
	}
	for y := 0; y < h; y++ {
		f.rows[y] = make([]byte, w)
	}
	return f
}

func (f *floor) want(x, y int) byte {
	ch := f.get(x, y)
	if ch == '.' {
		return ch
	}
	numOccupied := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if f.get(x+dx, y+dy) == '#' {
				numOccupied++
			}
		}
	}
	if ch == '#' {
		if numOccupied >= 4 {
			return 'L'
		}
		return '#'
	}
	if ch == 'L' && numOccupied == 0 {
		return '#'
	}
	return ch
}

func (f *floor) get(x, y int) byte {
	if x >= f.width || y >= len(f.rows) || y < 0 || x < 0 {
		return '.'
	}
	return f.rows[y][x]
}

func (f *floor) set(x, y int, b byte) {
	if x >= f.width || y >= len(f.rows) || y < 0 || x < 0 {
		return
	}
	f.rows[y][x] = b
}

func (f *floor) equals(other floor) bool {
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			if f.get(x, y) != other.get(x, y) {
				return false
			}
		}
	}
	return true
}

func (f floor) numOccupied() int {
	count := 0
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			if f.get(x, y) == '#' {
				count++
			}
		}
	}
	return count
}

func (f floor) String() string {
	lines := make([]string, f.height)
	for i, row := range f.rows {
		lines[i] = string(row)
	}
	return strings.Join(lines, "\n")
}

func read(fname string) (floor, error) {
	ret := floor{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		ret.rows = append(ret.rows, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	ret.height = len(ret.rows)
	ret.width = len(ret.rows[0])
	return ret, nil
}

func main() {
	flag.Parse()

	fmt.Printf("Day 11\n")

	infile := "day11.input"
	if *testFileFlag {
		infile = "day11-sample.input"
	}
	floor, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	w := newWorld(floor)
	round := 0
	for {
		fmt.Printf("\nRound %d\n", round)
		round++
		fmt.Printf("%s\n", w)
		if w.cycle() {
			break
		}
	}
	// 920  too low
	fmt.Printf("Num occupied %d\n", w.numOccupied())
}
