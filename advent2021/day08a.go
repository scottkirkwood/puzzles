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

	segments = map[int]string{
		1: "cf",
		7: "acf",
		4: "bcdf",
		2: "acdeg",
		3: "acdfg",
		5: "abdfg",
		0: "abcefg",
		6: "abdefg",
		9: "abcdfg",
		8: "abcdefg",
	}
)

type lineInfo struct {
	lhs []string
	rhs []string
}

type puzzle struct {
	lines []lineInfo
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
		leftRight := strings.Split(lineStr, "|")
		left := strings.TrimSpace(leftRight[0])
		right := strings.TrimSpace(leftRight[1])
		ret.lines = append(ret.lines, lineInfo{
			splitSort(left),
			splitSort(right),
		})
	}
	return &ret, nil
}

func splitSort(line string) []string {
	parts := strings.Split(line, " ")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		split := strings.Split(part, "")
		sort.Strings(split)
		ret = append(ret, strings.Join(split, ""))
	}
	return ret
}

func (p *puzzle) Len() int {
	return len(p.lines)
}

func (p *puzzle) String() string {
	return ""
}

func (p *puzzle) calculate() int {
	count := 0
	for _, line := range p.lines {
		for _, part := range line.rhs {
			l := len(part)
			if l == 2 || l == 3 || l == 4 || l == 7 {
				count++
			}
		}
	}
	return count
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

	fmt.Printf("Day 08a\n")
	infile := "day08.input"
	if *testFileFlag {
		infile = "day08-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
