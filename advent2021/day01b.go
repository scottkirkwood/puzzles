package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

func read(fname string) ([]int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []int{}
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error on line %d: %v", err)
		}
		lines = append(lines, int(num))
	}
	return lines, nil
}

func countUp(lines []int) int {
	increased := 0
	prevSum := 0
	for i := range lines {
		if i > 2 {
			curSum := sum3(lines[i-3 : i])
			if prevSum != 0 && curSum > prevSum {
				fmt.Printf("%d Sum %d (increased)\n", i, curSum)
				increased++
			} else {
				fmt.Printf("%d Sum %d (decreased)\n", i, curSum)
			}
			prevSum = curSum
		}
	}
	i := len(lines)
	curSum := sum3(lines[i-3 : i])
	if prevSum != 0 && curSum > prevSum {
		fmt.Printf("%d Sum %d (increased)\n", i, curSum)
		increased++
	}

	return increased
}

func sum3(toSum []int) int {
	sum := 0
	for _, n := range toSum {
		sum += n
	}
	return sum
}

func main() {
	flag.Parse()

	fmt.Printf("Day 01\n")
	infile := "day01.input"
	if *testFileFlag {
		infile = "day01-sample.input"
	}
	lines, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", len(lines))
	fmt.Printf("Count up %d\n", countUp(lines))
}
