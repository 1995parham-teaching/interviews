def number_of_1_bits(n: int) -> tuple[int, int]:
    input = n
    ones = 0
    while n >= 2:
        ones += n % 2
        n //= 2
    ones += n
    return (ones, input)


if __name__ == "__main__":
    assert number_of_1_bits(4) == (1, 4)
    assert number_of_1_bits(7) == (3, 7)

    assert sorted([0, 1, 2, 3, 4, 5, 6, 7, 8], key=number_of_1_bits) == [
        0,
        1,
        2,
        4,
        8,
        3,
        5,
        6,
        7,
    ]
    assert sorted(
        [1024, 512, 256, 128, 64, 32, 16, 8, 4, 2, 1], key=number_of_1_bits
    ) == [
        1,
        2,
        4,
        8,
        16,
        32,
        64,
        128,
        256,
        512,
        1024,
    ]
