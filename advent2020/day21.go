// See https://adventofcode.com/2020/day/21 for problem decscription
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
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
	fmt.Printf("\n")
	changed := true
	for changed {
		changed = false
		for k, v := range m.alergenFoods {
			onlyFood := getOnlyKey(v)
			if onlyFood != "" {
				changed = m.zapOthers(k, onlyFood)
			}
		}
	}
	finalMap := map[string]string{}
	for k, v := range m.alergenFoods {
		foods := getTrueKeys(v)
		fmt.Printf("%s: %s\n", k, strings.Join(foods, ","))
		finalMap[k] = foods[0]
	}
	toSort := []string{}
	for k := range m.alergenFoods {
		toSort = append(toSort, k)
	}
	sort.Strings(toSort)
	final := []string{}
	for _, k := range toSort {
		final = append(final, finalMap[k])
	}
	fmt.Printf("%s\n", strings.Join(final, ","))

}

// Zap food for all other alergens except `alergen`
// returns true is something needed to be zapped
func (m *menu) zapOthers(alergen string, food string) bool {
	zapped := false
	for k := range m.alergenFoods {
		if k != alergen && m.alergenFoods[k][food] {
			m.alergenFoods[k][food] = false
			zapped = true
		}
	}
	return zapped
}

func getOnlyKey(m map[string]bool) string {
	key := ""
	for k, v := range m {
		if v {
			if key != "" {
				// more than one key with bool true
				return ""
			}
			key = k
		}
	}
	return key
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
			return nil, fmt.Errorf("unable to parse %q", line)
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
