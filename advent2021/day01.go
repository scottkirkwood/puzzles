package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

func read(fname string) ([]int32, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []int32{}
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error on line %d: %v", err)
		}
		lines = append(lines, int32(num))
	}
	return lines, nil
}

func countUp(lines []int32) int {
	increased := 0
	prev := int32(0)
	for i, num := range lines {
		if i > 0 {
			if num > prev {
				increased++
			}
		}
		prev = num
	}
	return increased
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
