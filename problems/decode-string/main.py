def evaluate(expr: str) -> str:
    current_expr = ""
    string = ""
    i = 0
    numbers_stack = []
    strings_stack = []

    while i < len(expr):
        # print(f"expr: {expr}, string: {string}, at i {i}: {expr[i]}")

        if (current_expr + expr[i]).isnumeric():
            current_expr += expr[i]
        elif current_expr != "":
            number = int(current_expr)
            numbers_stack.insert(0, number)
            current_expr = expr[i]
        else:
            current_expr = expr[i]

        if current_expr == "[":
            strings_stack.insert(0, string)
            string = ""
            current_expr = ""
        elif current_expr == "]":
            number = numbers_stack.pop(0)
            string *= number
            string = strings_stack.pop(0) + string
            current_expr = ""
        elif current_expr.isalpha():
            string += current_expr
            current_expr = ""

        i += 1

    return string


if __name__ == "__main__":
    cases = [
        ("ab3[a]2[bc]", "abaaabcbc"),
        ("3[3[a]]", "aaaaaaaaa"),
        ("3[a2[c]]", "accaccacc"),
    ]

    for case in cases:
        assert evaluate(case[0]) == case[1]
