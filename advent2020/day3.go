// See https://adventofcode.com/2020/day/3 for problem decscription
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type rows []string

func read(fname string) (rows, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make(rows, 0, 323)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func (r rows) get(x, y int) string {
	if x >= len(r[0]) {
		x = x % len(r[0])
	}
	return r[y][x : x+1]
}

func (r rows) walk(dx, dy int) int {
	trees := 0
	rows := len(r)
	x := 0
	for y := 0; y < rows; y += dy {
		if r.get(x, y) == "#" {
			trees++
		}
		x += dx
	}
	return trees
}

func main() {
	flag.Parse()

	fmt.Printf("Day 3\n")
	infile := "day3.input"
	if *testFileFlag {
		infile = "day3-sample.input"
	}
	trees, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(trees))

	tries := []struct {
		dx, dy int
	}{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	mult := 1
	for _, trial := range tries {
		got := trees.walk(trial.dx, trial.dy)
		fmt.Printf("%d, %d = %d\n", trial.dx, trial.dy, got)
		mult *= got
	}
	fmt.Printf("Mult = %d\n", mult)
}
