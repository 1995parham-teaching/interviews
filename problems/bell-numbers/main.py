import functools


@functools.cache
def C(n: int, k: int) -> int:
    """
    the number of combinations of k objects from a set with n objects.
    """
    if k == 0:
        return 1
    if k == 1:
        return n
    if k == n:
        return 1
    return C(n - 1, k - 1) + C(n - 1, k)


@functools.cache
def bell(n: int) -> int:
    """
    recursive formulation for the bell number,
    which discusses the set containing the nth element.

    https://en.wikipedia.org/wiki/Bell_number
    """
    if n == 0:
        return 1
    count = 0
    for k in range(n):
        count += C(n - 1, k) * bell(k)
    return count


if __name__ == "__main__":
    assert C(10, 0) == 1
    assert C(5, 2) == 10
    assert C(3, 1) == 3
    assert C(3, 2) == 3
    assert C(4, 2) == 6

    assert bell(2) == 2
    assert bell(3) == 5
