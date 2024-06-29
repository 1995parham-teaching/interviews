/*
You are given a 0-indexed integer array nums having length n, and an integer k.

You can perform the following increment operation any number of times (including zero):

Choose an index i in the range [0, n - 1], and increase nums[i] by 1.
An array is considered beautiful if, for any subarray with a size of 3 or more,
its maximum element is greater than or equal to k.

Return an integer denoting the minimum number of increment operations needed to make nums beautiful.

A subarray is a contiguous non-empty sequence of elements within an array.



Example 1:

Input: nums = [2,3,0,0,2], k = 4
[i = 1, i = 1, = 4] -> [1, 4] -> [3, 2] -> 3
21012 k = 5
[2, 1, 0] -> 2 -> 3
[1, 0, 1] -> 1 -> 4
[0, 1, 2] -> 2 -> 3
Output: 3
f(n-1)
f(n-2)
f(n-3)
A = [ .... ] -> X
[ A , k ] -> X  10, 4, 0,0, 10

min(k - max(y, z, k), 0)
[ ..... ]
[a,b,c,X]
    x f1
   xx f2 f3
  xxxx
*/

package main

func main() {
}

func f(n int, nums []int, k int) int {
	if n == 0 || n == 1 || n == 2 {
		return 0
	}

	p1 := max(k-nums[0], 0) + f(n-1, nums, k)
	p2 := max(k-nums[1], 0) + f(n-2, nums, k)
	p3 := max(k-nums[2], 0) + f(n-3, nums, k)

	return min(min(p1, p2), p3)
}
