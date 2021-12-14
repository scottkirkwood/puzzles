package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	testFileFlag = flag.Bool("t", false, "Use the test file")
)

type puzzle struct {
	points []point
	folds  []point

	paper *grid
}

type grid struct {
	w, h int
	dots [][]bool
}

func makeGrid(w, h int) *grid {
	g := grid{
		w:    w,
		h:    h,
		dots: make([][]bool, h),
	}
	for y := 0; y < g.h; y++ {
		g.dots[y] = make([]bool, w)
	}
	return &g
}

func (g *grid) onGrid(pt point) bool {
	return pt.x >= 0 && pt.y >= 0 &&
		pt.x < g.w && pt.y < g.h
}

func (g *grid) get(pt point) bool {
	if !g.onGrid(pt) {
		return false
	}
	return g.dots[pt.y][pt.x]
}

func (g *grid) set(pt point, t bool) {
	if !g.onGrid(pt) {
		return
	}
	g.dots[pt.y][pt.x] = t
}

func (g *grid) union(pt point, t bool) {
	if !t || !g.onGrid(pt) {
		return
	}
	g.set(pt, true)
}

func (g *grid) setPoints(points []point) {
	for _, pt := range points {
		g.set(pt, true)
	}
}

func (g *grid) foldY(foldY int) *grid {
	newGrid := makeGrid(g.w, foldY)
	// Copy the top part
	for y := 0; y < foldY; y++ {
		for x := 0; x < g.w; x++ {
			pt := point{x, y}
			newGrid.set(pt, g.get(pt))
		}
	}
	for y := foldY; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			pt := point{x, 2*foldY - y}
			newGrid.union(pt, g.get(point{x, y}))
		}
	}
	return newGrid
}

func (g *grid) foldX(foldX int) *grid {
	newGrid := makeGrid(foldX, g.h)
	// Copy the left part
	for y := 0; y < g.h; y++ {
		for x := 0; x < foldX; x++ {
			pt := point{x, y}
			newGrid.set(pt, g.get(pt))
		}
	}
	for y := 0; y < g.h; y++ {
		for x := foldX; x < g.w; x++ {
			pt := point{2*foldX - x, y}
			newGrid.union(pt, g.get(point{x, y}))
		}
	}
	return newGrid
}

func (g *grid) countDots() int {
	count := 0
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.get(point{x, y}) {
				count++
			}
		}
	}
	return count
}

func (g *grid) Print() {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.get(point{x, y}) {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

type point struct {
	x, y int
}

func read(fname string) (*puzzle, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	p := puzzle{}
	for scanner.Scan() {
		lineStr := scanner.Text()
		if strings.Contains(lineStr, ",") {
			nums := parseInts(strings.Split(lineStr, ","))
			p.points = append(p.points, point{nums[0], nums[1]})
		} else if strings.Contains(lineStr, "=") {
			parts := strings.Split(lineStr, "=")
			val := parseInt(parts[1])
			if parts[0] == "fold along x" {
				p.folds = append(p.folds, point{val, 0})
			} else {
				p.folds = append(p.folds, point{0, val})
			}
		}
	}
	w, h := p.maxXY()
	p.paper = makeGrid(w, h)
	p.paper.setPoints(p.points)
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.points)
}

func (p *puzzle) maxXY() (int, int) {
	w, h := 0, 0
	for _, pt := range p.points {
		if pt.x > w {
			w = pt.x
		}
		if pt.y > h {
			h = pt.y
		}
	}
	return w + 1, h + 1
}

func (p *puzzle) Print() {
	for _, pt := range p.points {
		fmt.Printf("%d, %d\n", pt.x, pt.y)
	}
	fmt.Printf("Folds\n")
	for _, fold := range p.folds {
		fmt.Printf("%d, %d\n", fold.x, fold.y)
	}
}

func (p *puzzle) calculate() int {
	// p.Print()
	p.paper.Print()
	for _, fold := range p.folds {
		if fold.x == 0 {
			p.paper = p.paper.foldY(fold.y)
		} else {
			p.paper = p.paper.foldX(fold.x)
		}
		p.paper.Print()
		break
	}
	return p.paper.countDots()
}

func parseInts(nums []string) []int {
	ret := make([]int, 0, len(nums))
	for _, num := range nums {
		ret = append(ret, parseInt(num))
	}
	return ret
}

func parseInt(txt string) int {
	num, err := strconv.Atoi(strings.TrimSpace(txt))
	if err != nil {
		fmt.Printf("=========== Bad number %q ==============\n", txt)
	}
	return int(num)
}

func main() {
	flag.Parse()

	fmt.Printf("Day 13a\n")
	infile := "day13.input"
	if *testFileFlag {
		infile = "day13-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
