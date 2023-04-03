package main

func sortColors(nums []int) {
	var colors [3]int

	for _, c := range nums {
		colors[c] += 1
	}

	for i := 0; i < colors[0]; i++ {
		nums[i] = 0
	}

	for i := colors[0]; i < colors[0]+colors[1]; i++ {
		nums[i] = 1
	}

	for i := colors[0] + colors[1]; i < colors[0]+colors[1]+colors[2]; i++ {
		nums[i] = 2
	}
}
