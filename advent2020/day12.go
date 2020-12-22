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
	x, y int
	wp   *waypoint
}

type waypoint struct {
	x, y int
}

type action struct {
	action byte
	dist   int
}

func (s ship) String() string {
	return fmt.Sprintf("(%d, %d) wp (%d, %d)", s.x, s.y, s.wp.x, s.wp.y)
}

func (s ship) manhattanDist() float64 {
	return math.Abs(float64(s.x)) + math.Abs(float64(s.y))
}

func (s *ship) move(a action) {
	switch a.action {
	case 'N':
		s.wp.y += a.dist
	case 'S':
		s.wp.y -= a.dist
	case 'E':
		s.wp.x += a.dist
	case 'W':
		s.wp.x -= a.dist
	case 'L':
		s.wp.rotate(a.dist)
	case 'R':
		s.wp.rotate(-a.dist)
	case 'F':
		s.x += a.dist * s.wp.x
		s.y += a.dist * s.wp.y
	}
}

func (w *waypoint) rotate(angle int) {
	switch angle {
	case -90, 270:
		w.x, w.y = w.y, -w.x
	case 180, -180:
		w.x, w.y = -w.x, -w.y
	case -270, 90:
		w.x, w.y = -w.y, w.x
	case 360, -360:
		w.x, w.y = w.x, w.y
	default:
		fmt.Printf("Bad angle %d\n", angle)
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
	s := ship{wp: &waypoint{x: 10, y: 1}}
	for _, action := range actions {
		s.move(action)
		fmt.Printf("%s %s\n", action, s)
	}
	fmt.Printf("%d, %d = %g\n", s.x, s.y, s.manhattanDist())
}
