package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

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

func valid(nums []int, prev int) int {
	for i := prev; i < len(nums); i++ {
		if !hasSum(nums[i], nums[i-prev:i]) {
			fmt.Printf("Failed at finding sum for %d\n", nums[i])
			return nums[i]
		}
	}
	return 0
}

func hasSum(sum int, nums []int) bool {
	sorted := make([]int, len(nums))
	copy(sorted, nums)
	sort.Ints(sorted)
	low := 0
	high := len(sorted) - 1
	for low < high {
		curSum := sorted[low] + sorted[high]
		if curSum == sum {
			// fmt.Printf("%d + %d = %d\n", sorted[low], sorted[high], sum)
			return true
		}
		if curSum > sum {
			high--
		} else {
			low++
		}
	}
	return false
}

func findContiguousSum(invalidSum int, nums []int) (sum int) {
	low := 0
	high := low + 1
	for low < high {
		curRange := nums[low : high+1]
		curSum := calcSum(curRange)
		if curSum == invalidSum {
			sorted := make([]int, len(curRange))
			copy(sorted, curRange)
			sort.Ints(sorted)
			sum = sorted[0] + sorted[len(sorted)-1]
			fmt.Printf("%d + %d = %d\n", sorted[0], sorted[len(sorted)-1], sum)
			low++
			high++
		} else if curSum < invalidSum {
			high++
		} else {
			low++
		}
		if low >= high {
			high = low + 1
		}
		if high >= len(nums) {
			fmt.Print("Not found\n")
			return sum
		}
	}
	return sum
}

func calcSum(nums []int) int {
	sum := 0
	for _, x := range nums {
		sum += x
	}
	return sum
}

func main() {
	flag.Parse()

	fmt.Printf("Day 9\n")

	infile := "day9.input"
	prev := 25
	if *testFileFlag {
		prev = 5
		infile = "day9-sample.input"
	}
	nums, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(nums))

	invalidSum := valid(nums, prev)
	fmt.Printf("Sum to find %d\n", invalidSum)
	sum := findContiguousSum(invalidSum, nums)
	fmt.Printf("Sum %d\n", sum)
}
