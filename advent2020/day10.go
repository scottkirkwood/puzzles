package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

func deviceJoltage(nums []int) int {
	max := 0
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max + 3
}

type jolty struct {
	deviceJolts int
	done        map[int]bool
	nums        []int
	quit        bool
	sequence    []int
}

func newJolty(nums []int) jolty {
	return jolty{
		nums:        nums,
		done:        make(map[int]bool, len(nums)),
		deviceJolts: deviceJoltage(nums),
	}
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
		diff := j.nums[p] - need

		j.done[p] = true
		j.sequence = append(j.sequence, p)
		j.doChain(need + diff)
		j.sequence = j.sequence[:len(j.sequence)-1]
		j.done[p] = false
	}
}

func (j *jolty) visitedAll() bool {
	return len(j.sequence) == len(j.nums)
}

func (j *jolty) findPossible(need int) (ret []int) {
	for i := 0; i < len(j.nums); i++ {
		if j.done[i] {
			continue
		}
		if j.nums[i] == need+1 ||
			j.nums[i] == need+2 ||
			j.nums[i] == need+3 {
			ret = append(ret, i)
		}
	}
	return ret
}

func (j *jolty) printSummary() {
	jolts := 0
	numOnes := 0
	numThrees := 0
	for _, s := range j.sequence {
		val := j.nums[s]
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
