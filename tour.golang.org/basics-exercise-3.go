// Exercise: Maps
package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	fields := strings.Fields(s)
	m := make(map[string]int, len(fields))
	for _, field := range fields {
		m[field]++
	}
	return m
}

func main() {
	wc.Test(WordCount)
}
