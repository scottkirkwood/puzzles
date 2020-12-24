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

type game struct {
	moveNum int
	pos     int
	numbers []int
}

func (g game) String() string {
	lines := []string{}
	lines = append(lines, fmt.Sprintf("-- move %d --", g.moveNum))
	lines = append(lines, fmt.Sprintf("cups: %s", g.formatCups()))
	lines = append(lines, fmt.Sprintf("pick up: %s", joinInts(g.pickUp(), ", ")))
	lines = append(lines, fmt.Sprintf("destination: %d", g.destination()))
	return strings.Join(lines, "\n")
}

func (g *game) move() {
	prevVal := g.numbers[g.pos]
	dest := g.destination()
	pickup := g.pickUp()
	newList := g.removeThese(pickup)
	destIndex := indexOf(newList, dest)
	g.numbers = insertAfter(newList, destIndex, pickup)

	g.pos = indexOf(g.numbers, prevVal)
	g.pos = g.inc(g.pos)
	g.moveNum++
}

func (g game) removeThese(pickup []int) []int {
	m := intToMap(pickup)
	ret := make([]int, 0, len(g.numbers)-3)
	for i := 0; i < len(g.numbers); i++ {
		if m[g.numbers[i]] {
			continue
		}
		ret = append(ret, g.numbers[i])
	}
	return ret
}

func insertAfter(lst []int, after int, toInsert []int) []int {
	ret := make([]int, 0)
	for i := 0; i <= after; i++ {
		ret = append(ret, lst[i])
	}
	for _, x := range toInsert {
		ret = append(ret, x)
	}
	for i := after + 1; i < len(lst); i++ {
		ret = append(ret, lst[i])
	}
	return ret
}

func intToMap(nums []int) map[int]bool {
	m := make(map[int]bool, len(nums))
	for _, num := range nums {
		m[num] = true
	}
	return m
}

func (g game) inc(index int) int {
	return (index + 1) % len(g.numbers)
}

func (g game) afterOne() string {
	ret := []string{}
	index := indexOf(g.numbers, 1)
	for {
		index = g.inc(index)
		if g.numbers[index] == 1 {
			break
		}
		ret = append(ret, fmt.Sprintf("%d", g.numbers[index]))
	}
	return strings.Join(ret, "")
}

func indexOf(numbers []int, num int) int {
	for i, n := range numbers {
		if n == num {
			return i
		}
	}
	return 0
}

func (g game) formatCups() string {
	ret := []string{}
	for i, num := range g.numbers {
		if i == g.pos {
			ret = append(ret, fmt.Sprintf("(%d)", num))
		} else {
			ret = append(ret, fmt.Sprintf("%d ", num))
		}
	}
	return strings.Join(ret, " ")
}

func (g game) pickUp() []int {
	ret := make([]int, 3)
	p := g.pos
	for i := 0; i < 3; i++ {
		p = g.inc(p)
		ret[i] = g.numbers[p]
	}
	return ret
}

func (g game) destination() int {
	m := make(map[int]bool, 3)
	for _, n := range g.pickUp() {
		m[n] = true
	}
	num := g.numbers[g.pos]
	for {
		num--
		if num <= 0 {
			num = len(g.numbers)
		}
		if !m[num] {
			break
		}
	}
	return num
}

func joinInts(nums []int, joiner string) string {
	ret := make([]string, len(nums))
	for i, n := range nums {
		ret[i] = strconv.FormatInt(int64(n), 10)
	}
	return strings.Join(ret, joiner)
}

var playerRx = regexp.MustCompile(`Player (\d):`)

func read(fname string) (game, error) {
	ret := game{moveNum: 1}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		for _, ch := range line {
			ret.numbers = append(ret.numbers, parseNum(string(ch)))
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("Unable to parse decimal %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func main() {
	flag.Parse()

	fmt.Printf("Day 23\n")

	infile := "day23.input"
	if *testFileFlag {
		infile = "day23-sample.input"
	}
	game, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	for i := 0; i < 100; i++ {
		fmt.Printf("%s\n", game.String())
		game.move()
	}
	fmt.Printf("%s\n", game.afterOne())
}
