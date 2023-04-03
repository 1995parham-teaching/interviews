class Solution:
    def search_matrix(self, matrix: list[list[int]], target: int) -> bool:
        row = self.find_row(0, len(matrix), matrix, target)
        if row == -1:
            return False
        return self.find_in_row(0, len(matrix[row]), matrix[row], target)

    @staticmethod
    def find_row(
        first_row: int, last_row: int, matrix: list[list[int]], target: int
    ) -> int:
        """
        use binary search between first elements of rows to find
        a row that may contains the given target.
        """
        index = (first_row + last_row) // 2

        if last_row < first_row or index >= len(matrix):
            return -1

        if matrix[index][0] > target:
            return Solution.find_row(first_row, index - 1, matrix, target)
        if matrix[index][-1] < target:
            return Solution.find_row(index + 1, last_row, matrix, target)
        else:
            return index

    @staticmethod
    def find_in_row(start: int, end: int, row: list[int], target: int) -> bool:
        """
        use binary search in a given row to find a target.
        """
        index = (start + end) // 2

        if index >= len(row) or end < start:
            return False

        if row[index] > target:
            return Solution.find_in_row(start, index - 1, row, target)
        elif row[index] == target:
            return True
        else:
            return Solution.find_in_row(index + 1, end, row, target)


if __name__ == "__main__":
    matrix = [[1, 3, 5, 7], [10, 11, 16, 20], [23, 30, 34, 60]]
    target = 13

    assert Solution().search_matrix(matrix, target) is False

    matrix = [[1, 3, 5, 7], [10, 11, 16, 20], [23, 30, 34, 60]]
    target = 3

    assert Solution().search_matrix(matrix, target) is True

    matrix = [[1, 3, 5, 7], [10, 11, 16, 20], [23, 30, 34, 60]]
    target = 11

    assert Solution().search_matrix(matrix, target) is True

    matrix = [[1]]
    target = 1

    assert Solution().search_matrix(matrix, target) is True

    matrix = [[1]]
    target = 2

    assert Solution().search_matrix(matrix, target) is False

    matrix = [[1], [3]]
    target = 3

    assert Solution().search_matrix(matrix, target) is True
