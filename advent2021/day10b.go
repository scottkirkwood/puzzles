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

	scoring = map[string]int{
		")": 1,
		"]": 2,
		"}": 3,
		">": 4,
	}
	openCloseMap = map[string]string{
		"(": ")",
		"[": "]",
		"{": "}",
		"<": ">",
	}
)

type point struct {
	x, y int
}

type puzzle struct {
	lines []string
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
		ret.lines = append(ret.lines, lineStr)
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.lines)
}

func (p *puzzle) Print() {
}

func (p *puzzle) calculate() int {
	scores := []int{}
	for lineNum, line := range p.lines {
		stack := []string{}
		isCorrupted := false
		for i := 0; i < len(line); i++ {
			ch := string(line[i])
			toClose, ok := openCloseMap[ch]
			if ok {
				stack = append(stack, toClose)
			} else {
				pop := stack[len(stack)-1]
				if pop != ch {
					isCorrupted = true
					break
				}
				stack = stack[0 : len(stack)-1]
			}
		}
		if !isCorrupted {
			fmt.Printf("%d: Want: ", lineNum)
			curScore := 0
			for j := len(stack) - 1; j >= 0; j-- {
				fmt.Printf("%s", stack[j])
				curScore *= 5
				curScore += scoring[stack[j]]
			}
			fmt.Printf(" Score %d\n", curScore)
			scores = append(scores, curScore)
		}
	}
	sort.Ints(scores)
	return scores[len(scores)/2]
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

	fmt.Printf("Day 10a\n")
	infile := "day10.input"
	if *testFileFlag {
		infile = "day10-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
