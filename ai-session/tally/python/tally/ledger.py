"""The ledger holds a group's members, expenses and payments.

Balances are kept up to date as entries are recorded. A positive balance means
the group owes that member money; a negative one means the member owes the group.
"""

from dataclasses import dataclass


class LedgerError(Exception):
    pass


@dataclass(frozen=True)
class Expense:
    id: int
    payer: str
    amount: int
    shares: dict[str, int]
    description: str


@dataclass(frozen=True)
class Payment:
    sender: str
    receiver: str
    amount: int


class Ledger:
    def __init__(self) -> None:
        self._members: list[str] = []
        self._expenses: list[Expense] = []
        self._payments: list[Payment] = []
        self._balances: dict[str, int] = {}

    def add_member(self, name: str) -> None:
        if name in self._balances:
            raise LedgerError(f"member already exists: {name}")
        self._members.append(name)
        self._balances[name] = 0

    def members(self) -> list[str]:
        return list(self._members)

    def add_expense(
        self, payer: str, amount: int, shares: dict[str, int], description: str
    ) -> Expense:
        self._require(payer, *shares)
        expense = Expense(len(self._expenses) + 1, payer, amount, shares, description)
        self._expenses.append(expense)
        self._apply(expense, 1)
        return expense

    def delete_expense(self, expense_id: int) -> Expense:
        for expense in self._expenses:
            if expense.id == expense_id:
                self._expenses.remove(expense)
                self._apply(expense, -1)
                return expense
        raise LedgerError(f"no expense #{expense_id}")

    def expenses(self) -> list[Expense]:
        return list(self._expenses)

    def record_payment(self, sender: str, receiver: str, amount: int) -> Payment:
        self._require(sender, receiver)
        if sender == receiver:
            raise LedgerError(f"{sender} cannot pay themselves")
        payment = Payment(sender, receiver, amount)
        self._payments.append(payment)
        self._balances[sender] += amount
        self._balances[receiver] -= amount
        return payment

    def payments(self) -> list[Payment]:
        return list(self._payments)

    def balances(self) -> dict[str, int]:
        return self._balances

    def _apply(self, expense: Expense, sign: int) -> None:
        self._balances[expense.payer] += sign * expense.amount
        for name, share in expense.shares.items():
            self._balances[name] -= sign * share

    def _require(self, *names: str) -> None:
        for name in names:
            if name not in self._balances:
                raise LedgerError(f"unknown member: {name}")
