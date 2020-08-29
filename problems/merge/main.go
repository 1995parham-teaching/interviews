package main

import "fmt"

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
