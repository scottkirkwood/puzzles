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
	positions []int
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
		ret.positions = parseInts(strings.Split(lineStr, ","))
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.positions)
}

func (p *puzzle) String() string {
	strList := make([]string, len(p.positions))
	for i, t := range p.positions {
		strList[i] = strconv.Itoa(t)
	}
	return strings.Join(strList, ",")
}

func (p *puzzle) getMinMax() (min, max int) {
	for _, p := range p.positions {
		if p > max {
			max = p
		}
		if p < min {
			min = p
		}
	}
	return
}

func (p *puzzle) calcCostForPosition(pos int) int {
	total := 0
	for _, p := range p.positions {
		diff := pos - p
		if diff < 0 {
			diff = -diff
		}
		total += diff
	}
	return total
}

func (p *puzzle) calculate() int {
	min, max := p.getMinMax()
	minCost := 2147483647
	minp := 0
	for pos := min; pos <= max; pos++ {
		cost := p.calcCostForPosition(pos)
		if cost < minCost {
			minCost = cost
			minp = pos
			fmt.Printf("Pos %d, cost %d\n", pos, cost)
		}
	}
	fmt.Printf("Pos %d, cost %d\n", minp, minCost)
	return minCost
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

	fmt.Printf("Day 07a\n")
	infile := "day07.input"
	if *testFileFlag {
		infile = "day07-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
