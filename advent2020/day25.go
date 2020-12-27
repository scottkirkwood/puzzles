// See https://adventofcode.com/2020/day/25 for problem decscription
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var testFileFlag = flag.Bool("t", false, "Use the test file")

var directionRx = regexp.MustCompile(`(e|se|sw|w|nw|ne)`)

type puzzle struct {
	value         int
	subjectNumber int
	door          crypt
	card          crypt
}

type crypt struct {
	secretLoopSize int
	publicKey      int
	encryptionKey  int
}

func (c *crypt) findLoopSize() int {
	sn := 7
	val := sn
	for loopSize := 1; loopSize < 9999999999; loopSize++ {
		val *= sn
		val %= divider
		if val == c.publicKey {
			c.secretLoopSize = loopSize + 1
			break
		}
		if loopSize%10000 == 0 {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
	return c.secretLoopSize
}

func (c *crypt) calcEncryptionKey(subjectNumber int) int {
	sn := subjectNumber
	val := sn
	for i := 1; i < c.secretLoopSize; i++ {
		val *= sn
		val %= divider
	}
	c.encryptionKey = val
	return c.encryptionKey
}

const divider = 20201227

func read(fname string) (puzzle, error) {
	ret := puzzle{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rowNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if rowNum == 0 {
			ret.card.publicKey = parseNum(line)
		} else {
			ret.door.publicKey = parseNum(line)
		}
		rowNum++
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

func findCalc(wg *sync.WaitGroup, name string, c crypt, other crypt) {
	fmt.Printf("Loop size %s %d\n", name, c.findLoopSize())
	fmt.Printf("Calc %s key %d\n", name, c.calcEncryptionKey(other.publicKey))
	wg.Done()
}

func main() {
	flag.Parse()

	fmt.Printf("Day 25\n")

	infile := "day25.input"
	if *testFileFlag {
		infile = "day25-sample.input"
	}
	puzzle, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Card public key %d\n", puzzle.card.publicKey)
	fmt.Printf("Door public key %d\n", puzzle.door.publicKey)
	var wg sync.WaitGroup

	wg.Add(1)
	go findCalc(&wg, "card", puzzle.card, puzzle.door)
	wg.Add(1)
	go findCalc(&wg, "door", puzzle.door, puzzle.card)
	wg.Wait()
}
