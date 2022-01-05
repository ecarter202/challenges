package main

import (
	"testing"
)

func TestColumnLetterToIndex(t *testing.T) {
	cases := map[string]int{
		"a":  0,
		"z":  25,
		"aa": 26,
		"zz": 701,
	}

	for letter, value := range cases {
		if returned := columnLetterToIndex(letter); returned != value {
			t.Errorf("supplied %s and expected %d but got %d", letter, value, returned)
		} else {
			t.Logf("supplied %s and got %d", letter, returned)
		}
	}
}
