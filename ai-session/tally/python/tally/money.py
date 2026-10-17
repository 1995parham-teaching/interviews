"""Amounts are kept as integer cents so that every sum is exact."""

import re

_AMOUNT = re.compile(r"[0-9]+(\.[0-9]+)?")


def parse_amount(text: str) -> int:
    """Parse a decimal amount such as "12.50" or "7" into cents."""
    if not _AMOUNT.fullmatch(text):
        raise ValueError(f"invalid amount: {text}")
    if "." in text:
        whole, fraction = text.split(".")
        cents = int(whole) * 100 + int(fraction)
    else:
        cents = int(text) * 100
    if cents == 0:
        raise ValueError(f"amount must be positive: {text}")
    return cents


def format_amount(cents: int, signed: bool = False) -> str:
    """Format cents as a decimal amount; with signed, positive values get a "+"."""
    sign = ""
    if cents < 0:
        sign = "-"
    elif signed and cents > 0:
        sign = "+"
    cents = abs(cents)
    return f"{sign}{cents // 100}.{cents % 100:02d}"
