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

var testFileFlag = flag.Bool("t", false, "Use the test file")

type bag struct {
	color    string
	contains map[string]int
}

func read(fname string) ([]bag, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []bag{}
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

var (
	colorRx  = regexp.MustCompile(`(\w+ \w+) bags? contain `)
	color2Rx = regexp.MustCompile(`(\d+)? ?(\w+ \w+) bags?[.]?`)
)

func parseLine(line string) (b bag, err error) {
	b.contains = map[string]int{}

	parts := colorRx.FindStringSubmatch(line)
	if len(parts) != 2 {
		return b, fmt.Errorf("bad key %q", line)
	}
	b.color = parts[1]
	length := len(parts[0])
	splits := strings.Split(line[length:], ", ")
	for _, split := range splits {
		parts := color2Rx.FindStringSubmatch(split)
		if len(parts) != 3 {
			return b, fmt.Errorf("bad part %q", split)
		}
		if parts[2] == "no other" {
			continue
		}
		count, err := strconv.Atoi(parts[1])
		if err != nil || count == 0 {
			return b, fmt.Errorf("bad number %q", parts[1])
		}
		b.contains[parts[2]] = count
	}
	return b, err
}

func trace(bags []bag, toTrace string) int {
	couldContain := map[string]bool{}
	changed := true
	for changed {
		changed = false
		for _, bag := range bags {
			if couldContain[bag.color] {
				continue
			}
			if bag.color == toTrace {
				continue
			}
			if bag.contains[toTrace] > 0 {
				couldContain[bag.color] = true
				changed = true
			}
			for key, _ := range bag.contains {
				if couldContain[key] {
					couldContain[bag.color] = true
					changed = true
					break
				}
			}
		}
	}
	return len(couldContain)
}

func makeBagMap(bags []bag) map[string]bag {
	ret := make(map[string]bag, len(bags))
	for _, bag := range bags {
		ret[bag.color] = bag
	}
	return ret
}

func recurseBags(bags map[string]bag, cur string) int {
	if len(bags[cur].contains) == 0 {
		fmt.Printf("%s has 0 other bags\n", cur)
		return 0
	}
	count := 0
	for color, num := range bags[cur].contains {
		childCount := recurseBags(bags, color) + 1
		fmt.Printf("> %s has %d*%d = %d bags\n", color, num, childCount, num*childCount)
		count += num * childCount
	}
	return count
}

func main() {
	flag.Parse()

	fmt.Printf("Day 7\n")

	infile := "day7.input"
	if *testFileFlag {
		infile = "day7-sample.input"
	}
	bags, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(bags))
	fmt.Printf("Number of bag colors that can contain %d\n", trace(bags, "shiny gold"))

	bagMap := makeBagMap(bags)

	fmt.Printf("Number of bags inside shiny gold: %d\n", recurseBags(bagMap, "shiny gold"))
}
