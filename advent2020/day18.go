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

func read(fname string) ([]string, error) {
	ret := []string{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ret = append(ret, line)
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

var exprRx = regexp.MustCompile(`(\d+|[+*]|[()])\s*`)

func recurseParts(expr string) int {
	if !strings.Contains(expr, "(") {
		return evalLtoR(expr)
	}
	left, parens, right := outerParens(expr)
	val := recurseParts(parens)
	return recurseParts(fmt.Sprintf("%s %d %s", left, val, right))
}

func recurseParts2(expr string) int {
	if !strings.Contains(expr, "(") {
		return evalLtoR2(expr)
	}
	left, parens, right := outerParens(expr)
	val := recurseParts2(parens)
	return recurseParts2(fmt.Sprintf("%s %d %s", left, val, right))
}

func outerParens(expr string) (left, middle, right string) {
	l, r := outerParensIndexes(expr)
	return expr[:l], expr[l+1 : r], expr[r+1:]
}

func outerParensIndexes(expr string) (left, right int) {
	depth := 0
	for i, tok := range expr {
		switch tok {
		case '(':
			depth++
			if depth == 1 {
				left = i
			}
		case ')':
			depth--
			if depth == 0 {
				return left, i
			}
		}
	}
	return left, right
}

func evalLtoR(expr string) int {
	val := 0
	op := ""
	for _, part := range tokenize(expr) {
		if part == "*" || part == "+" {
			op = part
		} else {
			num := parseNum(part)
			switch op {
			case "":
				val = num
			case "*":
				val *= num
			case "+":
				val += num
			}
		}
	}
	return val
}

// This one + is higher precedence than *
func evalLtoR2(expr string) int {
	if !strings.Contains(expr, "*") {
		return evalLtoR(expr)
	}
	parts := strings.SplitN(expr, "*", 2)
	return evalLtoR2(parts[0]) * evalLtoR2(parts[1])
}

func tokenize(expr string) []string {
	parts := exprRx.FindAllStringSubmatch(expr, -1)
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		ret = append(ret, part[1])
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

	fmt.Printf("Day 18\n")

	infile := "day18.input"
	if *testFileFlag {
		infile = "day18-sample.input"
	}
	lines, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	sum := 0
	for i, line := range lines {
		val := recurseParts2(line)
		fmt.Printf("%d: %d = %q\n", i, val, line)
		sum += val
	}
	fmt.Printf("Sum %d\n", sum)
}
