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
	sums  map[string]int
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
		sums:  make(map[string]int),
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

	//p.sums[" "+p.cur[0]] = 1
	//p.sums[p.cur[len(p.cur)-1]+" "] = 1
	for i := 0; i < len(p.cur)-1; i++ {
		p.sums[p.cur[i]+p.cur[i+1]]++
	}
	return &p, nil
}

func (p *puzzle) Len() int {
	return len(p.pairs)
}

func (p *puzzle) Print() {
	fmt.Printf("%v\n", p.sums)
}

func copyMap(cur map[string]int) map[string]int {
	ret := make(map[string]int, len(cur))
	for chch, _ := range cur {
		ret[chch] = 0
	}
	return ret
}

func (p *puzzle) generate() {
	newSums := copyMap(p.sums)
	for chch, count := range p.sums {
		toCh, ok := p.pairs[chch]
		if ok {
			//newSums[chch] -= count
			leftKey := string(chch[0]) + toCh
			rightKey := toCh + string(chch[1])
			newSums[leftKey] += count
			newSums[rightKey] += count
		} else {
			//		newSums[chch] = count
		}
	}
	p.sums = newSums
}

func (p *puzzle) score() int {
	freq := map[string]int{}
	for chch, count := range p.sums {
		freq[string(chch[0])] += count
		freq[string(chch[1])] += count
	}
	min, max := 1<<60, 0
	minCh, maxCh := "", ""
	for ch, v := range freq {
		v /= 2
		if ch == " " {
			continue
		}
		if v < min {
			min = v
			minCh = ch
		}
		if v > max {
			max = v
			maxCh = ch
		}
	}
	fmt.Printf("Max %q: %d, Min %q: %d\n", maxCh, max, minCh, min)
	// 3390034818251 too high
	// 3390034818250 too high
	// 3390034818245 too low
	// 3390034818249 answer!
	return int(max - min)+1
}

func (p *puzzle) calculate() int {
	p.Print()
	for i := 0; i < *stepsFlag; i++ {
		fmt.Printf("Starting step %d\n", i)
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

	fmt.Printf("Day 14b\n")
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
