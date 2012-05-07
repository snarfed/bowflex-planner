// Misc utility functions.
package app

import (
	"net/url"
)

func Copy(input url.Values) url.Values {
	copied := make(url.Values, len(input))
	for key, val := range input {
		copied[key] = val
	}
	return copied
}

func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

func Permutations(n int) [][]int {
	elems := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		elems = append(elems, i)
	}
	return perms_helper(elems)
}

func perms_helper(elems []int) [][]int {
	// base case
	if len(elems) == 1 {
		return [][]int{elems}
	}

	n := len(elems)
	perms := make([][]int, 0, Factorial(n))

	// recursive step
	for i, e := range elems {
		// make a copy of elems with everything except e
		rest := make([]int, 0, n-1)
		rest = append(rest, elems[:i]...)
		if i < n-1 {
			rest = append(rest, elems[i+1:]...)
		}

		// populate the n permutations that end with e
		for _, subperm := range perms_helper(rest) {
			perms = append(perms, append(subperm, e))
		}
	}

	return perms
}
