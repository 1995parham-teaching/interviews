def find_duplicates(nums: list[int]) -> list[int]:
    result = []
    for i in range(len(nums)):
        index = abs(nums[i]) - 1
        if nums[index] > 0:
            nums[index] = -nums[index]
        else:
            result.append(abs(nums[i]))
    return result


if __name__ == "__main__":
    assert find_duplicates([1, 1, 2, 2]) == [1, 2]
    assert find_duplicates([1, 2, 3, 4]) == []
    assert find_duplicates([1, 2, 2, 4]) == [2]
