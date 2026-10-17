"""Turn a split specification into the amount each member owes for an expense.

A specification has a kind and a list of members:

    equal:ali,sara,reza            everyone owes the same
    shares:ali=1,sara=2            owed in proportion to the weights
    exact:ali=20.00,sara=30.00     the amounts are given directly
"""

import re

from .money import format_amount, parse_amount


def parse_split(spec: str, amount: int) -> dict[str, int]:
    kind, sep, rest = spec.partition(":")
    if not sep or not rest:
        raise ValueError(f"invalid split: {spec}")
    if kind == "equal":
        return equal(amount, rest.split(","))
    if kind == "shares":
        return by_shares(amount, {name: _parse_weight(w) for name, w in _pairs(rest)})
    if kind == "exact":
        return exact(amount, {name: parse_amount(a) for name, a in _pairs(rest)})
    raise ValueError(f"unknown split kind: {kind}")


def equal(amount: int, names: list[str]) -> dict[str, int]:
    each = amount // len(names)
    return {name: each for name in names}


def by_shares(amount: int, weights: dict[str, int]) -> dict[str, int]:
    total = sum(weights.values())
    return {name: amount * weight // total for name, weight in weights.items()}


def exact(amount: int, parts: dict[str, int]) -> dict[str, int]:
    total = sum(parts.values())
    if total > amount:
        raise ValueError(
            f"exact parts add up to {format_amount(total)}, more than {format_amount(amount)}"
        )
    return dict(parts)


def _pairs(text: str) -> list[tuple[str, str]]:
    pairs = []
    for item in text.split(","):
        name, sep, value = item.partition("=")
        if not sep or not name or not value:
            raise ValueError(f"expected name=value, got: {item}")
        pairs.append((name, value))
    return pairs


def _parse_weight(text: str) -> int:
    if not re.fullmatch(r"[0-9]+", text) or int(text) == 0:
        raise ValueError(f"invalid weight: {text}")
    return int(text)
