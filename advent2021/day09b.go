package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	testFileFlag = flag.Bool("t", false, "Use the test file")
)

type point struct {
	x, y int
}

type puzzle struct {
	w, h          int
	grid          [][]int
	pointsVisited map[int]bool
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
	ret.pointsVisited = make(map[int]bool)
	return &ret, nil
}

func (p *puzzle) get(pt point) int {
	if pt.x < 0 || pt.x >= p.w {
		return -1
	}
	if pt.y < 0 || pt.y >= p.h {
		return -1
	}
	return p.grid[pt.y][pt.x]
}

func (p *puzzle) getXY(x, y int) int {
	return p.get(point{x, y})
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

func (p *puzzle) lowestPoints() []point {
	adjacent := [4]int{}

	ret := []point{}
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
				ret = append(ret, point{x, y})
			}
		}
	}
	return ret
}

func (p *puzzle) visited(pt point) bool {
	index := pt.y*p.w + pt.x
	v := p.pointsVisited[index]
	if !v {
		p.pointsVisited[index] = true
	}
	return v
}

func (p *puzzle) floodFill(pt point) int {
	if p.get(pt) == 9 || p.get(pt) == -1 {
		return 0
	}
	if p.visited(pt) {
		return 0
	}
	sum := p.floodFill(point{pt.x, pt.y - 1}) // top
	sum += p.floodFill(point{pt.x + 1, pt.y}) // right
	sum += p.floodFill(point{pt.x, pt.y + 1}) // bottom
	sum += p.floodFill(point{pt.x - 1, pt.y}) // left
	sum += 1
	return sum
}

func (p *puzzle) calculate() int {
	basins := []int{}
	points := p.lowestPoints()
	fmt.Printf("Num seeds: %d\n", len(points))
	for _, point := range points {
		val := p.floodFill(point)
		fmt.Printf("Value %d\n", val)
		basins = append(basins, val)
	}
	sort.Ints(basins)
	mult := 1
	for i := len(basins) - 3; i < len(basins); i++ {
		fmt.Printf("Largest %d\n", basins[i])
		mult *= basins[i]
	}
	return mult
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

	fmt.Printf("Day 09b\n")
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
