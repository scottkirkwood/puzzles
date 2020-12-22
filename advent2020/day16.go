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

type ticket struct {
	names         map[string]ranges
	yourTicket    []int
	nearbyTickets [][]int
	possible      []map[string]bool
	order         []string
}

func (t *ticket) calculateOrder() {
	t.order = make([]string, len(t.yourTicket))
	t.possible = make([]map[string]bool, len(t.yourTicket))
	for i := 0; i < len(t.yourTicket); i++ {
		t.possible[i] = map[string]bool{}
		for key, ranges := range t.names {
			if t.matchRanges(i, ranges) {
				t.possible[i][key] = true
			}
		}
	}
	changed := false
	for {
		changed = false
		for i, m := range t.possible {
			key := theOnlyKey(m)
			if key != "" {
				t.order[i] = key
				t.zeroOut(key)
				changed = true
				break
			}
		}
		if !changed {
			break
		}
	}
	product := 1
	for i, key := range t.order {
		fmt.Printf("%d: %v\n", i, key)
		if strings.HasPrefix(key, "departure ") {
			fmt.Printf("%d\n", t.yourTicket[i])
			product *= t.yourTicket[i]
		}
	}
	fmt.Printf("%d\n", product)
}

func theOnlyKey(m map[string]bool) string {
	key := ""
	for k, v := range m {
		if v {
			if key != "" {
				return ""
			}
			key = k
		}
	}
	return key
}

func (t *ticket) zeroOut(key string) {
	for _, m := range t.possible {
		m[key] = false
	}
}

func (t ticket) matchRanges(index int, r ranges) bool {
	if !r.inRange(t.yourTicket[index]) {
		return false
	}
	for _, t := range t.nearbyTickets {
		if t != nil && !r.inRange(t[index]) {
			return false
		}
	}
	return true
}

func (t ticket) badValues() []int {
	ret := []int{}
	for _, values := range t.nearbyTickets {
		for _, v := range values {
			if !t.inRange(v) {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

func (t *ticket) removeBadTickets() {
	for i, values := range t.nearbyTickets {
		for _, v := range values {
			if !t.inRange(v) {
				t.nearbyTickets[i] = nil
				continue
			}
		}
	}
}

func (t ticket) inRange(x int) bool {
	for _, r := range t.names {
		if r.inRange(x) {
			return true
		}
	}
	return false
}

type ranges struct {
	first  fromTo
	second fromTo
}

func (r ranges) inRange(x int) bool {
	return r.first.inRange(x) || r.second.inRange(x)
}

type fromTo struct {
	from, to int
}

func (f fromTo) inRange(x int) bool {
	return x >= f.from && x <= f.to
}

func newRanges(from1, to1, from2, to2 string) ranges {
	return ranges{
		first:  fromTo{parseNum(from1), parseNum(to1)},
		second: fromTo{parseNum(from2), parseNum(to2)},
	}
}

var (
	nameFromToRx    = regexp.MustCompile(`(.+): (\d+)-(\d+) or (\d+)-(\d+)`)
	yourTicketRx    = regexp.MustCompile(`your ticket:`)
	nearbyTicketsRx = regexp.MustCompile(`nearby tickets:`)
)

func read(fname string) (ticket, error) {
	ret := ticket{names: map[string]ranges{}}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	phase := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		parts := nameFromToRx.FindStringSubmatch(line)
		if len(parts) == 6 {
			ret.names[parts[1]] = newRanges(parts[2], parts[3], parts[4], parts[5])

		}
		if yourTicketRx.MatchString(line) || nearbyTicketsRx.MatchString(line) {
			phase++
			continue
		}
		if phase == 1 {
			ret.yourTicket = parseNums(line)
		} else if phase == 2 {
			ret.nearbyTickets = append(ret.nearbyTickets, parseNums(line))
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func parseNums(txt string) []int {
	parts := strings.Split(txt, ",")
	ret := make([]int, len(parts))
	for i, part := range parts {
		ret[i] = parseNum(part)
	}
	return ret
}

func parseNum(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("Unable to parse decimal %q\n", txt)
		os.Exit(-1)
	}
	return num
}

func sum(vals []int) int {
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return sum
}

func main() {
	flag.Parse()

	fmt.Printf("Day 16\n")

	infile := "day16.input"
	if *testFileFlag {
		infile = "day16-sample.input"
	}
	tickets, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Tickets %d\n", len(tickets.names))
	fmt.Printf("Ticket error scanning rate %d\n", sum(tickets.badValues()))
	tickets.removeBadTickets()
	tickets.calculateOrder()
}
