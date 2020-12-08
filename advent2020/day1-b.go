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
	for a := 0; a < len(nums)-2; a++ {
		for b := a + 1; b < len(nums)-1; b++ {
			for c := b + 1; c < len(nums); c++ {
				x, y, z := nums[a], nums[b], nums[c]
				if x+y+z == 2020 {
					fmt.Printf("%d+%d+%d=2020, %d*%d*%d = %d\n", x, y, z, x, y, z, x*y*z)
					return
				}
			}
		}
	}
}

func main() {
	fmt.Printf("Day2\n")
	nums, err := read("day1.input")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("%d values\n", len(nums))
	bruteForce(nums)
}
