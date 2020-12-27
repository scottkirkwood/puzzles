// See https://adventofcode.com/2020/day/23 for problem decscription
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

// In order to do part 2 this needs to be a (singly) linked list
// Each node has two pointers one that points to the next element which
// finally points back to the first element.
// The other points to the next lower number.  With these two changes
// the algorithm will be fast
type game struct {
	moveNum    int
	countNodes int   // number of nodes in loop
	start      *node // start of loop
	pos        *node // current position, also used for add()
	max        *node // node with maximum number value
	min        *node // node with value 1
}

type node struct {
	num  int
	next *node
	dec  *node
}

func (g game) String() string {
	lines := []string{}
	lines = append(lines, fmt.Sprintf("-- move %d --", g.moveNum))
	lines = append(lines, fmt.Sprintf("cups: %s", g.formatCups()))
	pickup := g.pickupCopy()
	lines = append(lines, fmt.Sprintf("pick up: %s", joinInts(pickup, ", ")))
	// lines = append(lines, fmt.Sprintf("count down: %s", joinDecInts(g.max, ", ")))
	lines = append(lines, fmt.Sprintf("destination: %d", g.findDestination(g.pos, pickup).num))
	return strings.Join(lines, "\n")
}

func (g *game) move() {
	g.moveNum++
	pickup, last := g.pickUp()
	destination := g.findDestination(g.pos, pickup)
	last.next = destination.next
	destination.next = pickup
	g.pos = g.pos.next
}

func (g *game) init() {
	g.pos = g.start
}

func (g *game) add(num int) {
	newNode := &node{
		num: num,
	}
	if num == 1 {
		g.min = newNode
	}
	g.countNodes++
	if g.start == nil {
		g.start = newNode
		g.pos = g.start
		g.max = newNode
		newNode.next = g.start
		newNode.dec = nil
		return
	}
	newNode.next = g.pos.next
	g.pos.next = newNode
	g.pos = newNode
	if num > g.max.num {
		newNode.dec = g.max
		g.max = newNode
	} else {
		oneUp := g.findOneUp(num)
		newNode.dec = oneUp.dec
		oneUp.dec = newNode
	}
}

func (g game) findOneUp(num int) *node {
	p := g.max
	prev := p
	for {
		if p == nil || p.num < num {
			return prev
		}
		prev = p
		p = p.dec
	}
	return prev
}

func (g game) afterOne(count int) {
	p := g.min.next
	prod := 1
	for i := 0; i < count; i++ {
		fmt.Printf("%d ", p.num)
		prod *= p.num
		p = p.next
	}
	fmt.Printf("\nProduct: %d\n", prod)
}

func (g game) formatCups() string {
	ret := []string{}
	p := g.start
	i := 0
	for {
		if p == g.pos {
			ret = append(ret, fmt.Sprintf("(%d)", p.num))
		} else {
			ret = append(ret, fmt.Sprintf("%d ", p.num))
		}
		if i > 10 {
			ret = append(ret, "...")
			break
		}
		i++
		p = p.next
		if p == g.start || p == nil {
			if p == nil {
				fmt.Printf("Unexpected nil\n")
			}
			break
		}
	}
	return strings.Join(ret, " ")
}

func (g *game) pickupCopy() *node {
	p := g.pos.next

	var start *node
	var last *node
	for i := 0; i < 3; i++ {
		cur := &node{
			num: p.num,
		}
		if start == nil {
			start = cur
		} else {
			last.next = cur
		}
		p = p.next
		last = cur
	}
	return start
}
func (g *game) pickUp() (start, last *node) {
	start = g.pos.next
	last = start
	for i := 0; i < 2; i++ {
		last = last.next
	}
	// remove from loop
	g.pos.next = last.next
	last.next = nil
	return start, last
}

func (g game) findDestination(p *node, notIn *node) *node {
	p = g.decrement(p)
	for notIn.contains(p) {
		p = g.decrement(p)
	}
	return p
}

func (g game) decrement(p *node) *node {
	p = p.dec
	if p == nil {
		return g.max
	}
	return p
}

func (n *node) contains(other *node) bool {
	p := n
	for {
		if p.num == other.num {
			return true
		}
		p = p.next
		if p == nil || p == n {
			return false
		}
	}
	return false
}

func joinInts(start *node, joiner string) string {
	ret := []string{}
	p := start
	for {
		ret = append(ret, strconv.FormatInt(int64(p.num), 10))
		p = p.next
		if p == nil || p == start {
			break
		}
	}
	return strings.Join(ret, joiner)
}

func joinDecInts(start *node, joiner string) string {
	ret := []string{}
	p := start
	for {
		ret = append(ret, strconv.FormatInt(int64(p.num), 10))
		p = p.dec
		if p == nil {
			break
		}
	}
	return strings.Join(ret, joiner)
}

var playerRx = regexp.MustCompile(`Player (\d):`)

func read(fname string, extra int) (game, error) {
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
			ret.add(parseNum(string(ch)))
		}
		maxVal := ret.max.num
		for i := maxVal + 1; i <= extra; i++ {
			ret.add(i)
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
	extra := 1000000
	if *testFileFlag {
		infile = "day23-sample.input"
		//extra = 0
	}
	game, err := read(infile, extra)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	game.init()
	moves := 10000000
	for i := 0; i < moves; i++ {
		//fmt.Printf("%s\n", game.String())
		game.move()
	}
	game.afterOne(2)
}
