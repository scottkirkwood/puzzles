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

type puzzle struct {
	// edges is the from-to directional edges
	edges map[string][]string

	paths   int
	didNode map[string]bool
}

func read(fname string) (*puzzle, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ret := puzzle{
		edges:   make(map[string][]string),
		paths:   0,
		didNode: make(map[string]bool),
	}
	for scanner.Scan() {
		lineStr := scanner.Text()
		parts := strings.Split(lineStr, "-")
		ret.edges[parts[0]] = append(ret.edges[parts[0]], parts[1])
		if parts[0] != "start" && parts[1] != "end" {
			ret.edges[parts[1]] = append(ret.edges[parts[1]], parts[0])
		}
	}
	return &ret, nil
}

func (p *puzzle) Len() int {
	return len(p.edges)
}

func (p *puzzle) Print() {
	for k, v := range p.edges {
		fmt.Printf("%s->%s\n", k, v)
	}
}

func (p *puzzle) traverse(node string) {
	if node == "end" {
		p.paths++
		return
	}
	if p.didNode[node] {
		return
	}
	if strings.ToLower(node) == node {
		p.didNode[node] = true
	}
	for _, childNode := range p.edges[node] {
		p.traverse(childNode)
	}
	p.didNode[node] = false
}

func (p *puzzle) calculate() int {
	p.Print()
	p.traverse("start")
	return p.paths
}

func parseInts(nums []string) []int {
	ret := make([]int, 0, len(nums))
	for _, num := range nums {
		ret = append(ret, parseInt(num))
	}
	return ret
}

func parseInt(txt string) int {
	num, err := strconv.Atoi(strings.TrimSpace(txt))
	if err != nil {
		fmt.Printf("=========== Bad number %q ==============\n", txt)
	}
	return int(num)
}

func main() {
	flag.Parse()

	fmt.Printf("Day 12a\n")
	infile := "day12.input"
	if *testFileFlag {
		infile = "day12-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", p.Len())
	fmt.Printf("Answer %d\n", p.calculate())
}
