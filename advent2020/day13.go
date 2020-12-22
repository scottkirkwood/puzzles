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

type busCalc struct {
	ts    int
	buses []int
}

func (b busCalc) calcClosest() (int, int) {
	ts := b.ts
	for {
		for _, bus := range b.buses {
			if ts%bus == 0 {
				return bus, ts
			}
		}
		ts++
	}
	return 0, 0
}

func read(fname string) (busCalc, error) {
	ret := busCalc{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if ret.ts == 0 {
			ret.ts = parseNum(line)
		} else {
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if part == "x" {
					continue
				}
				ret.buses = append(ret.buses, parseNum(part))
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
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Buses %d\n", len(b.buses))
	bus, ts := b.calcClosest()
	delta := ts - b.ts
	fmt.Printf("Bus %d, ts %d, delta %d mult %d\n", bus, ts, delta, bus*delta)
}
