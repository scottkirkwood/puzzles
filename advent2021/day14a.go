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
	pairs map[string]string
	cur []string
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
	}
	w, h := p.maxXY()
	p.paper = makeGrid(w, h)
	p.paper.setPoints(p.points)
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.points)
}

func (p *puzzle) Print() {
}

func (p *puzzle) calculate() int {
	// p.Print()
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

	fmt.Printf("Day 14a\n")
	infile := "day14.input"
	if *testFileFlag {
		infile = "day14-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
