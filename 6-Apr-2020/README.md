You have an array with n-items (A).
We want to partition it into k-subarrays that each of them has n/k item, and each element of A appears precisely once.
The order of these subarrays must not be the same as the A.

**we know that: n % k == 0**

with duplication or without duplication?

For example:

A = {1, 2, 3, 4}

k = 2

we don't accept the following solution:
A1 = {1, 2}
A2 = {3, 4}

but we accept the following solution:
A1 = {1, 3}
A2 = {2, 4}
