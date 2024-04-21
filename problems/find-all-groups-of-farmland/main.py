def find_farmland(land: list[list[int]]) -> list[list[int]]:
    n = len(land)
    m = len(land[0])

    farms = []

    for r in range(n):
        for c in range(m):
            if land[r][c] == 1:
                length = 1
                org = (r, c)
                dst = (r, c)

                land[r][c] = 0

                while True:
                    i, j = org
                    if (
                        i + length < n
                        and j + length < m
                        and {land[i + k][j + length] for k in range(length + 1)} == {1}
                        and {land[i + length][j + k] for k in range(length + 1)} == {1}
                    ):
                        dst = (i + length, j + length)

                        for k in range(length + 1):
                            land[i + k][j + length] = 0

                        for k in range(length + 1):
                            land[i + length][j + k] = 0

                        length += 1
                    else:
                        break

                while True:
                    i, _ = org
                    _, j = dst

                    if j + 1 < m and {land[i + k][j + 1] for k in range(length)} == {1}:
                        for k in range(length):
                            land[i + k][j + 1] = 0

                        dst = (dst[0], dst[1] + 1)
                    else:
                        break

                while True:
                    _, j = org
                    i, _ = dst

                    if i + 1 < n and {land[i + 1][j + k] for k in range(length)} == {1}:
                        for k in range(length):
                            land[i + 1][j + k] = 0

                        dst = (dst[0] + 1, dst[1])
                    else:
                        break

                farms.append(
                    [
                        *org,
                        *dst,
                    ]
                )

    return farms


if __name__ == "__main__":
    assert find_farmland([[1, 0, 0], [0, 1, 1], [0, 1, 1]]) == [
        [0, 0, 0, 0],
        [1, 1, 2, 2],
    ], find_farmland([[1, 0, 0], [0, 1, 1], [0, 1, 1]])

    assert find_farmland([[1, 1], [1, 1]]) == [[0, 0, 1, 1]], find_farmland(
        [[1, 1], [1, 1]]
    )

    assert find_farmland([[1, 1, 0, 0, 0, 1], [1, 1, 0, 0, 0, 0]]) == [
        [0, 0, 1, 1],
        [0, 5, 0, 5],
    ], find_farmland([[1, 1, 0, 0, 0, 1], [1, 1, 0, 0, 0, 0]])
