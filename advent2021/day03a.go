package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type bits struct {
	numBits int
	values  []uint64
}

func read(fname string) (*bits, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ret := bits{
		numBits: 0,
		values:  []uint64{},
	}
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.ParseUint(line, 2, 64)
		if err != nil {
			return &ret, err
		}
		if ret.numBits == 0 {
			ret.numBits = len(line)
			fmt.Printf("Bit size %d\n", len(line))
		}
		ret.values = append(ret.values, num)
	}
	return &ret, nil
}

func (b *bits) Len() int {
	return len(b.values)
}

func (b *bits) countEm() uint64 {
	gammaRate := uint64(0)
	epsilonRate := uint64(0)
	for bitIndex := 0; bitIndex < b.numBits; bitIndex++ {
		mask := uint64(1) << bitIndex
		sum0, sum1 := 0, 0
		for _, v := range b.values {
			if v&mask != 0 {
				sum1++
			} else {
				sum0++
			}
		}
		fmt.Printf("Bit %d, sum0: %d, sum1: %d\n", bitIndex, sum0, sum1)
		if sum1 > sum0 {
			gammaRate |= 1 << bitIndex
		}
		if sum1 < sum0 {
			epsilonRate |= 1 << bitIndex
		}
	}
	fmt.Printf("Gamma rate %d, Epsilon rate %d\n", gammaRate, epsilonRate)
	return gammaRate * epsilonRate
}

func main() {
	flag.Parse()

	fmt.Printf("Day 03a\n")
	infile := "day03.input"
	if *testFileFlag {
		infile = "day03-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.countEm())
}
