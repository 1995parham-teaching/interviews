def num_of_islands(grid: list[list[str]]) -> int:
    n = len(grid)
    m = len(grid[0])

    islands = 0

    for i in range(n):
        for j in range(m):
            if grid[i][j] == "1":
                islands += 1
                queue: list[tuple[int, int]] = [(i, j)]
                grid[i][j] = "0"

                while len(queue) != 0:
                    r, c = queue.pop()
                    if r + 1 < n:
                        if grid[r + 1][c] == "1":
                            queue.append((r + 1, c))
                            grid[r + 1][c] = "0"
                    if r - 1 >= 0:
                        if grid[r - 1][c] == "1":
                            queue.append((r - 1, c))
                            grid[r - 1][c] = "0"
                    if c + 1 < m:
                        if grid[r][c + 1] == "1":
                            queue.append((r, c + 1))
                            grid[r][c + 1] = "0"
                    if c - 1 >= 0:
                        if grid[r][c - 1] == "1":
                            queue.append((r, c - 1))
                            grid[r][c - 1] = "0"

    return islands


if __name__ == "__main__":
    assert (
        num_of_islands(
            [
                ["1", "1", "1", "1", "0"],
                ["1", "1", "0", "1", "0"],
                ["1", "1", "0", "0", "0"],
                ["0", "0", "0", "0", "0"],
            ]
        )
        == 1
    )

    assert (
        num_of_islands(
            [
                ["1", "1", "0", "0", "0"],
                ["1", "1", "0", "0", "0"],
                ["0", "0", "1", "0", "0"],
                ["0", "0", "0", "1", "1"],
            ],
        )
        == 3
    )
