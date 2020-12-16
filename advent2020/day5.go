package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type seat string

var rowMapper = func(r rune) rune {
	switch r {
	case 'F':
		return '0'
	case 'B':
		return '1'
	}
	return ' '
}

func (s seat) Row() int {
	binary := strings.Map(rowMapper, string(s[0:7]))
	num, err := strconv.ParseInt(binary, 2, 32)
	if err != nil {
		fmt.Printf("Problem parseing %v\n", err)
	}
	return int(num)
}

var colMapper = func(r rune) rune {
	switch r {
	case 'L':
		return '0'
	case 'R':
		return '1'
	}
	return ' '
}

func (s seat) Col() int {
	binary := strings.Map(colMapper, string(s[7:]))
	num, err := strconv.ParseInt(binary, 2, 32)
	if err != nil {
		fmt.Printf("Problem parseing %v\n", err)
	}
	return int(num)
}

func (s seat) ID() int {
	return s.Row()*8 + s.Col()
}

func read(fname string) ([]seat, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []seat{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, seat(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func test(in string, wantRow, wantCol, wantID int) {
	s := seat(in)
	if s.Row() != wantRow {
		fmt.Printf("%q, Row = %d, want %d\n", in, s.Row(), wantRow)
	}
	if s.Col() != wantCol {
		fmt.Printf("%q, Col = %d, want %d\n", in, s.Col(), wantCol)
	}
	if s.ID() != wantID {
		fmt.Printf("%q, ID = %d, want %d\n", in, s.ID(), wantID)
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Day 5\n")
	test("BFFFBBFRRR", 70, 7, 567)
	test("FFFBBBFRRR", 14, 7, 119)
	test("BBFFBBFRLL", 102, 4, 820)
	infile := "day5.input"
	if *testFileFlag {
		infile = "day5-sample.input"
	}
	seats, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Printf("%d seats\n", len(seats))
	max := 0
	ids := []int{}
	for _, seat := range seats {
		ids = append(ids, seat.ID())
		if seat.ID() > max {
			max = seat.ID()
		}
	}
	fmt.Printf("Max ID %d\n", max)
	sort.Ints(ids)
	prev := 0
	for _, id := range ids {
		if id == prev+2 {
			fmt.Printf("Your seat %d\n", id-1)
		}
		prev = id
	}
}
