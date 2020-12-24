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

type floor struct {
	tiles map[hex]bool
}

type hex struct {
	q, r int
}

const (
	E = iota
	NE
	NW
	W
	SW
	SE
	NUM_DIRECTIONS
)

var axialDirections = []hex{
	hex{+1, 0}, hex{+1, -1}, hex{0, -1},
	hex{-1, 0}, hex{-1, +1}, hex{0, +1},
}

func (f *floor) flip(hx hex) {
	f.tiles[hx] = !f.tiles[hx]
}

func (f floor) countBlack() int {
	count := 0
	for _, black := range f.tiles {
		if black {
			count++
		}
	}
	return count
}

func (h *hex) move(hx hex) {
	h.q += hx.q
	h.r += hx.r
}

func letterToDirection(dir string) int {
	switch dir {
	case "e":
		return E
	case "ne":
		return NE
	case "nw":
		return NW
	case "w":
		return W
	case "sw":
		return SW
	case "se":
		return SE
	}
	fmt.Printf("Invalid direction %q\n", dir)
	return NUM_DIRECTIONS
}

var directionRx = regexp.MustCompile(`(e|se|sw|w|nw|ne)`)

func read(fname string) ([][]hex, error) {
	ret := [][]hex{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ret = append(ret, []hex{})
		for _, ch := range directionRx.FindAllString(line, -1) {
			ret[len(ret)-1] = append(ret[len(ret)-1], axialDirections[letterToDirection(ch)])
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

	fmt.Printf("Day 24\n")

	infile := "day24.input"
	if *testFileFlag {
		infile = "day24-sample.input"
	}
	directions, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("%d\n", len(directions))
	f := floor{
		tiles: map[hex]bool{},
	}
	for _, moves := range directions {
		start := hex{0, 0}
		for _, move := range moves {
			start.move(move)
		}
		f.flip(start)
	}
	fmt.Printf("Black tiles %d\n", f.countBlack())
}
