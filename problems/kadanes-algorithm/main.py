def maximum_sequence(arr: list[int]) -> int:
    n = len(arr)
    if n == 0:
        return 0

    # d[i] contains the largest continues sequence that ends
    # with arr[i]
    d = [0 for _ in range(n)]
    d[0] = arr[0]
    for i in range(1, n):
        d[i] = max(d[i - 1] + arr[i], arr[i])

    max_sum = d[0]
    for i in range(n):
        max_sum = max(max_sum, d[i])
    return max_sum


if __name__ == "__main__":
    assert maximum_sequence([1, 2, 3, -2, 5]) == 9
    assert maximum_sequence([-1, -2, -3, -4]) == -1
