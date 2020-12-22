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

type problem struct {
	rules    map[int]subMatch
	messages []string
	rx       *regexp.Regexp
}

type subMatch struct {
	char   string
	match1 []int
	match2 []int
}

var (
	rowRx    = regexp.MustCompile(`(\d+): (.+)`)
	charRx   = regexp.MustCompile(`"([a-z])"`)
	optionRx = regexp.MustCompile(`(.+) \| (.+)`)
)

func (p problem) match(message string) bool {
	return p.rx.MatchString(message)
}

func (p *problem) setRegex(regex string) (err error) {
	p.rx, err = regexp.Compile("^" + regex + "$")
	return err
}

func (p *problem) setRule(key int, s subMatch) {
	p.rules[key] = s
}

func (p problem) recurse(rule int) string {
	if p.rules[rule].char != "" {
		return p.rules[rule].char
	}
	if len(p.rules[rule].match2) > 0 {
		return "(" + p.recurseAll(p.rules[rule].match1) + "|" + p.recurseAll(p.rules[rule].match2) + ")"
	}
	return p.recurseAll(p.rules[rule].match1)
}

func (p problem) recurseAll(rules []int) string {
	parts := []string{}
	for _, rule := range rules {
		parts = append(parts, p.recurse(rule))
	}
	return strings.Join(parts, "")
}

func (p *problem) addRule(line string) error {
	parts := rowRx.FindStringSubmatch(line)
	if len(parts) != 3 {
		return fmt.Errorf("bad rule %q\n", line)
	}
	key, remainder := parseNum(parts[1]), parts[2]
	if _, ok := p.rules[key]; ok {
		return fmt.Errorf("two rules with same key %d\n", key)
	}

	parts = charRx.FindStringSubmatch(remainder)
	if len(parts) == 2 {
		p.rules[key] = subMatch{char: parts[1]}
		return nil
	}
	parts = optionRx.FindStringSubmatch(remainder)
	if len(parts) == 3 {
		p.rules[key] = subMatch{
			match1: parseNums(parts[1]),
			match2: parseNums(parts[2]),
		}
		return nil
	}
	p.rules[key] = subMatch{match1: parseNums(remainder)}
	return nil
}

func (p *problem) addMessage(line string) {
	p.messages = append(p.messages, line)
}

func read(fname string) (problem, error) {
	ret := problem{rules: map[int]subMatch{}}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			part = 1
			continue
		}
		if part == 0 {
			if err := ret.addRule(line); err != nil {
				return ret, err
			}
		} else {
			ret.addMessage(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func parseNums(txt string) []int {
	parts := strings.Split(txt, " ")
	ret := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		ret = append(ret, parseNum(part))
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

func main() {
	flag.Parse()

	fmt.Printf("Day 19\n")

	infile := "day19.input"
	if *testFileFlag {
		infile = "day19-sample.input"
	}
	p, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Problem rules %d\n", len(p.rules))
	p.setRule(8, subMatch{match1: []int{42}, match2: []int{42, 8}})
	p.setRule(11, subMatch{match1: []int{42, 31}, match2: []int{42, 11, 31}})
	regexStr := p.recurse(0)
	fmt.Printf("Calculated regex %q\n", regexStr)
	if err := p.setRegex(regexStr); err != nil {
		fmt.Printf("Problem compiling %q: %v\n", regexStr, err)
	}
	countMatches := 0
	for _, message := range p.messages {
		if !p.match(message) {
			fmt.Printf("Problem %q does not match\n", message)
		} else {
			countMatches++
		}
	}
	fmt.Printf("%d messages matches\n", countMatches)
}
