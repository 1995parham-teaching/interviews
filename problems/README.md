# Problems

These are some hands-on problem, and you can use a shared Google Doc to read and write code together.
They don't require the interviewee to have an IDE or coding environment, and the main part of your discussion
should be Algorithm.

## [KNN](./knn)

Implement KNN Algorithm.

## [Decode String](./decode-string)

Given an encoded string, return its decoded string.

The encoding rule is: `k[encoded_string]`, where the `encoded_string` inside the square brackets is being
repeated exactly `k` times. Note that `k` is guaranteed to be a positive integer.

You may assume that the input string is always valid; there are no extra white spaces, square brackets are well-formed,
etc. Furthermore, you may assume that the original data does not contain any digits and that digits are only
for those repeat numbers, k. For example, there will not be input like `3a` or `2[4]`.

The test cases are generated so that the length of the output will never exceed `10^5`.

Example 1:

```
Input: s = "3[a]2[bc]"
Output: "aaabcbc"
```

Example 2:

```
Input: s = "3[a2[c]]"
Output: "accaccacc"
```

Example 3:

```
Input: s = "2[abc]3[cd]ef"
Output: "abcabccdcdcdef"
```

Constraints:

- `1 <= s.length <= 30`
- `s` consists of lowercase English letters, digits, and square brackets `'[]'`.
- `s` is guaranteed to be a valid input.
- All the integers in `s` are in the range `[1, 300]`.

[LeetCode](https://leetcode.com/problems/decode-string/description/)

## [Find All Groups of Farmland](./find-all-groups-of-farmland)

You are given a **0-indexed** `m x n` binary matrix `land` where a `0` represents a hectare of forested land and
a `1` represents a hectare of farmland.

To keep the land organized, there are designated rectangular areas of hectares that consist **entirely** of farmland.
These rectangular areas are called **groups**. No two groups are adjacent, meaning farmland in one group
is not four-directionally adjacent to another farmland in a different group.

`land` can be represented by a coordinate system where the top left corner of `land` is `(0, 0)` and the bottom right
corner of `land` is `(m-1, n-1)`. Find the coordinates of the top left and bottom right corner of each **group**
of farmland. A **group** of farmland with a top left corner at `(r1, c1)` and a bottom right corner at
`(r2, c2)` is represented by the 4-length array `[r1, c1, r2, c2]`.

Return a 2D array containing the 4-length arrays described above for each **group** of farmland in `land`.
If there are no groups of farmland, return an empty array. You may return the answer in **any order**.

[LeetCode](https://leetcode.com/problems/find-all-groups-of-farmland/)

## [Number of Islands](./number-of-islands)

Given an m x n 2D binary grid which represents a map of '1's (land) and '0's (water), return the number of islands.

An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically. You may assume all four edges of the grid are all surrounded by water.

Example 1:

```
Input: grid = [
  ["1","1","1","1","0"],
  ["1","1","0","1","0"],
  ["1","1","0","0","0"],
  ["0","0","0","0","0"]
]
Output: 1
```

Example 2:

```
Input: grid = [
  ["1","1","0","0","0"],
  ["1","1","0","0","0"],
  ["0","0","1","0","0"],
  ["0","0","0","1","1"]
]
Output: 3
```

Constraints:

- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 300`
- `grid[i][j] is '0' or '1'`

[LeetCode](https://leetcode.com/problems/number-of-islands/)

## [Shuffle](./shuffle)

You have an array with n-items (A).
We want to partition it into k-subarrays that each of them has n/k items, and each element of A appears precisely once.
The order of these subarrays must not be the same as the A.

**We know that: n % k == 0**

- With duplication or without duplication?

For example:

```python
A = [1, 2, 3, 4]

k = 2
```

we don't accept the following solution:

```python
A1 = {1, 2}
A2 = {3, 4}
```

but we accept the following solution:

```python
A1 = [1, 3]
A2 = [2, 4]
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

```text
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

A happy number is a number defined by the following process: Starting with any positive integer,
replace the number by the sum of the squares of its digits, and repeat the process
until the number equals 1 (where it will stay), or it loops endlessly in a cycle which does not include 1.
Those numbers for which this process ends in 1 are happy numbers.

Return True if n is a happy number, and False if not.

Example:

```text
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

```text
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

```text
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

## [Search a 2D Matrix](./search-a-2d-matrix)

You are given a `m x n` integer matrix `matrix` with the following two properties:

- Each row is sorted in non-decreasing order.
- The first integer of each row is greater than the last integer of the previous row.

Given an integer `target`, return `true` if `target` is in `matrix` or `false` otherwise.

You must write a solution in `O(log(m * n))` time complexity.

Example 1:

```text
Input: matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
Output: true
```

Example 2:

```text
Input: matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 13
Output: false
```

Constraints:

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 100`
- `-10^4 <= matrix[i][j], target <= 10^4`

## [Search a 2D Matrix II](./search-a-2d-matrix-ii)

Write an efficient algorithm that searches for a value `target` in an `m x n` integer matrix.
This matrix has the following properties:

- Integers in each row are sorted in ascending from left to right.
- Integers in each column are sorted in ascending from top to bottom.

```text
Input: matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
Output: true
```

```text
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

```text
Input: s = "babad"
Output: "bab"
Explanation: "aba" is also a valid answer.
```

Example 2:

```text
Input: s = "cbbd"
Output: "bb"
```

[LeetCode](https://leetcode.com/problems/longest-palindromic-substring/)

## K-th Smallest Element in a Sorted Matrix

[LeetCode](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/)

## [Merge k Sorted Lists](./merge-k-sorted-lists)

You are given an array of `k` linked-lists `lists`, each linked-list is sorted in ascending order.
Merge all the linked-lists into one sorted linked-list and return it.

Example 1:

```text
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

```text
Input: lists = []
Output: []
```

Example 3:

```text
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

## [Sort Colors](./sort-colors)

Given an array `nums` with `n` objects colored red, white, or blue,
sort them in-place so that objects of the same color are adjacent,
with the colors in the order red, white, and blue.

We will use the integers `0`, `1`, and `2` to represent the color red, white, and blue, respectively.

You must solve this problem without using the library's sort function

Example 1:

```text
Input: nums = [2,0,2,1,1,0]
Output: [0,0,1,1,2,2]
```

Example 2:

```text
Input: nums = [2,0,1]
Output: [0,1,2]
```

Constraints:

- `n == nums.length`
- `1 <= n <= 300`
- `nums[i]` is either `0`, `1`, or `2`.

Follow up: Could you come up with a one-pass algorithm using only constant extra space?

[LeetCode](https://leetcode.com/problems/sort-colors)

## [Generate Parentheses](./generate-parentheses)

Given `n` pairs of parentheses, write a function to generate all combinations of well-formed parentheses.

Example 1:

```text
Input: n = 3
Output: ["((()))","(()())","(())()","()(())","()()()"]
```

Example 2:

```text
Input: n = 1
Output: ["()"]
```

Constraints:

- `1 <= n <= 8`

[LeetCode](https://leetcode.com/problems/generate-parentheses)

## [Longest Valid Parentheses](./longest-valid-parentheses)

Given a string containing just the characters '(' and ')', return the length of the longest valid (well-formed) parentheses
substring.

Example 1:

```text
Input: s = "(()"
Output: 2
Explanation: The longest valid parentheses substring is "()".
```

Example 2:

```text
Input: s = ")()())"
Output: 4
Explanation: The longest valid parentheses substring is "()()".
```

Example 3:

```text
Input: s = ""
Output: 0
```

Constraints:

- `0 <= s.length <= 3 * 10^4`
- `s[i] is '(', or ')'`

[LeetCode](https://leetcode.com/problems/longest-valid-parentheses/)

## Min By Column

### Part 1 of 2

Imagine that we are working with a simple database.
Each row associates column names (strings) with integer values.
Here's a table with three rows:

```text
a b c d
1 0 0 0
0 2 3 0
0 0 0 4
```

We might choose to represent a database table in JSON, as an array of objects.
For example, the previous table could be written as:

```json
[
  { "a": 1, "b": 0, "c": 0, "d": 0 },
  { "a": 0, "b": 2, "c": 3, "d": 0 },
  { "a": 0, "b": 0, "c": 0, "d": 4 }
]
```

Write a function, `min_by_column`, that takes a database table (as above),
along with a column name, and returns the row that contains the minimum value for the given column.
If a row doesn't have any value for the column, your function should behave as though the value for that column was zero.

#### Examples

```python
table_1 = [
{"a": 1},
{"a": 2},
{"a": 3}
]
assert min_by_column(table_1, "a") == {"a": 1}
```

```python
table_2 = [
{"a": 1, "b": 2},
{"a": 3, "b": 0}
]
assert min_by_column(table_2, "b") == {"a": 3, "b": 0}
```

```python
table_3 = [
{"a": 1, "b": -2},
{"a": 3}
]
assert min_by_column(table_3, "b") == {"a": 1, "b": -2}
```

### Part 2 of 2

In Part 1 you may have noticed that it's possible for two rows to be "tied",
meaning that either would be an acceptable return value from `min_by_column`.

Consider:

```python
table_4 = [
{"a": 1, "b": 2},
{"a": 1, "b": 3},
{"a": 1, "b": 4}
]
assert min_by_column(table_4, "a") == '???'
```

Since all three rows have the same value for column "a",
all three rows are acceptable candidates to be returned by `min_by_column(table, "a")`.

In these cases, it would be nice if users could specify additional columns (e.g. "b") to use as tie-breakers.
A tie-breaker would only apply in cases where multiple rows share the same minimum value.
In `table_4` above, the row `{"a": 1, "b": 2}` is tied for the smallest "a" value (1) and of all the tied candidates,
it has the smallest "b" value (2). If two records had equal values for "a" and also for "b" then another
tie-breaker (e.g. "c") could be used.
When records are tied with respect to all columns, any of the tied records may be considered the minimum.

Write a function `min_by_columns` that takes a database table and an ordered list of column names,
and returns the row with the minimum column values using the tie-breaking logic above.
Refactor `min_by_column` to use `min_by_columns` to produce its result.

#### Examples

```python
table_5 = [
{"x": 1, "y": 3},
{"x": 1, "y": 0}
]
assert min_by_columns(table_5, ["x", "y"]) == {"x": 1, "y": 0}
```

```python
table_6 = [
{"x": 2, "y": 3},
{"x": 2, "y": 1},
{"x": 1, "y": 10}
]
assert min_by_columns(table_6, ["x", "y"]) == {"x": 1, "y": 10}
```

```python
table_7 = [
{"x": 3, "y": -1, "z": 0},
{"x": 1, "y": 10, "z": 1},
{"x": 1, "y": 10, "z": 0}
]
assert min_by_columns(table_7, ["x", "y", "z"]) == {"x": 1, "y": 10, "z": 0}
```

```python
table_8 = [
{"x": 1, "y": 2, "z": 3},
{"x": 2, "y": 2, "z": 2}
]
assert min_by_columns(table_8, ["x", "y", "z"]) == {"x": 1, "y": 2, "z": 3}
```

## Prefix Search

We have a database, and we'd like it to support these operations:

- `insert(word)`: Inset a word to database
- `look up(prefix)`: Return all the words starting with the given prefix
- `delete(prefix)`: Delete all the words starting with the given prefix
- `count(prefix)`Count the number of words starting with the given prefix

## Package Delivery

### Part 1

Write a Delivery class (or object) that represents a delivery with a destination
and distance.
Deliveries require different sensors, depending on their distance.

Add a method, `getNeededSensors`, that returns a mapping of sensor name to the
count of that sensor needed to complete the delivery according to these rules:

- If distance < 10 miles, require 1 gps and 1 temp sensor.
- If 10 <= distance < 100 miles require 1 gps, 2 temp, and 1 weight sensor.
- If distance >= 100 miles require 2 gps, 4 temp, and 2 weight sensors.

### Part 2

Write a Scheduler class (or object) that represents a daily delivery scheduler with a set of available sensors.
Add a method, `scheduleDeliveries`, that given a list of deliveries as an argument,
returns a list of deliveries that can be made that day.
Assume all deliveries will be leaving at the same time every day, so sensors can only
be used once.

Test Cases to consider:

```python
deliveryA = Delivery("A", 9)
deliveryB = Delivery("B", 15)
deliveryC = Delivery("C", 100)
scheduler = Scheduler({"gps": 2, "temp": 4, "weight": 2})
scheduler.scheduleDeliveries([deliveryA, deliveryB, deliveryC]) == [deliveryA, deliveryB]
scheduler.scheduleDeliveries([deliveryA, deliveryC, deliveryB]) == [deliveryA, deliveryB]
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) == [deliveryC]
```

### Part 3

We get paid a flat fee for all deliveries. Modify the `scheduleDeliveries`
function to maximize the number of deliveries that will be made in a day.

Our previous test:

```python
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryC]
```

Should now return:

```python
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryA, deliveryB]
```

### Part 4

We recently purchased a new type of sensor, `doorSensor`.
A `doorSensor` can be used in place of 1 weight sensor or in place of 2 temperature
sensors at any time.
Modify our existing functions to maximize the day's deliveries with the new sensor.

## K-th biggest number

time complexity of retrieving the biggest number in a list: O(n)
time complexity of retrieving the second-biggest number in a list: 2\*O(n) = O(n)
time complexity of retrieving the k-th biggest number in a list:
if k is smaller than lg(n) we can retrieve the element in `O(kn)` and if k is bigger than `lg(n)` we can retrieve the
element in `O(nlg(n))` by sorting the list and returning the k-th element

## [Prefix and Suffix Search](./prefix-and-suffix-search/)

[LeetCode](https://leetcode.com/problems/prefix-and-suffix-search/)

## [Print N-bit binary numbers having more 1s than 0s](./print-n-bit-binary-numbers-having-more-1s-than-0s/)

Given a positive integer `n`.
Your task is to generate a string list of all n-bit binary numbers where, for any prefix of the number,
there are more or an equal number of 1's than 0's. The numbers should be sorted in decreasing order of magnitude.

Example 1:

```text
Input:
n = 2
Output:
"11, 10"
Explanation: Valid numbers are those where each prefix has more 1s than 0s:
11: all its prefixes (1 and 11) have more 1s than 0s.
10: all its prefixes (1 and 10) have more 1s than 0s.
So, the output is "11, 10".
```

Example 2:

```text
Input:
n = 3
Output:
"111, 110, 101"
Explanation: Valid numbers are those where each prefix has more 1s than 0s.
111: all its prefixes (1, 11, and 111) have more 1s than 0s.
110: all its prefixes (1, 11, and 110) have more 1s than 0s.
101: all its prefixes (1, 10, and 101) have more 1s than 0s.
So, the output is "111, 110, 101".
```

User Task:
Your task is to complete the function `NBitBinary()` which takes a single number as input `n` and
returns the list of strings in **decreasing** order. You need not take any input or print anything.

```python
class Solution:
  def NBitBinary(self, n):
    pass
```

Expected Time Complexity: `O(|2n|)`
Expected Auxiliary Space: `O(2n)`

Constraints:

```
1 <= n <= 15
```

[GeeksForGeeks](https://www.geeksforgeeks.org/problems/print-n-bit-binary-numbers-having-more-1s-than-0s0252)

## [Kadane's Algorithm](./kadanes-algorithm/)

Given an array `Arr[]` of `N` integers.
Find the contiguous sub-array (containing at least one number) which has the maximum sum and return its sum.

Example 1:

```text
Input:
N = 5
Arr[] = {1,2,3,-2,5}
Output:
9
Explanation:
Max subarray sum is 9
of elements (1, 2, 3, -2, 5) which
is a contiguous subarray.
```

Example 2:

```text
Input:
N = 4
Arr[] = {-1,-2,-3,-4}
Output:
-1
Explanation:
Max subarray sum is -1
of element (-1)
```

Your Task:
You don't need to read input or print anything.
The task is to complete the function `maxSubarraySum()` which takes `Arr[]` and `N` as input parameters
and returns the sum of subarray with maximum sum.

Expected Time Complexity: `O(N)`
Expected Auxiliary Space: `O(1)`

Constraints:

```
1 ≤ N ≤ 10^6
-10^7 ≤ A[i] ≤ 10^7
```

[GeeksForGeeks](https://www.geeksforgeeks.org/problems/kadanes-algorithm-1587115620/)

## [Find All Duplicates in an Array](./find-all-duplicates-in-an-array)

Given an integer array `nums` of length `n` where all the integers of `nums` are in the range `[1, n]`
and each integer appears once or twice, return an array of all the integers that appears twice.

You must write an algorithm that runs in `O(n)` time and uses only constant extra space.

Example 1:

```
Input: nums = [4,3,2,7,8,2,3,1]
Output: [2,3]
```

Example 2:

```
Input: nums = [1,1,2]
Output: [1]
```

Example 3:

```
Input: nums = [1]
Output: []
```

Constraints:

- `n == nums.length`
- `1 <= n <= 105`
- `1 <= nums[i] <= n`

Each element in `nums` appears once or twice.

[LeetCode](https://leetcode.com/problems/find-all-duplicates-in-an-array/)

## [First Missing Positive](./first-missing-positive)

Given an unsorted integer array `nums`.
Return the smallest positive integer that is not present in `nums`.

You must implement an algorithm that runs in `O(n)` time and uses `O(1)` auxiliary space.

Example 1:

```
Input: nums = [1,2,0]
Output: 3
Explanation: The numbers in the range [1,2] are all in the array.
```

Example 2:

```
Input: nums = [3,4,-1,1]
Output: 2
Explanation: 1 is in the array but 2 is missing.
```

Example 3:

```
Input: nums = [7,8,9,11,12]
Output: 1
Explanation: The smallest positive integer 1 is missing.
```

Constraints:

- `1 <= nums.length <= 105`
- `231 <= nums[i] <= 231 - 1`

[LeetCode](https://leetcode.com/problems/first-missing-positive)

## [Excel Sheet Column Title](./excel-sheet-column-title)

Given an integer `columnNumber`, return its corresponding column title as it appears in an Excel sheet.

For example:

```
A -> 1
B -> 2
C -> 3
...
Z -> 26
AA -> 27
AB -> 28
...
```

Example 1:

```
Input: columnNumber = 1
Output: "A"
```

Example 2:

```
Input: columnNumber = 28
Output: "AB"
```

Example 3:

```
Input: columnNumber = 701
Output: "ZY"
```

Constraints:

- `1 <= columnNumber <= 231 - 1`

[LeetCode](https://leetcode.com/problems/excel-sheet-column-title/)

## [Permutation Sequence](./permutation-sequence)

The set `[1, 2, 3, ..., n]` contains a total of `n!` unique permutations.

By listing and labeling all the permutations in order, we get the following sequence for n = 3:

```
"123"
"132"
"213"
"231"
"312"
"321"
```

Given `n` and `k`, return the `k`th permutation sequence.

Example 1:

```
Input: n = 3, k = 3
Output: "213"
```

Example 2:

```
Input: n = 4, k = 9
Output: "2314"
```

Example 3:

```
Input: n = 3, k = 1
Output: "123"
```

Constraints:

- `1 <= n <= 9`
- `1 <= k <= n!`

## [Length of Longest Subarray With at Most K Frequency](./length-of-longest-subarray-with-at-most-k-frequency)

You are given an integer array `nums` and an integer `k`.
The frequency of an element `x` is the number of times it occurs in an array.
An array is called good if the frequency of each element in this array is less than or equal to `k`.
Return the length of the longest good subarray of `nums`.
A subarray is a contiguous non-empty sequence of elements within an array.

Example 1:

```
Input: nums = [1,2,3,1,2,3,1,2], k = 2
Output: 6
Explanation: The longest possible good subarray is [1,2,3,1,2,3] since the values 1, 2, and 3 occur at most twice in this subarray. Note that the subarrays [2,3,1,2,3,1] and [3,1,2,3,1,2] are also good.
It can be shown that there are no good subarrays with length more than 6.
```

Example 2:

```
Input: nums = [1,2,1,2,1,2,1,2], k = 1
Output: 2
Explanation: The longest possible good subarray is [1,2] since the values 1 and 2 occur at most once in this subarray. Note that the subarray [2,1] is also good.
It can be shown that there are no good subarrays with length more than 2.
```

Example 3:

```
Input: nums = [5,5,5,5,5,5,5], k = 4
Output: 4
Explanation: The longest possible good subarray is [5,5,5,5] since the value 5 occurs 4 times in this subarray.
It can be shown that there are no good subarrays with length more than 4.
```

Constraints:

- `1 <= nums.length <= 10^5`
- `1 <= nums[i] <= 10^9`
- `1 <= k <= nums.length`

## [Count Subarrays Where Max Element Appears at Least K Times](./count-subarrays-where-max-element-appears-at-least-k-times)

You are given an integer array `nums` and a positive integer `k`.
Return the number of subarrays where the maximum element of `nums` appears at least `k` times in that subarray.
A subarray is a contiguous sequence of elements within an array.

Example 1:

```
Input: nums = [1,3,2,3,3], k = 2
Output: 6
Explanation: The subarrays that contain the element 3 at least 2 times are: [1,3,2,3], [1,3,2,3,3], [3,2,3], [3,2,3,3], [2,3,3] and [3,3].
```

Example 2:

```
Input: nums = [1,4,2,1], k = 3
Output: 0
Explanation: No subarray contains the element 4 at least 3 times.
```

Constraints:

- `1 <= nums.length <= 10^5`
- `1 <= nums[i] <= 10^6`
- `1 <= k <= 105`
