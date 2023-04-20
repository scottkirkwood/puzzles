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
	cave grid
}

type grid struct {
	w, h int
	vals [][]int
}

type point struct {
	x, y int
}

type traverse struct {
	cave       *grid
	visited    map[point]bool
	points     []point
	bestPoints []point
	minScore   int
	curScore   int
}

func (g *grid) get(pt point) int {
	if !g.inBounds(pt) {
		return -1
	}
	return g.vals[pt.y][pt.x]
}

func (g *grid) inBounds(pt point) bool {
	return pt.x >= 0 && pt.y >= 0 && pt.x < g.w && pt.y < g.h
}

func (g *grid) calcScore(points []point) int {
	score := 0
	for _, pt := range points {
		score += g.get(pt)
	}
	return score
}

func (g *grid) Print() {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			fmt.Printf("%d", g.get(point{x, y}))
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func newTraverse(cave *grid) *traverse {
	return &traverse{
		cave:     cave,
		visited:  make(map[point]bool),
		minScore: 1 << 62,
	}
}

func (t *traverse) atExit(pt point) bool {
	return pt.x == t.cave.w-1 && pt.y == t.cave.h-1
}

func (t *traverse) generate(pt point) {
	if t.atExit(pt) {
		score := t.curScore // t.cave.calcScore(t.points)
		if score < t.minScore {
			t.minScore = score
			// t.bestPoints = make([]point, len(t.points))
			//copy(t.bestPoints, t.points)
			fmt.Printf("%d\n", t.minScore)
		}
		return
	}
	for _, dxy := range []point{point{1, 0}, point{0, 1}} {
		pt := point{pt.x + dxy.x, pt.y + dxy.y}
		if !t.cave.inBounds(pt) { // || t.visited[pt] {
			continue
		}
		// t.points = append(t.points, pt)
		//t.visited[pt] = true
		val := t.cave.get(pt)
		t.curScore += val
		if t.curScore <= t.minScore {
			t.generate(pt)
		}
		t.curScore -= val
		// t.visited[pt] = false
		// t.points = t.points[0 : len(t.points)-1]
	}
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
		p.cave.vals = append(p.cave.vals, parseInts(strings.Split(lineStr, "")))
	}
	p.cave.h = len(p.cave.vals)
	p.cave.w = len(p.cave.vals[0])
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.cave.vals)
}

func (p *puzzle) Print() {
}

func (p *puzzle) calculate() int {
	tr := newTraverse(&p.cave)
	p.cave.Print()
	tr.generate(point{0, 0})
	// 1077 too high
	// 1076 too high
	// 1069 too high
	// 1058 wrong
	return tr.minScore
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

	fmt.Printf("Day 15a\n")
	infile := "day15.input"
	if *testFileFlag {
		infile = "day15-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
