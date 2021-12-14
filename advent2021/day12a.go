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
	// edges is the from-to directional edges
	edges map [string]string
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
		parts := strings.Split(lineStr, "-")
		ret.edges[parts[0]]=parts[1]
		if parts[0] != "start" || parts[1] != "end"
		  ret.edges[parts[1]] = parts[0]
	  }
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.edges)
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

func (p *puzzle) calculate() int {
	p.Print()
	for step := 0; step < *numStepsFlag; step++ {
		p.generation()
		p.Print()
	}
	return p.flashes
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

	fmt.Printf("Day 12a\n")
	infile := "day12.input"
	if *testFileFlag {
		infile = "day12-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
