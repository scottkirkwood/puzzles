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

func sum0s1s(values []uint64, bitIndex int) (int, int) {
	mask := uint64(1) << bitIndex
	sum0, sum1 := 0, 0
	for _, v := range values {
		if v&mask != 0 {
			sum1++
		} else {
			sum0++
		}
	}
	return sum0, sum1
}

func filter(values []uint64, mask, match uint64) []uint64 {
	ret := make([]uint64, 0, len(values))
	for _, v := range values {
		if mask&v == match {
			ret = append(ret, v)
		}
	}
	return ret
}

func (b *bits) printBinary(values []uint64) {
	lines := []string{}
	for _, v := range values {
		bits := strconv.FormatUint(v, 2)
		for len(bits) < b.numBits {
			bits = "0" + bits
		}
		lines = append(lines, bits)
	}
	fmt.Printf("%s\n", strings.Join(lines, ", "))
}

func (b *bits) countEm(oxygen bool) uint64 {
	remaining := make([]uint64, len(b.values))
	copy(remaining, b.values)
	for bitIndex := b.numBits - 1; bitIndex >= 0; bitIndex-- {
		mask := uint64(1) << bitIndex
		match := uint64(0)
		sum0, sum1 := sum0s1s(remaining, bitIndex)
		if oxygen {
			match = 0
			if sum1 >= sum0 {
				match = mask
			}
		} else {
			match = 0
			if sum1 < sum0 {
				match = mask
			}
		}
		remaining = filter(remaining, mask, match)
		b.printBinary(remaining)
		if len(remaining) == 1 {
			return remaining[0]
		}
	}
	return 0
}

func (b *bits) calculate() uint64 {
	oxygen := b.countEm(true)
	fmt.Printf("Oxygen %d\n", oxygen)
	co2 := b.countEm(false)
	fmt.Printf("CO2 %d\n", co2)
	return oxygen * co2
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
	fmt.Printf("Answer %d\n", b.calculate())
}
