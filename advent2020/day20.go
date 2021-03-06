// See https://adventofcode.com/2020/day/20 for problem decscription
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

const (
	width  = 10
	height = 10
)

const (
	deg0Right = iota
	deg0Bottom
	deg0Left
	deg0Top
	flip0Right
	flip0Bottom
	flip0Left
	flip0Top
	sigCount
)

type image struct {
	tiles []tile
}

type tile struct {
	id        int
	pixels    [][]byte
	signature [sigCount]string
}

func (im *image) addTile(num int) *tile {
	im.tiles = append(im.tiles, newTile(num))
	return &im.tiles[len(im.tiles)-1]
}

func (im *image) calcSignatures() {
	for i := 0; i < len(im.tiles); i++ {
		im.tiles[i].calcSignatures()
	}
}

func (im image) findCorners() {
	counts := map[string]int{}
	for _, tile := range im.tiles {
		for _, sig := range tile.signature {
			counts[sig]++
		}
	}
	m := map[int][]string{}
	for key, val := range counts {
		if val == 1 {
			id, pos := im.findSig(key)
			m[id] = append(m[id], posToTxt(pos))
		}
	}
	product := 1
	for id, val := range m {
		if len(val) == 4 {
			fmt.Printf("%d, %s\n", id, strings.Join(val, ", "))
			product *= id
		}
	}
	fmt.Printf("Product %d\n", product)
}

func (im image) findSig(sig string) (id, pos int) {
	for _, tile := range im.tiles {
		for i, curSig := range tile.signature {
			if curSig == sig {
				return tile.id, i
			}
		}
	}
	return 0, 0
}

func newTile(id int) tile {
	return tile{
		id:     id,
		pixels: make([][]byte, 0, height),
	}
}

// String returns a string representation of the tile
func (t tile) String() string {
	lines := []string{}
	lines = append(lines, fmt.Sprintf("Tile %d:", t.id))
	for _, row := range t.pixels {
		lines = append(lines, "")
		for _, col := range row {
			lines[len(lines)-1] += string(col)
		}
	}
	return strings.Join(lines, "\n")
}

func (t *tile) addRow(txt string) {
	t.pixels = append(t.pixels, make([]byte, width))
	row := len(t.pixels) - 1
	for i, ch := range txt {
		t.pixels[row][i] = byte(ch)
	}
}

func (t *tile) calcSignatures() {
	t.signature[deg0Top] = t.fromTo(0, 0, 1, 0)
	t.signature[deg0Right] = t.fromTo(1, 0, 0, 1)
	t.signature[deg0Bottom] = t.fromTo(1, 1, -1, 0)
	t.signature[deg0Left] = t.fromTo(0, 1, 0, -1)

	t.signature[flip0Top] = t.fromTo(1, 0, -1, 0)
	t.signature[flip0Right] = t.fromTo(1, 1, 0, -1)
	t.signature[flip0Bottom] = t.fromTo(0, 1, 1, 0)
	t.signature[flip0Left] = t.fromTo(0, 0, 0, 1)
}

func formatBits(val int) string {
	return fmt.Sprintf("%010s (%d)", strconv.FormatInt(int64(val), 2), val)
}

func sig(txt []byte) int {
	val := 0
	for i := 0; i < len(txt); i++ {
		val <<= 1
		if txt[i] == '#' {
			val |= 1
		} else if txt[i] != '.' {
			fmt.Printf("Unexpected string %q\n", txt)
		}
	}
	return val
}

func (t tile) fromTo(x, y, dx, dy int) string {
	x *= (width - 1)
	y *= (height - 1)
	txt := make([]byte, 0, width)
	for !t.isOutside(x, y) {
		txt = append(txt, t.get(x, y))
		x += dx
		y += dy
	}
	if len(txt) == 0 {
		fmt.Printf("Wrong parameters to fromTo(%d, %d, %d, %d)\n", x, y, dx, dy)
	}
	return string(txt)
}

func (t tile) get(x, y int) byte {
	return t.pixels[y][x]
}

func (t tile) isOutside(x, y int) bool {
	return x < 0 || y < 0 || x >= width || y >= height
}

var tileRx = regexp.MustCompile(`Tile (\d+):`)

func read(fname string) (image, error) {
	ret := image{}
	file, err := os.Open(fname)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lastTile *tile
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := tileRx.FindStringSubmatch(line)
		if len(parts) == 2 {
			id := parseNum(parts[1])
			lastTile = ret.addTile(id)
		} else {
			lastTile.addRow(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return ret, err
	}
	return ret, nil
}

func posToTxt(pos int) string {
	switch pos {
	case deg0Right:
		return "deg0Right"
	case deg0Bottom:
		return "deg0Bottom"
	case deg0Left:
		return "deg0Left"
	case deg0Top:
		return "deg0Top"
	case flip0Right:
		return "flip0Right"
	case flip0Bottom:
		return "flip0Bottom"
	case flip0Left:
		return "flip0Left"
	case flip0Top:
		return "flip0Top"
	}
	return "bad pos"
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

	fmt.Printf("Day 20\n")

	infile := "day20.input"
	if *testFileFlag {
		infile = "day20-sample.input"
	}
	im, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Num tiles %d\n", len(im.tiles))
	im.calcSignatures()
	im.findCorners()
}
