package main

import "fmt"

func factorial(n int) int {
	if n == 0 || n == 1 {
		return 1
	}

	return factorial(n-1) * n
}

func getPermutationWithList(nums []int, k int) string {
	n := len(nums)

	if n == 1 {
		return fmt.Sprintf("%d", nums[0])
	}

	i := 0

	for i*factorial(n-1) < k {
		i++
	}

	i--

	v := nums[i]
	nums = append(nums[:i], nums[i+1:]...)

	return fmt.Sprintf("%d%s", v, getPermutationWithList(nums, k-i*factorial(n-1)))
}

func getPermutation(n int, k int) string {
	nums := make([]int, n)

	for i := 1; i <= n; i++ {
		nums[i-1] = i
	}

	return getPermutationWithList(nums, k)
}

func main() {
	fmt.Println(getPermutation(3, 3))
}
