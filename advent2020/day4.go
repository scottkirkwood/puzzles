// See https://adventofcode.com/2020/day/4 for problem decscription
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

type passport struct {
	fields map[string]string
}

func read(fname string) ([]passport, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]passport, 0, 1069)
	scanner := bufio.NewScanner(file)

	prev := passport{fields: map[string]string{}}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			ret = append(ret, prev)
			prev = passport{fields: map[string]string{}}
			continue
		}
		parts := strings.Split(line, " ")
		for _, part := range parts {
			colon := strings.Split(part, ":")
			if len(colon) > 1 {
				prev.fields[colon[0]] = colon[1]
			}
		}
	}
	ret = append(ret, prev)
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

var fields = []string{
	"byr", // (Birth Year)
	"iyr", // (Issue Year)
	"eyr", // (Expiration Year)
	"hgt", // (Height)
	"hcl", // (Hair Color)
	"ecl", // (Eye Color)
	"pid", // (Passport ID)
}

func (p passport) AllFields() bool {
	allFields := map[string]bool{}
	for _, field := range fields {
		allFields[field] = true
	}
	for _, field := range fields {
		if _, ok := p.fields[field]; !ok {
			return false
		}
		allFields[field] = false
	}
	for _, field := range fields {
		if allFields[field] {
			return false
		}
	}
	return true
}

func (p passport) Valid() bool {
	if !p.AllFields() {
		return false
	}
	for _, field := range fields {
		val := p.fields[field]
		switch field {
		case "byr": // (Birth Year)
			if !validIntRange(val, 1920, 2002) {
				verbose("byr", val)
				return false
			}
		case "iyr": // (Issue Year)
			if !validIntRange(val, 2010, 2020) {
				verbose("iyr", val)
				return false
			}
		case "eyr": // (Expiration Year)
			if !validIntRange(val, 2020, 2030) {
				verbose("eyr", val)
				return false
			}
		case "hgt": // (Height)
			if !validHeight(val) {
				verbose("hgt", val)
				return false
			}
		case "hcl": // (Hair Color)
			if !validHairColor(val) {
				verbose("hcl", val)
				return false
			}
		case "ecl": // (Eye Color)
			if !validEyeColor(val) {
				verbose("ecl", val)
				return false
			}
		case "pid": // (Passport ID)
			if !validPid(val) {
				verbose("pid", val)
				return false
			}
		default:
			fmt.Printf("%q = %q\n", field, val)
		}
	}
	return true
}

func verbose(context, val string) {
	//fmt.Printf("Invalid %s %q\n", context, val)
}

func validIntRange(txt string, low, high int) bool {
	num, err := strconv.Atoi(txt)
	if err != nil {
		fmt.Printf("Unexpected number %q\n", txt)
		return false
	}
	return num >= low && num <= high
}

var heightRx = regexp.MustCompile(`^(\d+)(cm|in)$`)

func validHeight(txt string) bool {
	parts := heightRx.FindStringSubmatch(txt)
	if len(parts) != 3 {
		return false
	}
	switch parts[2] {
	case "cm":
		return validIntRange(parts[1], 150, 193)
	case "in":
		return validIntRange(parts[1], 59, 76)
	default:
		fmt.Printf("Unexpected height %q\n", parts[2])
	}
	return false
}

var hairColorRx = regexp.MustCompile(`^#[0-9a-f]{6}$`)

func validHairColor(txt string) bool {
	return hairColorRx.MatchString(txt)
}

var eyeColors = map[string]bool{
	"amb": true,
	"blu": true,
	"brn": true,
	"gry": true,
	"grn": true,
	"hzl": true,
	"oth": true,
}

func validEyeColor(txt string) bool {
	return eyeColors[txt]
}

var validPidRx = regexp.MustCompile(`^[0-9]{9}$`)

func validPid(txt string) bool {
	return validPidRx.MatchString(txt)
}

func checkInvalidInputs() {
	allInvalid, err := read("day4-invalid.input")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	for _, passport := range allInvalid {
		if passport.Valid() {
			fmt.Printf("Expecting invalid %v\n", passport)
		}
	}
}

func checkValidInputs() {
	allValid, err := read("day4-valid.input")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	for _, passport := range allValid {
		if !passport.Valid() {
			fmt.Printf("Expecting valid %v\n", passport)
		}
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Day 4\n")

	checkInvalidInputs()
	checkValidInputs()

	infile := "day4.input"
	if *testFileFlag {
		infile = "day4-sample.input"
	}
	passports, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Rows %d\n", len(passports))
	valid := 0
	for _, passport := range passports {
		if passport.Valid() {
			valid++
		}
	}
	fmt.Printf("Valid %d\n", valid)
	fmt.Printf("It's not 148\n")
}
