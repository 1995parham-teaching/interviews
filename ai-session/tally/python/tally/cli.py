"""Read commands, one per line, and run them against a single ledger."""

import re
import sys
from typing import Iterable, TextIO

from . import settle
from .ledger import Ledger, LedgerError
from .money import format_amount, parse_amount
from .split import parse_split

USAGE = """commands:
  member <name>
  expense <payer> <amount> <split> <description>
  pay <sender> <receiver> <amount>
  delete <expense-id>
  balances
  settle"""


def run(lines: Iterable[str], out: TextIO) -> int:
    """Run every command and return the number of commands that failed."""
    ledger = Ledger()
    failures = 0
    for number, raw in enumerate(lines, start=1):
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        try:
            execute(ledger, line, out)
        except (LedgerError, ValueError) as err:
            print(f"error: line {number}: {err}", file=out)
            failures += 1
    return failures


def execute(ledger: Ledger, line: str, out: TextIO) -> None:
    command, *args = line.split()
    if command == "member":
        _arity(args, 1, "member <name>")
        ledger.add_member(args[0])
    elif command == "expense":
        if len(args) < 4:
            raise ValueError("usage: expense <payer> <amount> <split> <description>")
        payer, amount_text, spec = args[:3]
        amount = parse_amount(amount_text)
        description = " ".join(args[3:])
        expense = ledger.add_expense(payer, amount, parse_split(spec, amount), description)
        print(
            f"#{expense.id} {expense.description}: "
            f"{expense.payer} paid {format_amount(expense.amount)}",
            file=out,
        )
    elif command == "pay":
        _arity(args, 3, "pay <sender> <receiver> <amount>")
        payment = ledger.record_payment(args[0], args[1], parse_amount(args[2]))
        print(
            f"{payment.sender} paid {payment.receiver} {format_amount(payment.amount)}",
            file=out,
        )
    elif command == "delete":
        _arity(args, 1, "delete <expense-id>")
        if not re.fullmatch(r"[0-9]+", args[0]):
            raise ValueError(f"invalid expense id: {args[0]}")
        expense = ledger.delete_expense(int(args[0]))
        print(f"deleted #{expense.id} {expense.description}", file=out)
    elif command == "balances":
        _arity(args, 0, "balances")
        _print_balances(ledger.balances(), out)
    elif command == "settle":
        _arity(args, 0, "settle")
        transfers = settle.plan(ledger.balances())
        if not transfers:
            print("nothing to settle", file=out)
        for t in transfers:
            print(f"{t.sender} -> {t.receiver} {format_amount(t.amount)}", file=out)
    else:
        raise ValueError(f"unknown command: {command}\n{USAGE}")


def _print_balances(balances: dict[str, int], out: TextIO) -> None:
    width = max((len(name) for name in balances), default=0)
    for name in sorted(balances):
        print(f"{name.ljust(width)}  {format_amount(balances[name], signed=True)}", file=out)


def _arity(args: list[str], n: int, usage: str) -> None:
    if len(args) != n:
        raise ValueError(f"usage: {usage}")


def main(argv: list[str]) -> int:
    if len(argv) > 2:
        print("usage: python -m tally [file]", file=sys.stderr)
        return 2
    if len(argv) == 2:
        with open(argv[1], encoding="utf-8") as f:
            failures = run(f, sys.stdout)
    else:
        failures = run(sys.stdin, sys.stdout)
    return 1 if failures else 0
