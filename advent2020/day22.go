// See https://adventofcode.com/2020/day/22 for problem decscription
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
	player1 []int
	player2 []int
}

func (g game) String() string {
	lines := []string{}
	if len(g.player1) == 0 || len(g.player2) == 0 {
		lines = append(lines, "=== Post game results ====")
		lines = append(lines, fmt.Sprintf("Player 1's deck: %s", joinInt(g.player1, ", ")))
		lines = append(lines, fmt.Sprintf("Player 2's deck: %s", joinInt(g.player2, ", ")))
	} else {
		lines = append(lines, fmt.Sprintf("Player 1's deck: %s", joinInt(g.player1, ", ")))
		lines = append(lines, fmt.Sprintf("Player 2's deck: %s", joinInt(g.player2, ", ")))
		lines = append(lines, fmt.Sprintf("Player 1 plays: %d", g.player1[0]))
		lines = append(lines, fmt.Sprintf("Player 2 plays: %d", g.player2[0]))
		player := 1
		if g.player1[0] < g.player2[0] {
			player = 2
		}
		lines = append(lines, fmt.Sprintf("Player %d wins the round!", player))
	}
	return strings.Join(lines, "\n")
}

func (g game) score() int {
	if len(g.player1) > 0 {
		return score(g.player1)
	}
	return score(g.player2)
}

func score(cards []int) int {
	score := 0
	count := len(cards)
	for i, v := range cards {
		score += (count - i) * v
	}
	return score
}

func (g *game) round() bool {
	if len(g.player1) == 0 || len(g.player2) == 0 {
		return false
	}
	if g.player1[0] > g.player2[0] {
		g.player1 = append(g.player1, g.player1[0])
		g.player1 = append(g.player1, g.player2[0])
	} else {
		g.player2 = append(g.player2, g.player2[0])
		g.player2 = append(g.player2, g.player1[0])
	}
	g.player1 = g.player1[1:]
	g.player2 = g.player2[1:]
	return true
}

var playerRx = regexp.MustCompile(`Player (\d):`)

func joinInt(values []int, joiner string) string {
	strList := make([]string, len(values))
	for i, value := range values {
		strList[i] = strconv.FormatInt(int64(value), 10)
	}
	return strings.Join(strList, joiner)
}

func read(fname string) (game, error) {
	ret := game{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	player := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := playerRx.FindStringSubmatch(line)
		if len(parts) == 2 {
			player++
			continue
		}
		num := parseNum(line)
		if player == 1 {
			ret.player1 = append(ret.player1, num)
		} else {
			ret.player2 = append(ret.player2, num)
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

	fmt.Printf("Day 22\n")

	infile := "day22.input"
	if *testFileFlag {
		infile = "day22-sample.input"
	}
	game, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	for round := 1; ; round++ {
		fmt.Printf("---- Round %d ---\n", round)
		fmt.Printf("%s\n", game)
		if !game.round() {
			break
		}
	}
	fmt.Printf("Score %d\n", game.score())
}
