class Solution:
    def longestOddPalindrome(self, s):
        res = ""
        for i in range(len(s)):
            r = s[i]
            for d in range(min(len(s) - i, i + 1)):
                if s[i - d] != s[i + d]:
                    break

                r = s[i - d:i + d + 1]

            if len(r) > len(res):
                res = r

        return res

    def longestEvenPalindrome(self, s):
        res = ""
        r = ""
        for i in range(len(s) - 1):
            if s[i] == s[i + 1]:
                r = s[i: i + 2]
                for d in range(min(len(s) - i - 2, i) + 1):
                    if s[i - d] != s[i + d + 1]:
                        break

                    r = s[i - d:i + d + 2]

            if len(r) > len(res):
                res = r

        return res

    def longestPalindrome(self, s: str) -> str:
        ro = self.longestOddPalindrome(s)
        re = self.longestEvenPalindrome(s)

        return ro if len(ro) > len(re) else re


r = Solution().longestPalindrome("babad")
print(r)
