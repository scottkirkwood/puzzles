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
	cave   grid
}

type grid struct {
	w, h int
	vals [][]int
}

type point struct {
	x, y int
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

func read(fname string) (*puzzle, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	p := puzzle{
	}
	for scanner.Scan() {
		lineStr := scanner.Text()
		p.cave.vals = append(p.cave.vals, parseInts(strings.Split(lineStr, "")))
	}
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.cave.vals)
}

func (p *puzzle) Print() {
}

func (p *puzzle) generate() {
}

func (p *puzzle) calculate() int {
	return 0
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
