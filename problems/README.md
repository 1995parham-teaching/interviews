# Problems

These are some hands-on problem, and you can use a shared Google Doc to read and write code together. They don't require the interviewee to have an IDE or coding environment, and the main part of your discussion should be Algorithm.

## Levels

Every problem is tagged with a level. The level describes **which bar the problem is calibrated against**, not how clever it is. Picking a problem above the candidate's level does not give you a stronger signal, it just gives you a noisier one: you learn that someone struggled, but not whether they would struggle at the work you are actually hiring for.

### Level 1 — Frontend and SRE

Fundamentals. Arrays, strings, hash maps, straightforward iteration and direct simulation of the described process. No algorithm has to be *discovered*; the candidate has to translate a clear specification into correct code and handle the edge cases.

- **Budget:** 20–25 minutes.
- **You are looking for:** correct decomposition, edge cases (empty input, single element, ties), readable naming, and being able to say what the running time is.
- **Do not** penalise for not finding an asymptotically better solution. That is not the signal here.

### Level 2 — Backend engineers

Requires choosing the right data structure and justifying it. Stacks, tries, heaps, graphs on grids, sliding windows, recursion with state. There is a naive solution *and* a better one, and the gap between them is the interesting part of the conversation.

- **Budget:** 30–40 minutes.
- **You are looking for:** starting from the naive solution and improving it deliberately, correct complexity analysis of both, and choosing the data structure *before* writing code rather than discovering the need for it halfway through.
- **Do not** accept a correct answer without the complexity discussion — at this level the reasoning is worth more than the code.

### Level 3 — Seniors

The optimal solution depends on an invariant, a tight space constraint, or a design decision with no single right answer. Many of these problems ask the candidate to *design* something (an API, a data structure with a maintained invariant) rather than just compute an answer.

- **Budget:** 45–60 minutes, and use the layers — do not treat it as pass/fail.
- **You are looking for:** identifying the constraint that makes the naive solution unacceptable, proposing more than one approach and choosing between them out loud, and questioning the assumptions in the problem statement.
- **Do not** ask a Level 3 problem as the only coding problem in a loop. Pair it with a Level 2 so a candidate who does not crack it still produces usable signal.

## Interviewer Notes

Every problem below carries a collapsed **Interviewer notes** block with four things: the trap, the naive solution and what it costs, the intended solution, and one follow-up. It is collapsed so you can share the statement above it without showing the answers.

Read the notes before the session, not during it. If a candidate finds a route the notes do not mention, that is a fact about the notes, not about the candidate.

## Writing a New Problem

Use [Package Delivery](#package-delivery) and [Min By Column](#min-by-column) as the template. They are the two strongest problems here and neither came from a puzzle site: each starts from an ordinary task, specifies it precisely enough to code against, and then adds constraints that invalidate the obvious first answer. Package Delivery in particular escalates across four parts, with a final requirement that breaks whatever greedy rule was built in the third.

For how to design at each level, see [`docs/designing-interview-questions.typ`](../docs/designing-interview-questions.typ).

## Index by Level

Problems marked † have a solution in this repository but no problem statement in this file yet. Problems with no link have a statement here but no reference solution.

### Level 1

- [Diagonal Traverse](./diagonal-traverse) †
- [K Nearest Neighbor](./k-nearest-neighbour) — also in Python as [knn](./knn)
- [Longest Common Prefix](./longest_common_prefix) †
- [Merge Two Sorted Arrays](./merge) †
- [Min By Column](./min-by-column) — Part 1 only
- [Sort Integers by the Number of 1 Bits](./sort-integers-by-the-number-of-1-bits) †
- Happy Number

### Level 2

- [Bell Numbers](./bell-numbers) †
- [Coins](./coins)
- [Decode String](./decode-string)
- [Excel Sheet Column Title](./excel-sheet-column-title)
- [Find All Groups of Farmland](./find-all-groups-of-farmland)
- [Generate IPv4](./generate-ipv4) †
- [Generate Parentheses](./generate-parentheses)
- [Jump Game](./jump-game) †
- [Kadane's Algorithm](./kadanes-algorithm)
- [Length of Longest Subarray With at Most K Frequency](./length-of-longest-subarray-with-at-most-k-frequency)
- [Longest Palindromic Substring](./longest-palindrome) — duplicated in [longest-palindromic-substring](./longest-palindromic-substring)
- [Number of Islands](./number-of-islands)
- [Package Delivery](./package-delivery) — Parts 1 and 2
- [Prefix Search](./prefix-search)
- [Print N-bit Binary Numbers Having More 1s Than 0s](./print-n-bit-binary-numbers-having-more-1s-than-0s)
- [Search a 2D Matrix](./search-a-2d-matrix)
- [Search a 2D Matrix II](./search-a-2d-matrix-ii)
- [Shuffle](./shuffle)
- [Sort Colors](./sort-colors)
- K-th Biggest Number
- Rotate Image

### Level 3

- [Beautiful Nums](./beautiful-nums) †
- [Count Subarrays Where Max Element Appears at Least K Times](./count-subarrays-where-max-element-appears-at-least-k-times)
- [Find All Duplicates in an Array](./find-all-duplicates-in-an-array)
- [First Missing Positive](./first-missing-positive)
- [Longest Valid Parentheses](./longest-valid-parentheses)
- [Merge k Sorted Lists](./merge-k-sorted-lists)
- [Package Delivery](./package-delivery) — Parts 3 and 4
- [Permutation Sequence](./permutation-sequence)
- [Prefix and Suffix Search](./prefix-and-suffix-search)
- [Shuffler](./shuffler) †
- K-th Smallest Element in a Sorted Matrix
- Snappfood

One problem, [Bulb Switcher](#retired), has been retired. It is kept at the bottom of this file with the reason.

## [Decode String](./decode-string)

**Level 2**

Given an encoded string, return its decoded string.

The encoding rule is: `k[encoded_string]`, where the `encoded_string` inside the square brackets is being repeated exactly `k` times. Note that `k` is guaranteed to be a positive integer.

You may assume that the input string is always valid; there are no extra white spaces, square brackets are well-formed, etc. Furthermore, you may assume that the original data does not contain any digits and that digits are only for those repeat numbers, k. For example, there will not be input like `3a` or `2[4]`.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Nesting. `3[a2[c]]` cannot be handled by a regex or by a single counter — this is what separates a working answer from a plausible one. Multi-digit repeat counts like `12[a]` break naive character-at-a-time digit parsing.
- **Naive:** Recursive expansion or repeated rescanning of the string. Correct, and acceptable as a first answer.
- **Intended:** One stack holding `(count, prefix_so_far)`; push on `[`, pop and repeat on `]`. Linear in the length of the output.
- **Follow-up:** What would you do if the input could be malformed? What if the output were too large to hold in memory — could you stream it?

</details>

## [Find All Groups of Farmland](./find-all-groups-of-farmland)

**Level 2**

You are given a **0-indexed** `m x n` binary matrix `land` where a `0` represents a hectare of forested land and a `1` represents a hectare of farmland.

To keep the land organized, there are designated rectangular areas of hectares that consist **entirely** of farmland. These rectangular areas are called **groups**. No two groups are adjacent, meaning farmland in one group is not four-directionally adjacent to another farmland in a different group.

`land` can be represented by a coordinate system where the top left corner of `land` is `(0, 0)` and the bottom right corner of `land` is `(m-1, n-1)`. Find the coordinates of the top left and bottom right corner of each **group** of farmland. A **group** of farmland with a top left corner at `(r1, c1)` and a bottom right corner at `(r2, c2)` is represented by the 4-length array `[r1, c1, r2, c2]`.

Return a 2D array containing the 4-length arrays described above for each **group** of farmland in `land`. If there are no groups of farmland, return an empty array. You may return the answer in **any order**.

[LeetCode](https://leetcode.com/problems/find-all-groups-of-farmland/)

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The guarantees do the work. Groups are rectangles and are never adjacent, so no flood fill is needed — but the statement never says so outright. Watch whether they notice.
- **Naive:** BFS or DFS per group, tracking min and max coordinates. `O(mn)` and perfectly correct.
- **Intended:** Scan for a cell whose top and left neighbours are both `0` — that is a top-left corner — then walk right and down to find the opposite corner. `O(mn)` with no extra structure.
- **Follow-up:** What breaks if groups were allowed to touch? What if they were not rectangles?

</details>

## [Number of Islands](./number-of-islands)

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Whether they mutate the input grid to mark visited cells or keep a separate structure. Both are defensible; ask which they chose and whether a caller would mind.
- **Naive:** Recursive DFS. Note that a 300×300 all-land grid can overflow the call stack in some languages — worth raising if they do not.
- **Intended:** Iterative BFS or DFS, or union-find. `O(mn)`.
- **Follow-up:** The grid is too large for memory and arrives one row at a time. What changes? (This leads naturally to union-find over a two-row window.)

</details>

## [Shuffle](./shuffle)

**Level 2**

You have an array with n-items (A). We want to partition it into k-subarrays that each of them has n/k items, and each element of A appears precisely once. The order of these subarrays must not be the same as the A.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The statement is deliberately underspecified — it even asks "with duplication or without duplication?" out loud. The signal is whether they pin down what "not the same order" means before writing anything.
- **Naive:** Shuffle the array and deal it into chunks. As the reference solution notes in its own comment, this can still produce an ordering matching the original when values repeat.
- **Intended:** First define the requirement precisely, then deal round-robin so no subarray preserves adjacency from the input.
- **Follow-up:** How would you test a randomised function? What property would you assert?

</details>

## [Coins](./coins)

**Level 2**

We have `n` amount of money and our country have the following coins:

- coin-1
- coin-5
- coin-7
- coin-10

We want to have this money with minimum number of coins. What is the minimum? For example:

- 2 = 2 x coin-1
- 5 = 1 x coin-5
- 6 = 1 x coin-5 + 1 x coin-1

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Greedy fails on this coin set, which is the whole reason the set is `{1, 5, 7, 10}` and not `{1, 5, 10, 25}`. If they propose greedy, hand them `14`: greedy pays `10+1+1+1+1` for five coins, the answer is `7+7` for two.
- **Naive:** Greedy — wrong here. Or exhaustive search over combinations.
- **Intended:** Bottom-up DP over amounts: `dp[a] = 1 + min(dp[a - c])` for each coin `c <= a`. `O(n · |coins|)`.
- **Follow-up:** Return the coins used, not just the count. What if the coin set were arbitrary — how would you know greedy is safe?

</details>

## Happy Number

**Level 1**

Write an algorithm to determine if a number n is "happy".

A happy number is a number defined by the following process: Starting with any positive integer, replace the number by the sum of the squares of its digits, and repeat the process until the number equals 1 (where it will stay), or it loops endlessly in a cycle which does not include 1. Those numbers for which this process ends in 1 are happy numbers.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Knowing when to stop. No iteration bound is given, so termination has to come from cycle detection — that is the actual problem, not the digit arithmetic.
- **Naive:** A set of previously seen values. `O(log n)` space and entirely acceptable at this level.
- **Intended:** The same, or Floyd's tortoise and hare for constant space.
- **Follow-up:** Why can the values not grow without bound? (For any 3-digit or larger number the next value is smaller, so the sequence is eventually trapped in a small range.)

</details>

## Rotate Image

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** "In place." A candidate who allocates a second matrix has answered a different question. Let them finish, then hold them to the constraint.
- **Naive:** Build a new matrix by index mapping. `O(n²)` extra space.
- **Intended:** Transpose, then reverse each row. Or rotate four cells at a time, ring by ring.
- **Follow-up:** Counter-clockwise? What about a non-square `m x n` matrix — why can that one not be done in place?

</details>

## Snappfood

**Level 3**

We have motorcycles and restaurants. Motorcycles deliver foods to peoples from restaurants. How we can schedule this delivery process?

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** This is a seed, not a finished question — it has no constraints on purpose. Do not ask it without deciding beforehand which of the layers below you intend to reach.
- **Naive:** Assign each order to the nearest free courier. Fine as an opening move; falls apart under every layer below.
- **Intended:** There is no single intended answer. Expect a queue of orders, couriers with locations and state, and some batching rule. The design conversation is the deliverable.
- **Follow-up:** Layers, in order: (1) a courier can carry two orders if they are close — how do you decide? (2) food goes cold, so lateness is not linear in time. (3) demand spikes and there are not enough couriers — what do you drop? (4) what would you measure to know the scheduler is working?

</details>

## [Search a 2D Matrix](./search-a-2d-matrix)

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The row-boundary property means the matrix is one sorted array in disguise. The slip is in the index mapping, `mid / n` and `mid % n`, not in the binary search.
- **Naive:** Scan every cell, `O(mn)` — ruled out by the stated `O(log(m·n))` requirement, which is one of the few explicit complexity requirements in this bank. Point at it if they miss it.
- **Intended:** A single binary search over `[0, m·n)`, mapping the index back to a cell.
- **Follow-up:** Drop the guarantee that each row starts above the previous row's end. What now? (That is the next problem in this file.)

</details>

## [Search a 2D Matrix II](./search-a-2d-matrix-ii)

**Level 2**

Write an efficient algorithm that searches for a value `target` in an `m x n` integer matrix. This matrix has the following properties:

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The elegant answer is not the only acceptable one. Binary search per row is `O(m log n)` and fine — ask them to compare it with the staircase rather than demanding the staircase.
- **Naive:** Full scan `O(mn)`, or binary search each row `O(m log n)`.
- **Intended:** Start at the top-right corner: move left when the value is too large, down when too small. `O(m + n)`.
- **Follow-up:** Why does starting at the top-left corner not work? What does that tell you about which corners are usable?

</details>

## [Longest Palindromic Substring](./longest-palindrome)

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Even-length palindromes. Anyone expanding around single characters only will pass `"babad"` and silently fail `"cbbd"` — which is exactly why both examples are in the statement.
- **Naive:** Check every substring for palindromicity, `O(n³)`. Or DP over `is_pal[i][j]`, `O(n²)` time and `O(n²)` space.
- **Intended:** Expand around each of the `2n - 1` centres. `O(n²)` time, `O(1)` space.
- **Follow-up:** Manacher's algorithm gets `O(n)` — ask them to describe the idea, not to code it. Note this problem has two solution directories in this repo (`longest-palindrome` and `longest-palindromic-substring`); they are the same problem.

</details>

## K-th Smallest Element in a Sorted Matrix

**Level 3**

Given an `n x n` matrix where each row and each column is sorted in ascending order, return the `k`-th smallest element **in sorted order** — not the `k`-th distinct element.

Note that the rows and columns being sorted does **not** mean the flattened matrix is sorted.

Example:

```text
Input: matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8
Output: 13
Explanation: sorted, the elements are [1,5,9,10,11,12,13,13,15] and the 8th is 13.
```

Constraints:

- `n == matrix.length == matrix[i].length`
- `1 <= n <= 300`
- `1 <= k <= n^2`
- Each row and each column is sorted in ascending order.

[LeetCode](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/)

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Two of them. `k`-th smallest counts duplicates, and sorted rows plus sorted columns does *not* make the flattened matrix sorted — the example is chosen so that assuming otherwise gives the wrong answer.
- **Naive:** Flatten and sort, `O(n² log n)`. Correct, and a fine starting point to improve from.
- **Intended:** Either a min-heap seeded with the first column, popped `k` times — `O(k log n)`; or binary search on the *value* range, counting elements `<= mid` with the staircase walk — `O(n log(max - min))`.
- **Follow-up:** Which of the two is better when `k` is close to `n²`? What if the matrix did not fit in memory?

</details>

## [Merge k Sorted Lists](./merge-k-sorted-lists)

**Level 3**

You are given an array of `k` linked-lists `lists`, each linked-list is sorted in ascending order. Merge all the linked-lists into one sorted linked-list and return it.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The empty cases. `[]` and `[[]]` are both in the examples deliberately — an empty list of lists and a list containing an empty list are different inputs.
- **Naive:** Concatenate everything and sort, `O(N log N)`. Or merge one list at a time into an accumulator, `O(kN)` — ask them why that second one is worse than it looks.
- **Intended:** A min-heap holding the current head of each list, `O(N log k)`. Or divide and conquer with pairwise merges — same bound, no heap needed.
- **Follow-up:** The lists are streams too large to hold in memory. What still works?

</details>

## K Nearest Neighbor

**Level 1** · reference solutions: [Go](./k-nearest-neighbour), [Python](./knn)

We have `n` points and one reference point. Each point has `x` and `y` coordinates. We want to find `k` the nearest points to the reference point.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Ties. The example deliberately shows two acceptable answers, so a candidate who assumes the result is unique has missed something the statement is showing them.
- **Naive:** Sort all points by distance and take the first `k`. `O(n log n)` — completely acceptable at this level, do not push for better.
- **Intended:** The same. If they want to go further, a bounded max-heap of size `k` gives `O(n log k)`.
- **Follow-up:** Can you skip the square root, and why is that safe? What changes if the points arrive as a stream, or in many more than two dimensions?

</details>

## [Sort Colors](./sort-colors)

**Level 2**

Given an array `nums` with `n` objects colored red, white, or blue, sort them in-place so that objects of the same color are adjacent, with the colors in the order red, white, and blue.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The one-pass follow-up. Counting each value and rewriting the array is correct and uses constant space — accept it first, then ask for one pass.
- **Naive:** Count the `0`s, `1`s and `2`s, then overwrite. Two passes, `O(1)` space.
- **Intended:** Dutch national flag: three pointers, one pass. The subtle part is not advancing the middle pointer after swapping with the high pointer — ask them to justify it.
- **Follow-up:** What if there were four colours? What if there were `k`, given as a parameter?

</details>

## [Generate Parentheses](./generate-parentheses)

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Pruning. Generating all `2^(2n)` strings and filtering is the obvious idea and is exponentially wasteful; the insight is never placing `)` unless a `(` is still unmatched.
- **Naive:** Generate every string of length `2n` and validate each one, `O(2^(2n) · n)`.
- **Intended:** Backtracking with counts of open and closed brackets. The number of results is the `n`-th Catalan number.
- **Follow-up:** Count the results without generating them. Return only the `k`-th in lexicographic order — how does that change the approach?

</details>

## [Longest Valid Parentheses](./longest-valid-parentheses)

**Level 3**

Given a string containing just the characters '(' and ')', return the length of the longest valid (well-formed) parentheses substring.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The valid substring can start anywhere, so greedy counting of matched pairs fails. Give them `")()())"` early — it breaks most first attempts.
- **Naive:** Check every substring for validity, `O(n³)`, or `O(n²)` with incremental counting.
- **Intended:** A stack of indices seeded with `-1`, or DP over `dp[i]`, or two passes with left and right counters for `O(1)` space. All `O(n)`.
- **Follow-up:** Return the substring itself rather than its length. What changes with more than one kind of bracket?

</details>

## Min By Column

**Level 1**

### Part 1 of 2

Imagine that we are working with a simple database. Each row associates column names (strings) with integer values. Here's a table with three rows:

```text
a b c d
1 0 0 0
0 2 3 0
0 0 0 4
```

We might choose to represent a database table in JSON, as an array of objects. For example, the previous table could be written as:

```json
[
  { "a": 1, "b": 0, "c": 0, "d": 0 },
  { "a": 0, "b": 2, "c": 3, "d": 0 },
  { "a": 0, "b": 0, "c": 0, "d": 4 }
]
```

Write a function, `min_by_column`, that takes a database table (as above), along with a column name, and returns the row that contains the minimum value for the given column. If a row doesn't have any value for the column, your function should behave as though the value for that column was zero.

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

In Part 1 you may have noticed that it's possible for two rows to be "tied", meaning that either would be an acceptable return value from `min_by_column`.

Consider:

```python
table_4 = [
{"a": 1, "b": 2},
{"a": 1, "b": 3},
{"a": 1, "b": 4}
]
assert min_by_column(table_4, "a") == '???'
```

Since all three rows have the same value for column "a", all three rows are acceptable candidates to be returned by `min_by_column(table, "a")`.

In these cases, it would be nice if users could specify additional columns (e.g. "b") to use as tie-breakers. A tie-breaker would only apply in cases where multiple rows share the same minimum value. In `table_4` above, the row `{"a": 1, "b": 2}` is tied for the smallest "a" value (1) and of all the tied candidates, it has the smallest "b" value (2). If two records had equal values for "a" and also for "b" then another tie-breaker (e.g. "c") could be used. When records are tied with respect to all columns, any of the tied records may be considered the minimum.

Write a function `min_by_columns` that takes a database table and an ordered list of column names, and returns the row with the minimum column values using the tie-breaking logic above. Refactor `min_by_column` to use `min_by_columns` to produce its result.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Part 2 asks them to *refactor* `min_by_column` to use `min_by_columns`, not to write a second independent function. This is the real signal in the question and it is easy to miss.
- **Naive:** A single pass with a comparison. Nothing more is needed at this level.
- **Intended:** Part 2 is a comparison over a tuple of column values, in the order given.
- **Follow-up:** Descending order. A column name that appears in no row. What should happen if the list of tie-breaker columns is empty?

</details>

## Prefix Search

**Level 2**

We have a database, and we'd like it to support these operations:

- `insert(word)`: Inset a word to database
- `look up(prefix)`: Return all the words starting with the given prefix
- `delete(prefix)`: Delete all the words starting with the given prefix
- `count(prefix)`Count the number of words starting with the given prefix

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** `delete(prefix)` is the interesting operation and the reason this beats a stock trie question — most people have rehearsed insert and lookup but not subtree removal with counts kept correct.
- **Naive:** A list of words scanned per query, `O(N · M)`. Fine as an opening answer.
- **Intended:** A trie with a word count stored at every node, so `count(prefix)` is `O(|prefix|)` and `delete(prefix)` detaches a subtree.
- **Follow-up:** After a delete, how do the counts on the ancestor nodes stay right? What about case-insensitivity, or ranking the results?

</details>

## Package Delivery

**Level 2 (Parts 1-2) / Level 3 (Parts 3-4)**

### Part 1

Write a Delivery class (or object) that represents a delivery with a destination and distance. Deliveries require different sensors, depending on their distance.

Add a method, `getNeededSensors`, that returns a mapping of sensor name to the count of that sensor needed to complete the delivery according to these rules:

- If distance < 10 miles, require 1 gps and 1 temp sensor.
- If 10 <= distance < 100 miles require 1 gps, 2 temp, and 1 weight sensor.
- If distance >= 100 miles require 2 gps, 4 temp, and 2 weight sensors.

### Part 2

Write a Scheduler class (or object) that represents a daily delivery scheduler with a set of available sensors. Add a method, `scheduleDeliveries`, that given a list of deliveries as an argument, returns a list of deliveries that can be made that day. Assume all deliveries will be leaving at the same time every day, so sensors can only be used once.

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

We get paid a flat fee for all deliveries. Modify the `scheduleDeliveries` function to maximize the number of deliveries that will be made in a day.

Our previous test:

```python
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryC]
```

Should now return:

```python
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryA, deliveryB]
```

### Part 4

We recently purchased a new type of sensor, `doorSensor`. A `doorSensor` can be used in place of 1 weight sensor or in place of 2 temperature sensors at any time. Modify our existing functions to maximize the day's deliveries with the new sensor.

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Part 4 is the point. The substitutable `doorSensor` invalidates whatever greedy rule they built in Part 3, and watching them repair the design is worth more than the code.
- **Naive:** Parts 1–2 are modelling and first-fit — note that Part 2's expected outputs are order-dependent on purpose.
- **Intended:** Part 3 is a sort plus greedy, or a small knapsack. Part 4 is genuinely combinatorial; a correct exhaustive answer with a stated cost beats a confident wrong greedy.
- **Follow-up:** What would you measure to know your scheduler is actually good? This is the best question in this bank — use it as the template when writing new ones.

</details>

## K-th biggest number

**Level 2**

Given an unsorted array of `n` integers and an integer `k`, return the `k`-th largest element. Duplicates count: in `[3, 2, 3, 1]` with `k = 2` the answer is `3`, not `2`.

Do not sort the whole array.

Constraints:

- `1 <= n <= 10^5`
- `1 <= k <= n`
- `-10^9 <= nums[i] <= 10^9`

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Duplicates count. In `[3, 2, 3, 1]` with `k = 2` the answer is `3`, not `2` — a candidate deduplicating first has answered a different question.
- **Naive:** Sort and index, `O(n log n)` — ruled out by the statement. Or scan for the maximum `k` times, `O(kn)`, which is fine only for tiny `k`.
- **Intended:** Quickselect: `O(n)` expected, for any `k`. Or a bounded min-heap of size `k` for `O(n log k)`, which is the better choice when the data streams.
- **Follow-up:** What is quickselect's worst case, and how do you avoid it? Which would you pick if the numbers arrived one at a time and you could not store them all?

</details>

## [Prefix and Suffix Search](./prefix-and-suffix-search/)

**Level 3**

Design a `WordFilter` that is initialised with a list of words and supports one query:

- `f(prefix, suffix)`: return the **largest index** of a word that has both the given prefix and the given suffix. If no such word exists, return `-1`.

Example:

```text
Input:  words = ["apple"]
        f("a", "e")  -> 0
        f("b", "")   -> -1
```

Constraints:

- `1 <= words.length <= 10^4`
- `1 <= words[i].length <= 7`
- Up to `10^4` calls to `f`.
- `prefix` and `suffix` may be empty.

[LeetCode](https://leetcode.com/problems/prefix-and-suffix-search/)

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** It asks for the **largest** index, so a later word must overwrite an earlier one with the same prefix and suffix. Building the structure without that is the common failure.
- **Naive:** Scan every word per query, `O(N · M)` each. At `10^4` words and `10^4` queries that is far too slow — point at the constraints.
- **Intended:** Insert every `suffix + separator + word` combination into one trie, `O(M²)` per word, so each query is a single walk. Two separate tries plus set intersection also works and is worth comparing.
- **Follow-up:** What does the combined-key trie cost in memory? Why does the word length limit of 7 make it viable, and what would you do if words were 1000 characters?

</details>

## [Print N-bit binary numbers having more 1s than 0s](./print-n-bit-binary-numbers-having-more-1s-than-0s/)

**Level 2**

Given a positive integer `n`. Your task is to generate a string list of all n-bit binary numbers where, for any prefix of the number, there are more or an equal number of 1's than 0's. The numbers should be sorted in decreasing order of magnitude.

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

User Task: Your task is to complete the function `NBitBinary()` which takes a single number as input `n` and returns the list of strings in **decreasing** order. You need not take any input or print anything.

```python
class Solution:
  def NBitBinary(self, n):
    pass
```

Expected Time Complexity: `O(|2n|)` Expected Auxiliary Space: `O(2n)`

Constraints:

```
1 <= n <= 15
```

[GeeksForGeeks](https://www.geeksforgeeks.org/problems/print-n-bit-binary-numbers-having-more-1s-than-0s0252)

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The condition applies to **every prefix**, not just the whole string. Also, the required decreasing order falls out for free if you always try `1` before `0`.
- **Naive:** Generate all `2^n` strings and check each prefix of each, `O(2^n · n)`.
- **Intended:** Backtracking carrying the counts of ones and zeros; only place a `0` while ones exceed zeros.
- **Follow-up:** Count the valid strings without generating them — these are the ballot numbers, and they are Catalan numbers when `n` is even.

</details>

## [Kadane's Algorithm](./kadanes-algorithm/)

**Level 2**

Given an array `Arr[]` of `N` integers. Find the contiguous sub-array (containing at least one number) which has the maximum sum and return its sum.

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

Your Task: You don't need to read input or print anything. The task is to complete the function `maxSubarraySum()` which takes `Arr[]` and `N` as input parameters and returns the sum of subarray with maximum sum.

Expected Time Complexity: `O(N)` Expected Auxiliary Space: `O(1)`

Constraints:

```
1 ≤ N ≤ 10^6
-10^7 ≤ A[i] ≤ 10^7
```

[GeeksForGeeks](https://www.geeksforgeeks.org/problems/kadanes-algorithm-1587115620/)

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The all-negative case. Initialising the best sum to `0` returns `0` instead of `-1` for `{-1,-2,-3,-4}` — the second example exists precisely to catch this, so do not skip it.
- **Naive:** Every subarray, `O(n²)`.
- **Intended:** Carry a running sum, resetting it to the current element whenever the element alone is larger. `O(n)` time, `O(1)` space.
- **Follow-up:** Return the start and end indices, not just the sum. What if the array were circular?

</details>

## [Find All Duplicates in an Array](./find-all-duplicates-in-an-array)

**Level 3**

Given an integer array `nums` of length `n` where all the integers of `nums` are in the range `[1, n]` and each integer appears once or twice, return an array of all the integers that appears twice.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** "Constant extra space" is the entire problem. The output array does not count against it, but a set of seen values does — say so if they ask, because it is a fair question.
- **Naive:** A hash set, `O(n)` time and `O(n)` space. Correct, but does not meet the stated bound.
- **Intended:** Use the value range `[1, n]` as indices into the array itself: negate the entry at `abs(v) - 1` and treat an already-negative entry as a second sighting. Mutates the input.
- **Follow-up:** Restore the array afterwards. What would change if the values were in `[0, n]` instead?

</details>

## [First Missing Positive](./first-missing-positive)

**Level 3**

Given an unsorted integer array `nums`. Return the smallest positive integer that is not present in `nums`.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** The answer always lies in `[1, n+1]`, and that bound is what makes constant space possible. Without noticing it, the constraint looks impossible. (The constraint line below is mangled — the value range should read `-2^31` to `2^31 - 1`.)
- **Naive:** Sort then scan, `O(n log n)`. Or a hash set, `O(n)` time but `O(n)` space.
- **Intended:** Cyclic sort: repeatedly place each value `v` in `[1, n]` at index `v - 1`, then scan for the first index whose value is wrong.
- **Follow-up:** What if you were not allowed to modify the input at all? What does that cost you?

</details>

## [Excel Sheet Column Title](./excel-sheet-column-title)

**Level 2**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** Bijective base-26. There is no zero digit, so you must subtract 1 before each modulo. This off-by-one *is* the question — check them on `26` (`"Z"`) and `701` (`"ZY"`), not on `28`.
- **Naive:** There is no meaningfully worse approach here; the question is short by design.
- **Intended:** Loop while `n > 0`: decrement `n`, take `n % 26` as the letter, divide by 26, then reverse.
- **Follow-up:** Write the inverse, title to number. Where does this show up for real? (Spreadsheet addressing, and short human-readable IDs.)

</details>

## [Permutation Sequence](./permutation-sequence)

**Level 3**

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** `k` is 1-indexed. Converting to 0-indexed before doing the factorial arithmetic is where most attempts break, and the error is invisible for `k = 1`.
- **Naive:** Generate permutations in lexicographic order until reaching the `k`-th. At `n = 9` that is up to `362880` steps — it terminates, so push on why it is still the wrong idea.
- **Intended:** The factorial number system: the first digit is index `(k-1) / (n-1)!` of the remaining values; remove it and recurse on the remainder.
- **Follow-up:** Write the inverse — given a permutation, return its rank.

</details>

## [Length of Longest Subarray With at Most K Frequency](./length-of-longest-subarray-with-at-most-k-frequency)

**Level 2**

You are given an integer array `nums` and an integer `k`. The frequency of an element `x` is the number of times it occurs in an array. An array is called good if the frequency of each element in this array is less than or equal to `k`. Return the length of the longest good subarray of `nums`. A subarray is a contiguous non-empty sequence of elements within an array.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** A single counter is not enough; the window must shrink while the *offending* element exceeds `k`, which needs a frequency map. Watch for shrinking on the wrong condition.
- **Naive:** Every subarray with a frequency check, `O(n²)` at best.
- **Intended:** Sliding window with a frequency map: extend right, and while `freq[nums[right]] > k` advance left. `O(n)`.
- **Follow-up:** Return the subarray rather than its length. What should happen if `k` were `0`?

</details>

## [Count Subarrays Where Max Element Appears at Least K Times](./count-subarrays-where-max-element-appears-at-least-k-times)

**Level 3**

You are given an integer array `nums` and a positive integer `k`. Return the number of subarrays where the maximum element of `nums` appears at least `k` times in that subarray. A subarray is a contiguous sequence of elements within an array.

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

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** It counts *subarrays*, not a length — so the sliding-window habit of tracking a best size does not apply. The leap is: for each right endpoint, add the number of valid left endpoints.
- **Naive:** Enumerate every subarray and count occurrences, `O(n²)`.
- **Intended:** Find the global maximum first, then slide a window; whenever it contains `k` occurrences, add `left + 1` to the running total. `O(n)`.
- **Follow-up:** What if it asked for at least `k` occurrences of *any* element rather than the maximum?

</details>

## Retired

Kept for reference. Do not ask these.

### Bulb Switcher

There are n bulbs that are initially off. You first turn on all the bulbs. Then, you turn off every second bulb. On the third round, you toggle every third bulb (turning on if it's off or turning off if it's on). For the i-th round, you toggle every i bulb. For the n-th round, you only toggle the last bulb. Find how many bulbs are on after n rounds.

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

**Why it is retired.** The answer is `floor(sqrt(n))`, because a bulb ends up on only when its index has an odd number of divisors, which happens only for perfect squares. There is no intermediate step between not seeing that and seeing it — the candidate either has the insight or produces nothing, which is the single-flash anti-pattern the design guide rules out. It is a good puzzle and a bad interview question.
