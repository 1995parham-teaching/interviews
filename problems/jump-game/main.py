from typing import List


class Solution:
    def jumpOver(self, nums: List[int]) -> bool:
        for i in range(len(nums)):
            if nums[i] > len(nums) - i:
                return True
        return False

    def canJump(self, nums: List[int]) -> bool:
        if len(nums) == 1:
            return True
        for i in range(len(nums) - 1):
            if nums[i] == 0:
                if not self.jumpOver(nums[:i]):
                    return False

        return True


print(Solution().canJump([2,0,0]))
