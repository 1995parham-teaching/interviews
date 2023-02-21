# Problems

## [Shuffle](./shuffle)

You have an array with n-items (A).
We want to partition it into k-subarrays that each of them has n/k items, and each element of A appears precisely once.
The order of these subarrays must not be the same as the A.

**We know that: n % k == 0**

- With duplication or without duplication?

For example:

A = {1, 2, 3, 4}

k = 2

we don't accept the following solution:
A1 = {1, 2}
A2 = {3, 4}

but we accept the following solution:
A1 = {1, 3}
A2 = {2, 4}

## [Coins](./coins)

We have n amount of money and our country have the following coins:

coin-1 coin-5 coin-7 coin-10

We want to have this money with minimum number of coins. what is the minimum? for example:

- 2 = 2 \* coin-1
- 5 = 1 \* coin-5
- 6 = 1 _coin-5 + 1_ coin-1

## Bulb Switcher

There are n bulbs that are initially off. You first turn on all the bulbs.
Then, you turn off every second bulb. On the third round,
you toggle every third bulb (turning on if it's off or turning off if it's on). For the i-th round, you toggle every i bulb. For the n-th round, you only toggle the last bulb. Find how many bulbs are on after n rounds.

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

## Binary Search

Consider an _m x n_ matrix that each of its rows is sorted.
Write an efficient algorithm to find a given number into it.
