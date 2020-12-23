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

var rowRx = regexp.MustCompile(`(.+) \(contains (.+)\)`)

type menu struct {
	all          []ingredients
	foods        map[string]int
	foodAlergens map[string]map[string]bool
	alergenFoods map[string]map[string]bool
}

type ingredients struct {
	foods    map[string]bool
	contains map[string]bool
}

func newMenu(all []ingredients) menu {
	return menu{
		all:          all,
		foods:        map[string]int{},
		foodAlergens: map[string]map[string]bool{},
		alergenFoods: map[string]map[string]bool{},
	}
}

func newIngredients(foodsStr, containsStr string) ingredients {
	ret := ingredients{}
	ret.foods = listToBoolMap(strings.Split(foodsStr, " "))
	ret.contains = listToBoolMap(strings.Split(containsStr, ", "))
	return ret
}

func (m *menu) buildMaps() {
	for _, i := range m.all {
		for food := range i.foods {
			m.foods[food]++
		}
	}
	somethingChanged := true
	for somethingChanged {
		somethingChanged = false
		for _, i := range m.all {
			for alergen := range i.contains {
				if _, ok := m.alergenFoods[alergen]; !ok {
					m.alergenFoods[alergen] = map[string]bool{}
				}
				for food := range i.foods {
					if _, ok := m.alergenFoods[alergen][food]; ok {
						// Either already set or set to false
						// leave as is
						continue
					} else {
						if !m.alergenFoods[alergen][food] {
							m.alergenFoods[alergen][food] = true
							somethingChanged = true
						}
					}
				}
				for oldFood := range m.alergenFoods[alergen] {
					if !i.foods[oldFood] {
						// Doesn't appear in this list therefore not
						// the the food in question
						if m.alergenFoods[alergen][oldFood] {
							m.alergenFoods[alergen][oldFood] = false
							somethingChanged = true
						}
					}
				}
			}
		}
		for _, i := range m.all {
			for food := range i.foods {
				if _, ok := m.foodAlergens[food]; !ok {
					m.foodAlergens[food] = map[string]bool{}
				}
				for alergen := range i.contains {
					if _, ok := m.foodAlergens[food][alergen]; ok {
						// Either already set or set to false
						// leave as is
						continue
					} else {
						// first time seen
						if !m.foodAlergens[food][alergen] {
							m.foodAlergens[food][alergen] = true
							somethingChanged = true
						}
					}
				}
				for oldAlergen := range m.foodAlergens[food] {
					if !i.contains[oldAlergen] {
						// Doesn't appear in this list therefore not
						// the the alergen in question
						if m.foodAlergens[food][oldAlergen] {
							m.foodAlergens[food][oldAlergen] = false
							somethingChanged = true
						}
					}
				}
			}
		}
	}
}

func (m menu) lowCounts() {
	for k, v := range m.foodAlergens {
		alergens := getTrueKeys(v)
		fmt.Printf("%s: %s\n", k, strings.Join(alergens, ","))
	}
	fmt.Printf("\n")

	possibles := map[string]bool{}
	for k, v := range m.alergenFoods {
		foods := getTrueKeys(v)
		for _, food := range foods {
			possibles[food] = true
		}
		fmt.Printf("%s: %s\n", k, strings.Join(foods, ","))
	}

	fmt.Printf("\n")

	count := 0
	impossibles := map[string]int{}

	for food := range m.foodAlergens {
		if possibles[food] {
			continue
		}
		count++
		impossibles[food] = 0
	}
	fmt.Printf("Count %d\n", count)

	count = 0
	for _, ing := range m.all {
		for food := range impossibles {
			if ing.foods[food] {
				impossibles[food]++
				count++
			}
		}
	}
	fmt.Printf("Final count %d\n", count)
}

func (i ingredients) String() string {
	return fmt.Sprintf("%s (contains %s)", strings.Join(getTrueKeys(i.foods), " "), strings.Join(getTrueKeys(i.contains), ", "))
}

func getTrueKeys(kv map[string]bool) []string {
	keys := make([]string, 0, len(kv))
	for k, v := range kv {
		if v {
			keys = append(keys, k)
		}
	}
	return keys
}

func listToBoolMap(lst []string) (m map[string]bool) {
	m = make(map[string]bool, len(lst))
	for _, key := range lst {
		m[key] = true
	}
	return m
}

func read(fname string) ([]ingredients, error) {
	ret := []ingredients{}
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

		parts := rowRx.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("Unable to parse %q\n", line)
		}
		ret = append(ret, newIngredients(parts[1], parts[2]))
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

	fmt.Printf("Day 21\n")

	infile := "day21.input"
	if *testFileFlag {
		infile = "day21-sample.input"
	}
	input, err := read(infile)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	fmt.Printf("Num rows %d\n", len(input))
	m := newMenu(input)
	m.buildMaps()
	m.lowCounts()
}
