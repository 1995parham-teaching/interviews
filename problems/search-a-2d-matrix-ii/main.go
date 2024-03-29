package main

import "fmt"

func searchMatrix(matrix [][]int, target int) bool {
	for i := 0; i < len(matrix); i++ {
		row := matrix[i]
		if j := search(row, target, 0); j != -1 {
			return true
		}
	}

	return false
}

func search(row []int, wanted, offset int) int {
	n := len(row)

	if n == 1 && row[0] != wanted {
		return -1
	}

	if row[n/2] == wanted {
		return offset + n/2
	}

	if row[n/2] < wanted {
		return search(row[n/2:], wanted, offset+n/2)
	}

	return search(row[:n/2], wanted, offset)
}

func main() {
	matrix := [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}}
	wanted := 20

	fmt.Println(searchMatrix(matrix, wanted))
}
