package app

import (
	"net/url"
	"testing"
)

func TestCopy(t *testing.T) {
	orig := url.Values{"foo": {"bar"}}
	copied := Copy(orig)

	if len(copied) != 1 || len(copied["foo"]) != 1 || copied["foo"][0] != "bar" {
		t.Error(orig, copied)
	}

	orig["baz"] = []string{"baj"}

	if copied["baz"] != nil {
		t.Error(orig, copied)
	}
}

func TestFactorial(t *testing.T) {
	for n, expected := range map[int]int{0: 1, 1: 1, 2: 2, 3: 6, 4: 24, 5: 120} {
		fact := Factorial(n)
		if fact != expected {
			t.Error(n, fact)
		}
	}
}

func Equal(this [][]int, that [][]int) bool {
	for i, this_i := range this {
		for j, this_i_j := range this_i {
			if this_i_j != that[i][j] {
				return false
			}
		}
	}
	return true
}

func TestPermutations(t *testing.T) {
	for n, expected := range map[int][][]int{
		0: [][]int{},
		1: [][]int{{0}},
		2: [][]int{{0, 1}, {1, 0}},
		3: [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}},
	} {
		perms := Permutations(n)
		println(perms)
		if !Equal(perms, expected) {
			t.Error(n, perms)
		}
	}
}
