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
	daysFlag     = flag.Int("days", 18, "Number of days")
)

type puzzle struct {
	timers [9]int
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
		initial := parseInts(strings.Split(lineStr, ","))
		for _, j := range initial {
			ret.timers[j]++
		}
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.timers)
}

func printTimers(timers [9]int) string {
	strList := make([]string, len(timers))
	for i, t := range timers {
		strList[i] = strconv.Itoa(t)
	}
	return strings.Join(strList, ",")
}

func (p *puzzle) cycle() {
	prev := [9]int{}
	for i := 0; i <= 8; i++ {
		prev[i] = p.timers[i]
	}
	p.timers[8] = prev[0]
	for i := 0; i < 8; i++ {
		p.timers[i] = prev[i+1]
	}
	p.timers[6] += prev[0]
}

func (p *puzzle) printDay(dayNum int) {
	fmt.Printf("After %2d days: %s\n", dayNum, printTimers(p.timers))
}

func (p *puzzle) sum() int {
	sum := 0
	for i := 0; i <= 8; i++ {
		sum += p.timers[i]
	}
	return sum
}

func (p *puzzle) calculate(days int) int {
	fmt.Printf("Initial State: %s\n", printTimers(p.timers))
	for day := 1; day <= days; day++ {
		p.cycle()
		fmt.Printf("Day %d: %s\n", day, printTimers(p.timers))
	}
	return p.sum()
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

	fmt.Printf("Day 06b\n")
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
	fmt.Printf("Answer %d\n", b.calculate(*daysFlag))
}
