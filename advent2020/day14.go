package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	testFileFlag = flag.Bool("t", false, "Use the test file")
)

type machine struct {
	mem     map[int64]int64
	curMask string
}

func (m *machine) evaluate(inst instruction) {
	if inst.memLoc > 0 {
		oldVal := m.mem[inst.memLoc]
		val := oldVal & m.bitsToKeep()
		val |= (val&m.bitsToZero() | (m.bitsToSet() & inst.val))
		m.mem[inst.memLoc] = val
	} else {
		m.curMask = inst.mask
	}
}

func (m machine) String(inst instruction) string {
	return fmt.Sprintf("value:  %s\nmask:   %s\nresult: %s", formatBits(inst.val), m.curMask, formatBits(m.mem[inst.memLoc]))
}

func formatBits(val int64) string {
	return fmt.Sprintf("%036s", strconv.FormatInt(val, 2))
}

type instruction struct {
	mask   string
	memLoc int64
	val    int64
}

var (
	maskRx = regexp.MustCompile(`mask = ([01X]{36})`)
	memRx  = regexp.MustCompile(`mem\[(\d+)\] = (\d+)`)
)

func read(fname string) ([]instruction, error) {
	ret := []instruction{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := maskRx.FindStringSubmatch(line)
		if len(parts) == 2 {
			ret = append(ret, instruction{
				mask: parts[1],
			})
		} else {
			parts = memRx.FindStringSubmatch(line)
			if len(parts) != 3 {
				return nil, fmt.Errorf("Unable to parse %q\n", line)
			}
			ret = append(ret, instruction{
				memLoc: parseNum(parts[1]),
				val:    parseNum(parts[2]),
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func (m machine) bitsToSet() int64 {
	n, err := parseBits(strings.Map(func(r rune) rune {
		if r == 'X' {
			return '0'
		}
		return r
	}, m.curMask))
	if err != nil {
		fmt.Printf("Unable to parse bits to set %q\n", m.curMask)
	}
	return n
}

func (m machine) bitsToZero() int64 {
	n, err := parseBits(strings.Map(func(r rune) rune {
		if r == '0' {
			return '0'
		}
		return '1'
	}, m.curMask))
	if err != nil {
		fmt.Printf("Unable to parse bits to zero %q\n", m.curMask)
	}
	return n
}

func (m machine) bitsToKeep() int64 {
	n, err := parseBits(strings.Map(func(r rune) rune {
		if r == 'X' {
			return '1'
		}
		return '0'
	}, m.curMask))
	if err != nil {
		fmt.Printf("Unable to parse bits to keep %q\n", m.curMask)
	}
	return n
}

func parseNum(txt string) int64 {
	num, err := strconv.ParseInt(txt, 10, 64)
	if err != nil {
		fmt.Printf("Unable to parse decimal %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func parseBits(txt string) (int64, error) {
	num, err := strconv.ParseInt(txt, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse bits %", txt)
	}
	return num, nil
}

func main() {
	flag.Parse()

	fmt.Printf("Day 14\n")

	infile := "day14.input"
	if *testFileFlag {
		infile = "day14-sample.input"
	}
	instructions, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Instructions %d\n", len(instructions))
	m := machine{mem: map[int64]int64{}}

	for _, inst := range instructions {
		m.evaluate(inst)
		if inst.memLoc != 0 {
			fmt.Printf("%s\n\n", m.String(inst))
		}
	}
}
