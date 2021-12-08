package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type puzzle struct {
	timers []int
}

func toIntMap(numbers []int) map[int]bool {
	ret := make(map[int]bool, len(numbers))
	for _, n := range numbers {
		ret[n] = true
	}
	return ret
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
		ret.timers = parseInts(strings.Split(lineStr, ","))
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.timers)
}

func printTimers(timers []int) string {
	strList := make([]string, len(timers))
	for i, t := range timers {
		strList[i] = strconv.Itoa(t)
	}
	return strings.Join(strList, ",")
}

func (p *puzzle) cycle() {
	newChildren := 0
	for i := 0; i < len(p.timers); i++ {
		if p.timers[i] == 0 {
			newChildren++
			p.timers[i] = 6
		} else {
			p.timers[i]--
		}
	}
	for i := 0; i < newChildren; i++ {
		p.timers = append(p.timers, 8)
	}
}

func (p *puzzle) printDay(dayNum int) {
	fmt.Printf("After %2d days: %s\n", dayNum, printTimers(p.timers))
}

func (p *puzzle) calculate() int {
	fmt.Printf("Initial State: %s\n", printTimers(p.timers))
	for day := 1; day <= 80; day++ {
		p.cycle()
	}
	return len(p.timers)
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

	fmt.Printf("Day 06a\n")
	infile := "day06.input"
	if *testFileFlag {
		infile = "day06-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
