def atoi(input: str) -> int:
    index = 0
    number = 0

    while index < len(input) and input[index] == " ":
        index += 1

    if index == len(input):
        return 0

    match input[index]:
        case "+":
            sign = 1
            index += 1
        case "-":
            sign = -1
            index += 1
        case _:
            sign = 1

    while index < len(input):
        match input[index]:
            case "0":
                number = number * 10 + 0
            case "1":
                number = number * 10 + 1
            case "2":
                number = number * 10 + 2
            case "3":
                number = number * 10 + 3
            case "4":
                number = number * 10 + 4
            case "5":
                number = number * 10 + 5
            case "6":
                number = number * 10 + 6
            case "7":
                number = number * 10 + 7
            case "8":
                number = number * 10 + 8
            case "9":
                number = number * 10 + 9
            case _:
                return sign * number
        if sign == 1 and number >= 2**31 - 1:
            return 2**31 - 1
        if sign == -1 and -number <= -(2**31):
            return -(2**31)
        index += 1

    return sign * number


if __name__ == "__main__":
    assert atoi("") == 0
    assert atoi(" ") == 0
    assert atoi("-1") == -1
    assert atoi("   -12abac") == -12
    assert atoi("   -12") == -12
    assert atoi("   12") == 12
    assert atoi("   12   ") == 12
