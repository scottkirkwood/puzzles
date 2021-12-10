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
	w, h int
	grid [][]int
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
		nums := parseInts(strings.Split(lineStr, ""))
		ret.grid = append(ret.grid, nums)
	}
	ret.w = len(ret.grid[0])
	ret.h = len(ret.grid)
	return &ret, nil
}

func (p *puzzle) getXY(x, y int) int {
	if x < 0 || x >= p.w {
		return -1
	}
	if y < 0 || y >= p.h {
		return -1
	}
	return p.grid[y][x]
}

func (p *puzzle) Len() int {
	return len(p.grid)
}

func (p *puzzle) Print() {
	for y := 0; y < p.Len(); y++ {
		for x := 0; x < len(p.grid[y]); x++ {
			fmt.Printf("%d", p.grid[y][x])
		}
		fmt.Printf("\n")
	}
}

func (p *puzzle) calculate() int {
	risk := 0
	adjacent := [4]int{}

	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			cur := p.getXY(x, y)
			adjacent[0] = p.getXY(x, y-1) // top
			adjacent[1] = p.getXY(x+1, y) // right
			adjacent[2] = p.getXY(x, y+1) // bottom
			adjacent[3] = p.getXY(x-1, y) // left
			lowest := true
			for _, adj := range adjacent {
				if adj == -1 {
					continue
				}
				if adj <= cur {
					lowest = false
					break
				}
			}
			if lowest {
				fmt.Printf("%d, %d = %d risk %d\n", x, y, cur, 1+cur)
				risk += 1 + cur // 1674 too high
			}
		}
	}
	return risk
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

	fmt.Printf("Day 09a\n")
	infile := "day09.input"
	if *testFileFlag {
		infile = "day09-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
