// See https://adventofcode.com/2020/day/6 for problem decscription
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type group struct {
	question map[string]bool
}

func read(fname string) ([]group, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []group{}
	scanner := bufio.NewScanner(file)

	prev := group{question: map[string]bool{}}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			ret = append(ret, prev)
			prev = group{question: map[string]bool{}}
			continue
		}
		cur := map[string]bool{}
		if len(prev.question) == 0 {
			for _, part := range line {
				prev.question[string(part)] = true
			}
		} else {
			for _, part := range line {
				key := string(part)
				cur[key] = true
				if val, ok := prev.question[key]; !ok || !val {
					prev.question[key] = false
				}
			}
			// intersect
			for key, val := range prev.question {
				if val && !cur[key] {
					prev.question[key] = false
				}
			}
		}
	}
	ret = append(ret, prev)
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func sumCounts(groups []group) int {
	counts := 0
	for _, group := range groups {
		count := 0
		for _, val := range group.question {
			if val {
				count++
			}
		}
		fmt.Printf("Count %d\n", count)
		counts += count
	}
	return counts
}

func sumCounts2(groups []group) int {
	count := 0
	for _, group := range groups {
		count += len(group.question)
	}
	return count
}

func main() {
	flag.Parse()

	fmt.Printf("Day 6\n")
	infile := "day6.input"
	if *testFileFlag {
		infile = "day6-sample.input"
	}
	groups, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(groups))
	fmt.Printf("Sum: %d\n", sumCounts(groups))
}
