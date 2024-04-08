def longest_palindromic_substring(s):
    """
    Return lognest plandromic substring.
    """
    if len(s) == 0:
        return ""

    lg = len(s)
    res = s[0]
    for j in range(2):
        offset = 0
        candidate = ""
        i = 0
        while i + j < lg:
            # When j is 0, the following condition is obvious and in case
            # of 1 it causes moving on when two neighbor cells are equal.
            if s[i] != s[i + j]:
                i += 1
                continue

            if (
                i - offset >= 0
                and i + j + offset < lg
                and s[i - offset] == s[i + j + offset]
            ):
                candidate = s[i - offset : i + j + offset + 1]
                offset += 1
            else:
                i += 1
                if len(candidate) > len(res):
                    res = candidate
                offset = 0
                candidate = ""

    return res


if __name__ == "__main__":
    assert longest_palindromic_substring("babad") == "bab"
    assert longest_palindromic_substring("cbbd") == "bb"
    assert longest_palindromic_substring("bb") == "bb"
