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
	numStepsFlag = flag.Int("n", 10, "Number of steps")
)

type puzzle struct {
	w, h    int
	step    int
	flashes int
	grid    [][]int
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
	ret := puzzle{}
	for scanner.Scan() {
		lineStr := scanner.Text()
		ret.grid = append(ret.grid, parseInts(strings.Split(lineStr, "")))
	}
	ret.h = len(ret.grid)
	ret.w = len(ret.grid[0])
	return &ret, nil
}

func (p *puzzle) getXY(pt point) int {
	if !p.inBounds(pt) {
		return -1
	}
	return p.grid[pt.y][pt.x]
}

func (p *puzzle) inBounds(pt point) bool {
	return pt.x >= 0 && pt.y >= 0 && pt.x < p.w && pt.y < p.h
}

func (p *puzzle) setXY(pt point, val int) {
	if !p.inBounds(pt) {
		return
	}
	p.grid[pt.y][pt.x] = val
}

func (p *puzzle) incXY(pt point) int {
	if !p.inBounds(pt) {
		return -1
	}
	p.grid[pt.y][pt.x]++
	return p.grid[pt.y][pt.x]
}

func (p *puzzle) Len() int {
	return len(p.grid)
}

func (p *puzzle) Print() {
	fmt.Printf("After step %d:\n", p.step)
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			pt := point{x, y}
			val := p.getXY(pt)
			if val == 0 {
				fmt.Printf("\u001b[1m0\u001b[0m")
			} else {
				fmt.Printf("%d", val)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (p *puzzle) incNeighbours(pt point) (toZero []point) {
	toZero = append(toZero, pt)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			pt := point{pt.x + dx, pt.y + dy}
			if p.incXY(pt) > 9 {
				p.setXY(pt, 0)
				p.flashes++
				toZero = append(toZero, p.incNeighbours(pt)...)
			}
		}
	}
	return toZero
}

func (p *puzzle) propFlashes() {
	for {
		toZero := []point{}
		for y := 0; y < p.h; y++ {
			for x := 0; x < p.w; x++ {
				pt := point{x, y}
				if p.getXY(pt) > 9 {
					p.setXY(pt, 0)
					p.flashes++
					toZero = append(toZero, p.incNeighbours(pt)...)
				}
			}
		}
		if len(toZero) == 0 {
			return
		}
		for _, pt := range toZero {
			p.setXY(pt, 0)
		}
	}
}

func (p *puzzle) generation() {
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			p.incXY(point{x, y})
		}
	}
	p.propFlashes()
	p.step++
}

func (p *puzzle) allFlash() bool {
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			if p.getXY(point{x, y}) != 0 {
				return false
			}
		}
	}
	return true
}

func (p *puzzle) calculate() int {
	p.Print()
	stepAllFlash := 0
	for step := 0; step < *numStepsFlag; step++ {
		p.generation()
		if p.allFlash() {
			p.Print()
			stepAllFlash = p.step
			break
		}
	}
	fmt.Printf("Step %d\n", stepAllFlash)
	return stepAllFlash
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

	fmt.Printf("Day 11b\n")
	infile := "day11.input"
	if *testFileFlag {
		infile = "day11-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
