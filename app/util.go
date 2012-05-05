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

// TODO: write map copy fn and implement recursively by passing down remaining elements
// actually, use array/slice and -1 as tombstone and fill that in for removed elements
func Permutations(n int) [][]int {
	num := Factorial(n)
	perms := make([][]int, num)
	for i := range perms {
		perms[i] = make([]int, n)
	}

	for j := 0; j < n; j++ {
		perms[i][j] = (i*n*j/num + j) % n
	}

	return perms
}

// var x = []int{3, nil}

// func copy(map[int]

// func Permutations(n int) [][]int {
// 	num := Factorial(n)
// 	perms := make([][]int, num)
// 	for i := range perms {
// 		perms[i] = make([]int, n)
// 		for j := 0; j < n; j++ {
// 			perms[i][j] = (i*n*j/num + j) % n
// 		}
// 	}

// 	return perms
// }
