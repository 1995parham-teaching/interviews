# Code Session

These questions are hard to implement, so they require your candidate to have coding environment at hands, and you need to check edge cases or even ask him/her to write some test cases.

## Loan

Design a loan system which gives people loan, and they can pay it back. They must pay back them in time, and they can have more than one loan at the time.

## Record Appender

There is a data file as below. We want to read it and then insert it into database.

data.txt:

```
John, Doe, 0939 123 1234
Jane, Doe, 0399 123 1234
```

- API for retrieving data from database
- Improve performance of Insert phase to support thousands of records

P.S. [Solution](https://github.com/1995parham-teaching/record-appender)

## Hangman

Implement Hangman!

```
> Computer: _ _ _ _ _

< Player: P

> Computer _ _ _ _ _ (player lose some score)

< Player: E

> Computer E _ _ _ E

< Player: L

> Computer: E L _ _ E

< Player: H

> Computer E L _ H E

< Player A

> Computer E L A H E (player win)
```

If the player runs out of the score he or she will die.

## String to Integer (`atoi`)

Implement the `myAtoi(string s)` function, which converts a string to a 32-bit signed integer (similar to C/C++'s `atoi` function).

The algorithm for myAtoi(string s) is as follows:

1. Read in and ignore any leading whitespace.
2. Check if the next character (if not already at the end of the string) is `'-'` or `'+'`. Read this character in if it is either. This determines if the final result is negative or positive respectively. Assume the result is positive if neither is present.
3. Read in next the characters until the next non-digit character or the end of the input is reached. The rest of the string is ignored.
4. Convert these digits into an integer (i.e. `"123" -> 123`, `"0032" -> 32`). If no digits were read, then the integer is `0`. Change the sign as necessary (from step 2).
5. If the integer is out of the 32-bit signed integer range `[-2**31, 2**31 - 1]`, then clamp the integer so that it remains in the range. Specifically, integers less than `-2**31` should be clamped to `-2**31`, and integers greater than `2**31 - 1` should be clamped to `2**31 - 1`.
6. Return the integer as the final result.

Note:

- Only the space character `' '` is considered a whitespace character.
- Do not ignore any characters other than the leading whitespace or the rest of the string after the digits.

Example 1:

```
Input: s = "42"
Output: 42
Explanation: The underlined characters are what is read in, the caret is the current reader position.
Step 1: "42" (no characters read because there is no leading whitespace)
         ^
Step 2: "42" (no characters read because there is neither a '-' nor '+')
         ^
Step 3: "42" ("42" is read in)
           ^
The parsed integer is 42.
Since 42 is in the range [-231, 231 - 1], the final result is 42.
```

Example 2:

```
Input: s = "   -42"
Output: -42
Explanation:
Step 1: "   -42" (leading whitespace is read and ignored)
            ^
Step 2: "   -42" ('-' is read, so the result should be negative)
             ^
Step 3: "   -42" ("42" is read in)
               ^
The parsed integer is -42.
Since -42 is in the range [-231, 231 - 1], the final result is -42.
```

Example 3:

```
Input: s = "4193 with words"
Output: 4193
Explanation:
Step 1: "4193 with words" (no characters read because there is no leading whitespace)
         ^
Step 2: "4193 with words" (no characters read because there is neither a '-' nor '+')
         ^
Step 3: "4193 with words" ("4193" is read in; reading stops because the next character is a non-digit)
             ^
The parsed integer is 4193.
Since 4193 is in the range [-231, 231 - 1], the final result is 4193.
```

Constraints:

- `0 <= s.length <= 200`
- `s` consists of English letters (lower-case and upper-case), digits (`0-9`), `' '`, `'+'`, `'-'`, and `'.'`.

[LeetCode](https://leetcode.com/problems/string-to-integer-atoi/)

## LRU Cache

Write an `LRUCache` class that is constructed with a fixed `capacity` and supports the two following operations in O(1):

- `get(key)`: return the value stored for `key`, and report a miss out of band — either raise `KeyError`, or accept a `default` to return instead.
- `put(key, value)`: store or update `value` for `key`. When the cache is full, inserting a new key must evict the least recently used entry.

Both `get` and a `put` on an existing key count as a use, and make that key the most recently used one.

Do not signal a miss with an in-band value such as `-1` or `None`. Both are perfectly good things to *store*, so a cache that answers `None` cannot tell "absent" from "present and empty" — and that is a bug the caller gets to debug, not the cache. Asking the candidate what their `get` returns for a cached `None` is a cheap way to find out whether they design interfaces or just satisfy examples.

```python
cache = LRUCache(2)
cache.put(1, 1)
cache.put(2, 2)
cache.get(1, -1)    # -> 1
cache.put(3, 3)     # evicts key 2
cache.get(2, -1)    # -> -1
cache.put(4, 4)     # evicts key 1
cache.get(1, -1)    # -> -1
cache.get(3, -1)    # -> 3
cache.get(4, -1)    # -> 4
cache[2]            # KeyError
```

Constraints:

- `1 <= capacity <= 10^5`
- Both operations must run in O(1) on average.
- Extra memory for auxiliary structures (nodes plus a hash map) is allowed, but the memory has to stay bounded by the capacity — no entry may outlive its eviction.

Follow-ups, in the order they are worth asking:

1. Make it thread-safe.
2. Add a per-entry TTL.
3. Generalise it to LFU (Least Frequently Used) instead of LRU.
4. Bound the cache by total item size rather than item count, so a single `put` may evict several entries.

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** `put` on a key that is **already present**. It is neither an insert nor a plain update: the existing node has to be moved, not appended, and no eviction should happen. Most broken implementations pass the LeetCode example and fail here — either leaking a duplicate node into the list, or evicting an entry while the cache is not actually full. Ask for this case explicitly if their own tests do not cover it.
- **Naive:** a single list or array — `get` scans it, so both operations are O(n). Also common: a dict plus a timestamp per entry, where eviction has to scan every entry for the minimum, again O(n). A `list.remove(key)` on a Python list is O(n) too, even next to a dict.
- **Intended:** a doubly linked list ordered most-recently-used to least-recently-used, plus a hash map from key to *node*. The node must carry its own key, otherwise eviction cannot find the dict entry to delete. Sentinel head and tail nodes remove every `if node.prev is None` branch — a candidate who reaches for them has done this before.
- **Edge cases to push on:** capacity 1; `put` on an existing key (above); a cached value that is legitimately `None` or `-1`, which is what the miss signalling above is there to survive; eviction from a list that just became empty; `capacity` of zero or less, which should be rejected in the constructor rather than crashing on the first eviction.
- **Follow-up:** which of the four extensions above fits the time left. TTL is the most productive one — it forces the question of *when* expired entries are removed (lazily on read, or by a sweeper), and whether an expired-but-present entry still counts as used.

</details>

P.S. [Solution](./cache-lru)
