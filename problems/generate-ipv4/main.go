package main

import "fmt"

func main() {
	GenerateIPv4([]int{0, 1, 0, 0, 0, 0}) // 0.100.0.0
}

// nums = [1,2,3,4]
// GenerateIPv4 -> Return all IPV4 that can be generated.
// 0-255 x 4
// [1,2,3,4] -> 1.2.3.4 This is the only example that can be generated.
// [1,2,1,2,2,3,1,2,3] -> 121.223.12.3, 121.223.1.23, 12.122.31.23, etc...
// 0 <= n <= 9
// - All number must used.
// 0,0,0,0,0 -> 0.0.0.0
// 0,1,0,0,0,0 -> 0.100.0.0
// 01.00.0.0
// 0.1.0.0, 0.10.0.0

func sliceToNum(num []int) (int, bool) {
	n := len(num)
	s := 0

	leadingZero := false
	if num[0] == 0 && n > 1 {
		leadingZero = true
	}

	for i := 0; i < n; i++ {
		s = s*10 + num[i]
	}

	return s, leadingZero
}

func GenerateIPv4(nums []int) {
	n := len(nums)

	// O(n ^ 3)
	// if we consider 10^6 is the most instructions that we can run in 1 seconds
	// we only can handle 100 elements in our input during 1 second.
	for i := 1; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				p1, l1 := sliceToNum(nums[:i])
				p2, l2 := sliceToNum(nums[i:j])
				p3, l3 := sliceToNum(nums[j:k])
				p4, l4 := sliceToNum(nums[k:])

				if l1 || l2 || l3 || l4 {
					continue
				}

				if p1 >= 0 && p1 <= 255 && p2 >= 0 && p2 <= 255 && p3 >= 0 && p3 <= 255 && p4 >= 0 && p4 <= 255 {
					fmt.Printf("%d.%d.%d.%d\n", p1, p2, p3, p4)
				}
			}
		}
	}
}
