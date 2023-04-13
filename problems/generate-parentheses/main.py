def generate(n) -> list[str]:
    if n == 0:
        return []
    if n == 1:
        return ["()"]

    res = []
    for i in range(n):
        inner = generate(n - 1 - i)
        outer = generate(i)

        if len(outer) > 0 and len(inner) > 0:
            for iparn in inner:
                for oparn in outer:
                    res.append(f"({iparn}){oparn}")
        elif len(outer) > 0:
            for oparn in outer:
                res.append(f"(){oparn}")
        elif len(inner) > 0:
            for iparn in inner:
                res.append(f"({iparn})")
    return res


if __name__ == "__main__":
    assert generate(1) == ["()"]
    assert generate(3) == ["((()))", "(()())", "(())()", "()(())", "()()()"]
