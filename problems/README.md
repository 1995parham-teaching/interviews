# Problems

## [Shuffle](./shuffle)

You have an array with n-items (A).
We want to partition it into k-subarrays that each of them has n/k items, and each element of A appears precisely once.
The order of these subarrays must not be the same as the A.

**We know that: n % k == 0**

- With duplication or without duplication?

For example:

```
A = {1, 2, 3, 4}

k = 2
```

we don't accept the following solution:

```
A1 = {1, 2}
A2 = {3, 4}
```

but we accept the following solution:

```
A1 = {1, 3}
A2 = {2, 4}
```

## [Coins](./coins)

We have `n` amount of money and our country have the following coins:

- coin-1
- coin-5
- coin-7
- coin-10

We want to have this money with minimum number of coins. What is the minimum? For example:

- 2 = 2 x coin-1
- 5 = 1 x coin-5
- 6 = 1 x coin-5 + 1 x coin-1

## Bulb Switcher

There are n bulbs that are initially off. You first turn on all the bulbs.
Then, you turn off every second bulb. On the third round,
you toggle every third bulb (turning on if it's off or turning off if it's on).
For the i-th round, you toggle every i bulb.
For the n-th round, you only toggle the last bulb.
Find how many bulbs are on after n rounds.

Example:

```
Input: 3
Output: 1
Explanation:
At first, the three bulbs are [off, off, off].
After first round, the three bulbs are [on, on, on].
After second round, the three bulbs are [on, off, on].
After third round, the three bulbs are [on, off, off].
```

So you should return 1, because there is only one bulb is on.

[LeetCode](https://leetcode.com/problems/bulb-switcher/)

## Happy Number

Write an algorithm to determine if a number n is "happy".

A happy number is a number defined by the following process: Starting with any positive integer, replace the number by the sum of the squares of its digits, and repeat the process until the number equals 1 (where it will stay), or it loops endlessly in a cycle which does not include 1. Those numbers for which this process ends in 1 are happy numbers.

Return True if n is a happy number, and False if not.

Example:

```
Input: 19
Output: true
Explanation:
1^2 + 9^2 = 82
8^2 + 2^2 = 68
6^2 + 8^2 = 100
1^2 + 0^2 + 0^2 = 1
```

[LeetCode](https://leetcode.com/problems/happy-number/)

## Rotate Image

You are given an n x n 2D matrix representing an image. Rotate the image by 90 degrees (clockwise).

```
Given input matrix =
[
  [1,2,3],
  [4,5,6],
  [7,8,9]
],

rotate the input matrix in-place such that it becomes:
[
  [7,4,1],
  [8,5,2],
  [9,6,3]
]
```

```
Given input matrix =
[
  [ 5, 1, 9,11],
  [ 2, 4, 8,10],
  [13, 3, 6, 7],
  [15,14,12,16]
],

rotate the input matrix in-place such that it becomes:
[
  [15,13, 2, 5],
  [14, 3, 4, 1],
  [12, 6, 8, 9],
  [16, 7,10,11]
]
```

[LeetCode](https://leetcode.com/problems/rotate-image/)

## Snappfood

We have motorcycles and restaurants. Motorcycles deliver foods to peoples from restaurants.
How we can schedule this delivery process?

## [Search a 2D Matrix]()

You are given a `m x n` integer matrix `matrix` with the following two properties:

- Each row is sorted in non-decreasing order.
- The first integer of each row is greater than the last integer of the previous row.

Given an integer `target`, return `true` if `target` is in `matrix` or `false` otherwise.

You must write a solution in `O(log(m * n))` time complexity.

Example 1:

```
Input: matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
Output: true
```

Example 2:

```
Input: matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 13
Output: false
```

Constraints:

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 100`
- `-10^4 <= matrix[i][j], target <= 10^4`

## [Search a 2D Matrix II](./search-a-2d-matrix-ii)

Write an efficient algorithm that searches for a value `target` in an `m x n` integer matrix matrix.
This matrix has the following properties:

- Integers in each row are sorted in ascending from left to right.
- Integers in each column are sorted in ascending from top to bottom.

```
Input: matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
Output: true
```

```
Input: matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 20
Output: false
```

Constraints:

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= n, m <= 300`
- `-109 <= matrix[i][j] <= 109`
- All the integers in each row are sorted in ascending order.
- All the integers in each column are sorted in ascending order.
- `-109 <= target <= 109`

[LeetCode](https://leetcode.com/problems/search-a-2d-matrix-ii/)

## [Longest Palindromic Substring](./longest-palindrome)

Given a string `s`, return the longest palindromic substring in `s`.

Example 1:

```
Input: s = "babad"
Output: "bab"
Explanation: "aba" is also a valid answer.
```

Example 2:

```
Input: s = "cbbd"
Output: "bb"
```

[LeetCode](https://leetcode.com/problems/longest-palindromic-substring/)

## Kth Smallest Element in a Sorted Matrix

[LeetCode](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/)

## [Merge k Sorted Lists](./merge-k-sorted-lists)

You are given an array of `k` linked-lists `lists`, each linked-list is sorted in ascending order.
Merge all the linked-lists into one sorted linked-list and return it.

Example 1:

```
Input: lists = [[1,4,5],[1,3,4],[2,6]]
Output: [1,1,2,3,4,4,5,6]
Explanation: The linked-lists are:
[
  1->4->5,
  1->3->4,
  2->6
]
merging them into one sorted list:
1->1->2->3->4->4->5->6
```

Example 2:

```
Input: lists = []
Output: []
```

Example 3:

```
Input: lists = [[]]
Output: []
```

[LeetCode](https://leetcode.com/problems/merge-k-sorted-lists/)

## K Nearest Neighbor

We have `n` points and one reference point.
Each point has `x` and `y` coordinates.
We want to find `k` the nearest points to the reference point.

For example:

```python
import dataclasses

@dataclasses.dataclass()
class Point:
  x: float
  y: float

points = [
  Point(0, 0), Point(0, 1), Point(1, 1), Point(1, 0),
  Point(-1, -1), Point(0, -1), Point(-1, 0),
]
reference = Point(-1, -1)
n = len(points)
k = 2

k_nearest_points = [Point(-1, -1), Point(-1, 0)]
# or
k_nearest_points = [Point(-1, -1), Point(0, -1)]
```
