package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type dir struct {
	dir    string
	amount int
}

var lineRx = regexp.MustCompile(`(forward|down|up) (\d+)`)

func read(fname string) ([]dir, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []dir{}
	for scanner.Scan() {
		match := lineRx.FindStringSubmatch(scanner.Text())
		if len(match) != 3 {
			return nil, fmt.Errorf("error on line %d: %v", len(lines), scanner.Text())
		}

		lines = append(lines, dir{
			dir:    match[1],
			amount: parseNum(match[2]),
		})
	}
	return lines, nil
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("=========== Bad number %q ==============\n", txt)
	}
	return int(num)
}

func traverse(lines []dir) int {
	x, y := 0, 0
	for _, line := range lines {
		switch line.dir {
		case "forward":
			x += line.amount
		case "up":
			y -= line.amount
		case "down":
			y += line.amount
		}
	}
	fmt.Printf("x:%d, y:%d\n", x, y)
	return x * y
}

func main() {
	flag.Parse()

	fmt.Printf("Day 02\n")
	infile := "day02.input"
	if *testFileFlag {
		infile = "day02-sample.input"
	}
	lines, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", len(lines))
	fmt.Printf("Answer %d\n", traverse(lines))
}
