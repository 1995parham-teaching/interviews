package main

import "fmt"

func maxSubarrayLength(nums []int, k int) int {
	length := 0

	start := 0
	end := 0

	frequency := map[int]int{}

	for start != len(nums) {
		for {
			frequency[nums[end]] += 1
			if frequency[nums[end]] <= k {
				if end-start+1 > length {
					length = end - start + 1
				}
			} else {
				break
			}
			if end < len(nums)-1 {
				end += 1
			}
		}

		start += 1
		if start == len(nums) {
			break
		}
		frequency[nums[start-1]] -= 1
		frequency[nums[end]] -= 1
	}

	return length
}

func main() {
	fmt.Println(maxSubarrayLength([]int{1, 2, 3, 1, 2, 3, 1, 2}, 2))
	fmt.Println(maxSubarrayLength([]int{1, 4, 4, 3}, 1))
}
