package main

import (
	"fmt"
)

func main() {
	// In the following example my code will answer wrongly whether using shuffle or not
	a := []int{1, 1, 2, 2, 3, 3, 4, 4}
	k := 4

	// Wrong answer
	r := subArray(a, k)
	fmt.Println(r)

	// Correct answer
	r = subArrayPrecisely(a, k)
	fmt.Println(r)
}

func subArray(a []int, k int) [][]int {
	result := make([][]int, k)
	l := len(a) / k

	index := 0

	for i := range result {
		sub := a[index : index+l]
		result[i] = sub

		index += l
	}

	first := result[0]
	temp := first[0]
	first[0] = first[1]
	first[1] = temp

	return result
}

func subArrayPrecisely(a []int, k int) [][]int {
	result := make([][]int, k)
	l := len(a) / k

	fSub := make([]int, 0)
	fSub = append(fSub, a[0])
	fSub = append(fSub, a[len(a)-l+1:]...)
	result[0] = fSub

	index := 1

	for i := 1; i < len(result); i++ {
		sub := a[index : index+l]
		result[i] = sub

		index += l
	}

	return result
}

// This code may cause no shuffle
//a := []int{1, 1, 3}
//rand.Seed(time.Now().UnixNano())
//rand.Shuffle(len(a), func(i, j int) {
//	a[i], a[j] = a[j], a[i]
//})
