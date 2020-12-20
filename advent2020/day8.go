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

type machine struct {
	p     int
	acc   int
	ops   []inst
	visit map[int]bool
}

type inst struct {
	op  string
	arg int
}

func read(fname string) ([]inst, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []inst{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		cur, err := parseLine(line)
		if err != nil {
			return ret, err
		}
		ret = append(ret, cur)
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

var instRx = regexp.MustCompile(`(\w+) ([+-]\d+)`)

func parseLine(line string) (op inst, err error) {
	parts := instRx.FindStringSubmatch(line)
	if len(parts) != 3 {
		return op, fmt.Errorf("bad instruction %q", line)
	}
	op.op = parts[1]
	arg, err := strconv.Atoi(parts[2])
	if err != nil {
		return op, fmt.Errorf("bad number %q", parts[2])
	}
	op.arg = arg
	return op, err
}

func (c *machine) execute() bool {
	for {
		if c.p == len(c.ops) {
			return true
		}
		if c.visit[c.p] {
			return false
		}
		c.visit[c.p] = true
		op, arg := c.ops[c.p].op, c.ops[c.p].arg
		switch op {
		case "acc":
			c.acc += arg
			c.p++
		case "jmp":
			c.p += arg
		case "nop":
			c.p++
		}
	}
	return false
}

func (c *machine) toggle(i int) {
	if i < 0 {
		return
	}
	if c.ops[i].op == "jmp" {
		c.ops[i].op = "nop"
	} else {
		c.ops[i].op = "jmp"
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Day 8\n")

	infile := "day8.input"
	if *testFileFlag {
		infile = "day8-sample.input"
	}
	ops, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(ops))
	lastI := -1
	for i := 0; i < len(ops); i++ {
		if ops[i].op == "acc" {
			continue
		}
		cpu := machine{
			ops:   ops[:],
			visit: make(map[int]bool),
		}
		cpu.toggle(lastI)
		lastI = i
		cpu.toggle(i)
		if cpu.execute() {
			inst := cpu.ops[i]
			fmt.Printf("Change instruction %d (to %s %d), val %d\n", i, inst.op, inst.arg, cpu.acc)
			break
		}
	}
}
