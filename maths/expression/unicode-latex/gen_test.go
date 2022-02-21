package main

import "testing"

func Test_importData(t *testing.T) {
	chars, err := importData()
	if err != nil {
		t.Fatal(err)
	}

	if chars['a'] != "a" {
		t.Fatal()
	}
	if chars['\u03C0'] != "\\pi" {
		t.Fatal()
	}
}
