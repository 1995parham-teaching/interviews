class Solution:

    def myLongestCommonPrefix(self, strs) -> str:
        m = {}
        for s in strs:
            if m.get(s[0]) is None:
                m[s[0]] = 1
            else:
                m[s[0]] = m[s[0]] + 1

        maxValue = 0
        max_key = ""

        for k, v in m.items():
            if v > maxValue:
                maxValue = v
                max_key = k

        if maxValue == 1:
            return ""

        new_strs = []
        for s in strs:
            if s[0] == max_key:
                new_strs.append(s[1:])

        if len(new_strs) < len(strs):
            return ""

        return max_key + self.myLongestCommonPrefix(new_strs)

    def longestCommonPrefix(self, strs) -> str:
        m = {}
        for s in strs:
            if m.get(s[0]) is None:
                m[s[0]] = 1
            else:
                m[s[0]] = m[s[0]] + 1

        maxValue = 0
        max_key = ""

        for k, v in m.items():
            if v > maxValue:
                maxValue = v
                max_key = k

        if maxValue == 1:
            return ""

        new_strs = []
        for s in strs:
            if s[0] == max_key:
                new_strs.append(s[1:])

        if len(new_strs) < len(strs):
            return ""

        return max_key + self.myLongestCommonPrefix(new_strs)


s = Solution()
print(s.longestCommonPrefix(["flower", "flow", "flight"]))

# a = {"A": 0, "B": 1}
# for k, v in a.items():
#     print(k)
#     print(v)
# if a.get("A") == None:
#     print("hi")
