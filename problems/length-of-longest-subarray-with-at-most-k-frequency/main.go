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
	tests := []struct {
		nums   []int
		k      int
		answer int
	}{
		{
			nums:   []int{1, 2, 3, 1, 2, 3, 1, 2},
			k:      2,
			answer: 6,
		}, {
			nums:   []int{1, 4, 4, 3},
			k:      1,
			answer: 2,
		},
	}

	for _, t := range tests {
		if maxSubarrayLength(t.nums, t.k) != t.answer {
			fmt.Printf("maximum length should be %d for %v and %d but it is %d\n",
				t.answer,
				t.nums,
				t.k,
				maxSubarrayLength(t.nums, t.k),
			)
		}
	}
}
