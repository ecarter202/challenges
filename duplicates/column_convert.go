package main

import (
	"math"
	"strings"
)

func columnLetterToIndex(letter string) int {
	var column = 0

	for i := 0; i < len(letter); i++ {
		if value, ok := isLetter(letter); ok {
			column += (value - 64) * int(math.Pow(26, float64(len(letter)-i-1)))
		}
	}

	return column - 1 // start indexing at zero
}

func isLetter(letter string) (value int, is bool) {
	x := []rune(strings.ToUpper(letter))
	l := x[0]

	if len(x) > 0 {
		return int(l), l > 64 && l < 91
	}

	return int(l), false
}
