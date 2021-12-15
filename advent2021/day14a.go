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
	stepsFlag    = flag.Int("n", 10, "Number of steps")
)

type puzzle struct {
	pairs map[string]string
	cur   []string
}

func read(fname string) (*puzzle, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	p := puzzle{
		pairs: make(map[string]string),
	}
	for scanner.Scan() {
		lineStr := scanner.Text()
		if len(lineStr) == 0 {
			continue
		}
		if len(p.cur) == 0 {
			p.cur = strings.Split(lineStr, "")
		} else {
			parts := strings.Split(lineStr, " -> ")
			key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if len(p.pairs[key]) != 0 {
				return nil, fmt.Errorf("unexpected duplicate pair %q", key)
			}
			p.pairs[key] = val
		}
	}
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.pairs)
}

func (p *puzzle) Print() {
	fmt.Printf("%s\n", strings.Join(p.cur, ""))
}

func (p *puzzle) generate() {
	next := make([]string, 0, len(p.cur)*2)
	next = append(next, p.cur[0])
	for i := 0; i < len(p.cur)-1; i++ {
		pair := p.cur[i] + p.cur[i+1]
		insert, ok := p.pairs[pair]
		if ok {
			next = append(next, insert, p.cur[i+1])
		} else {
			next = append(next, p.cur[i+1])
		}
	}
	p.cur = next
}

func (p *puzzle) score() int {
	freq := map[string]int{}
	for _, ch := range p.cur {
		freq[ch]++
	}
	min, max := 1<<60, 0
	minCh, maxCh := "", ""
	for k, v := range freq {
		if v < min {
			min = v
			minCh = k
		}
		if v > max {
			max = v
			maxCh = k
		}
	}
	fmt.Printf("Min %s: %d, Max %s: %d\n", minCh, min, maxCh, max)
	return max - min
}

func (p *puzzle) calculate() int {
	p.Print()
	for i := 0; i < *stepsFlag; i++ {
		p.generate()
		p.Print()
	}
	return p.score()
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
