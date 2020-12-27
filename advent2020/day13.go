// See https://adventofcode.com/2020/day/13 for problem decscription
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

type busDelta struct {
	bus   int
	delta int
}

func calc(bds []busDelta) {
	ts := 0
	gcf := 1
	for _, bd := range bds {
		for {
			if (ts+bd.delta)%bd.bus == 0 {
				break
			}
			ts += gcf
		}
		gcf *= bd.bus
	}
	fmt.Printf("ts %d\n", ts)
}

func okTs(bds []busDelta, ts int) bool {
	for _, bd := range bds {
		if (ts+bd.delta)%bd.bus != 0 {
			return false
		}
	}
	return true
}

func read(fname string) ([]busDelta, error) {
	ret := []busDelta{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if lineNum == 1 {
			continue
		} else {
			parts := strings.Split(line, ",")
			for i, part := range parts {
				if part == "x" {
					continue
				}
				ret = append(ret, busDelta{parseNum(part), i})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("Unable to parse %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func main() {
	flag.Parse()

	fmt.Printf("Day 13\n")

	infile := "day13.input"
	if *testFileFlag {
		infile = "day13-sample.input"
	}
	bds, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Buses %d\n", len(bds))
	calc(bds)
}
