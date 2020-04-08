package main

import (
	"fmt"
)

// There are two sorted arrays. You can safely assume these arrays are in ascending order.
// These arrays can have different sizes and we want to merge them into new sorted arrays.
// You cannot use any sort's algorithms or function in this problem.
// [1, 2, 3]
// [4, 5, 6, 7]
// => [1, 2, 3, 4, 5, 6, 7]
// [1, 2]
// [1, 2]
// => [1, 1, 2, 2]
// [1]
// [1, 2]
// => [1, 1, 2]

func main() {
	f := []int{-1}
	s := []int{-2, -2}

	fmt.Println(sort(f, s))
}

func sort(first []int, second []int) []int {
	result := make([]int, 0)
	fIndex, sIndex := 0, 0

	for {
		if fIndex == len(first) {
			result = append(result, second[sIndex:]...)
			break
		} else if sIndex == len(second) {
			result = append(result, first[fIndex:]...)
			break
		}

		if first[fIndex] <= second[sIndex] {
			result = append(result, first[fIndex])
			fIndex++
		} else {
			result = append(result, second[sIndex])
			sIndex++
		}
	}

	return result
}
