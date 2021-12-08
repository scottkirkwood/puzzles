package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

var lineRx = regexp.MustCompile(`(\d+),(\d+) -> (\d+),(\d+)`)

type vec struct {
	x, y int
}

type line struct {
	from vec
	to   vec
}

type grid [][]int

type puzzle struct {
	lines []line
	board grid
}

func toIntMap(numbers []int) map[int]bool {
	ret := make(map[int]bool, len(numbers))
	for _, n := range numbers {
		ret[n] = true
	}
	return ret
}

func read(fname string) (*puzzle, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ret := puzzle{
		lines: []line{},
	}
	for scanner.Scan() {
		lineStr := scanner.Text()
		parts := lineRx.FindStringSubmatch(lineStr)
		if len(parts) != 5 {
			return &ret, fmt.Errorf("Unexpected line %q\n", lineStr)
		}
		ret.lines = append(ret.lines, line{
			vec{parseNum(parts[1]), parseNum(parts[2])},
			vec{parseNum(parts[3]), parseNum(parts[4])},
		})
	}
	ret.makeGrid()
	return &ret, nil
}

func (g grid) printGrid() {
	h := len(g)
	w := len(g[0])

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ch := "."
			if g[y][x] > 0 {
				ch = fmt.Sprintf("%d", g[y][x])
			}
			fmt.Printf(ch)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (g grid) countOverlaps() int {
	h := len(g)
	w := len(g[0])

	count := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if g[y][x] > 1 {
				count++
			}
		}
	}
	return count
}

func (p *puzzle) getWidthHeight() (int, int) {
	w, h := 0, 0
	for _, l := range p.lines {
		if l.from.x > w {
			w = l.from.x
		}
		if l.from.y > h {
			h = l.from.y
		}
		if l.to.x > w {
			w = l.to.x
		}
		if l.to.y > h {
			h = l.to.y
		}
	}
	return w + 1, h + 1
}

func (p *puzzle) makeGrid() {
	width, height := p.getWidthHeight()
	p.board = make([][]int, height)
	for y := 0; y < height; y++ {
		p.board[y] = make([]int, width)
	}
}

func sign(x1, x2 int) int {
	if x1 < x2 {
		return 1
	} else if x1 > x2 {
		return -1
	}
	return 0
}

func leq(x, y, sign int) bool {
	if sign == 1 {
		return x <= y
	} else if sign == -1 {
		return y <= x
	}
	return true
}

func (p *puzzle) drawLinesOnGrid() {
	for _, l := range p.lines {
		x1, x2 := l.from.x, l.to.x
		y1, y2 := l.from.y, l.to.y
		dx := sign(x1, x2)
		dy := sign(y1, y2)
		x, y := x1, y1

		if dx == 0 {
			for y := y1; leq(y, y2, dy); y += dy {
				p.inc(x, y)
			}
		} else if dy == 0 {
			for x := x1; leq(x, x2, dx); x += dx {
				p.inc(x, y)
			}
		} else {
			for x := x1; leq(x, x2, dx); x += dx {
				p.inc(x, y)
				y += dy
			}
		}
	}
}

func (p *puzzle) inc(x, y int) {
	p.board[y][x]++
}

func (p *puzzle) Len() int {
	return len(p.lines)
}

func (p *puzzle) calculate() int {
	p.drawLinesOnGrid()
	p.board.printGrid()
	return p.board.countOverlaps()
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(strings.TrimSpace(txt))
	if err != nil {
		fmt.Printf("=========== Bad number %q ==============\n", txt)
	}
	return int(num)
}

func main() {
	flag.Parse()

	fmt.Printf("Day 05a\n")
	infile := "day05.input"
	if *testFileFlag {
		infile = "day05-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
