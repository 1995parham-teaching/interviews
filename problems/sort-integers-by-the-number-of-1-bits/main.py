def number_of_1_bits(n: int) -> int:
    ones = 0
    while n >= 2:
        ones += n % 2
        n //= 2
    ones += n
    return ones


if __name__ == "__main__":
    assert number_of_1_bits(4) == 1
    assert number_of_1_bits(7) == 3
