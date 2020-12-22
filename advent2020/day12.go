package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

var (
	testFileFlag = flag.Bool("t", false, "Use the test file")
)

type ship struct {
	angle int     // 0=east, 90=north
	x, y  float64 // +y is north, +x=east
}

type action struct {
	action byte
	dist   int
}

func (s ship) String() string {
	return fmt.Sprintf("(%g, %g) %d deg", s.x, s.y, s.angle)
}

func (s ship) manhattanDist() float64 {
	return math.Abs(s.x) + math.Abs(s.y)
}

func (s *ship) move(a action) {
	switch a.action {
	case 'N':
		s.y += float64(a.dist)
	case 'S':
		s.y -= float64(a.dist)
	case 'E':
		s.x += float64(a.dist)
	case 'W':
		s.x -= float64(a.dist)
	case 'L':
		s.angle += a.dist
	case 'R':
		s.angle -= a.dist
	case 'F':
		s.x += float64(a.dist) * math.Cos(degToRads(s.angle))
		s.y += float64(a.dist) * math.Sin(degToRads(s.angle))
	}
}

func (a action) String() string {
	return fmt.Sprintf("%c%d", a.action, a.dist)
}

func degToRads(deg int) float64 {
	return (float64(deg) * math.Pi) / 180
}

var actionRx = regexp.MustCompile(`([NSEWLRF])(\d+)`)

func read(fname string) ([]action, error) {
	ret := []action{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := actionRx.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("unable to parse %q\n", line)
		}
		num, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("unable to parse number %q\n", parts[2])
		}
		ret = append(ret, action{byte(parts[1][0]), num})
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func main() {
	flag.Parse()

	fmt.Printf("Day 12\n")

	infile := "day12.input"
	if *testFileFlag {
		infile = "day12-sample.input"
	}
	actions, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Num actions %d\n", len(actions))
	s := ship{}
	for _, action := range actions {
		s.move(action)
		fmt.Printf("%s %s\n", action, s)
	}
	fmt.Printf("%g, %g = %g\n", s.x, s.y, s.manhattanDist())
}
