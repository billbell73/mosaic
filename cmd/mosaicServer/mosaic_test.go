package main

import "testing"

func TestSq(t *testing.T) {
	var v int
	v = sq(4)
	if v != 16 {
		t.Error("Expected 16, got ", v)
	}
}

func TestDistance(t *testing.T) {
	var v int
	p1 := [3]int{1, 1, 2}
	p2 := [3]int{3, 4, 8}
	v = distance(p1, p2)
	if v != 49 {
		t.Error("Expected 49, got ", v)
	}
}

func TestNearest(t *testing.T) {
	target := [3]int{1, 1, 1}
	db := map[string][3]int{
		"image2": [3]int{10, 10, 10},
		"image1": [3]int{3, 2, 2},
		"image3": [3]int{3, 4, 19},
	}
	filename := nearest(target, &db)
	if filename != "image1" {
		t.Error("Expected \"image1\", got ", filename)
	}
}
