package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type rows struct {
}

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

func main() {
	flag.Parse()

	fmt.Printf("Day 6\n")
	infile := "day6.input"
	if *testFileFlag {
		infile = "day6-sample.input"
	}
}
