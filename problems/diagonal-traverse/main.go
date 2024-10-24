// Given an MxN matrix, write code which prints out the diagonals (from upper right to lower left)
// of the matrix. In this example where M = 4, N = 3:
// {{9, 3, 2},
//  {8, 6, 1},
//  {5, 5, 6},
//  {1, 2, 8}}

// Your code should print out:
// 9
// 3 8
// 2 6 5
// 1 5 1
// 6 2
// 8

package main

import "fmt"

type tuple struct {
	x int
	y int
}

func main() {
	arr := [][]int{
		{9, 3, 2},
		{8, 6, 1},
		{5, 5, 6},
		{1, 2, 8},
	}
	M := len(arr)
	N := len(arr[0])

	diags := make([]tuple, 0)

	for i := 0; i < N; i++ {
		diags = append(diags, tuple{0, i})
	}
	for i := 1; i < M; i++ {
		diags = append(diags, tuple{i, N - 1})
	}

	for _, t := range diags {
		i, j := t.x, t.y
		for i < M && j >= 0 {
			fmt.Printf("%d ", arr[i][j])
			j--
			i++
		}
		fmt.Println()
	}
}
