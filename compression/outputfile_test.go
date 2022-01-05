package main

import "testing"

func TestOutputFile(t *testing.T) {
	cases := map[string][]string{
		"":                     {"main"},
		"testing_out.txt.gzip": {"main", "testing.txt"},
		"newtesting.gzip":      {"main", "testing", "newtesting"},
		"newtesting.txt.gzip":  {"main", "testing.txt", "newtesting.txt"},
	}

	for expected, args := range cases {
		if output := outputFile(args); output != expected {
			t.Errorf("supplied %s and expected %s but got %s", args, expected, output)
		}
	}
}
