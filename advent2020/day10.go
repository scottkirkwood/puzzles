// See https://adventofcode.com/2020/day/10 for problem decscription
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var (
	testFileFlag  = flag.Bool("t", false, "Use the test file")
	testFile2Flag = flag.Bool("t2", false, "Use the test2 file")
)

func read(fname string) ([]int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []int{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		cur, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		ret = append(ret, cur)
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

type jolty struct {
	deviceJolts int
	orig        []int
	nums        []bool
	quit        bool
	sequence    []int
}

func newJolty(nums []int) jolty {
	sort.Ints(nums)
	max := nums[len(nums)-1]

	j := jolty{
		orig:        nums,
		deviceJolts: max + 3,
		nums:        make([]bool, max+1),
	}
	for _, v := range nums {
		if j.nums[v] {
			fmt.Printf("Duplicate value at %d\n", v)
		}
		j.nums[v] = true
	}
	return j
}

func (j *jolty) doChain(need int) {
	if j.quit {
		return
	}
	if need == j.deviceJolts-3 {
		if j.visitedAll() {
			fmt.Printf("Done!\n")
			j.printSummary()
			j.quit = true
		}
		return
	}
	possible := j.findPossible(need)
	if len(possible) == 0 {
		fmt.Printf("Couldn't find for %d jolts\n", need)
		return
	}
	for _, p := range possible {
		if j.quit {
			return
		}
		diff := p - need
		j.nums[p] = false
		j.sequence = append(j.sequence, p)
		j.doChain(need + diff)
		j.sequence = j.sequence[:len(j.sequence)-1]
		j.nums[p] = true
	}
}

func (j *jolty) visitedAll() bool {
	return len(j.sequence) == len(j.orig)
}

func (j *jolty) findPossible(need int) []int {
	ret := make([]int, 0, 3)
	if need+1 < len(j.nums) && j.nums[need+1] {
		ret = append(ret, need+1)
	}
	if need+2 < len(j.nums) && j.nums[need+2] {
		ret = append(ret, need+2)
	}
	if need+3 < len(j.nums) && j.nums[need+3] {
		ret = append(ret, need+3)
	}
	return ret
}

func (j *jolty) printSummary() {
	jolts := 0
	numOnes := 0
	numThrees := 0
	for _, s := range j.sequence {
		val := s
		diff := val - jolts
		if diff == 1 {
			numOnes++
		} else if diff == 3 {
			numThrees++
		}
		jolts = jolts + diff
		fmt.Printf("%d jolts, diff %d\n", val, diff)
	}
	numThrees++ // Because of final one?
	fmt.Printf("ones %d, threes %d = %d\n", numOnes, numThrees, numOnes*numThrees)
}

func main() {
	flag.Parse()

	fmt.Printf("Day 10\n")

	infile := "day10.input"
	if *testFileFlag {
		infile = "day10-sample.input"
	} else if *testFile2Flag {
		infile = "day10-sample2.input"
	}
	nums, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(nums))
	j := newJolty(nums)
	fmt.Printf("Device Joltage %d\n", j.deviceJolts)
	j.doChain(0)
}
