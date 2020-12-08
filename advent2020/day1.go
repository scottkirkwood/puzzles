package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func read(fname string) ([]int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]int, 0, 200)
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return ret, fmt.Errorf("line %d: %v", line, err)
		}
		ret = append(ret, num)
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func bruteForce(nums []int) {
	fmt.Printf("Brute forcing\n")
	for a := 0; a < len(nums)-1; a++ {
		for b := a + 1; b < len(nums); b++ {
			x, y := nums[a], nums[b]
			if x+y == 2020 {
				fmt.Printf("%d + %d = 200, %d * %d = %d\n", x, y, x, y, x*y)
				return
			}
		}
	}
}

func main() {
	fmt.Printf("Day1\n")
	nums, err := read("day1.input")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("%d values\n", len(nums))
	bruteForce(nums)
}
