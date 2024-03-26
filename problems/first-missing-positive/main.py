def first_missing_positive(nums: list[int]) -> int:
    """
    Use the same array to mark numbers that we seen.
    By marking, I mean converting them into the negative number,
    for example when I see 4, I convert nums[3] to negative and
    this will indicates we have 4 in our array.
    """

    # by converting negative and zero numbers
    # to big enough positive numbers we remove
    # them from our algorithm.
    for i in range(len(nums)):
        if nums[i] == 0:
            nums[i] = len(nums) + 1
        elif nums[i] < 0:
            nums[i] = len(nums) + 1

    sp = 1

    for i in range(len(nums)):
        if sp == abs(nums[i]):
            sp += 1
            while sp <= len(nums) and nums[sp - 1] < 0:
                sp += 1
        index = abs(nums[i]) - 1
        if index < len(nums) and nums[index] > 0:
            nums[index] *= -1
    return sp


if __name__ == "__main__":
    assert first_missing_positive([1, 2, 0]) == 3
    assert first_missing_positive([3, 4, -1, 1]) == 2
    assert first_missing_positive([7, 8, 9, 11, 12]) == 1
    assert first_missing_positive([2, 1, 3, 4, 4, 3, 1, 1]) == 5
