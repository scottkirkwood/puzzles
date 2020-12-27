// See https://adventofcode.com/2020/day/2 for problem decscription
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type pwd struct {
	min, max int
	char     string
	password string
}

// ex. "4-5 l: rllllj"
var lineRx = regexp.MustCompile(`(\d+)-(\d+) ([a-z]): ([a-z]+)`)

func read(fname string) ([]pwd, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]pwd, 0, 1000)
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		parts := lineRx.FindStringSubmatch(scanner.Text())
		if len(parts) != 5 {
			return ret, fmt.Errorf("invalid line %d: %q", line, scanner.Text())
		}
		min, err := strconv.Atoi(parts[1])
		if err != nil {
			return ret, fmt.Errorf("invalid min %q, line %d; %v", parts[1], line, err)
		}
		max, err := strconv.Atoi(parts[2])
		if err != nil {
			return ret, fmt.Errorf("invalid max %q, line %d: %v", parts[2], line, err)
		}
		ret = append(ret, pwd{
			min:      min,
			max:      max,
			char:     parts[3],
			password: parts[4],
		})
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func numCorrect1(pwds []pwd) int {
	count := 0
	for _, pwd := range pwds {
		if pwd.Valid1() {
			count++
		}
	}
	return count
}

func numCorrect2(pwds []pwd) int {
	count := 0
	for _, pwd := range pwds {
		if pwd.Valid2() {
			count++
		}
	}
	return count
}

func (p *pwd) Valid1() bool {
	count := 0
	for _, ch := range p.password {
		if p.char == string(ch) {
			count++
		}
	}
	return count >= p.min && count <= p.max
}

func (p *pwd) Valid2() bool {
	x, y := p.min-1, p.max-1
	count := 0
	if p.char == p.password[x:x+1] {
		count++
	}
	if p.char == p.password[y:y+1] {
		count++
	}
	return count == 1
}

func main() {
	fmt.Printf("Day2\n")
	//pwds, err := read("day2-sample.input")
	pwds, err := read("day2.input")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", len(pwds))
	fmt.Printf("Num correct1 %d\n", numCorrect1(pwds))
	fmt.Printf("Num correct2 %d\n", numCorrect2(pwds))
}
