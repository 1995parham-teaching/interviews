package main

import "fmt"

func countSubarraysWithMax(nums []int, k int, max int) int64 {
	if len(nums) < k {
		return 0
	}

	firstMaxIndex := -1
	i := 0
	for firstMaxIndex == -1 {
		if i == len(nums) {
			return 0
		}
		if nums[i] == max {
			firstMaxIndex = i
		}
		i++
	}

	c := countSubarraysWithMax(nums[firstMaxIndex+1:], k, max)

	n := 0
	lastMaxIndex := -1
	for i := firstMaxIndex; i < len(nums); i++ {
		if nums[i] == max {
			n++
		}
		if n == k {
			lastMaxIndex = i
			break
		}
	}

	if lastMaxIndex == -1 {
		return c
	}

	return c + int64(firstMaxIndex+1)*int64(len(nums)-lastMaxIndex)
}

func countSubarrays(nums []int, k int) int64 {
	max := -1
	for _, n := range nums {
		if n > max {
			max = n
		}
	}

	return countSubarraysWithMax(nums, k, max)
}

func main() {
	tests := []struct {
		nums   []int
		k      int
		answer int64
	}{
		{
			nums:   []int{1, 3, 2, 3, 3},
			k:      2,
			answer: 6,
		},
		{
			nums:   []int{1, 4, 2, 1},
			k:      3,
			answer: 0,
		},
	}

	for _, t := range tests {
		if t.answer != countSubarrays(t.nums, t.k) {
			fmt.Printf("The number of subarrays should be %d but it is %d for %v %d\n",
				t.answer,
				countSubarrays(t.nums, t.k),
				t.nums,
				t.k,
			)
		}
	}
}
