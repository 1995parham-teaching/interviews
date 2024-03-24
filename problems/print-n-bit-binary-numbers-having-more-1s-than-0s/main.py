def n_bit_binaries(n: int) -> list[str]:
    if n == 1:
        return ["1"]

    result: list[str] = []
    binaries = n_bit_binaries(n - 1)
    for binary in binaries:
        nunber_of_1s = binary.count("1")
        nunber_of_0s = binary.count("0")

        if nunber_of_1s - nunber_of_0s > 0:
            result.append(binary + "0")
        result.append(binary + "1")

    return sorted(result, reverse=True)


if __name__ == "__main__":
    assert n_bit_binaries(1) == ["1"]
    assert n_bit_binaries(2) == ["11", "10"]
    assert n_bit_binaries(3) == ["111", "110", "101"]
