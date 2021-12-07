package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

type board [5][5]int

type bingo struct {
	numbers []int
	boards  []board
}

func (b board) Print(selected map[int]bool) {
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			val := b[row][col]
			if selected[val] {
				fmt.Printf("[%2d]", val)
			} else {
				fmt.Printf(" %2d ", val)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (b board) winner(selected map[int]bool) bool {
	for row := 0; row < 5; row++ {
		count := 0
		for col := 0; col < 5; col++ {
			if selected[b[row][col]] {
				count++
			}
		}
		if count == 5 {
			return true
		}
	}
	for col := 0; col < 5; col++ {
		count := 0
		for row := 0; row < 5; row++ {
			if selected[b[row][col]] {
				count++
			}
		}
		if count == 5 {
			return true
		}
	}
	return false
}

func (b board) sumRemaining(selected map[int]bool) int {
	sum := 0
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if !selected[b[row][col]] {
				sum += b[row][col]
			}
		}
	}
	return sum
}

func toIntMap(numbers []int) map[int]bool {
	ret := make(map[int]bool, len(numbers))
	for _, n := range numbers {
		ret[n] = true
	}
	return ret
}

func read(fname string) (*bingo, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ret := bingo{
		numbers: []int{},
		boards:  []board{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ",") {
			ret.numbers = parseNums(strings.Split(line, ","))
		} else if len(line) == 14 {
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "  ", " ")
			vals := parseNums(strings.Split(line, " "))
			if row == 0 {
				ret.boards = append(ret.boards, board{})
			}
			lastIndex := len(ret.boards) - 1
			for col := 0; col < 5; col++ {
				ret.boards[lastIndex][row][col] = vals[col]
			}
			row++
			if row > 4 {
				row = 0
			}
		}
	}
	return &ret, nil
}

func (b *bingo) Len() int {
	return len(b.boards)
}

func (b *bingo) calculateAt(i int) int {
	slice := b.numbers[0:i]
	selected := toIntMap(slice)
	for _, board := range b.boards {
		if !board.winner(selected) {
			board.Print(selected)
			boardSum := board.sumRemaining(selected)
			lastValue := slice[len(slice)-1]
			fmt.Printf("Board sum %d, lastValue %d\n", boardSum, lastValue)
			return lastValue * boardSum
		}
	}
	return 0
}

func (b *bingo) calculate() int {
	winningBoards := map[int]bool{}
	lastNumber := 0
	lastSelected := map[int]bool{}
	boardRemaining := board{}
	for i := 1; i <= len(b.numbers); i++ {
		slice := b.numbers[0:i]
		selected := toIntMap(slice)
		for i, board := range b.boards {
			if board.winner(selected) {
				winningBoards[i] = true
				if len(winningBoards) == len(b.boards)-1 {
					lastNumber = slice[len(selected)-1]
					lastSelected = selected
					fmt.Printf("Last number %d\n", lastNumber)
					for j, board := range b.boards {
						if !winningBoards[j] {
							boardRemaining = board
						}
					}
				}
			}
		}
	}
	sum := boardRemaining.sumRemaining(lastSelected)
	fmt.Printf("Sum %d\n", sum)
	return sum * lastNumber
}

func parseNums(texts []string) []int {
	ret := make([]int, 0, len(texts))
	for _, txt := range texts {
		ret = append(ret, parseNum(txt))
	}
	return ret
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(strings.TrimSpace(txt))
	if err != nil {
		fmt.Printf("=========== Bad number %q ==============\n", txt)
	}
	return int(num)
}

func main() {
	flag.Parse()

	fmt.Printf("Day 04b\n")
	infile := "day04.input"
	if *testFileFlag {
		infile = "day04-sample.input"
	}
	b, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Lines %d\n", b.Len())
	fmt.Printf("Answer %d\n", b.calculate())
}
