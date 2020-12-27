// See https://adventofcode.com/2020/day/15 for problem decscription
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

type game struct {
	turn   int
	last   int
	spoken map[int][]int
}

func (g *game) initStarting(val int) (int, int) {
	g.spoken[val] = append(g.spoken[val], g.turn)
	g.turn++
	g.last = val
	return g.turn - 1, g.last
}

func (g *game) set(val int) (int, int) {
	g.spoken[val] = append(g.spoken[val], g.turn)
	g.turn++
	g.last = val
	return g.turn - 1, g.last
}

func (g *game) eval() (int, int) {
	if len(g.spoken[g.last]) == 1 {
		g.set(0)
	} else {
		diff := g.diffLastTwo()
		g.set(diff)
	}
	return g.turn - 1, g.last
}

func (g *game) diffLastTwo() int {
	lst := g.spoken[g.last]
	ln := len(lst)
	return lst[ln-1] - lst[ln-2]
}

func read(fname string) ([]int, error) {
	ret := []int{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		for _, part := range parts {
			ret = append(ret, parseNum(part))
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
		fmt.Printf("Unable to parse decimal %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func main() {
	flag.Parse()

	fmt.Printf("Day 15\n")

	infile := "day15.input"
	if *testFileFlag {
		infile = "day15-sample.input"
	}
	nums, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Nums %d\n", len(nums))
	g := game{
		turn:   1,
		spoken: map[int][]int{},
	}
	for _, x := range nums {
		turn, val := g.initStarting(x)
		fmt.Printf("Turn %d: %d\n", turn, val)
	}
	for i := len(nums); i < 30000000; i++ {
		g.eval()
	}
	fmt.Printf("Turn %d: %d\n", g.turn-1, g.last)
}
