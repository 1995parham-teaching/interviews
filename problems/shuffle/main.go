package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// in the following example my code will answer wrongly whether using shuffle or not
	a := []int{1, 1, 2, 2, 3, 3, 4, 4}
	k := 4

	sa := Shuffle(a)
	r := Partition(sa, k)
	fmt.Println(r)
}

// Shuffle shuffles a copy of given array in-place.
func Shuffle(arr []int) []int {
	result := make([]int, len(arr))

	copy(result, arr)

	for i := 0; i < len(result); i++ {
		j := rand.Intn(len(result)-i) + i
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// Partition partions a copy of given array into k equal size sub arrays.
func Partition(arr []int, k int) [][]int {
	result := make([][]int, k)

	n := len(arr)
	l := n / k

	for i := 0; i < k; i++ {
		result[i] = make([]int, l)
		copy(result[i], arr[i*l:(i+1)*l])
	}

	return result
}
